package aiservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type ExtractFaceResult struct {
	Success          bool      `json:"success"`
	FaceDetected     bool      `json:"face_detected"`
	LivenessScore    float64   `json:"liveness_score"`
	IsLive           bool      `json:"is_live"`
	Embedding        []float64 `json:"embedding"`
	ProcessingTimeMs float64   `json:"processing_time_ms"`
	Message          string    `json:"message"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) ExtractFaceAndVerifyLiveness(ctx context.Context, photoBytes []byte, filename string) (*ExtractFaceResult, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart form: %w", err)
	}

	if _, err := part.Write(photoBytes); err != nil {
		return nil, fmt.Errorf("failed to write photo bytes: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/extract-face", &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to ai-face-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai-face-service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result ExtractFaceResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ai service response: %w", err)
	}

	return &result, nil
}
