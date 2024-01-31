package chat

import (
	"fmt"
)

type Service struct {
	repo DBRepository
	ai   AIChat
}

func NewService(repo DBRepository, ai AIChat) *Service {
	return &Service{
		repo,
		ai,
	}
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
