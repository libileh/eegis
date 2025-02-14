package json_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func JsonResponse(w http.ResponseWriter, status int, data interface{}) error {
	type Response struct {
		Data interface{} `json:"data"`
	}
	return WriteJSON(w, status, &Response{Data: data})
}

// Helper: Decode JSON response
func DecodeJSONResponse(resp *http.Response, target interface{}) error {
	if resp.Body == nil {
		return errors.New("response body is nil")
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return nil
}
