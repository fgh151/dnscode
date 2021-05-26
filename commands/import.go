package commands

import (
	"encoding/json"
	"io/ioutil"
	"openitstudio.ru/dnscode/providers"
	"openitstudio.ru/dnscode/utils"
)

func GetRemoteZones(configFile string) providers.Zones {
	var localZones = utils.GetZonesFromConfig(configFile)
	var remoteZones = providers.Zones{}

	for i := 0; i < len(localZones.Zones); i++ {
		currentZone := providers.ZoneProvider{Name: localZones.Zones[i].Name}
		provider := localZones.Zones[i].GetProvider()
		records := provider.GetRecords(localZones.Zones[i].Name)
		currentZone.Records = records
		remoteZones.Zones = append(remoteZones.Zones, currentZone)
	}

	return remoteZones
}

func ImportDomains(configFile string, fileToSave string, useImport bool) {

	var zones = utils.GetZonesFromConfig(configFile)

	for i := 0; i < len(zones.Zones); i++ {
		currentZone := &zones.Zones[i]
		provider := currentZone.GetProvider()
		records := provider.GetRecords(currentZone.Name)
		currentZone.Records = append(currentZone.Records, records...)

		if useImport {
			currentZone.Include = ""
			currentZoneFileContent, _ := json.MarshalIndent(currentZone, "", " ")
			_ = ioutil.WriteFile(currentZone.Name+".json", currentZoneFileContent, 0644)
			currentZone.Include = currentZone.Name + ".json"
			currentZone.Records = nil
		}
	}

	file, _ := json.MarshalIndent(zones, "", " ")
	_ = ioutil.WriteFile(fileToSave, file, 0644)
}
