package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"openitstudio.ru/dnscode/utils"
	"os"
)

const FILENAME = "dnscode.json"

func main() {

	if len(os.Args) < 2 {
		fmt.Println("list or count subcommand is required")
		os.Exit(1)
	}

	importCommand := flag.NewFlagSet("import", flag.ExitOnError)

	importTextPtr := importCommand.String("filename", "", "File name to save. If empty it will override "+FILENAME)
	importDirectiveImportPtr := importCommand.Bool("useImport", true, "Use import directive, default true")

	switch os.Args[1] {
	case "import":
		importCommand.Parse(os.Args[2:])

	//case "count":
	//	countCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if importCommand.Parsed() {

		filename := *importTextPtr
		if filename == "" {
			filename = FILENAME
		}

		fmt.Println(filename)
		fmt.Println("importing")
		importDomains(filename, *importDirectiveImportPtr)
	}

	//jsonFile, err := os.Open(FILENAME)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//defer jsonFile.Close()
	//
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//var zones providers.Zones;
	//
	//json.Unmarshal(byteValue, &zones)
	//for i := 0; i < len(zones.Zones); i++ {
	//
	//	currentZone := zones.Zones[i]
	//
	//	provider := currentZone.GetProvider()
	//
	//
	//	fmt.Println(provider.GetRecords(currentZone.Name))
	//	for j :=0; j < len(zones.Zones[i].Records); j++ {
	//		fmt.Println(zones.Zones[i].Records[j].Value)
	//	}
	//}

}

func importDomains(fileToSave string, useImport bool) {

	var zones = utils.GetZonesFromConfig(FILENAME)

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
