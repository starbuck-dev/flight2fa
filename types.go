package main

import "net/http"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// LoginRequest represents the authentication request structure
type LoginRequest struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	RememberMe       bool   `json:"rememberMe"`
	Eng              bool   `json:"eng"`
	ConfirmationCode string `json:"confirmationCode,omitempty"`
}

// LoginResponse represents the server's response to authentication
type LoginResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	User    User   `json:"user"`
	Info    Info   `json:"info"`
}

type User struct {
	Name        string `json:"name"`
	AllowAccess bool   `json:"allowAccess"`
}

type Info struct {
	Type        int    `json:"type"`
	SendingInfo string `json:"sendingInfo"`
}
