package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AgentClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type SendMessageRequest struct {
	Message  string  `json:"message"`
	ThreadID *string `json:"threadId,omitempty"`
}

type AgentResponse struct {
	ThreadID  string    `json:"threadId"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type UserProfile struct {
	PreferredAgentInstructions string `json:"preferredAgentInstructions,omitempty"`
	CustomWorkflowsJson        string `json:"customWorkflowsJson,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewAgentClient(baseURL, accessToken string) *AgentClient {
	return &AgentClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: accessToken,
	}
}

func (c *AgentClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Message)
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

func (c *AgentClient) SendMessage(ctx context.Context, message string, threadID *string) (*AgentResponse, error) {
	req := SendMessageRequest{
		Message:  message,
		ThreadID: threadID,
	}

	var resp AgentResponse
	if err := c.doRequest(ctx, "POST", "/api/agent/send", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *AgentClient) GetProfile(ctx context.Context) (*UserProfile, error) {
	var profile UserProfile
	if err := c.doRequest(ctx, "GET", "/api/profile", nil, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (c *AgentClient) UpdateProfile(ctx context.Context, profile *UserProfile) error {
	return c.doRequest(ctx, "PUT", "/api/profile", profile, nil)
}

func (c *AgentClient) ListConversations(ctx context.Context) ([]ConversationMetadata, error) {
	var conversations []ConversationMetadata
	if err := c.doRequest(ctx, "GET", "/api/agent/conversations", nil, &conversations); err != nil {
		return nil, err
	}
	return conversations, nil
}

type ConversationMetadata struct {
	ID            int       `json:"id"`
	ThreadID      string    `json:"threadId"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"createdAt"`
	LastMessageAt time.Time `json:"lastMessageAt"`
}
