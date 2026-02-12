package dataedge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ViberService provides access to Viber-related API methods.
type ViberService struct {
	client *Client
}

// NewViberService creates a new ViberService.
func NewViberService(client *Client) *ViberService {
	return &ViberService{client: client}
}

// CreateViberMessage creates a new Viber message template.
func (s *ViberService) CreateViberMessage(ctx context.Context, message string) (interface{}, error) {
	params := map[string]interface{}{
		"message":      message,
		"vibername_id": s.client.viberNameID,
	}

	resp, err := s.client.sendRequest(ctx, http.MethodPost, "createViberMessage", params, "")
	if err != nil {
		return nil, err
	}

	// Assuming response returns an ID or object
	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SendViberMessage sends a created Viber message.
func (s *ViberService) SendViberMessage(ctx context.Context, phone string, viberMessageID int) (interface{}, error) {
	params := map[string]interface{}{
		"phone":            phone,
		"viber_message_id": viberMessageID,
	}

	resp, err := s.client.sendRequest(ctx, http.MethodPost, "sendViberMessage", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SendQuickViberMessage sends a Viber message without creating it first.
func (s *ViberService) SendQuickViberMessage(ctx context.Context, phone string, message string) (interface{}, error) {
	params := map[string]interface{}{
		"phone":        phone,
		"vibername_id": s.client.viberNameID,
		"message":      message,
	}

	resp, err := s.client.sendRequest(ctx, http.MethodPost, "sendQuickViberMessage", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SendQuickViberMessageWithImage sends a Viber message with an image.
// This function handles multipart form data upload.
func (s *ViberService) SendQuickViberMessageWithImage(ctx context.Context, phone string, message string, typeMessage string, imagePath string, buttonText string, buttonLink string) (interface{}, error) {
	// PHP passes `v=2` explicitly here.
	v := "2"

	fields := map[string]string{
		"phone":        phone,
		"vibername_id": fmt.Sprintf("%d", s.client.viberNameID),
		"message":      message,
		"type_message": typeMessage,
	}

	if buttonText != "" {
		fields["button"] = buttonText
	}
	if buttonLink != "" {
		fields["button_link"] = buttonLink
	}

	files := make(map[string]string)
	if imagePath != "" {
		// Verify file exists
		if _, err := os.Stat(imagePath); err != nil {
			return nil, fmt.Errorf("image file not found: %w", err)
		}
		files["image"] = imagePath
	} else if typeMessage == "IMAGE" {
		return nil, fmt.Errorf("image path required for IMAGE message type")
	}

	resp, err := s.client.sendMultipartRequest(ctx, "sendQuickViberMessage", fields, files, v)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// SendViberMessageList sends a predictable Viber message list.
func (s *ViberService) SendViberMessageList(ctx context.Context, name string, listID int, dSchedule string, extraParams map[string]interface{}) (interface{}, error) {
	params := map[string]interface{}{
		"name":         name,
		"vibername_id": s.client.viberNameID,
		"list_id":      listID,
		"d_schedule":   dSchedule,
	}
	for k, v := range extraParams {
		params[k] = v
	}

	resp, err := s.client.sendRequest(ctx, http.MethodPost, "sendViberMessageList", params, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetViberMessageList returns the list of Viber messages.
func (s *ViberService) GetViberMessageList(ctx context.Context) (interface{}, error) {
	resp, err := s.client.sendRequest(ctx, http.MethodGet, "getViberMessageList", nil, "")
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}
