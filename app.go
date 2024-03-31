package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cp "github.com/otiai10/copy"
	"github.com/sqweek/dialog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DomReady(ctx context.Context) {
	debugLog("Dom ready!")
	runtime.EventsEmit(ctx, "AVAILABLE_CONVERTERS", ConvertersMap)
	runtime.EventsEmit(ctx, "DOM_READY", map[string]interface{}{
		"game": getGameInstallDirectory(),
	})

	var updateState, updateMessage = checkForUpdate(a)
	var updateType = ""

	debugLog(fmt.Sprintf("Updatestate: %d, updatemessage: %s", updateState, updateMessage))

	if updateState == 0 {
		updateType = "none"
	} else if updateState == 1 {
		updateType = "available"
	} else {
		updateType = "error"
	}

	var emitData = map[string]interface{}{
		"state":   updateType,
		"message": updateMessage,
	}

	runtime.EventsEmit(ctx, "UPDATE_CHECK_INFO", map[string]interface{}{
		"data": emitData,
	})
}

func (a *App) ChooseManifest() map[string]interface{} {
	filename, err := dialog.File().Filter("Mod Manifest (manifest.json)", "json").Title("Choose Input Content Pack manifest.json").Load()
	if err != nil {
		debugLog("failed to select file")
		return map[string]interface{}{
			"filename": nil,
			"content":  nil,
		}
	}
	manifest, err := os.ReadFile(filename)
	if err != nil {
		debugLog("failed to os.ReadFile")
		return map[string]interface{}{
			"filename": nil,
			"content":  nil,
		}
	}

	return map[string]interface{}{
		"filename": filename,
		"content":  string(manifest),
	}
}

func (a *App) ConvertMod(manifest map[string]interface{}, converterName string, manifestPath string) string {
	var converter = ConvertersMap[converterName]
	if !CheckCompatibleGame(converter) {
		return "error|This converter is not compatible with your game install!"
	}
	if !converter.SupportsManifest(manifest["content"].(map[string]interface{})) {
		return "error|This converter is not compatible with this content pack!"
	}

	var downloadDirName = fmt.Sprintf("%s_%s", converter.Name, randomString(5))
	var downloadDir = filepath.Join(localFilesRootDirectory, downloadDirName)

	debugLog(converter.GitFile)
	runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", fmt.Sprintf("[GUI] Cloning %s", converter.Name))

	var gitCloneCmd = runSimpleCommand(fmt.Sprintf(
		"& \"%s\" clone %s \"%s\" --recurse-submodules", gitExecutable, converter.GitFile, downloadDir,
	))

	if err := gitCloneCmd.Start(); err != nil {
		return "error|git clone failed to start with error " + err.Error()
	}

	// Wait for the command to finish
	if err := gitCloneCmd.Wait(); err != nil {
		return "error|git clone completed with error " + err.Error()
	}

	var modFolder = filepath.Dir(manifestPath)
	var inputFolder = createFolderIfNeeded(filepath.Join(downloadDir, converter.InputDirectory))
	var outputFolder = createFolderIfNeeded(filepath.Join(downloadDir, converter.OutputDirectory))

	err := updatePackageResolution(filepath.Join(downloadDir, converter.MainFile))
	if err != nil {
		debugLog(fmt.Sprintf("failed to update package resolution %s", err.Error()))

		return fmt.Sprintf("error|Failed change directory with error %s", err.Error())
	}

	runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", fmt.Sprintf("[GUI] Copying %s to input directory",
		manifest["content"].(map[string]interface{})["Name"].(string)))
	_ = cp.Copy(modFolder, inputFolder)

	runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", "[GUI] Installing dependencies...")

	var reqsCmd = runSimpleCommand(fmt.Sprintf(
		"& \"%s\" -m pip install -r \"%s\"", pyExecutable, filepath.Join(downloadDir, converter.RequirementsFile),
	))

	reqsStdout, err := reqsCmd.StdoutPipe()
	if err != nil {
		debugLog(fmt.Sprintf("Error creating stdout pipe: %v\n", err))
	}

	if err := reqsCmd.Start(); err != nil {
		var r = "error|Dependencies installation failed to start with error " + err.Error()
		debugLog(r)
		return r
	}

	reqsScanner := bufio.NewScanner(reqsStdout)

	for reqsScanner.Scan() {
		debugLog(reqsScanner.Text())
		if strings.Contains(reqsScanner.Text(), "Successfully installed") {
			runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", "[GUI] "+reqsScanner.Text())
		}
	}

	converter.ModifyConfig(filepath.Join(downloadDir, "config.json"))

	if err := reqsCmd.Wait(); err != nil {
		var r = "error|Dependencies installation completed with error " + err.Error()
		debugLog(fmt.Sprintf("Failed to install dependencies %v // %s", err, err.Error()))
		debugLog(r)
		return r
	}

	runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", fmt.Sprintf("[GUI] Running %s", converter.Name))

	old_wd, _ := os.Getwd()
	var chdir1_err = os.Chdir(downloadDir)
	if chdir1_err != nil {
		return fmt.Sprintf("error|Failed change directory with error %s", err.Error())
	}

	wd, wde := os.Getwd()
	debugLog(fmt.Sprintf("PWD: %s - possible err %s", wd, wde))

	var mainRunCmd = runSimpleCommand(fmt.Sprintf(
		"& \"%s\" -u \"%s\" %s", pyExecutable, converter.MainFile, converter.ExtraArgs,
	))

	mainStdout, err := mainRunCmd.StdoutPipe()
	if err != nil {
		debugLog(fmt.Sprintf("Error creating stdout pipe: %v\n", err))
	}

	mainStdoutScanner := bufio.NewScanner(mainStdout)

	mainSterr, err := mainRunCmd.StderrPipe()
	if err != nil {
		debugLog(fmt.Sprintf("Error creating sterr pipe: %v\n", err))
	}

	mainStderrScanner := bufio.NewScanner(mainSterr)

	if err := mainRunCmd.Start(); err != nil {
		var r = fmt.Sprintf("error|%s failed to start with error %s", converter.Name, err.Error())
		debugLog(r)
		return r
	}

	for mainStdoutScanner.Scan() {
		var txt = mainStdoutScanner.Text()
		debugLog(txt)
		runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", txt)
	}

	for mainStderrScanner.Scan() {
		var txt = mainStderrScanner.Text()
		debugLog(txt)
		runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", "ERR: "+txt)
	}

	if err := mainRunCmd.Wait(); err != nil {
		var r = fmt.Sprintf("error|%s completed with error %s", converter.Name, err.Error())
		debugLog(r)
		return r
	}

	var chdir2_err = os.Chdir(old_wd)
	if chdir2_err != nil {
		var r = fmt.Sprintf("error|Failed change directory with error %s", err.Error())
		debugLog(r)
		return r
	}

	var newOutputFolder = createFolderIfNeeded(filepath.Join(localFilesRootDirectory, "conversions",
		fmt.Sprintf("%s-converted-%s_%s",
			manifest["content"].(map[string]interface{})["UniqueID"].(string), converter.Name, randomString(8),
		),
	))

	showFolder(newOutputFolder)

	_ = cp.Copy(outputFolder, newOutputFolder)

	var zipPath = newOutputFolder + ".zip"
	var zipErr = zipFolder(outputFolder, zipPath)

	if zipErr != nil {
		debugLog("Failed to zip mod")
		debugLog(zipErr.Error())
	} else {
		debugLog(fmt.Sprintf("Completed mod conversion to zip: %s", zipPath))
		runtime.EventsEmit(a.ctx, "CONVERTER_MOD_DONE", zipPath)
	}

	runtime.EventsEmit(a.ctx, "OUTPUT_CONVERTER_TEXT", fmt.Sprintf("[GUI] %s completed!", converter.Name))

	os.RemoveAll(downloadDir) // Clean up afterwards
	return "success|Conversion complete! The converted mod has been opened in explorer.."
}

func (a *App) ShowFolderInExplorer(path string) {
	showFolder(path)
}
