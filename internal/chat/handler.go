package chat

import (
	"admin/pkg/apperror"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	s Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s: s,
	}
}

// RegisterRoutes registers JWT-protected routes on the provided fiber App.
func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/sellie/chat/:organisation", h.Chat)
}

func (h *Handler) Chat(c *fiber.Ctx) error {
	org := 1
	body := struct {
		Input    string `json:"Input"`
		SenderID string `json:"SenderID"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		apperror.NewBadRequestError(err.Error())
	}

	h.s.Chat(body.Input, "IG", body.SenderID, org)
	//Logic to get the identifier
	//get the organisation
	//call chat
	return c.JSON(fiber.Map{"message": "ok"})
}
