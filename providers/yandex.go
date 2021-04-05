package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//https://yandex.ru/dev/pdd/doc/reference/dns-list.html

type YandexProvider struct {
	PddToken string
}

type YandexErrorResponse struct {
	Domain  string
	success string
	error   string
}

type YandexResponse struct {
	Domain  string         `json:"domain"`
	Records []YandexRecord `json:"records"`
}

type YandexRecord struct {
	Content   string `json:"content"`
	RecordId  int    `json:"record_id"`
	Fqdn      string `json:"fqdn"`
	Ttl       int    `json:"ttl"`
	Domain    string `json:"domain"`
	Priority  int    `json:"priority"`
	Port      int    `json:"port"`
	Weight    int    `json:"weight"`
	Target    string `json:"target"`
	Subdomain string `json:"subdomain"`
	Type      string `json:"type"`
}

func (p YandexProvider) GetRecords(domain string) []DnsRecord {

	url := "https://pddimp.yandex.ru/api2/admin/dns/list?domain=" + domain

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PddToken", p.getToken())
	res, _ := client.Do(req)

	var yaResp YandexResponse
	json.NewDecoder(res.Body).Decode(&yaResp)
	defer res.Body.Close()

	var returnAr []DnsRecord

	for _, r := range yaResp.Records {
		returnAr = append(returnAr, DnsRecord{Value: r.Content, Type: r.Type, Host: r.Domain, Ttl: r.Ttl, ExternalId: strconv.Itoa(r.RecordId)})
	}

	return returnAr
}

func (p YandexProvider) AddRecord(record DnsRecord) {
	url := fmt.Sprintf("https://pddimp.yandex.ru/api2/admin/dns/list?domain=%s&type=%s&content%s&ttl=%s", record.Host, record.Type, record.Value, strconv.Itoa(record.Ttl))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("PddToken", p.getToken())
	client.Do(req)
}

func (p YandexProvider) DeleteRecord(record DnsRecord) {

	url := fmt.Sprintf("https://pddimp.yandex.ru/api2/admin/dns/del?domain=%s&record_id=%s", record.Host, record.ExternalId)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("PddToken", p.getToken())
	client.Do(req)
}

func (p YandexProvider) getToken() string {
	token := p.PddToken
	if strings.HasPrefix(p.PddToken, "ENV_") {
		token = os.Getenv(p.PddToken)
	}
	return token
}
