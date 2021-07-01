package commands

import (
	"fmt"
	"openitstudio.ru/dnscode/providers"
	"openitstudio.ru/dnscode/utils"
	"reflect"
)

type callback func(record providers.DnsRecord, provider providers.ZoneProvider)

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
							callbackDelete(remoteRecord, localZone)
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
						callbackAdd(localRecord, localZone)
					}
				}
			}
		}
	}
}

func printDelete(record providers.DnsRecord, provider providers.ZoneProvider) {
	utils.PrintlnWarning(fmt.Sprintf(provider.Name+" / DELETE %s %s %s %d", record.Subdomain, record.Type, record.Value, record.Ttl))
}

func printAdd(record providers.DnsRecord, provider providers.ZoneProvider) {
	utils.PrintInfo(fmt.Sprintf(provider.Name+" / ADD %s %s %s %d", record.Host, record.Type, record.Value, record.Ttl))
}

func deleteRecord(record providers.DnsRecord, provider providers.ZoneProvider) {
	printDelete(record, provider)
	provider.GetProvider().DeleteRecord(record)
}

func addRecord(record providers.DnsRecord, provider providers.ZoneProvider) {
	printAdd(record, provider)
	provider.GetProvider().AddRecord(record)
}

func Plan(configFile string, force *bool) {
	analyze(configFile, force, printDelete, printAdd)
}

func Apply(configFile string, force *bool) {
	analyze(configFile, force, deleteRecord, addRecord)
}
