package client

import (
	"net/http"
	neturl "net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var scs spew.ConfigState = spew.ConfigState{
	Indent:         "    ",
	DisableMethods: true,
	// DisablePointerAddresses: true,
}

type Client struct {
	BaseUrl  string
	Err      error
	ReqResp  *http.Response
	BeginReq time.Time
	EndReq   time.Time
	ReqTime  time.Duration

	httpClient *http.Client
}

func NewClient(url string, opts ...ClientOpt) *Client {
	client := &Client{BaseUrl: url}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

type ClientOpt func(*Client)

func HttpClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func (client *Client) Get(rel string) {
	client.BeginReq = time.Now()

	defer func(reqBegan time.Time) {
		client.EndReq = time.Now()
		client.ReqTime = time.Since(reqBegan)
	}(client.BeginReq)

	var base *neturl.URL
	base, client.Err = neturl.Parse(client.BaseUrl)
	if client.Err != nil {
		return
	}

	var url *neturl.URL
	url, client.Err = base.Parse(rel)
	if client.Err != nil {
		return
	}

	// add any req query params

	httpReq := &http.Request{
		Method: "GET",
		URL:    url,
	}

	// set any req headers

	client.ReqResp, client.Err = client.httpClient.Do(httpReq)
}
