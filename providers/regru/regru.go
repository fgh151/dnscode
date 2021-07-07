package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	http2 "openitstudio.ru/dnscode/http"
	"openitstudio.ru/dnscode/providers"
	"os"
	"strings"
)

type RequestParams struct {
	Username          string          `json:"username"`
	Password          string          `json:"password"`
	OutputContentType string          `json:"output_content_type"`
	Domains           []RequestDomain `json:"domains"`
}

type RequestDomain struct {
	Dname string `json:"dname"`
}

type RegruProvider struct {
	Username string
	Password string
}

func (r RegruProvider) AddRecord(record providers.DnsRecord) {
	var endpoint = ""

	params := map[string]interface{}{
		"username": r.Username,
		"password": r.Password,
	}

	switch record.Type {
	case "A":
		endpoint = "https://api.reg.ru/api/regru2/zone/add_alias"
		params["ipaddr"] = record.Value
		break
	case "AAAA":
		endpoint = "https://api.reg.ru/api/regru2/zone/add_aaaa"
		params["ipaddr"] = record.Value
		break
	case "CNAME":
		endpoint = "https://api.reg.ru/api/regru2/zone/add_cname"
		params["canonical_name"] = record.Value
		break
	case "MX":
		endpoint = "https://api.reg.ru/api/regru2/zone/add_mx"
		params["mail_server"] = record.Value
		params["priority"] = record.Ttl
		break
	case "NS":
		endpoint = "https://api.reg.ru/api/regru2/zone/add_ns"
		params["dns_server"] = record.Value
		params["priority"] = record.Ttl
		break
	case "TXT":
		endpoint = "https://api.reg.ru/api/regru2/zone/add_txt"
		params["text"] = record.Value
		break
	default:
		endpoint = "https://api.reg.ru/api/regru2/zone/add_txt"
		params["text"] = record.Value
		break
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	b, _ := json.Marshal(params)

	q.Add("input_data", string(b))
	q.Add("input_format", "json")

	req.URL.RawQuery = q.Encode()

	c, _ := http2.CreateHttpClient()
	c.Do(req)
}

type RegruRecord struct {
	Content string `json:"content"`
	Prio    int    `json:"prio"`
	Rectype string `json:"rectype"`
	State   string `json:"state"`
	Subname string `json:"subname"`
}

type RegruDomain struct {
	Dname string        `json:"dname"`
	Reult string        `json:"reult"`
	Rrs   []RegruRecord `json:"rrs"`
}

type RegruAnswer struct {
	Domains []RegruDomain `json:"domains"`
}

type RegruResponse struct {
	Answer RegruAnswer `json:"answer"`
	Result string      `json:"result"`
}

//https://www.reg.ru/support/help/api2#zone_get_resource_records

func (r RegruProvider) GetRecords(domain string) []providers.DnsRecord {

	req, err := http.NewRequest("GET", "https://api.reg.ru/api/regru2/zone/get_resource_records", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	q.Add("input_data", string(r.crateParams(domain)))
	q.Add("input_format", "json")

	req.URL.RawQuery = q.Encode()

	c, _ := http2.CreateHttpClient()

	resp, _ := c.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var rResp RegruResponse
	json.Unmarshal(body, &rResp)

	var returnAr []providers.DnsRecord

	for _, d := range rResp.Answer.Domains {
		if d.Dname == domain {
			for _, rrs := range d.Rrs {
				returnAr = append(returnAr, providers.DnsRecord{
					Host:  rrs.Subname,
					Type:  rrs.Rectype,
					Value: rrs.Content,
					Ttl:   10,
				})
			}
		}
	}

	return returnAr
}

func (r RegruProvider) DeleteRecord(record providers.DnsRecord) {
	req, err := http.NewRequest("GET", "https://api.reg.ru/api/regru2/zone/remove_record", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	params := struct {
		Username          string `json:"username"`
		Password          string `json:"password"`
		Subdomain         string `json:"subdomain"`
		Content           string `json:"content"`
		RecordType        string `json:"record_type"`
		OutputContentType string `json:"output_content_type"`
	}{
		Username:          r.Username,
		Password:          r.Password,
		Subdomain:         record.Host,
		Content:           record.Value,
		RecordType:        record.Type,
		OutputContentType: "plain",
	}

	b, _ := json.Marshal(params)

	q.Add("input_data", string(b))
	q.Add("input_format", "json")

	req.URL.RawQuery = q.Encode()

	c, _ := http2.CreateHttpClient()
	c.Do(req)
}

func (r RegruProvider) crateParams(domain string) []byte {

	var d []RequestDomain

	p := r.getRequestParams()

	d = append(p.Domains, RequestDomain{Dname: domain})

	p.Domains = d

	b, _ := json.Marshal(p)

	return b
}

func (r RegruProvider) getRequestParams() RequestParams {

	user := r.Username
	pass := r.Password

	if strings.HasPrefix(pass, "ENV_") {
		pass = os.Getenv(pass)
	}

	if strings.HasPrefix(user, "ENV_") {
		user = os.Getenv(user)
	}

	return RequestParams{
		Username:          user,
		Password:          pass,
		OutputContentType: "json",
	}
}

//goland:noinspection GoUnusedGlobalVariable
var Provider RegruProvider
