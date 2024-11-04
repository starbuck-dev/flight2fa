package flight2fa

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name        string
		responses   []string
		statusCodes []int
		wantErr     bool
		errContains string
	}{
		{
			name:        "successful login without 2FA",
			responses:   []string{`{"user":{"name":"Test User","allowAccess":true}}`},
			statusCodes: []int{200},
			wantErr:     false,
		},
		{
			name: "successful login with 2FA",
			responses: []string{
				`{"code":"CHECK_AUTH_FAILED","info":{"sendingInfo":"test@email.com"}}`,
				`{"user":{"name":"Test User","allowAccess":true}}`,
			},
			statusCodes: []int{200, 200},
			wantErr:     false,
		},
		{
			name:        "invalid credentials",
			responses:   []string{`{"user":{"allowAccess":false},"message":"Invalid credentials"}`},
			statusCodes: []int{401},
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseIndex := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCodes[responseIndex])
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.responses[responseIndex]))
				if responseIndex < len(tt.responses)-1 {
					responseIndex++
				}
			}))
			defer server.Close()

			// Mock stdin for 2FA code input if needed
			if len(tt.responses) > 1 {
				cleanup := mockStdin("123456")
				defer cleanup()
			}

			_, err := Authenticate("testuser", "testpass", server.URL)

			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Authenticate() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

func TestNetworkError(t *testing.T) {
	// Create a test server that will be immediately closed to simulate network error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	// Test authentication with network error
	_, err := Authenticate("testuser", "testpass", server.URL)
	if err == nil || !strings.Contains(err.Error(), "network error") {
		t.Errorf("Authenticate() expected network error, got %v", err)
	}
}

// Helper functions
func mockStdin(input string) func() {
	old := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.Write([]byte(input + "\n"))
	w.Close()
	return func() {
		os.Stdin = old
	}
}

type errorTransport struct{}

func (t *errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("network error")
}

// Helper function to create test server with delayed response
func newTestServerWithDelay(delay time.Duration, statusCode int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

// Helper function to check error messages
func containsError(err error, want string) bool {
	if err == nil {
		return want == ""
	}
	return strings.Contains(err.Error(), want)
}

func TestNetworkScenarios(t *testing.T) {
	tests := []struct {
		name        string
		transport   http.RoundTripper
		timeout     time.Duration
		wantErr     bool
		errContains string
	}{
		{
			name:        "network error",
			transport:   &errorTransport{},
			timeout:     1 * time.Second,
			wantErr:     true,
			errContains: "network error",
		},
		{
			name:        "timeout error",
			transport:   &timeoutTransport{},
			timeout:     100 * time.Millisecond,
			wantErr:     true,
			errContains: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Server implementation
			}))
			defer server.Close()

			// Override the default HTTP client for testing
			originalClient := http.DefaultClient
			http.DefaultClient = &http.Client{
				Transport: tt.transport,
				Timeout:   tt.timeout,
			}
			defer func() { http.DefaultClient = originalClient }()

			// Call the package-level Authenticate function
			_, err := Authenticate("testuser", "testpass", server.URL)

			// Check if error occurred as expected
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check error message if error was expected
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Authenticate() error = %v, want error containing %q", err, tt.errContains)
			}
		})
	}
}

// Add a timeout transport for testing timeouts
type timeoutTransport struct{}

func (t *timeoutTransport) RoundTrip(*http.Request) (*http.Response, error) {
	time.Sleep(200 * time.Millisecond)
	return nil, nil
}
