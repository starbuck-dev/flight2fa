package flight2fa

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Authenticate performs the complete authentication process with automatic 2FA handling
func Authenticate(username, password, baseURL string) (*Client, error) {

	client, err := NewClient(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	resp, err := client.attemptLogin(username, password, "")
	if err != nil {
		return nil, fmt.Errorf("login attempt failed: %v", err)
	}

	if resp.Code == "CHECK_AUTH_FAILED" {
		fmt.Print("Enter 2FA code: ")
		reader := bufio.NewReader(os.Stdin)
		code, _ := reader.ReadString('\n')
		code = strings.TrimSpace(code)

		resp, err = client.attemptLogin(username, password, code)
		if err != nil {
			return nil, fmt.Errorf("2FA verification failed: %v", err)
		}
	}

	if !resp.User.AllowAccess && resp.Code != "CHECK_AUTH_FAILED" {
		if resp.Message != "" {
			return nil, fmt.Errorf("access denied: %s", resp.Message)
		}
		return nil, fmt.Errorf("access denied")
	}

	fmt.Println("Successfully authenticated!")
	return client, nil
}

// attemptLogin handles the login request (now unexported as it's an internal helper)
func (c *Client) attemptLogin(username, password, code string) (*LoginResponse, error) {
	req := LoginRequest{
		Username:   username,
		Password:   password,
		RememberMe: true,
		Eng:        false,
	}

	if code != "" {
		req.ConfirmationCode = code
	}

	payloadBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	request, err := http.NewRequest("POST", c.BaseURL+"/api/auth", strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	println(c.BaseURL)
	println(c.BaseURL + "/Login")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Origin", c.BaseURL)
	request.Header.Set("Referer", c.BaseURL+"/Login")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var loginResp LoginResponse
	if err := json.Unmarshal(body, &loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if resp.StatusCode == http.StatusUnauthorized && loginResp.Code != "CHECK_AUTH_FAILED" {
		return nil, fmt.Errorf("invalid credentials")
	}

	if resp.StatusCode != http.StatusOK && loginResp.Code != "CHECK_AUTH_FAILED" {
		return nil, fmt.Errorf("server error: status code %d", resp.StatusCode)
	}

	return &loginResp, nil
}
