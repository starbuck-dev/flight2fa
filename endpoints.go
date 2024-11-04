package flight2fa

import "os"

var Days = "30"

func GetAPIEndpoint(url string) string {
	if endpoint := os.Getenv("API_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return url + "/api"
}
