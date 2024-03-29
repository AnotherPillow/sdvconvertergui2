package main

import (
	"embed"
	"math/rand"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	initDebugLogging()
	rand.Seed(time.Now().UnixNano())
	debugLog("Started application...")
	debugLog("Application version: " + APP_VERSION)

	createFolderIfNeeded(localFilesRootDirectory)

	if !checkExists(pythonDownloadFilePath) {
		debugLog("Downloading python.")
		var err = downloadFile(pythonDownloadFilePath, pythonDownloadLink)
		if err != nil {
			debugLog("Failed to download Python! " + err.Error())
			return
		}

		var e = extractZipFile(pythonDownloadFilePath, pythonFolderPath)
		if e != nil {
			debugLog("Failed to extract Python! " + e.Error())
			return
		}

		debugLog("Installing pip.")
		installPip()
	} else {
		debugLog("Python already downloaded.")
	}

	if !checkExists(gitDownloadFilePath) {
		debugLog("Downloading git!")
		var err = downloadFile(gitDownloadFilePath, gitDownloadLink)
		if err != nil {
			debugLog("Failed to download Git! " + err.Error())
		}
		var e = extractSelfExtracing7z(gitDownloadFilePath, gitFolderPath)
		if e != nil {
			debugLog("Failed to extract Git! " + e.Error())
		}
	} else {
		debugLog("Git already downloaded.")
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "sdvconvertergui2 " + APP_VERSION,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnDomReady:       app.DomReady,
		OnStartup:        app.startup,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		debugLog("Failed to run/create wails app/window " + err.Error())
	}

}
