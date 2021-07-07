package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"openitstudio.ru/dnscode/providers"
	"openitstudio.ru/dnscode/utils/fs"
	"os"
)

func GetZonesFromConfig(filename string) providers.Zones {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var zones providers.Zones

	json.Unmarshal(byteValue, &zones)

	for i, z := range zones.Zones {
		if z.Include != "" && fs.FileExists(z.Include) {
			incFileName := z.Include
			zones.Zones = remove(zones.Zones, i)
			zones.Zones = append(zones.Zones, readIncludeFile(incFileName))
		}
	}

	return zones
}

func readIncludeFile(fileName string) providers.ZoneProvider {
	zoneFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	defer zoneFile.Close()
	byteZoneValue, _ := ioutil.ReadAll(zoneFile)
	var provider providers.ZoneProvider
	json.Unmarshal(byteZoneValue, &provider)

	return provider
}

func remove(s []providers.ZoneProvider, i int) []providers.ZoneProvider {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
