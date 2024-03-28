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

	return strings.ToLower(contentPackFor) == strings.ToLower(c.SupportedUniqueID)
}

// https://github.com/AnotherPillow/TMXL2CP
// https://github.com/AnotherPillow/BFAV2CP
// https://github.com/AnotherPillow/CM2CP
// https://github.com/AnotherPillow/STF2CP
// https://github.com/elizabethcd/FurnitureConverter
// https://github.com/elizabethcd/SkirtConverter
// https://github.com/holy-the-sea/CP2AT

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
}

var ConvertersMap = map[string]Converter{
	"TMXL2CP":            TMXL2CP,
	"BFAV2CP":            BFAV2CP,
	"CM2CP":              CM2CP,
	"STF2CP":             STF2CP,
	"FurnitureConverter": FurnitureConverter,
	"CP2AT":              CP2AT,
}
