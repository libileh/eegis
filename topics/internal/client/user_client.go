package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/json_utils"
	"github.com/libileh/eegis/topics/internal/infra/config"
	requests "github.com/libileh/eegis/topics/internal/infra/request"
	"io"
	"net/http"
)

type UserClient struct {
	AuthProps  *config.AuthClientProperties
	BaseURL    string
	HTTPClient *http.Client
}

// NewUserClient creates a new UserClient instance with the specified base URL and authentication properties.
func NewUserClient(baseURL string, props *config.TopicProperties) *UserClient {
	return &UserClient{
		AuthProps:  &props.AuthProps,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// GetUserById retrieves a user by their ID, ensuring an auth token is loaded or fetched.
func (u *UserClient) GetUserById(id uuid.UUID) (*requests.UserDTO, error) {

	// Load or fetch the auth token
	token, err := u.LoadOrFetchAuthToken(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load or fetch auth token: %w", err)
	}

	// Create the request
	req, err := u.NewRequest("GET", fmt.Sprintf("%s/v1/users/%s", u.BaseURL, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the Authorization header with the fetched token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Execute the request
	resp, err := u.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Decode the response
	var user requests.UserDTO
	if err := json_utils.DecodeJSONResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// NewRequest creates a new HTTP request and adds the authorization header when an auth token is available.
func (u *UserClient) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	// Serialize the request body if provided.
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header if token exists.
	if u.AuthProps.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.AuthProps.AuthToken))
	}

	// Add content type header for JSON if body is provided.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
