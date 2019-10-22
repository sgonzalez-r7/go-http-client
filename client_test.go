package client

import (
	"context"
	"fmt"
	"gotest.tools/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockHttpClient(handler http.Handler) (*http.Client, func()) {
	server := httptest.NewServer(handler)

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, server.Listener.Addr().String())
			},
		},
	}

	return client, server.Close
}

func testReqHandler(t *testing.T, tc testCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqUrl := fmt.Sprintf("http://%s%s", r.Host, r.URL.Path)

		assert.Equal(t, r.Method, tc.ExpectMethod)
		assert.Equal(t, reqUrl, tc.ExpectUrl)
	})
}

type testCase struct {
	Descr    string
	Method   string
	BaseUrl  string
	Endpoint string

	ExpectMethod string
	ExpectUrl    string
	ExpectErr    string
}

func TestClientGet(t *testing.T) {
	testCases := []testCase{
		testCase{
			Descr:    "with endoint",
			Method:   "GET",
			BaseUrl:  "http://testing.example.com",
			Endpoint: "/rel/path",

			ExpectMethod: "GET",
			ExpectUrl:    "http://testing.example.com/rel/path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Descr, func(t *testing.T) {
			httpClient, closeClient := mockHttpClient(testReqHandler(t, tc))

			client := NewClient(tc.BaseUrl, HttpClient(httpClient))
			client.Get(tc.Endpoint)
			closeClient()

			if tc.ExpectErr != "" {
				assert.Error(t, client.Err, tc.ExpectErr)
			} else {
				assert.NilError(t, client.Err)
			}

		})
	}
}
