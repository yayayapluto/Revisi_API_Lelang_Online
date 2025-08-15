package subdistrict

import "time"

type Subdistrict struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Nama      string    `gorm:"not null" json:"nama"`
	FullCode  string    `gorm:"not null" json:"full_code"`
	Code      string    `gorm:"not null" json:"code"`
}
