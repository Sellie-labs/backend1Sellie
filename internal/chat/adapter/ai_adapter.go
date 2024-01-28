package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ChatAdapter implements the AIChat interface.
type ChatAdapter struct {
	serverURL string
}

// NewChatAdapter creates a new instance of ChatAdapter with the specified server URL.
func NewChatAdapter(serverURL string) *ChatAdapter {
	return &ChatAdapter{
		serverURL: serverURL,
	}
}

// GenerateResponse sends a message to the chat server and returns the response.
func (c *ChatAdapter) GenerateResponse(message string, chatIdentifier string, organisationID int) (string, error) {
	// Construct the request payload.
	payload := map[string]string{
		"content":         message,
		"chat_identifier": chatIdentifier,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// Create the HTTP request.
	req, err := http.NewRequest("POST", c.serverURL+"/v1/api/chat/session", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check for a non-200 status code.
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("received non-200 status code from chat server")
	}

	// Return the response as a string.
	return string(body), nil
}
