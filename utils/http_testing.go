package utils

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func MakeRequest(method string, t *testing.T, app *fiber.App, url string, payload interface{}, headers map[string]string) (map[string]string, *http.Response) {
	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	var req *http.Request

	// Create a new HTTP request
	if method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		req, err = http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
		req.Header.Add("Content-Type", "application/json")

	}
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Test the request using Fiber's testing framework
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		t.Errorf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	// Decode the response body into a map
	respBody := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	return respBody, resp
}