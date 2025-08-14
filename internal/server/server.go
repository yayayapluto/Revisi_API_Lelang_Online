package server

import (
	"github.com/API_Lelang_Online_Go/internal/config"
	"github.com/API_Lelang_Online_Go/internal/modules/city"
	"github.com/API_Lelang_Online_Go/internal/modules/country"
	"github.com/API_Lelang_Online_Go/internal/modules/province"
	"github.com/API_Lelang_Online_Go/pkg/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
)

type FiberServer struct {
	app *fiber.App
	db  *gorm.DB
}

func New() *FiberServer {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config")
	}

	db, err := database.StartDB(conf.DBHost, conf.DBUser, conf.DBPassword, conf.DBName, conf.DBPort)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	if err := db.AutoMigrate(
		&country.Country{},
		&province.Province{},
		&city.City{},
	); err != nil {
		log.Fatal("Failed to migrate: ", err)
	}

	return &FiberServer{
		app: fiber.New(fiber.Config{
			ServerHeader: "api-lelang-online",
			AppName:      "api-lelang-online",
		}),
		db: db,
	}
}

func (s *FiberServer) RegisterFiberServer() {
	RegisterFiberRoutes(s.app, s.db)
}

func (s *FiberServer) Listen(addr string) {
	if err := s.app.Listen(addr); err != nil {
		log.Fatalf("Failed to start Fiber server: %v", err)
	}
}
