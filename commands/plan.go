package commands

import (
	"fmt"
	"openitstudio.ru/dnscode/providers"
	"openitstudio.ru/dnscode/utils"
	"reflect"
)

type callback func(record providers.DnsRecord)

func analyze(configFile string, force *bool, callbackDelete callback, callbackAdd callback) {
	localZones := utils.GetZonesFromConfig(configFile)
	remoteZones := GetRemoteZones(configFile)

	for _, remoteZone := range remoteZones.Zones {
		for _, localZone := range localZones.Zones {
			if localZone.Name == remoteZone.Name {
				if *force {
					for _, remoteRecord := range remoteZone.Records {
						del := true
						for _, localRecord := range localZone.Records {
							if reflect.DeepEqual(remoteRecord, localRecord) {
								del = false
							}
						}
						if del {
							callbackDelete(remoteRecord)
						}
					}
				}

				for _, localRecord := range localZone.Records {
					exist := false
					for _, remoteRecord := range remoteZone.Records {

						if reflect.DeepEqual(remoteRecord, localRecord) {
							exist = true
						}
					}
					if !exist {
						callbackAdd(localRecord)
					}
				}
			}
		}
	}
}

func printDelete(record providers.DnsRecord) {
	fmt.Println(fmt.Sprintf("DELETE %s %s %s %d", record.Host, record.Type, record.Value, record.Ttl))
}

func printAdd(record providers.DnsRecord) {
	fmt.Println(fmt.Sprintf("ADD %s %s %s %d", record.Host, record.Type, record.Value, record.Ttl))
}

func Plan(configFile string, force *bool) {
	analyze(configFile, force, printDelete, printAdd)
}
