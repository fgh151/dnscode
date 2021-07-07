package fs

import (
	"log"
	"os"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetWorkDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	return path
}
