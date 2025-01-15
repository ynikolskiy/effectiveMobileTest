// main.go
package main

import (
	"github.com/gorilla/mux"
	"main.go/db"
	_ "main.go/docs"
	"main.go/handlers"

	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Инициализация базы данных
	db.InitDB()

	r := mux.NewRouter()

	r.Handle("/swagger/{any:.*}", httpSwagger.WrapHandler)

	// Роуты

	r.HandleFunc("/songs", handlers.GetSongs).Methods("GET")                  // Для получения всех песен
	r.HandleFunc("/songs", handlers.AddSong).Methods("POST")                  // Для добавления новой песни
	r.HandleFunc("/songs/{id}", handlers.UpdateSong).Methods("PUT")           // Для обновления песни по ID
	r.HandleFunc("/songs/{id}", handlers.DeleteSong).Methods("DELETE")        // Для удаления песни по ID
	r.HandleFunc("/songs/lyrics/{id}", handlers.GetSongLyrics).Methods("GET") // Для получения текста песни по ID

	// Поднимаем сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер поднимется на порту: %s\n", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}

}
