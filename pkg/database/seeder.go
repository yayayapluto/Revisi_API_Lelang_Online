package database

import (
	"fmt"
	country2 "github.com/API_Lelang_Online_Go/internal/modules/country"
	province2 "github.com/API_Lelang_Online_Go/internal/modules/province"
	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
	"log"
)

type entitySeeder struct {
	db *gorm.DB
}

type EntitySeeder interface {
	CountrySeeder(total int)
	ProvinceSeeder(total int)
}

func NewEntitySeeder(db *gorm.DB) EntitySeeder {
	return &entitySeeder{db: db}
}

func (e *entitySeeder) CountrySeeder(total int) {
	if err := e.db.Exec("TRUNCATE TABLE countries RESTART IDENTITY").Error; err != nil {
		log.Fatal("Cannot truncate countries: ", err)
	}

	for i := 0; i < total; i++ {
		country := country2.Country{
			Kode:  gofakeit.CountryAbr(),
			Nama:  gofakeit.Country(),
			Nomor: gofakeit.Numerify("##"),
		}
		if err := e.db.Create(&country).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding country done!")
}

func (e *entitySeeder) ProvinceSeeder(total int) {
	if err := e.db.Exec("TRUNCATE TABLE provinces RESTART IDENTITY").Error; err != nil {
		log.Fatal("Cannot truncate provinces: ", err)
	}

	for i := 0; i < total; i++ {
		cd := gofakeit.Numerify("##")

		province := province2.Province{
			Code:     cd,
			Nama:     fmt.Sprintf("Province %v", gofakeit.LoremIpsumWord()),
			FullCode: cd,
		}
		if err := e.db.Create(&province).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding province done!")
}
