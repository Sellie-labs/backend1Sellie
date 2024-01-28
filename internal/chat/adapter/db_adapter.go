package adapter

import (
	"admin/internal/chat"
	"admin/pkg/apperror"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

type SQLChatSessionRepository struct {
	db *sql.DB
}

func NewSQLChatSessionRepository(db *sql.DB) chat.DBRepository {
	return &SQLChatSessionRepository{db: db}
}

func (repo *SQLChatSessionRepository) FindById(id int) (*chat.ChatSession, error) {
	query := `SELECT id, organization_id, source, history, identifier, ai_responder FROM chat_session WHERE id = $1`
	var cs chat.ChatSession
	var historyJson []byte
	err := repo.db.QueryRow(query, id).Scan(&cs.ID, &cs.OrganizationID, &cs.Source, &historyJson, &cs.Identifier, &cs.AIResponder)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(fmt.Sprintf("Chat session with ID %d not found", id))
		}
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching chat session by ID: %v", err))
	}

	if err := json.Unmarshal(historyJson, &cs.History); err != nil {
		return nil, apperror.NewInternalError("Error unmarshalling chat session history")
	}

	return &cs, nil
}

func (repo *SQLChatSessionRepository) FindByIdentifier(identifier string) (*chat.ChatSession, error) {
	log.Info("FindByIdentifier")
	query := `SELECT id, organization_id, source, history, identifier, ai_responder FROM chat_session WHERE identifier = $1`
	var cs chat.ChatSession
	var historyJson []byte
	err := repo.db.QueryRow(query, identifier).Scan(&cs.ID, &cs.OrganizationID, &cs.Source, &historyJson, &cs.Identifier, &cs.AIResponder)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperror.NewNotFoundError(fmt.Sprintf("Chat session with identifier %s not found", identifier))
		}
		return nil, apperror.NewInternalError(fmt.Sprintf("Error fetching chat session by identifier: %v", err))
	}

	if err := json.Unmarshal(historyJson, &cs.History); err != nil {
		return nil, apperror.NewInternalError("Error unmarshalling chat session history")
	}
	log.Info("FindByIdentifier2")

	return &cs, nil
}

func (repo *SQLChatSessionRepository) Update(cs *chat.ChatSession) error {
	historyJson, err := json.Marshal(cs.History)
	if err != nil {
		return apperror.NewInternalError("Error marshalling chat session history")
	}

	query := `UPDATE chat_session SET organization_id = $1, source = $2, history = $3, identifier = $4, ai_responder = $5 WHERE id = $6`
	_, err = repo.db.Exec(query, cs.OrganizationID, cs.Source, historyJson, cs.Identifier, cs.AIResponder, cs.ID)
	if err != nil {
		return apperror.NewInternalError(fmt.Sprintf("Error updating chat session: %v", err))
	}

	return nil
}

func (repo *SQLChatSessionRepository) Save(cs *chat.ChatSession) error {
	historyJson, err := json.Marshal(cs.History)
	if err != nil {
		return apperror.NewInternalError("Error marshalling chat session history")
	}

	query := `INSERT INTO chat_session (organization_id, source, history, identifier, ai_responder) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = repo.db.QueryRow(query, cs.OrganizationID, cs.Source, historyJson, cs.Identifier, cs.AIResponder).Scan(&cs.ID)
	if err != nil {
		return apperror.NewInternalError(fmt.Sprintf("Error saving chat session: %v", err))
	}

	return nil
}
