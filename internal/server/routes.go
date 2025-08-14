package server

import (
	"github.com/API_Lelang_Online_Go/internal/modules/country"
	"github.com/API_Lelang_Online_Go/internal/modules/province"
	"github.com/API_Lelang_Online_Go/internal/shared/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func RegisterFiberRoutes(app *fiber.App, db *gorm.DB) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("hi there, the correct url is within /api")
	})

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.JSON("hi there, yeah it is the right url, u need the module name tho")
	})

	// All Routes goes here
	country.RegisterRoutes(app, db)
	province.RegisterRoutes(app, db)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(shared.ErrorResponse("Route not found", nil))
	})
}
