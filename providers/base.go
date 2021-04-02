package providers

import "github.com/mitchellh/mapstructure"

type ZoneProvider struct {
	Provider     string      `json:"provider"`
	Name         string      `json:"name"`
	Records      []DnsRecord `json:"records"`
	BaseProvider interface{} `json:"connection"`
	Include      string      `json:"include"`
}

func (z ZoneProvider) GetProvider() BaseProvider {
	var provider BaseProvider

	//
	switch z.Provider {
	case "yandex":
		{
			provider = YandexProvider{}
		}
	case "adman":
		{
			provider = AdmanProvider{}
		}
	case "regru":
		{
			provider = RegruProvider{}
		}
	}

	mapstructure.Decode(z.BaseProvider, &provider)

	return provider
}

type Zones struct {
	Zones []ZoneProvider `json:"zones"`
}

type DnsRecord struct {
	Host  string `json:"host"`
	Type  string `json:"type"`
	Value string `json:"value"`
	Ttl   int    `json:"ttl"`
}

type BaseProvider interface {
	GetRecords(domain string) []DnsRecord
}
