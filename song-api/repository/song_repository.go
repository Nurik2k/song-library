package repository

import (
	"gorm.io/gorm"
	"song-library/models"
)

type SongRepository interface {
	GetAllSongs(filter map[string]string, page int, limit int) ([]models.Song, error)
	GetSongByID(id uint) (models.Song, error)
	AddNewSong(song models.Song) (models.Song, error)
	UpdateSong(id uint, updatedSong models.Song) error
	DeleteSong(id uint) error
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) GetAllSongs(filter map[string]string, page int, limit int) ([]models.Song, error) {
	var songs []models.Song
	query := r.db.Model(&models.Song{})

	if group, ok := filter["group"]; ok {
		query = query.Where("group = ?", group)
	}
	if song, ok := filter["song"]; ok {
		query = query.Where("song = ?", song)
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&songs).Error
	return songs, err
}

func (r *songRepository) GetSongByID(id uint) (models.Song, error) {
	var song models.Song
	err := r.db.First(&song, id).Error
	return song, err
}

func (r *songRepository) AddNewSong(song models.Song) (models.Song, error) {
	err := r.db.Create(&song).Error
	return song, err
}

func (r *songRepository) UpdateSong(id uint, updatedSong models.Song) error {
	return r.db.Model(&models.Song{}).Where("id = ?", id).Updates(updatedSong).Error
}

func (r *songRepository) DeleteSong(id uint) error {
	return r.db.Delete(&models.Song{}, id).Error
}
