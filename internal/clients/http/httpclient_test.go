package httpclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	baseURL := &url.URL{Scheme: "https", Host: "example.com"}
	client := NewClient(baseURL)

	assert.NotNil(t, client)
	assert.Equal(t, baseURL, client.BaseURL)
	assert.NotNil(t, client.c)
	assert.Equal(t, 10*time.Second, client.c.Timeout)
}

func TestClient_Do(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		path           string
		body           any
		args           map[string]string
		serverResponse string
		serverStatus   int
		expectedBody   string
		expectedError  bool
	}{
		{
			name:           "Successful GET request",
			method:         http.MethodGet,
			path:           "/api/data",
			serverResponse: `{"message": "Success"}`,
			serverStatus:   http.StatusOK,
			expectedBody:   `{"message": "Success"}`,
			expectedError:  false,
		},
		{
			name:           "Successful POST request with JSON body",
			method:         http.MethodPost,
			path:           "/api/create",
			body:           map[string]string{"key": "value"},
			serverResponse: `{"status": "created"}`,
			serverStatus:   http.StatusCreated,
			expectedBody:   `{"status": "created"}`,
			expectedError:  false,
		},
		{
			name:           "Server error",
			method:         http.MethodGet,
			path:           "/api/error",
			serverResponse: "Internal Server Error",
			serverStatus:   http.StatusInternalServerError,
			expectedBody:   "Internal Server Error",
			expectedError:  false,
		},
		{
			name:           "GET with query parameters",
			method:         http.MethodGet,
			path:           "/api/items",
			args:           map[string]string{"id": "123", "name": "test"},
			serverResponse: `{"items": []}`,
			serverStatus:   http.StatusOK,
			expectedBody:   `{"items": []}`,
			expectedError:  false,
		},
		{
			name:          "Invalid URL in request",
			method:        http.MethodGet,
			path:          ":invalid", // Invalid path that causes url.Parse error
			expectedError: true,
		},
		{
			name:           "Empty response body",
			method:         http.MethodGet,
			path:           "/api/empty",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectedBody:   "",
			expectedError:  false,
		},
		{
			name:           "invalid json body",
			method:         http.MethodPost,
			path:           "/api/create",
			body:           make(chan int), //invalid json
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectedBody:   "",
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.serverStatus)
				if tc.body != nil && r.Method == http.MethodPost {
					bodyBytes, _ := io.ReadAll(r.Body)
					assert.NotEmpty(t, bodyBytes)
					contentType := r.Header.Get("Content-Type")
					assert.Equal(t, "application/json", contentType)

					// Verify args are in the url
					for key, value := range tc.args {
						assert.Contains(t, r.URL.RawQuery, key+"="+value)
					}
				}
				if tc.args != nil && r.Method == http.MethodGet {
					for key, value := range tc.args {
						assert.Contains(t, r.URL.RawQuery, key+"="+value)
					}
				}

				if tc.serverResponse != "" {
					_, err := w.Write([]byte(tc.serverResponse))
					assert.NoError(t, err)
				}

			}))
			defer server.Close()

			baseURL, _ := url.Parse(server.URL)
			client := NewClient(baseURL)

			body, err := client.Do(context.Background(), tc.method, tc.path, tc.body, tc.args)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, string(body))
			}
		})
	}
}

func TestClient_newRequest(t *testing.T) {
	baseURL, _ := url.Parse("https://example.com")
	client := NewClient(baseURL)

	t.Run("GET request with no body or args", func(t *testing.T) {
		req, err := client.newRequest(context.Background(), http.MethodGet, "/api/data", nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "https://example.com/api/data", req.URL.String())
		assert.Nil(t, req.Body)
	})

	t.Run("POST request with JSON body", func(t *testing.T) {
		body := map[string]string{"key": "value"}
		req, err := client.newRequest(context.Background(), http.MethodPost, "/api/create", body, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "https://example.com/api/create", req.URL.String())
		assert.NotNil(t, req.Body)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var reqBody map[string]string
		err = json.NewDecoder(req.Body).Decode(&reqBody)
		assert.NoError(t, err)
		assert.Equal(t, body, reqBody)
	})

	t.Run("GET request with query arguments", func(t *testing.T) {
		args := map[string]string{"id": "123", "name": "test"}
		req, err := client.newRequest(context.Background(), http.MethodGet, "/api/items", nil, args)
		assert.NoError(t, err)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "https://example.com/api/items?id=123&name=test", req.URL.String())
		assert.Nil(t, req.Body)
		assert.Equal(t, args["id"], req.URL.Query().Get("id"))
		assert.Equal(t, args["name"], req.URL.Query().Get("name"))
	})
	t.Run("invalid json", func(t *testing.T) {
		body := make(chan int)
		_, err := client.newRequest(context.Background(), http.MethodPost, "/api/create", body, nil)
		assert.Error(t, err)
	})
}
