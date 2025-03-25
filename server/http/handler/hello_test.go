package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	httpclient "go-template/internal/clients/httpClient"
)

func TestHandler_Hello(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	client := httpclient.NewClient(httpclient.ClientOptions{
		BaseURL: &url.URL{Scheme: "https", Host: "hacker-news.firebaseio.com", Path: "v0/"},
		// Sock5Proxy:         config.Get(config.SOCKS5_PROXY),
		InsecureSkipVerify: false,
	})

	h := Handler{
		HTTPClient: client,
	}

	// Assertions
	if assert.NoError(t, h.Hello(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello, World!", rec.Body.String())
	}
}
