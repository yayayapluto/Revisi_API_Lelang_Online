package main

import (
	"github.com/API_Lelang_Online_Go/internal/config"
	country2 "github.com/API_Lelang_Online_Go/internal/modules/country"
	"github.com/API_Lelang_Online_Go/pkg/database"
	"github.com/brianvoe/gofakeit/v7"
	"log"
	"time"
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

	if err := gofakeit.Seed(time.Now().UnixNano()); err != nil {
		return
	}

	if err := db.Exec("TRUNCATE TABLE countries RESTART IDENTITY").Error; err != nil {
		log.Fatal("Cannot truncate countries: ", err)
	}

	for i := 0; i < 100; i++ {
		country := country2.Country{
			Kode:  gofakeit.CountryAbr(),
			Nama:  gofakeit.Country(),
			Nomor: gofakeit.Numerify("##"),
		}
		if err := db.Create(&country).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding selesai")
}
