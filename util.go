package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"runtime"
)

const APP_VERSION = "1.1.1"
const APP_UNIQUEID = "fake.nexus.20986"
const APP_UPDATEKEY = "Nexus:20986"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func runSimpleCommand(cmd string) *exec.Cmd {
	debugLog("Running " + cmd)
	if runtime.GOOS == "windows" {
		return exec.Command(
			"powershell",
			"-NoProfile",
			"-Command",
			cmd,
		)
	} else {
		// TODO
		return nil
	}
}

func sendPOST(url string, data interface{}) ([]map[string]interface{}, error) {
	// Marshal the JSON data
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create a request object
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create a HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	debugLog(fmt.Sprintf("Response from %s with %s", url, string(responseBody)))

	// Parse the JSON response array
	var responseArray []map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseArray); err != nil {
		return nil, fmt.Errorf("error parsing JSON response: %v", err)
	}

	return responseArray, nil
}

func checkForUpdate(a *App) (int, string) {
	reqMods := []map[string]interface{}{{
		"id":               "fake.nexus.20986",
		"updateKeys":       []string{"Nexus:20986"},
		"installedVersion": APP_VERSION,
	}}
	reqData := map[string]interface{}{
		"mods":       reqMods,
		"apiVersion": "3.0.0",
	}
	var updates, err = sendPOST("https://smapi.io/api/v3.0/mods", reqData)

	if err != nil {
		return 2, "Failed to check for an update."
	} else if len(updates) == 0 {
		return 0, "No update available."
	} else if updates[0]["suggestedUpdate"] != nil {
		return 1, fmt.Sprintf("There is a new update - %s - available. ", updates[0]["suggestedUpdate"].(map[string]interface{})["version"].(string))
	} else {
		return 0, "No update available."
	}
}
