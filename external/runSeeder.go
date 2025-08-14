package main

import (
	"github.com/API_Lelang_Online_Go/internal/config"
	"github.com/API_Lelang_Online_Go/internal/modules/country"
	"github.com/API_Lelang_Online_Go/internal/modules/province"
	"github.com/API_Lelang_Online_Go/pkg/database"
	"log"
	"sync"
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
	); err != nil {
		log.Fatal("Failed to migrate: ", err)
	}

	wg := sync.WaitGroup{}

	entitySeeder := database.NewEntitySeeder(db)
	go func() {
		wg.Add(1)
		defer wg.Done()

		entitySeeder.CountrySeeder(100)
	}()

	wg.Wait()

	go func() {
		wg.Add(1)
		defer wg.Done()

		entitySeeder.ProvinceSeeder(87)
	}()

}
