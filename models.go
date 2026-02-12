package dataedge

// Common response structures

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// BalanceResponse represents the response for GetBalance
type BalanceResponse struct {
	Balance float64 `json:"balance"` // Assuming balance is numeric based on usage
	// If it's a string, we might need to adjust.
}

// IDResponse represents a response containing just an ID (common for create/send methods)
type IDResponse struct {
	ID int `json:"id"`
	// Sometimes it might be message_id, sms_id, etc.
	// We will define specific structs if needed or use tags.
}

// SmsMessageResponse for createSMSMessage
type SmsMessageResponse struct {
	MessageID int `json:"message_id"`
	// Add other fields if known
}

// SendSmsResponse for sendSms
type SendSmsResponse struct {
	ID int `json:"id"`
	// check actual response format from examples or docs if available
}

// ChecksmsResponse
type CheckSMSResponse struct {
	Status string `json:"status"`
	// Add timestamp etc
}

// Since we don't have the exact JSON schema, we will use map[string]interface{} for some generic returns
// or simpler structs where we are reasonably sure (like ID).
// For strict static typing, we'd need the docs.
// I will provide specific structs for the methods we implement based on PHP usage.

// SMSIDResponse creates a common struct for responses returning an ID
type SMSIDResponse struct {
	MessageID int `json:"message_id"`
}

// SendQuickSmsResponse
type SendQuickSmsResponse struct {
	IDs []int `json:"id"` // Usually returns an array of IDs if bulk, or single ID?
	// PHP example prints $result.
}
