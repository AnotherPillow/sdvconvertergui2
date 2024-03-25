package main

import (
	"math/rand"
	"os/exec"
	"runtime"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func runSimpleCommand(cmd string) *exec.Cmd {
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
