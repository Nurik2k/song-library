package main

import (
	"log"
	"net/http"
	"os"

	_ "song-library/docs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"song-library/handlers"
	"song-library/middleware"
	"song-library/models"
)

// @title        Song Library API
// @version      1.0
// @description  API сервиса онлайн библиотеки песен.
// @host         localhost:8080
// @BasePath     /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки .env файла")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	err = db.AutoMigrate(&models.Song{})
	if err != nil {
		log.Fatalf("Не удалось мигрировать базу данных: %v", err)
	}

	r := mux.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	handlers.RegisterHandlers(r, db)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
