package main

import (
	"fmt"
	"os"
	"strings"
)

func CheckCompatibleGame(converter Converter) bool {
	var is16 = checkGameVersion(getGameInstallDirectory()) == "1.6"

	if converter.Needs16 && is16 {
		return true
	} else if converter.Needs16 && !is16 {
		return false
	}
	return true
}

type Converter struct {
	Name              string
	Needs16           bool
	Repo              string
	GitFile           string
	InputDirectory    string
	OutputDirectory   string
	MainFile          string
	RequirementsFile  string
	SupportedUniqueID string
	ExtraArgs         string
	IsPython          bool
}

func (c Converter) ModifyConfig(configPath string) {
	if c.Name == "TMXL2CP" {
		var bytes = []byte(fmt.Sprintf(`{"game_folder": "%s","ran_before": false}`, strings.ReplaceAll(
			getGameInstallDirectory(), "\\", "\\\\")))
		os.WriteFile(configPath, bytes, 0644)
	} else if c.Name == "CP2AT" {
		var bytes = []byte(`{"mod_folder_path": "input","keywords": [""], "output_folder_path": "output"}`)
		os.WriteFile(configPath, bytes, 0644)
	}
}

func (c Converter) SupportsManifest(manifest map[string]interface{}) bool {
	var contentPackFor = manifest["ContentPackFor"].(map[string]interface{})["UniqueID"].(string)

	debugLog(fmt.Sprintf("checking if CPF %s is compatible with UID %s", contentPackFor, c.SupportedUniqueID))
	return strings.ToLower(contentPackFor) == strings.ToLower(c.SupportedUniqueID)
}

var TMXL2CP = Converter{
	Name:              "TMXL2CP",
	Needs16:           true,
	Repo:              "https://github.com/AnotherPillow/TMXL2CP",
	GitFile:           "https://github.com/AnotherPillow/TMXL2CP.git",
	InputDirectory:    "TMXL",
	OutputDirectory:   "CP",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "platonymous.tmxloader",
	ExtraArgs:         "",
	IsPython:          true,
}

var BFAV2CP = Converter{
	Name:              "BFAV2CP",
	Needs16:           false,
	Repo:              "https://github.com/AnotherPillow/BFAV2CP",
	GitFile:           "https://github.com/AnotherPillow/BFAV2CP.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "paritee.betterfarmanimalvariety",
	ExtraArgs:         "",
	IsPython:          true,
}

var CM2CP = Converter{
	Name:              "CM2CP",
	Needs16:           false,
	Repo:              "https://github.com/AnotherPillow/CM2CP",
	GitFile:           "https://github.com/AnotherPillow/CM2CP.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "platonymous.custommusic",
	ExtraArgs:         "",
	IsPython:          true,
}

var STF2CP = Converter{
	Name:              "STF2CP",
	Needs16:           false,
	Repo:              "https://github.com/AnotherPillow/STF2CP",
	GitFile:           "https://github.com/AnotherPillow/STF2CP.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "cherry.shoptileframework",
	ExtraArgs:         "",
	IsPython:          true,
}

var FurnitureConverter = Converter{
	Name:              "FurnitureConverter",
	Needs16:           false,
	Repo:              "https://github.com/elizabethcd/FurnitureConverter",
	GitFile:           "https://github.com/elizabethcd/FurnitureConverter.git",
	InputDirectory:    "input",
	OutputDirectory:   "output-1.6",
	MainFile:          "furniture_converter.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "platonymous.customfurniture",
	ExtraArgs:         "--inputDir input --outputDir output",
	IsPython:          true,
}

var CP2AT = Converter{
	Name:              "CP2AT",
	Needs16:           false,
	Repo:              "https://github.com/holy-the-sea/CP2AT",
	GitFile:           "https://github.com/holy-the-sea/CP2AT.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "pathoschild.contentpatcher",
	ExtraArgs:         "",
	IsPython:          true,
}

var CPA2SC = Converter{
	Name:              "CPA2SC",
	Needs16:           false,
	Repo:              "https://github.com/AnotherPillow/CPA2SC",
	GitFile:           "https://github.com/AnotherPillow/CPA2SC.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "pathoschild.contentpatcher",
	ExtraArgs:         "",
	IsPython:          true,
}

var SAAT2CP = Converter{
	Name:              "SAAT2CP",
	Needs16:           false,
	Repo:              "https://github.com/AnotherPillow/SAAT2CP",
	GitFile:           "https://github.com/AnotherPillow/SAAT2CP.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "zerometers.saat.mod",
	ExtraArgs:         "",
	IsPython:          true,
}

var MTN2CP = Converter{
	Name:              "SAAT2CP",
	Needs16:           false,
	Repo:              "https://github.com/AnotherPillow/MTN2CP",
	GitFile:           "https://github.com/AnotherPillow/MTN2CP.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "main.py",
	RequirementsFile:  "requirements.txt",
	SupportedUniqueID: "sgtpickles.mtn",
	ExtraArgs:         "",
	IsPython:          true,
}

var ConvertToFashionSenseJsonAssets = Converter{
	Name:              "Convert To Fashion Sense JA ONLY",
	Needs16:           false,
	Repo:              "https://github.com/anotherpillow/stardew-convert-to-fashion-sense",
	GitFile:           "https://github.com/anotherpillow/stardew-convert-to-fashion-sense.git",
	InputDirectory:    "input",
	OutputDirectory:   "output",
	MainFile:          "index.js",
	RequirementsFile:  "",
	SupportedUniqueID: "spacechase0.jsonassets",
	ExtraArgs:         "",
	IsPython:          false,
}

var ConvertersMap = map[string]Converter{
	"TMXL2CP":                          TMXL2CP,
	"BFAV2CP":                          BFAV2CP,
	"CM2CP":                            CM2CP,
	"STF2CP":                           STF2CP,
	"FurnitureConverter":               FurnitureConverter,
	"CP2AT":                            CP2AT,
	"CPA2SC":                           CPA2SC,
	"SAAT2CP":                          SAAT2CP,
	"MTN2CP":                           MTN2CP,
	"Convert To Fashion Sense JA ONLY": ConvertToFashionSenseJsonAssets,
}
