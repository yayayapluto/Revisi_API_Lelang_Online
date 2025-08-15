package database

import (
	"fmt"
	city2 "github.com/API_Lelang_Online_Go/internal/modules/city"
	country2 "github.com/API_Lelang_Online_Go/internal/modules/country"
	province2 "github.com/API_Lelang_Online_Go/internal/modules/province"
	subdistrict2 "github.com/API_Lelang_Online_Go/internal/modules/subdistrict"
	village2 "github.com/API_Lelang_Online_Go/internal/modules/village"
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
	CitySeeder(total int)
	SubdistrictSeeder(total int)
	VillageSeeder(total int)
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

func (e *entitySeeder) CitySeeder(total int) {
	if err := e.db.Exec("TRUNCATE TABLE cities RESTART IDENTITY").Error; err != nil {
		log.Fatal("Cannot truncate cities: ", err)
	}

	for i := 0; i < total; i++ {
		cd := gofakeit.Numerify("##")

		city := city2.City{
			Code:     cd,
			Nama:     fmt.Sprintf("City %v", gofakeit.LoremIpsumWord()),
			FullCode: cd,
		}
		if err := e.db.Create(&city).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding city done!")
}

func (e *entitySeeder) SubdistrictSeeder(total int) {
	if err := e.db.Exec("TRUNCATE TABLE cities RESTART IDENTITY").Error; err != nil {
		log.Fatal("Cannot truncate cities: ", err)
	}

	for i := 0; i < total; i++ {
		cd := gofakeit.Numerify("##")

		subdistrict := subdistrict2.Subdistrict{
			Code:     cd,
			Nama:     fmt.Sprintf("Subdistrict %v", gofakeit.LoremIpsumWord()),
			FullCode: cd,
		}
		if err := e.db.Create(&subdistrict).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding subdistrict done!")
}

func (e *entitySeeder) VillageSeeder(total int) {
	if err := e.db.Exec("TRUNCATE TABLE villages RESTART IDENTITY").Error; err != nil {
		log.Fatal("Cannot truncate villages: ", err)
	}

	for i := 0; i < total; i++ {
		cd := gofakeit.Numerify("##")

		village := village2.Village{
			Code:     cd,
			Nama:     fmt.Sprintf("Village %v", gofakeit.LoremIpsumWord()),
			FullCode: cd,
		}
		if err := e.db.Create(&village).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Seeding village done!")
}
