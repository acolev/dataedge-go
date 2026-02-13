package dataedge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://app.dataedge.md/api/v"
	defaultVersion = "1"
)

// Client is the DataEdge API client.
type Client struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	debug       bool
	alphaNameID int
	viberNameID int
	translit    bool
}

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithDebug enables or disables debug logging.
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithBaseURL sets a custom base URL for the client.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithAlphaNameID sets the default alpha name ID for SMS.
func WithAlphaNameID(id int) Option {
	return func(c *Client) {
		c.alphaNameID = id
	}
}

// WithViberNameID sets the default viber name ID for Viber.
func WithViberNameID(id int) Option {
	return func(c *Client) {
		c.viberNameID = id
	}
}

// WithTranslit enables or disables auto-transliteration by default.
func WithTranslit(translit bool) Option {
	return func(c *Client) {
		c.translit = translit
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// NewClient creates a new DataEdge API client.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		baseURL:  defaultBaseURL,
		token:    token,
		translit: true, // Default to true like in PHP SDK
		httpClient: &http.Client{
			Timeout: 20 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// APIError represents an error returned by the DataEdge API.
type APIError struct {
	Message string `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("DataEdge API error: %s", e.Message)
}

// sendRequest sends an HTTP request to the DataEdge API.
func (c *Client) sendRequest(ctx context.Context, method, command string, params map[string]interface{}, v string) ([]byte, error) {
	if v == "" {
		v = defaultVersion
	}

	var reqURL string
	var body io.Reader
	var contentType string

	if method == http.MethodGet {
		u, err := url.Parse(fmt.Sprintf("%s%s/%s", c.baseURL, v, command))
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}
		q := u.Query()
		q.Set("token", c.token)
		for k, v := range params {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		u.RawQuery = q.Encode()
		reqURL = u.String()
	} else {
		// POST request
		// Note: We include version in URL for consistency, matching GET requests and ViberService requirements.
		reqURL = fmt.Sprintf("%s%s/%s", c.baseURL, v, command)

		params["token"] = c.token

		// Check for multipart (images)
		if _, ok := params["_multipart_writer"]; ok {
			// Special handling if we passed a multipart writer (custom logic for Go)
			// But for a generic request, we usually map map[string]interface{} to fields.
			// Let's handle multipart separately or detect it.
			// For this implementation, let's stick to form-urlencoded for standard POST, as PHP uses `http_build_query`.
			// PHP: `curl_setopt($ch, CURLOPT_POSTFIELDS, http_build_query($params));` -> application/x-www-form-urlencoded.

			form := url.Values{}
			for k, v := range params {
				form.Set(k, fmt.Sprintf("%v", v))
			}
			body = strings.NewReader(form.Encode())
			contentType = "application/x-www-form-urlencoded"
		} else {
			form := url.Values{}
			for k, v := range params {
				form.Set(k, fmt.Sprintf("%v", v))
			}
			body = strings.NewReader(form.Encode())
			contentType = "application/x-www-form-urlencoded"
		}
	}

	// Handle multipart specifically for Viber image upload if needed later.
	// The PHP SDK `sendQuickViberMessageWithImage` does manual CURLFile stuff.
	// We will handle that in the specific service method or a helper, but `sendRequest` generic might need to support it.
	// For now, let's keep `sendRequest` simple for standard params.

	if c.debug {
		fmt.Printf("[DEBUG] Request: %s %s\n", method, reqURL)
		if contentType != "" {
			fmt.Printf("[DEBUG] Content-Type: %s\n", contentType)
		}
	}

	req, err := http.NewRequestWithContext(ctx, strings.ToUpper(method), reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response Status: %s\n", resp.Status)
		fmt.Printf("[DEBUG] Response Body: %s\n", string(respBody))
	}

	// Check logical error in JSON
	var errResp APIError
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
		return nil, &errResp
	}

	// Double check if response is false/null which might indicate error in PHP SDK logic,
	// but here we just return body.

	return respBody, nil
}

// Helper to handle Multipart requests for Viber image
func (c *Client) sendMultipartRequest(ctx context.Context, command string, fields map[string]string, files map[string]string, v string) ([]byte, error) {
	if v == "" {
		v = "2" // default for viber with image in PHP is 2
	}
	reqURL := fmt.Sprintf("%s%s/%s", c.baseURL, v, command)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add fields
	fields["token"] = c.token
	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	// Add files
	for key, path := range files {
		file, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(key, filepath.Base(path))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var errResp APIError
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
		return nil, &errResp
	}

	return respBody, nil
}
