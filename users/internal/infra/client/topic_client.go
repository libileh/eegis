package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/libileh/eegis/common/auth"
	"github.com/libileh/eegis/common/errors"
	commonProps "github.com/libileh/eegis/common/properties"
	"github.com/libileh/eegis/users/internal/config"
	"io"
	"net/http"
)

type TopicHTTPClient struct {
	HTTPClient *http.Client
	TopicProps config.TopicProperties
	AuthProps  commonProps.AuthProperties
	Auth       auth.Authenticator
}

func NewTopicHTTPClient(props *config.UserProperties, authenticator *auth.JWTAuthenticator) *TopicHTTPClient {
	return &TopicHTTPClient{
		HTTPClient: &http.Client{},
		TopicProps: props.Topics,
		AuthProps:  props.CommonProps.AuthProps,
		Auth:       authenticator,
	}
}

func (t *TopicHTTPClient) FollowTopic(topicId, userId uuid.UUID, claims jwt.MapClaims) *errors.CustomError {
	token, err := t.Auth.GenerateToken(claims)
	if err != nil {
		return errors.NewInternalServerError("failed to generate token: %v", err)
	}
	// Create a new request to the topics endpoint.
	req, err := t.NewRequest("POST", fmt.Sprintf("%s/v1/topics/%s/followers/%s", t.TopicProps.BaseUrl, topicId, userId), nil)
	if err != nil {
		return errors.NewInternalServerError("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return errors.NewInternalServerError("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.NewInternalServerError(fmt.Sprintf("received non-OK HTTP status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode)))
	}

	return nil
}

func (t *TopicHTTPClient) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
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

	// Add content type header for JSON if body is provided.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
