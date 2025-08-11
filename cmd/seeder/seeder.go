package main

import (
	"github.com/API_Lelang_Online_Go/internal/config"
	"github.com/API_Lelang_Online_Go/internal/modules/example"
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

	for i := 0; i < 10; i++ {
		country := example.Country{
			Kode: gofakeit.Country(),
			Nama: gofakeit.CountryAbr(),
		}
		if err := db.Create(&country).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding selesai")
}
