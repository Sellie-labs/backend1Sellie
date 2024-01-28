package chat

import "encoding/json"

const (
	UserType      = "user"
	AssistantType = "assistant"
)

type ChatSession struct {
	ID             int                 `json:"id"`
	OrganizationID int                 `json:"organization_id"`
	Source         string              `json:"source"`     // e.g., IG, WhatsApp
	Identifier     string              `json:"identifier"` // e.g., IG user,whatsap number
	History        []map[string]string `json:"history"`    // Flexible structure for JSON data
	AIResponder    bool                `json:"ai_responder"`
}

// NewChatSession creates a new ChatSession instance with the provided details.
func New(organizationID int, identifier, source string) ChatSession {

	ch := make([]map[string]string, 0)
	return ChatSession{
		OrganizationID: organizationID,
		Source:         source,
		History:        ch,
		Identifier:     identifier,
		AIResponder:    true,
	}
}

// ToJSON converts the history map to a JSON string for storage.
func (c *ChatSession) ToJSON() (string, error) {
	historyJson, err := json.Marshal(c.History)
	if err != nil {
		return "", err // Return an error if the map cannot be converted to JSON
	}
	return string(historyJson), nil
}

// UpdateHistory updates the chat history with new data.
// historyJson is a byte slice of the new history JSON you want to add.
func (c *ChatSession) UpdateHistory(historyJson []byte) error {
	var newMessages []map[string]string
	if err := json.Unmarshal(historyJson, &newMessages); err != nil {
		return err // Return an error if the JSON cannot be parsed
	}

	// Append new messages to the existing history
	c.History = append(c.History, newMessages...)

	return nil
}

// AddUserMessage adds a new message from the user to the chat history.
func (c *ChatSession) AddUserMessage(content string) {
	c.History = append(c.History, map[string]string{
		"type":    UserType,
		"content": content,
	})
}

// AddAssistantMessage adds a new message from the assistant to the chat history.
func (c *ChatSession) AddAssistantMessage(content string) {
	c.History = append(c.History, map[string]string{
		"type":    AssistantType,
		"content": content,
	})
}

func (c *ChatSession) ShoudAiRespond() bool {
	return c.AIResponder
}
