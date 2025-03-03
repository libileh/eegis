package client

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/json_utils"
	"net/http"
)

type HttpUserService struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHttpUserService(baseURL string) *HttpUserService {
	return &HttpUserService{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// CheckRolePrecedence Method to check role precedence
func (s *HttpUserService) CheckRolePrecedence(user *auth.CtxUser) (bool, error) {
	// Build the request URL
	url := fmt.Sprintf("%s/v1/users/role-precedence?userId=%s&role=%s", s.BaseURL, user.ID, user.ContextRole.Description)

	// Send HTTP GET request using the HTTP sender
	resp, err := s.sendGetRequest(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Decode JSON response
	var result struct {
		Allowed bool `json:"data"`
	}
	if err := json_utils.DecodeJSONResponse(resp, &result); err != nil {
		return false, err
	}

	return result.Allowed, nil
}

// Helper: Send an HTTP GET request
func (s *HttpUserService) sendGetRequest(url string) (*http.Response, error) {
	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}

	// Check for non-OK HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return resp, nil
}

func (s *HttpUserService) GetUser(uuid2 uuid.UUID) (*UserDTO, error) {
	resp, err := s.sendGetRequest(fmt.Sprintf("%s/v1/users/%s", s.BaseURL, uuid2))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var user UserDTO
	if err := json_utils.DecodeJSONResponse(resp, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
