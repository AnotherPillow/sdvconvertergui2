package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bodgit/sevenzip"
	"golang.org/x/sys/windows/registry"
)

func getRegistryValue(area registry.Key, path string, value string) string {
	k, err := registry.OpenKey(area, path, registry.QUERY_VALUE)
	if err != nil { // `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 413150` registry.LOCAL_MACHINE
		return ""
	}
	defer k.Close()

	s, _, err := k.GetStringValue(value) // "InstallLocation"
	if err != nil {
		return ""
	}
	return s
}

var home, _ = os.UserHomeDir()

// Game path lists taken from https://github.com/Pathoschild/SMAPI/blob/develop/build/find-game-folder.targets

var possibleUnixOSXGamePaths = []string{
	filepath.Join(home, "GOG Games/Stardew Valley/game"),
	filepath.Join(home, ".steam/steam/steamapps/common/Stardew Valley"),
	filepath.Join(home, ".local/share/Steam/steamapps/common/Stardew Valley"),
	filepath.Join(home, ".var/app/com.valvesoftware.Steam/data/Steam/steamapps/common/Stardew Valley"),
	"/Applications/Stardew Valley.app/Contents/MacOS",
	filepath.Join(home, "Library/Application Support/Steam/steamapps/common/Stardew Valley/Contents/MacOS"),
}

var possibleWindowsGamePaths = []string{
	getRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\GOG.com\Games\1453375253`, "PATH"),
	getRegistryValue(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 413150`, "InstallLocation"),
	filepath.Join(getRegistryValue(registry.CURRENT_USER, `SOFTWARE\Valve\Steam`, "SteamPath"), "steamapps\\common\\Stardew Valley"),
	`C:\Program Files\GalaxyClient\Games\Stardew Valley`,
	`C:\Program Files\GOG Galaxy\Games\Stardew Valley`,
	`C:\Program Files\GOG Games\Stardew Valley`,
	`C:\Program Files (x86)\GalaxyClient\Games\Stardew Valley`,
	`C:\Program Files (x86)\GOG Galaxy\Games\Stardew Valley`,
	`C:\Program Files (x86)\GOG Games\Stardew Valley`,
	`C:\Program Files\Steam\steamapps\common\Stardew Valley`,
	`C:\Program Files (x86)\Steam\steamapps\common\Stardew Valley`,
	`C:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`D:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`E:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`F:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`G:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`H:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`I:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`J:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`K:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`L:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`M:\Program Files\ModifiableWindowsApps\Stardew Valley`,
	`C:\GOG Games\Stardew Valley`,
}

func getGameInstallDirectory() string {
	var paths = make([]string, 0)

	if runtime.GOOS == "windows" {
		paths = possibleWindowsGamePaths
	} else {
		paths = possibleUnixOSXGamePaths
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil && path != "" {
			return path
		}
	}

	return ""
}

func getLocalFilesRootDirectory() string {
	// ~/.local/share
	var outDir = home

	if runtime.GOOS == "windows" {
		outDir = filepath.Join(outDir, "Appdata\\Local\\sdvconvertergui2")
	} else {
		outDir = filepath.Join(outDir, ".local/share/sdvconvertergui2")
	}

	return outDir
}

func getPythonDownloadLink() string {
	if runtime.GOOS == "windows" {
		return "https://www.python.org/ftp/python/3.11.8/python-3.11.8-embed-amd64.zip"
	} else if runtime.GOOS == "darwin" {
		// TODO: Can't find portable build for mac
		return ""
	} else {
		// TODO: Can't find portable build for linux
		return ""
	}
}

func getGitDownloadLink() string {
	if runtime.GOOS == "windows" {
		return "https://github.com/git-for-windows/git/releases/download/v2.44.0.windows.1/PortableGit-2.44.0-64-bit.7z.exe"
	} else if runtime.GOOS == "darwin" {
		// TODO: Can't find portable build for mac
		return ""
	} else {
		// TODO: Can't find portable build for linux
		return ""
	}
}

var localFilesRootDirectory = getLocalFilesRootDirectory()

var gitDownloadLink = getGitDownloadLink()
var pythonDownloadLink = getPythonDownloadLink()

var gitDownloadFilePath = filepath.Join(localFilesRootDirectory, "PortableGit-2.44.0-64-bit.7z.exe")
var pythonDownloadFilePath = filepath.Join(localFilesRootDirectory, "python-3.11.8-embed-amd64.zip")

var gitFolderPath = filepath.Join(localFilesRootDirectory, "git")
var pythonFolderPath = filepath.Join(localFilesRootDirectory, "python-311")

func getPythonExecutable() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(pythonFolderPath, "python.exe")
	} else if runtime.GOOS == "darwin" {
		// TODO
		return ""
	} else {
		// TODO
		return ""
	}
}

func getGitExecutable() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(gitFolderPath, "bin", "git.exe")
	} else if runtime.GOOS == "darwin" {
		// TODO
		return ""
	} else {
		// TODO
		return ""
	}
}

var pyExecutable = getPythonExecutable()
var gitExecutable = getGitExecutable()

func checkExists(path string) bool {
	_, error := os.Stat(path)
	return !os.IsNotExist(error)
}

// https://stackoverflow.com/a/33853856
func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func checkGameVersion(gamepath string) string {
	var doesExist = checkExists(filepath.Join(gamepath, "Content/Data/Objects.xnb"))
	if doesExist {
		return "1.6"
	} else {
		return "1.5"
	}
}

// https://stackoverflow.com/a/24792688
func extractZipFile(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractSelfExtracing7z(src string, dest string) error {
	r, err := sevenzip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		var destFolder = filepath.Join(dest, filepath.Dir(file.Name))
		if !checkExists(destFolder) {
			println("creating " + destFolder)
			e := os.MkdirAll(destFolder, os.ModePerm)
			if e != nil {
				println("failed to create folder " + destFolder + e.Error())
			}
		}
		var destFilePath = filepath.Join(dest, file.Name)
		// println(destFilePath)

		var rc, _ = file.Open()

		dest, e := os.Create(destFilePath)

		if e != nil {
			println("Failed to create " + destFilePath + e.Error())
		}

		defer dest.Close()

		_, e2 := io.Copy(dest, rc)
		if e2 != nil {
			println("Failed to copy file from git archive! " + e.Error())
		}

		_ = rc.Close()

	}

	return nil
}

func createFolderIfNeeded(fp string) string {
	if !checkExists(fp) {
		os.MkdirAll(fp, os.ModePerm)
	}

	return fp
}

func installPip() {
	var getpip = filepath.Join(pythonFolderPath, "get-pip.py")
	err := downloadFile(getpip, "https://bootstrap.pypa.io/get-pip.py")

	if err != nil {
		println("Failed to download pip: " + err.Error())
	}

	var pipInstallCmd = runSimpleCommand(fmt.Sprintf("& \"%s\" \"%s\"",
		pyExecutable, getpip,
	))

	if err := pipInstallCmd.Start(); err != nil {
		println("get-pip.py failed to start with error " + err.Error())
	}

	if err := pipInstallCmd.Wait(); err != nil {
		println("get-pip.py completed with error " + err.Error())
	}

	var pthContent = []byte(`
python311.zip
.

# Uncomment to run site.main() automatically
#import site
./Scripts
./Lib/site-packages`)

	os.WriteFile(filepath.Join(pythonFolderPath, "python311._pth"),
		pthContent, os.ModePerm)

}

func showFolder(path string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command(`explorer`, `/select,`, path)
		_ = cmd.Start()
		return
	} else if runtime.GOOS == "darwin" {
		// TODO
		return
	} else {
		// TODO
		return
	}
}

// https://gosamples.dev/zip-file/
func zipFolder(source string, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func updatePackageResolution(mf string) error {
	var text = `
# --- START ADDED BY SDVCONVERTERGUI2 ---
import os, sys
SCRIPT_DIR = (os.path.abspath(__file__))
sys.path.append(os.path.dirname(SCRIPT_DIR))
# --- END ADDED BY SDVCONVERTERGUI2 ---
`
	file, err := os.OpenFile(mf, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	existingContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	existingString := string(existingContent)

	if strings.Contains(existingString, "SCRIPT_DIR = (os") {
		return nil
	}

	updatedContent := text + existingString

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.WriteString(updatedContent)
	if err != nil {
		return err
	}

	err = file.Truncate(int64(len(updatedContent)))
	if err != nil {
		return err
	}

	return nil
}
