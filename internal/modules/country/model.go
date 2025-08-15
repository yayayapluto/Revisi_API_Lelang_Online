package country

import "time"

type Country struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Kode      string    `gorm:"not null" json:"kode"`
	Nama      string    `gorm:"not null" json:"nama"`
	Nomor     string    `gorm:"not null" json:"nomor"`
}
