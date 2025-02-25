package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
)

// tokenResponse represents the expected JSON structure from the token API.
type tokenResponse struct {
	Token string `json:"data"`
}

// LoadOrFetchAuthToken checks if an auth token is present in the environment,
// and if not, calls an external API endpoint to retrieve a new token.
func (u *UserClient) LoadOrFetchAuthToken(userId uuid.UUID) (string, error) {
	// First, try to load the token from environment variables.
	if token := u.AuthProps.AuthToken; token != "" && token != "secret" {
		return token, nil
	}

	// If no token is provided, attempt to fetch a new one from the token generation API.
	endpoint := fmt.Sprintf("%s?user_id=%s", u.AuthProps.AuthEndpoint, userId.String())
	if endpoint == "" {
		return "", errors.New("auth endpoint not provided in environment")
	}

	resp, err := u.HTTPClient.Get(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch token: non-OK HTTP status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result tokenResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Token == "" {
		return "", errors.New("received empty token from endpoint")
	}

	return result.Token, nil
}
