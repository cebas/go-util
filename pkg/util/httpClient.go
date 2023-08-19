package util

import (
	"time"

	"io"
	"net/http"
	"net/url"
)

func torTransport() http.RoundTripper {
	proxy, err := url.Parse("socks5://127.0.0.1:9050")
	FatalErrorCheck(err)

	return &http.Transport{Proxy: http.ProxyURL(proxy)}
}

func GetHttpContent(uri string, torified bool) (body []byte, err error) {
	httpClient := http.DefaultClient

	if torified {
		httpClient = &http.Client{
			Transport: torTransport(),
			Timeout:   time.Second * 42,
		}
	}

	resp, err := httpClient.Get(uri)
	if err != nil {
		return
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

var savedTransport http.RoundTripper

func StartTor() {
	savedTransport = http.DefaultTransport
	http.DefaultTransport = torTransport()
}

func StopTor() {
	http.DefaultTransport = savedTransport
}
