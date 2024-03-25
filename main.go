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
	rand.Seed(time.Now().UnixNano())

	createFolderIfNeeded(localFilesRootDirectory)

	if !checkExists(pythonDownloadFilePath) {
		println("Downloading python!")
		var err = downloadFile(pythonDownloadFilePath, pythonDownloadLink)
		if err != nil {
			println(err.Error())
			println("Failed to download Python!")
			return
		}

		var e = extractZipFile(pythonDownloadFilePath, pythonFolderPath)
		if e != nil {
			println(e.Error())
			println("Failed to extract Python!")
			return
		}

		installPip()
	}

	if !checkExists(gitDownloadFilePath) {
		println("Downloading git!")
		var err = downloadFile(gitDownloadFilePath, gitDownloadLink)
		if err != nil {
			println(err.Error())
		}
		var e = extractSelfExtracing7z(gitDownloadFilePath, gitFolderPath)
		if e != nil {
			println(e.Error())
		}
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "sdvconvertergui2",
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
		println("Error:", err.Error())
	}

}
