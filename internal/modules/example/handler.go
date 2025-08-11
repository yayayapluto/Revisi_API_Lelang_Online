package example

import "github.com/gofiber/fiber/v2"

type Handler struct {
	useCase UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{useCase: uc}
}

func (h *Handler) List(c *fiber.Ctx) error {
	panic("implement me")
}

func (h *Handler) Create(c *fiber.Ctx) error {
	panic("implement me")
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	panic("implement me")
}

func (h *Handler) Update(c *fiber.Ctx) error {
	panic("implement me")
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	panic("implement me")
}
