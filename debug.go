package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logDirectory = createFolderIfNeeded(filepath.Join(getLocalFilesRootDirectory(), "debug-logs"))

func generateLogFilename() string {
	var now = time.Now()

	var name = fmt.Sprintf("%s %s.log", now.Format("02-01-06 15;04;05"), randomString(4))

	return name
}

func getCurrentTime() string {
	var now = time.Now()
	return now.Format("15;04;05")
}

var logFilename = generateLogFilename()
var logPath = filepath.Join(logDirectory, logFilename)

func initDebugLogging() {

	os.WriteFile(logPath,
		[]byte(fmt.Sprintf("Started logging at %s", getCurrentTime())), os.ModePerm)
}

func debugLog(rawmsg string) {
	var message = fmt.Sprintf("%s: %s\n", getCurrentTime(), rawmsg)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		log.Fatal(err)
	}
}
