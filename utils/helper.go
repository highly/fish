package utils

import (
	"log"
	"os"
	"path/filepath"
)

var CurrentPath string

func init() {
	CurrentPath = GetBaseLocation()
}

func GetBaseLocation() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func IsDirectory(file string) bool {
	fi, err := os.Stat(file)
	if err == nil && fi.IsDir() {
		return true
	}
	return false
}

func IsFileExisted(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func GetPathDirectoryName(path string) string {
	return filepath.Base(path)
}

func GetProgramName(path string) string {
	return GetPathDirectoryName(path)
}

func DelFile(filePath string) {
	if IsFileExisted(filePath) {
		if err := os.Remove(filePath); err != nil {
			Red.Println(err.Error())
		}
	}
}
