package chat

type DBRepository interface {
	FindById(id int) (*ChatSession, error)
	FindByIdentifier(identifier string) (*ChatSession, error)
	Update(chatSession *ChatSession) error
	Save(chatSession *ChatSession) error
}

type AIChat interface {
	GenerateResponse(message string, chatIdentifier string, organisationID int) (string, error)
}
