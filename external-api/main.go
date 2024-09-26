package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")

	if group == "" || song == "" {
		http.Error(w, "Параметры 'group' и 'song' обязательны", http.StatusBadRequest)
		return
	}

	response := SongDetail{
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/info", infoHandler) // Регистрация обработчика
	fmt.Println("Сервер запущен на порту 8081")
	http.ListenAndServe(":8081", nil) // Запуск сервера на порту 8081
}
