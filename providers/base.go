package providers

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"openitstudio.ru/dnscode/utils/fs"
	"os"
	"plugin"
	//"openitstudio.ru/dnscode/providers/adman"
)

type ZoneProvider struct {
	Provider     string      `json:"provider"`
	Name         string      `json:"name"`
	Records      []DnsRecord `json:"records"`
	BaseProvider interface{} `json:"connection"`
	Include      string      `json:"include"`
}

func (z ZoneProvider) GetProvider() BaseProvider {
	var provider BaseProvider

	var providerName = z.Provider

	p, err := plugin.Open(fs.GetWorkDir() + "/.providers/" + providerName + ".so")
	if err != nil {
		fmt.Println("No provider " + providerName + " found! run init command")
		os.Exit(0)
	}

	importedProvider, err := p.Lookup("Provider")
	if err != nil {
		panic(err)
	}

	provider, _ = importedProvider.(BaseProvider)
	mapstructure.Decode(z.BaseProvider, &provider)

	return provider
}

type Zones struct {
	Zones []ZoneProvider `json:"zones"`
}

type DnsRecord struct {
	Host           string `json:"host"`
	Type           string `json:"type"`
	Value          string `json:"value"`
	Ttl            int    `json:"ttl"`
	Subdomain      string `json:"subdomain"`
	ExternalId     string
	AdditionalInfo string
}

type BaseProvider interface {
	GetRecords(domain string) []DnsRecord

	DeleteRecord(record DnsRecord)
	AddRecord(record DnsRecord)
}
