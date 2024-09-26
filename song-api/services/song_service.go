package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"song-library/models"
)

type SongService struct {
	db *gorm.DB
}

func NewSongService(db *gorm.DB) *SongService {
	return &SongService{db: db}
}

func (s *SongService) GetSongs(page, limit int) ([]models.Song, error) {
	var songs []models.Song

	query := s.db.Model(&models.Song{})

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&songs).Error; err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *SongService) GetSongByID(id string) (*models.Song, error) {
	var song models.Song
	if err := s.db.First(&song, id).Error; err != nil {
		return nil, err
	}
	return &song, nil
}

func (s *SongService) DeleteSong(id string) error {
	if err := s.db.Delete(&models.Song{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *SongService) UpdateSong(id string, updatedSong *models.Song) error {
	var song models.Song
	if err := s.db.First(&song, id).Error; err != nil {
		return err
	}

	song.Group = updatedSong.Group
	song.SongName = updatedSong.SongName
	song.ReleaseDate = updatedSong.ReleaseDate
	song.Text = updatedSong.Text
	song.Link = updatedSong.Link

	if err := s.db.Save(&song).Error; err != nil {
		return err
	}
	return nil
}

func (s *SongService) AddSong(group, songName string) (*models.Song, error) {
	songDetail, err := s.fetchSongDetail(group, songName)
	if err != nil {
		return nil, err
	}

	song := models.Song{
		Group:       group,
		SongName:    songName,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	if err := s.db.Create(&song).Error; err != nil {
		return nil, err
	}

	return &song, nil
}

func (s *SongService) fetchSongDetail(group, song string) (*models.SongDetail, error) {
	apiURL := os.Getenv("EXTERNAL_API_URL")
	url := fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, group, song)

	resp, err := http.Get(url)
	if err != nil {
		log.Error("Ошибка запроса к внешнему API: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Некорректный ответ от внешнего API: ", resp.Status)
		return nil, fmt.Errorf("внешний API вернул статус %d", resp.StatusCode)
	}

	var songDetail models.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		log.Error("Ошибка декодирования ответа внешнего API: ", err)
		return nil, err
	}

	return &songDetail, nil
}
