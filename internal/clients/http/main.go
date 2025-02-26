package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	_ "go-template/pkg/tracer"
	"net/http"
	"net/url"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	BaseURL *url.URL
	c       *http.Client
}

// func init() {
// 	client := NewClient(&url.URL{Scheme: "https", Host: "httpbin.org"})

// 	var test interface{}
// 	args := map[string]string{
// 		"name": "go-template",
// 		"age":  "22",
// 	}
// 	test, err := client.Do(context.Background(), "GET", "get", nil, args)
// 	fmt.Printf("%+v", err)
// 	fmt.Printf("%+v", test)
// 	test, err = client.Do(context.Background(), "POST", "post", args, args)

// 	test, _ = client.Do(context.Background(), "GET", "/get", args, args)

// }

func NewClient(baseUrl *url.URL) *Client {
	c := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}
	return &Client{
		c:       &c,
		BaseURL: baseUrl,
	}

}
func (client Client) Do(ctx context.Context, method, path string, body interface{}, args map[string]string) (interface{}, error) {
	request, err := client.newRequest(ctx, method, path, body, args)
	if err != nil {
		return nil, err
	}

	resp, err := client.c.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (client Client) newRequest(ctx context.Context, method, path string, body interface{}, args map[string]string) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := client.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range args {
		q := req.URL.Query()
		q.Add(key, value)
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}
