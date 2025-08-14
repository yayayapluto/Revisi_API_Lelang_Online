package server

import (
	"github.com/API_Lelang_Online_Go/internal/config"
	"github.com/API_Lelang_Online_Go/internal/modules/city"
	"github.com/API_Lelang_Online_Go/internal/modules/country"
	"github.com/API_Lelang_Online_Go/internal/modules/province"
	subdistrict2 "github.com/API_Lelang_Online_Go/internal/modules/subdistrict"
	"github.com/API_Lelang_Online_Go/internal/modules/village"
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
		&subdistrict2.Subdistrict{},
		&village.Village{},
	); err != nil {
		log.Fatal("Failed to migrate: ", err)
	}

	return &FiberServer{
		app: fiber.New(fiber.Config{
			ServerHeader: "api-lelang-online",
			AppName:      "api-lelang-online",
			BodyLimit:    100 * 1024 * 1024,
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
