package net

import (
	"github.com/avast/retry-go"
	"github.com/cebas/go-util/util"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	httpClient *http.Client
	cached     bool
}

type HttpHeader struct {
	Key   string
	Value string
}

var cache = util.NewCache()

func NewHttpClient(torified bool, cached bool) *HttpClient {
	var httpClient *http.Client

	if torified {
		httpClient = &http.Client{
			Transport: torTransport(),
			Timeout:   time.Second * 42,
		}
	} else {
		httpClient = http.DefaultClient
	}

	return &HttpClient{
		httpClient: httpClient,
		cached:     cached,
	}
}

func (c *HttpClient) httpCall(method string, path string, headers []HttpHeader, params url.Values, payload interface{}) ([]byte, error) {
	var reader io.Reader

	urlStruct, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	urlStruct.RawQuery = params.Encode()

	urlString := urlStruct.String()

	withCache := c.cached && method == http.MethodGet && payload == nil
	if withCache {
		body, ok := cache.Get(urlString)
		if ok {
			return body.([]byte), nil
		}
	}

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		reader = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(method, urlString, reader)
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}

	var resp *http.Response
	err = retry.Do(
		func() error {
			resp, err = c.httpClient.Do(req)
			if err == nil {
				if resp.StatusCode != 200 {
					return retry.Unrecoverable(fmt.Errorf("unexpected status code: %s", resp.Status))
				}
				return err
			}

			// err != nil: recoverable?
			return err
		},
		retry.Attempts(10),
		retry.Delay(1*time.Second),
	)

	if err != nil {
		return nil, err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if withCache {
		cache.Set(urlString, body)
	}

	return body, nil
}

func (c *HttpClient) Get(path string, headers []HttpHeader, params url.Values) (body []byte, err error) {
	return c.httpCall(http.MethodGet, path, headers, params, nil)
}

func (c *HttpClient) Patch(path string, headers []HttpHeader, params url.Values, payload interface{}) (err error) {
	_, err = c.httpCall(http.MethodPatch, path, headers, params, payload)
	return
}

func (c *HttpClient) Post(path string, headers []HttpHeader, params url.Values, payload interface{}) (body []byte, err error) {
	return c.httpCall(http.MethodPost, path, headers, params, payload)
}

func GetHttpContent(path string) (content []byte, err error) {
	httpClient := NewHttpClient(false, false)
	content, err = httpClient.Get(path, []HttpHeader{}, url.Values{})
	return
}
