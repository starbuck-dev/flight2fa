package main

import "os"

var global_url string

var defaultAPIEndpoint = global_url + "/api"
var origin = global_url
var referer = global_url + "/login"
var authEndpoint = global_url + "/auth"
var ScheduleEndpoint = GetAPIEndpoint(defaultAPIEndpoint) + "/crew-plan"
var Days = "30"

func SetGlobalUrl(url string) {
	global_url = url
}

func GetAPIEndpoint(url string) string {
	if endpoint := os.Getenv("API_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return defaultAPIEndpoint
}
