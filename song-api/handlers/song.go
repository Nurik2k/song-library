package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"song-library/models"
	"song-library/services"
)

func RegisterHandlers(r *mux.Router, db *gorm.DB) {
	songService := services.NewSongService(db)

	r.HandleFunc("/library", getLibrary(songService)).Methods("GET")
	r.HandleFunc("/song/{id}", getSong(songService)).Methods("GET")
	r.HandleFunc("/song/{id}", deleteSong(songService)).Methods("DELETE")
	r.HandleFunc("/song/{id}", updateSong(songService)).Methods("PUT")
	r.HandleFunc("/song", addSong(songService)).Methods("POST")
}

// getLibrary godoc
// @Summary      Получение данных библиотеки
// @Description  Получение данных библиотеки с фильтрацией по всем полям и пагинацией
// @Tags         library
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Номер страницы"
// @Param        limit     query     int     false  "Размер страницы"
// @Success      200       {array}   models.Song
// @Failure      400       {object}  map[string]string
// @Router       /library [get]
func getLibrary(songService *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Получение данных библиотеки")

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		if page == 0 {
			page = 1
		}
		if limit == 0 {
			limit = 10
		}

		songs, err := songService.GetSongs(page, limit)
		if err != nil {
			log.Error("Ошибка получения данных: ", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(songs)
	}
}

// getSong godoc
// @Summary      Получение текста песни
// @Description  Получение текста песни с пагинацией по куплетам
// @Tags         song
// @Accept       json
// @Produce      json
// @Param        id     path     int     true  "ID песни"
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object} map[string]string
// @Router       /song/{id} [get]
func getSong(songService *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Получение текста песни")
		vars := mux.Vars(r)
		id := vars["id"]

		song, err := songService.GetSongByID(id)
		if err != nil {
			log.Error("Песня не найдена: ", err)
			http.Error(w, "Песня не найдена", http.StatusNotFound)
			return
		}

		verses := strings.Split(song.Text, "\n\n")
		totalVerses := len(verses)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"verses": verses,
			"total":  totalVerses,
		})
	}
}

// deleteSong godoc
// @Summary      Удаление песни
// @Description  Удаление песни по ID
// @Tags         song
// @Accept       json
// @Produce      json
// @Param        id   path     int  true  "ID песни"
// @Success      204  "Песня удалена"
// @Failure      400  {object} map[string]string
// @Router       /song/{id} [delete]
func deleteSong(songService *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Удаление песни")
		vars := mux.Vars(r)
		id := vars["id"]

		if err := songService.DeleteSong(id); err != nil {
			log.Error("Ошибка удаления песни: ", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// updateSong godoc
// @Summary      Изменение данных песни
// @Description  Изменение данных песни по ID
// @Tags         song
// @Accept       json
// @Produce      json
// @Param        id    path     int            true  "ID песни"
// @Param        song  body     models.Song    true  "Данные песни"
// @Success      204   "Песня обновлена"
// @Failure      400   {object} map[string]string
// @Router       /song/{id} [put]
func updateSong(songService *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Обновление песни")
		vars := mux.Vars(r)
		id := vars["id"]

		var updatedSong models.Song
		if err := json.NewDecoder(r.Body).Decode(&updatedSong); err != nil {
			log.Error("Ошибка декодирования данных: ", err)
			http.Error(w, "Неверные данные", http.StatusBadRequest)
			return
		}

		if err := songService.UpdateSong(id, &updatedSong); err != nil {
			log.Error("Ошибка обновления песни: ", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// addSong godoc
// @Summary      Добавление новой песни
// @Description  Добавление новой песни с обогащением данных через внешний API
// @Tags         song
// @Accept       json
// @Produce      json
// @Param        group  query     string  true  "Название группы"
// @Param        song   query     string  true  "Название песни"
// @Success      201    {object} models.Song
// @Failure      400    {object} map[string]string
// @Router       /song [post]
func addSong(songService *services.SongService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Добавление новой песни")

		group := r.URL.Query().Get("group")
		songName := r.URL.Query().Get("song")

		if group == "" || songName == "" {
			http.Error(w, "group и song обязательны", http.StatusBadRequest)
			return
		}

		song, err := songService.AddSong(group, songName)
		if err != nil {
			log.Error("Ошибка добавления песни: ", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(song)
	}
}
