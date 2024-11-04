package flight2fa

import "os"

func GetAPIEndpoint(url string) string {
	if endpoint := os.Getenv("API_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	return url + "/api"
}
