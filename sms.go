package dataedge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SmsService provides access to SMS-related API methods.
type SmsService struct {
	client *Client
}

// NewSmsService creates a new SmsService.
// Note: translit option from PHP is handled per method if needed, or we can add it to struct.
// PHP: `public function __construct(string $token, bool $translit = true)`
// We will follow the pattern of creating a service attached to a client.
func NewSmsService(client *Client) *SmsService {
	return &SmsService{client: client}
}

// GetBalance returns the current balance.
func (s *SmsService) GetBalance(ctx context.Context) (float64, error) {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, "getBalance", nil, "")
	if err != nil {
		return 0, err
	}

	// The API might return {"balance": "100.50"} (string) or {"balance": 100.50} (number)
	// We use a custom struct to handle both or just try to unmarshal.
	var result struct {
		Balance interface{} `json:"balance"`
	}

	if err := json.Unmarshal(resp, &result); err == nil {
		switch v := result.Balance.(type) {
		case float64:
			return v, nil
		case string:
			var f float64
			if _, err := fmt.Sscanf(v, "%f", &f); err == nil {
				return f, nil
			}
		case int:
			return float64(v), nil
		}
	}

	// Fallback: try to parse the whole body as a number
	var balance float64
	if err := json.Unmarshal(resp, &balance); err == nil {
		return balance, nil
	}

	return 0, fmt.Errorf("failed to parse balance from response: %s (error: %v)", string(resp), err)
}

// GetLimit returns the current limit.
func (s *SmsService) GetLimit(ctx context.Context) (int, error) {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, "getLimit", nil, "")
	if err != nil {
		return 0, err
	}

	// Similar assumption to GetBalance
	var result struct {
		Limit int `json:"limit"`
	}
	if err := json.Unmarshal(resp, &result); err == nil {
		return result.Limit, nil
	}

	var limit int
	if err := json.Unmarshal(resp, &limit); err == nil {
		return limit, nil
	}

	return 0, fmt.Errorf("unexpected response format: %s", string(resp))
}

// CreateSMSMessage creates a new SMS message.
func (s *SmsService) CreateSMSMessage(ctx context.Context, message string, alphaNameID int, translit bool) (int, error) {
	params := map[string]interface{}{
		"message": message,
	}
	if translit {
		params["message"] = Transliterate(message)
	}
	if alphaNameID > 0 {
		params["alphaname_id"] = alphaNameID
	}

	resp, err := s.client.sendRequest(ctx, http.MethodGet, "createSmsMessage", params, "")
	if err != nil {
		return 0, err
	}

	// var result []struct {
	// 	ID int `json:"id"`
	// }
	// Actually PHP returns `$result` which is `json_decode`.
	// CreateSMSMessage implies creating ONE message.
	// But `sendRequest` parses JSON.

	// Let's try to parse as single object first
	var singleRes struct {
		ID int `json:"message_id"` // See SmsService.php line 18 in example: $message->message_id
	}
	if err := json.Unmarshal(resp, &singleRes); err == nil && singleRes.ID != 0 {
		return singleRes.ID, nil
	}

	// Try generic map if needed to debug
	return 0, fmt.Errorf("failed to parse response or invalid format: %s", string(resp))
}

// SendSms sends a created SMS message.
func (s *SmsService) SendSms(ctx context.Context, messageID int, phone string) (interface{}, error) {
	params := map[string]interface{}{
		"message_id": messageID,
		"phone":      phone,
	}

	resp, err := s.client.sendRequest(ctx, http.MethodGet, "sendSms", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SendQuickSms sends an SMS without creating it first.
func (s *SmsService) SendQuickSms(ctx context.Context, message string, phone string, translit bool) (interface{}, error) {
	if message == "" || phone == "" {
		return nil, fmt.Errorf("message and phone are required")
	}

	params := map[string]interface{}{
		"message": message,
		"phone":   phone,
	}
	if translit {
		params["message"] = Transliterate(message)
	}

	resp, err := s.client.sendRequest(ctx, http.MethodGet, "sendQuickSms", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CheckSMS checks the status of a sent SMS.
func (s *SmsService) CheckSMS(ctx context.Context, smsID int) (interface{}, error) {
	params := map[string]interface{}{
		"sms_id": smsID,
	}

	resp, err := s.client.sendRequest(ctx, http.MethodGet, "checkSMS", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetMessagesList returns the list of messages.
func (s *SmsService) GetMessagesList(ctx context.Context) (interface{}, error) {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, "getMessagesList", nil, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAlphaNames returns the list of alpha names.
func (s *SmsService) GetAlphaNames(ctx context.Context) (interface{}, error) {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, "getAlphanames", nil, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAlphaNameID returns the ID of an alpha name.
func (s *SmsService) GetAlphaNameID(ctx context.Context, name string) (interface{}, error) {
	params := map[string]interface{}{
		"name": name,
	}

	resp, err := s.client.sendRequest(ctx, http.MethodGet, "getAlphanameId", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}
