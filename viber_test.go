package dataedge

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestViberService_SendQuickViberMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expect POST to /api/v2/sendQuickViberMessage based on our implementation logic (v=2 for Viber usually in PHP SDK, but we pass empty string by default in method, wait.
		// In `SendQuickViberMessage`, we call `s.client.sendRequest(..., "sendQuickViberMessage", params, "")`
		// `sendRequest` defaults v to "1" if empty.
		// BUT wait, does PHP SDK usage of `SendQuickViberMessage` (not image) use v=1 or v=2?
		// `ViberService.php`:
		// `sendQuickViberMessage` -> `sendRequest(..., 'post')` -> defaults v=1.
		// `sendQuickViberMessageWithImage` -> `sendRequest(..., 'post', 2)` -> v=2.
		// So my Go implementation of `SendQuickViberMessage` using default v (which is 1) mimics PHP.
		// So URL should be /api/v1/sendQuickViberMessage (if implicit) or /api/vsendQuickViberMessage if strict PHP copy?
		// My client logic: `reqURL = fmt.Sprintf("%s%s/%s", c.baseURL, v, command)` -> `.../api/v1/sendQuickViberMessage`.

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if !strings.Contains(r.URL.Path, "sendQuickViberMessage") {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		// Check body params
		body, _ := io.ReadAll(r.Body)
		// For normal POST we use form-urlencoded in `sendRequest`?
		// "application/x-www-form-urlencoded"
		// Let's decode or just check string
		bs := string(body)
		if !strings.Contains(bs, "message=Hello") {
			t.Errorf("Body missing message: %s", bs)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id": [111]}`)) // Assuming array or single ID
	}))
	defer ts.Close()

	client := NewClient("token", WithBaseURL(ts.URL+"/api/v"), WithTranslit(true))
	viber := NewViberService(client)
	// Test with defaults from client
	_, err := viber.SendQuickViberMessage(context.Background(), "375290000000", "Hello")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
