package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	fmt.Println(rResp.Answer.Domains[0].Rrs)

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

	p := RequestParams{
		Username:          r.Username,
		Password:          r.Password,
		OutputContentType: "json",
	}

	d = append(p.Domains, RequestDomain{Dname: domain})

	p.Domains = d

	b, _ := json.Marshal(p)

	return b
}
