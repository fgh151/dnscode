package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	http2 "openitstudio.ru/dnscode/http"
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

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("PddToken", p.getToken())

	res, _ := p.getClient().Do(req)

	var yaResp YandexResponse
	json.NewDecoder(res.Body).Decode(&yaResp)
	defer res.Body.Close()

	var returnAr []DnsRecord

	for _, r := range yaResp.Records {
		returnAr = append(returnAr, DnsRecord{Value: r.Content, Type: r.Type, Host: r.Domain, Subdomain: r.Subdomain, Ttl: r.Ttl, ExternalId: strconv.Itoa(r.RecordId)})
	}

	return returnAr
}

func (p YandexProvider) AddRecord(record DnsRecord) {
	url := "https://pddimp.yandex.ru/api2/admin/dns/add"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("domain", record.Host)
	_ = writer.WriteField("type", record.Type)
	_ = writer.WriteField("content", record.Value)
	_ = writer.WriteField("subdomain", record.Subdomain)
	_ = writer.WriteField("ttl", strconv.Itoa(record.Ttl))

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("PddToken", p.getToken())

	p.getClient().Do(req)
}

func (p YandexProvider) DeleteRecord(record DnsRecord) {

	url := fmt.Sprintf("https://pddimp.yandex.ru/api2/admin/dns/del?domain=%s&record_id=%s", record.Host, record.ExternalId)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("PddToken", p.getToken())

	p.getClient().Do(req)
}

func (p YandexProvider) getClient() *http.Client {
	c, _ := http2.CreateHttpClient()
	return c
}

func (p YandexProvider) getToken() string {
	token := p.PddToken
	if strings.HasPrefix(p.PddToken, "ENV_") {
		token = os.Getenv(p.PddToken)
	}
	return token
}
