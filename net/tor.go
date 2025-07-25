package net

import (
	"github.com/cebas/go-util/util"
	"net/http"
	"net/url"
)

const torServerUrl = "socks5://127.0.0.1:9050"

var savedTransport http.RoundTripper

func StartTor() {
	savedTransport = http.DefaultTransport
	http.DefaultTransport = torTransport()
}

func StopTor() {
	http.DefaultTransport = savedTransport
}

func torTransport() http.RoundTripper {
	proxy, err := url.Parse(torServerUrl)
	util.FatalErrorCheck(err)

	return &http.Transport{Proxy: http.ProxyURL(proxy)}
}
