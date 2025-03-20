package httpClient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/proxy"

	"go-template/pkg/logger"
	_ "go-template/pkg/tracer" // Fix init
)

type ClientOptions struct {
	BaseURL            *url.URL
	Sock5Proxy         string
	InsecureSkipVerify bool
}
type Client struct {
	c       *http.Client
	BaseURL *url.URL
}

func NewClient(options ClientOptions) *Client {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
	}

	if options.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if options.Sock5Proxy != "" {
		dialSocksProxy, err := proxy.SOCKS5("tcp", "localhost:11337", nil, proxy.Direct)
		if err != nil {
			fmt.Println("Error connecting to proxy:", err)
		}

		if contextDialer, ok := dialSocksProxy.(proxy.ContextDialer); ok {
			transport.DialContext = contextDialer.DialContext
		} else {
			logger.FatalErr("Error connecting to proxy:", errors.New("Failed type assertion to DialContext"))

			return nil
		}
	}

	c := http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, //nolint:mnd // it is sane
	}

	return &Client{
		c:       &c,
		BaseURL: options.BaseURL,
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
