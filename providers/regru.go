package providers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

func (r RegruProvider) GetRecords(domain string) []DnsRecord {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.reg.ru/api/regru2/zone/get_resource_records", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	q.Add("input_data", string(r.crateParams(domain)))
	q.Add("input_format", "json")

	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var rResp RegruResponse
	json.Unmarshal(body, &rResp)

	var returnAr []DnsRecord

	for _, d := range rResp.Answer.Domains {
		if d.Dname == domain {
			for _, rrs := range d.Rrs {
				returnAr = append(returnAr, DnsRecord{
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
