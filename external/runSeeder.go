package main

import (
	"github.com/API_Lelang_Online_Go/internal/config"
	"github.com/API_Lelang_Online_Go/internal/modules/city"
	"github.com/API_Lelang_Online_Go/internal/modules/country"
	"github.com/API_Lelang_Online_Go/internal/modules/province"
	"github.com/API_Lelang_Online_Go/pkg/database"
	"log"
)

func main() {
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

	entitySeeder := database.NewEntitySeeder(db)
	entitySeeder.CountrySeeder(100)
	entitySeeder.ProvinceSeeder(87)
	entitySeeder.CitySeeder(92)
	entitySeeder.SubdistrictSeeder(54)
	entitySeeder.VillageSeeder(74)
}
