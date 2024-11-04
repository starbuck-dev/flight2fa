# Flight2FA

Flight2FA is a Go package that provides authentication functionality with 2FA support for accessing flight schedule information. It handles both regular authentication and two-factor authentication flows automatically.

## Installation

```bash
go get github.com/starbuck-dev/flight2fa
```

## Usage

Here's a basic example of how to use the Flight2FA package:

```go
package main

import (
    "fmt"
    "log"
    "strconv"
    "time"

    "github.com/starbuck-dev/flight2fa"
)

func main() {
    // Authenticate with credentials
    client, err := flight2fa.Authenticate("username", "password", "https://some.example.edu")
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    // Get schedule data
    dataResponse, err := GetSchedule(client)
    if err != nil {
        log.Fatalf("Failed to get schedule: %v", err)
    }
    fmt.Printf("Data response: %s\n", dataResponse)
}

func GetSchedule(client *flight2fa.Client) (string, error) {
    days := 30
    now := time.Now()
    url := fmt.Sprintf(
        "%s?dateBegin=%d&dateEnd=%d&eng=false&airportCode=3",
        "https://some.example.edu/api/crew-plan",
        now.Unix()*1000,
        now.Add(time.Duration(days)*24*time.Hour).Unix()*1000,
    )

    return client.Get(url)
}
```

## Features

- Automatic handling of two-factor authentication (2FA)
- Cookie management for maintaining session
- Configurable base URL
- Customizable number of days for schedule retrieval
- Error handling and retry logic
- Support for custom HTTP client configuration

## Configuration

You can also set a custom API endpoint using an environment variable:

```bash
export API_ENDPOINT=https://custom.api.endpoint
```

## Authentication Flow

1. Initial authentication attempt with username and password
2. If 2FA is required, the user will be prompted to enter the verification code
3. After successful authentication, a client instance is returned
4. The client maintains session cookies for subsequent requests

## Error Handling

The package includes comprehensive error handling for various scenarios:

- Network errors
- Invalid credentials
- Server errors
- Timeout errors
- Malformed responses

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

GPL2

## Disclaimer

This package is for educational purposes only. Make sure you have the necessary permissions to access the services you're connecting to.
