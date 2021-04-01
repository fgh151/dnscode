package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AdmanProvider struct {
	Login  string
	Mdpass string
}

type AdmanDomain struct {
	DomainId string `json:"domain_id"`
	Domain   string `json:"domain"`
}

type DomainListResponse struct {
	Data []AdmanDomain `json:"data"`
}

type AdmanZoneRecordsResponse struct {
	Data []AdmanZoneRecord `json:"data"`
}

type AdmanZoneRecord struct {
	RecId     string `json:"rec_id"`
	DomainId  string `json:"domain_id"`
	Type      string `json:"type"`
	Subdomain string `json:"subdomain"`
	Rec       string `json:"rec"`
	Prior     string `json:"prior"`
}

func (p AdmanProvider) GetRecords(domain string) []DnsRecord {
	domains := p.getZones()

	var returnAr []DnsRecord
	for _, d := range domains {
		if domain == d.Domain {
			records := p.getZoneRecords(d.DomainId)

			for _, r := range records {
				returnAr = append(returnAr, DnsRecord{Value: r.Rec, Type: r.Type, Host: r.Subdomain, Ttl: 10})
			}
		}
	}

	return returnAr
}

func (p AdmanProvider) getZoneRecords(domainId string) []AdmanZoneRecord {

	surl := "https://adman.com/api/domain/zonelist"
	authString := fmt.Sprintf("{\"login\":\"%s\",\"mdpass\":\"%s\", \"filter\":[{\"domain_id\": \"%s\"}]}", p.Login, p.Mdpass, domainId)
	params := url.QueryEscape(authString)
	payload := strings.NewReader(fmt.Sprintf("req=%s", params))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", surl, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	var aResp AdmanZoneRecordsResponse
	json.Unmarshal(body, &aResp)

	return aResp.Data
}

func (p AdmanProvider) getZones() []AdmanDomain {
	surl := "https://adman.com/api/domain/list"

	authString := fmt.Sprintf("{\"login\":\"%s\",\"mdpass\":\"%s\"}", p.Login, p.Mdpass)
	params := url.QueryEscape(authString)
	payload := strings.NewReader(fmt.Sprintf("req=%s", params))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", surl, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var aResp DomainListResponse

	json.Unmarshal(body, &aResp)

	defer resp.Body.Close()

	return aResp.Data
}
