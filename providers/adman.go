package providers

type AdmanProvider struct {
	Login  string
	Mdpass string
}

func (p AdmanProvider) GetRecords(domain string) []DnsRecord {
	return []DnsRecord{}
}
