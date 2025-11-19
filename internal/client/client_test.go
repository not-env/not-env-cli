package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGet(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Authorization header, got %s", r.Header.Get("Authorization"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	resp, err := client.Get("/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestClientPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "test" {
			t.Errorf("expected name=test, got %v", body["name"])
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	resp, err := client.Post("/test", map[string]interface{}{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		body           string
		wantErr        bool
		expectedResult map[string]interface{}
	}{
		{
			name:       "success response",
			statusCode: 200,
			body:       `{"key": "value"}`,
			wantErr:    false,
			expectedResult: map[string]interface{}{"key": "value"},
		},
		{
			name:       "error response",
			statusCode: 400,
			body:       `{"error": "Bad Request", "message": "invalid input"}`,
			wantErr:    true,
		},
		{
			name:       "not found",
			statusCode: 404,
			body:       `{"error": "Not Found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewClient(server.URL, "test-key")
			resp, err := client.Get("/test")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var result map[string]interface{}
			err = ParseResponse(resp, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResponse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && result["key"] != tt.expectedResult["key"] {
				t.Errorf("ParseResponse() result = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

