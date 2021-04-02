package commands

import (
	"encoding/json"
	"io/ioutil"
	"openitstudio.ru/dnscode/utils"
)

func ImportDomains(configFile string, fileToSave string, useImport bool) {

	var zones = utils.GetZonesFromConfig(configFile)

	for i := 0; i < len(zones.Zones); i++ {
		currentZone := &zones.Zones[i]
		provider := currentZone.GetProvider()
		records := provider.GetRecords(currentZone.Name)
		currentZone.Records = append(currentZone.Records, records...)

		if useImport {
			currentZone.Include = currentZone.Name + ".json"
			currentZoneFile, _ := json.MarshalIndent(currentZone, "", " ")
			_ = ioutil.WriteFile(currentZone.Name+".json", currentZoneFile, 0644)
			currentZone.Records = nil
		}
	}

	file, _ := json.MarshalIndent(zones, "", " ")
	_ = ioutil.WriteFile(fileToSave, file, 0644)
}
