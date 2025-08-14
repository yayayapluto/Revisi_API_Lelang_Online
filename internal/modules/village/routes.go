package village

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	repo := NewRepository(db)
	uc := NewUseCase(repo)
	h := NewHandler(uc)

	rg := app.Group("/api/villages")
	rg.Get("/", h.List)
	rg.Post("/", h.Create)
	rg.Get("/:id", h.GetByID)
	rg.Put("/:id", h.Update)
	rg.Delete("/:id", h.Delete)
	rg.Post("/uploadBatchData", h.Batch)
}
