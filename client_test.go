package dataedge

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SendRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/getBalance" {
			t.Errorf("Expected path /api/v1/getBalance, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("token") != "test-token" {
			t.Errorf("Expected token test-token, got %s", r.URL.Query().Get("token"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"balance": 100.50}`))
	}))
	defer ts.Close()

	client := NewClient("test-token", WithBaseURL(ts.URL+"/api/v"))
	sms := NewSmsService(client)

	balance, err := sms.GetBalance(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if balance != 100.50 {
		t.Errorf("Expected balance 100.50, got %f", balance)
	}
}

func TestSmsService_CreateSMSMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/createSmsMessage" {
			t.Errorf("Expected path /api/v1/createSmsMessage, got %s", r.URL.Path)
		}
		// Check transliteration
		msg := r.URL.Query().Get("message")
		if msg != "Privet" { // "Привет" -> "Privet"
			t.Errorf("Expected message 'Privet', got '%s'", msg)
		}

		w.Write([]byte(`{"message_id": 12345}`))
	}))
	defer ts.Close()

	client := NewClient("test-token", WithBaseURL(ts.URL+"/api/v"))
	sms := NewSmsService(client)

	// Test with defaults from client
	id, err := sms.CreateSMSMessage(context.Background(), "Привет")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if id != 12345 {
		t.Errorf("Expected ID 12345, got %d", id)
	}
}

func TestClient_APIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"error": "Invalid token"}`))
	}))
	defer ts.Close()

	client := NewClient("bad-token", WithBaseURL(ts.URL+"/api/v"))
	sms := NewSmsService(client)

	_, err := sms.GetBalance(context.Background())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}
	if apiErr.Message != "Invalid token" {
		t.Errorf("Expected error message 'Invalid token', got '%s'", apiErr.Message)
	}
}
