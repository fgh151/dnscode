package http

import (
	"fmt"
	"net/http"
	"net/url"
)

var proxyUrl string

func SetProxy(proxy string) {
	proxyUrl = proxy
}

func CreateHttpClient() (*http.Client, error) {
	if proxyUrl != "" {
		proxyUrl, err := url.Parse(proxyUrl)

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}
		fmt.Println(proxyUrl)

		return client, err
	}

	return &http.Client{}, nil
}
