package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Service struct {
	repo  DBRepository
	ai    AIChat
	token string
}

func NewService(repo DBRepository, ai AIChat, token string) *Service {
	return &Service{
		repo,
		ai,
		token,
	}
}

func (s *Service) sendFacebookTextMessage(recipientNumber, messageContent string) error {
	client := &http.Client{}
	url := "https://graph.facebook.com/v18.0/124306174110269/messages"

	reqBody, err := json.Marshal(map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                recipientNumber,
		"type":              "text",
		"text": map[string]string{
			"body": messageContent,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.token))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("Message sent successfully")
	} else {
		responseBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("Failed to send message. Status Code: %d. Response: %s\n", resp.StatusCode, string(responseBody))
	}

	return nil
}

// CreateChatSession saves a new ChatSession to the repository.
func (s *Service) CreateChatSession(chatIdentifier, source string, organizationID int) (*ChatSession, error) {
	fmt.Println("Chat created")
	c := New(organizationID, chatIdentifier, source)
	if err := s.repo.Save(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// Interact with the chat
func (s *Service) Chat(message, source, chatIdentifier string, organizationID int) error {
	c, check, err := s.repo.FindByIdentifier(chatIdentifier)
	if err != nil && check != true {
		return err
	}
	//If chat dont exit create it
	if check == false {
		c, err = s.CreateChatSession(chatIdentifier, source, organizationID)
	}
	if err != nil {
		return err
	}

	c.AddUserMessage(message)
	if c.ShoudAiRespond() {
		answer, err := s.ai.GenerateResponse(message, c.Identifier, c.OrganizationID)
		if err != nil {
			return err
		}
		if err := s.sendFacebookTextMessage(chatIdentifier, answer); err != nil {
			return fmt.Errorf("failed to send message through Facebook API: %v", err)
		}
		c.AddAssistantMessage(answer)
	}

	if err := s.repo.Update(c); err != nil {
		return err
	}

	fmt.Println(c.History)
	//TODO:send the messsage to the api, whatsapp, ig, etc

	return nil
}

// GetChatSessionByID retrieves a ChatSession by its identifier.
func (s *Service) GetChatSessionByIdentifier(identifier string) (*ChatSession, bool, error) {
	return s.repo.FindByIdentifier(identifier)
}

// GetChatSessionByID retrieves a ChatSession by its ID.
func (s *Service) GetChatSessionByID(id int) (*ChatSession, error) {
	return s.repo.FindById(id)
}
