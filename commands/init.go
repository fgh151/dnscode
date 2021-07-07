package commands

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"openitstudio.ru/dnscode/providers"
	"openitstudio.ru/dnscode/utils/fs"
	"os"
	"strconv"
)

func DownloadProviders(zones providers.Zones) {
	for _, zone := range zones.Zones {
		name := zone.Provider
		path := fs.GetWorkDir() + "/.providers/" + name + ".so"
		uri := "https://github.com/fgh151/dnscode/releases/latest/download/" + name + ".so"
		if false == fs.FileExists(path) {
			err := downloadFile(path, uri)
			if err != nil {
				fmt.Println("Cant find provider " + name + " to download from " + uri)
			}
		}
	}
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := errors.New("Status code " + strconv.Itoa(resp.StatusCode))

		return err
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
