package core

import (
	"os"
	"path/filepath"
	"runtime"
)

const DefaultFileName = "ikasbox.db"

func getHome() string {
	home := "HOME"
	if runtime.GOOS == "windows" {
		home = "USERPROFILE"
	}
	return os.Getenv(home)
}

func GetPATH() string {
	return filepath.Join(getHome(), ".ikascrew")
}

func GetDBFile() string {
	if _, err := os.Stat(DefaultFileName); err == nil {
		return DefaultFileName
	}
	return filepath.Join(GetPATH(), DefaultFileName)
}
