package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	_ "go-template/pkg/tracer" // Fix init
)

type Client struct {
	BaseURL *url.URL
	c       *http.Client
}

func NewClient(baseURL *url.URL) *Client {
	c := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second, //nolint:mnd // it is sane
	}

	return &Client{
		c:       &c,
		BaseURL: baseURL,
	}
}

func (client Client) Do(ctx context.Context, method, path string, body any, args map[string]string) ([]byte, error) {
	request, err := client.newRequest(ctx, method, path, body, args)
	if err != nil {
		return nil, err
	}

	resp, err := client.c.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (client Client) newRequest(ctx context.Context, method, path string, body any, args map[string]string) (*http.Request, error) {
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
