package models

import (
	"time"
)

type Song struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
	Group       string     `json:"group"`
	SongName    string     `json:"song"`
	ReleaseDate string     `json:"releaseDate"`
	Text        string     `json:"text"`
	Link        string     `json:"link"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
