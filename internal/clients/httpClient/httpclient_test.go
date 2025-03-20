package httpClient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		options ClientOptions
		wantErr bool
	}{
		{
			name: "basic client creation",
			options: ClientOptions{
				BaseURL:            &url.URL{Scheme: "https", Host: "example.com"},
				InsecureSkipVerify: false,
			},
			wantErr: false,
		},
		{
			name: "client with insecure skip verify",
			options: ClientOptions{
				BaseURL:            &url.URL{Scheme: "https", Host: "example.com"},
				InsecureSkipVerify: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.options)
			assert.NotNil(t, client)
			assert.Equal(t, tt.options.BaseURL, client.BaseURL)
		})
	}
}

func TestClient_Do(t *testing.T) {
	type testResponse struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		args           map[string]string
		serverResponse *testResponse
		serverStatus   int
		wantErr        bool
	}{
		{
			name:   "successful GET request",
			method: http.MethodGet,
			path:   "/test",
			serverResponse: &testResponse{
				Message: "success",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:   "successful POST request with body",
			method: http.MethodPost,
			path:   "/test",
			body: map[string]string{
				"key": "value",
			},
			serverResponse: &testResponse{
				Message: "created",
			},
			serverStatus: http.StatusCreated,
			wantErr:      false,
		},
		{
			name:   "successful GET request with query params",
			method: http.MethodGet,
			path:   "/test",
			args: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			serverResponse: &testResponse{
				Message: "success with params",
			},
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				assert.Equal(t, tt.method, r.Method)

				// Verify query parameters if provided
				if tt.args != nil {
					query := r.URL.Query()
					for key, value := range tt.args {
						assert.Equal(t, value, query.Get(key))
					}
				}

				// Verify request body if provided
				if tt.body != nil {
					var receivedBody map[string]string
					err := json.NewDecoder(r.Body).Decode(&receivedBody)
					require.NoError(t, err)
					assert.Equal(t, tt.body, receivedBody)
				}

				// Send response
				w.WriteHeader(tt.serverStatus)

				if tt.serverResponse != nil {
					json.NewEncoder(w).Encode(tt.serverResponse)
				}
			}))
			defer server.Close()

			// Parse server URL
			serverURL, err := url.Parse(server.URL)
			require.NoError(t, err)

			// Create client
			client := NewClient(ClientOptions{
				BaseURL: serverURL,
			})

			// Make request
			resp, err := client.Do(context.Background(), tt.method, tt.path, tt.body, tt.args)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			// Verify response
			if tt.serverResponse != nil {
				var response testResponse
				err = json.Unmarshal(resp, &response)
				require.NoError(t, err)
				assert.Equal(t, tt.serverResponse.Message, response.Message)
			}
		})
	}
}

func TestClient_newRequest(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		path    string
		body    interface{}
		args    map[string]string
		wantErr bool
	}{
		{
			name:    "valid request without body",
			method:  http.MethodGet,
			path:    "/test",
			wantErr: false,
		},
		{
			name:   "valid request with body",
			method: http.MethodPost,
			path:   "/test",
			body: map[string]string{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:   "valid request with query params",
			method: http.MethodGet,
			path:   "/test",
			args: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(ClientOptions{
				BaseURL: &url.URL{Scheme: "https", Host: "example.com"},
			})

			req, err := client.newRequest(context.Background(), tt.method, tt.path, tt.body, tt.args)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.NotNil(t, req)
			assert.Equal(t, tt.method, req.Method)

			// Verify query parameters
			if tt.args != nil {
				query := req.URL.Query()
				for key, value := range tt.args {
					assert.Equal(t, value, query.Get(key))
				}
			}

			// Verify request body
			if tt.body != nil {
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			}
		})
	}
}
