package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	http2 "openitstudio.ru/dnscode/http"
	"os"
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
				returnAr = append(returnAr, DnsRecord{Value: r.Rec, Type: r.Type, Host: r.Subdomain, Ttl: 10, ExternalId: r.RecId, AdditionalInfo: r.DomainId})
			}
		}
	}

	return returnAr
}

func (p AdmanProvider) getZoneRecords(domainId string) []AdmanZoneRecord {

	surl := "https://adman.com/api/domain/zonelist"
	authString := fmt.Sprintf("{\"login\":\"%s\",\"mdpass\":\"%s\", \"filter\":[{\"domain_id\": \"%s\"}]}", p.getLogin(), p.getPass(), domainId)
	params := url.QueryEscape(authString)
	payload := strings.NewReader(fmt.Sprintf("req=%s", params))
	req, _ := http.NewRequest("POST", surl, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c, _ := http2.CreateHttpClient()
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	var aResp AdmanZoneRecordsResponse
	json.Unmarshal(body, &aResp)

	return aResp.Data
}

func (p AdmanProvider) AddRecord(record DnsRecord) {
	surl := "https://adman.com/api/domain/zoneadd"

	rp := map[string]interface{}{
		"login":     p.getLogin(),
		"mdpass":    p.getPass(),
		"rec_id":    record.ExternalId,
		"domain_id": record.AdditionalInfo,
		"type":      record.Type,
		"subdomain": record.Host,
		"rec":       record.Value,
		"prior":     record.Ttl,
	}
	b, _ := json.Marshal(rp)
	params := string(b)

	payload := strings.NewReader(fmt.Sprintf("req=%s", params))
	req, _ := http.NewRequest("POST", surl, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c, _ := http2.CreateHttpClient()
	c.Do(req)
}

func (p AdmanProvider) DeleteRecord(record DnsRecord) {
	surl := "https://adman.com/api/domain/zonedelete"
	rp := map[string]interface{}{
		"login":     p.getLogin(),
		"mdpass":    p.getPass(),
		"rec_id":    record.ExternalId,
		"domain_id": record.AdditionalInfo,
	}
	b, _ := json.Marshal(rp)
	params := string(b)

	payload := strings.NewReader(fmt.Sprintf("req=%s", params))
	req, _ := http.NewRequest("POST", surl, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c, _ := http2.CreateHttpClient()
	c.Do(req)
}

func (p AdmanProvider) getZones() []AdmanDomain {
	surl := "https://adman.com/api/domain/list"

	authString := fmt.Sprintf("{\"login\":\"%s\",\"mdpass\":\"%s\"}", p.getLogin(), p.getPass())
	params := url.QueryEscape(authString)
	payload := strings.NewReader(fmt.Sprintf("req=%s", params))

	req, _ := http.NewRequest("POST", surl, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c, _ := http2.CreateHttpClient()
	resp, _ := c.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var aResp DomainListResponse

	json.Unmarshal(body, &aResp)

	defer resp.Body.Close()

	return aResp.Data
}

func (p AdmanProvider) getLogin() string {
	l := p.Login
	if strings.HasPrefix(p.Login, "ENV_") {
		l = os.Getenv(p.Login)
	}

	return l
}

func (p AdmanProvider) getPass() string {
	mdpass := p.Mdpass
	if strings.HasPrefix(p.Mdpass, "ENV_") {
		mdpass = os.Getenv(p.Mdpass)
	}

	return mdpass
}
