package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"main.go/db"
	"main.go/models"

	_ "github.com/lib/pq"
)

// @title Music API
// @version 1.0
// @description API для работы с музыкальной библиотекой
// @host localhost:8080
// @BasePath /api

// ErrorResponse описание структуры для ошибок
// @Description Структура ответа с ошибкой
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Получение данных о песнях с фильтрацией и пагинацией

// GetSongs godoc
// @Summary Получить список всех песен
// @Description Получить все песни из библиотеки с возможностью фильтрации и пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group name"
// @Param song query string false "Song name"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.Song "List of songs"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Router /songs [get]
func GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	offset := (page - 1) * limit

	query := "SELECT * FROM songs WHERE group_name LIKE $1 AND song_name LIKE $2 LIMIT $3 OFFSET $4"
	rows, err := db.DB.Query(query, "%"+group+"%", "%"+song+"%", limit, offset)
	if err != nil {
		http.Error(w, "Ошибка выбора песен", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			http.Error(w, "Не удалось прочесть", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		songs = append(songs, song)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}

// Получение текста песни с пагинацией по куплетам

// GetSongLyrics godoc
// @Summary Получить текст песни по ID
// @Description Получить текст песни с пагинацией по куплетам
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param limit query int false "Limit" default(5)
// @Param offset query int false "Offset" default(0)
// @Success 200 {string} string "Lyrics of the song"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Song not found"
// @Router /songs/lyrics/{id} [get]
func GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromPath(r.URL.Path)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 1
	}

	offset := (page - 1) * limit

	query := "SELECT text FROM songs WHERE id = $1 LIMIT $2 OFFSET $3"
	rows, err := db.DB.Query(query, id, limit, offset)
	if err != nil {
		http.Error(w, "Ошибка выбора текста", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()

	var lyrics []string
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			http.Error(w, "Не удалось прочесть", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		lyrics = append(lyrics, text)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lyrics)
}

// Удаление песни

// DeleteSong godoc
// @Summary Удалить песню по ID
// @Description Удалить песню из библиотеки по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Song not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /songs/{id} [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromPath(r.URL.Path)

	query := "DELETE FROM songs WHERE id = $1"
	_, err := db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Не удалось удалить песню", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Изменение данных песни

// UpdateSong godoc
// @Summary Обновить данные о песне
// @Description Обновить данные песни по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song to update"
// @Success 200 {object} models.Song "Updated song"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Song not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /songs/{id} [put]
func UpdateSong(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromPath(r.URL.Path)
	var song models.Song
	_ = json.NewDecoder(r.Body).Decode(&song)

	query := "UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6"
	_, err := db.DB.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, id)
	if err != nil {
		http.Error(w, "Не удалось обновить песню", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Добавление новой песни

// AddSong godoc
// @Summary Добавить новую песню
// @Description Добавить новую песню в библиотеку
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song to add"
// @Success 201 {object} models.Song "New song added"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /songs [post]
func AddSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	_ = json.NewDecoder(r.Body).Decode(&song)

	query := "INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)"
	_, err := db.DB.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		http.Error(w, "Не удалось добавить песню", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Вспомогательная функция для извлечения ID из URL
func extractIDFromPath(path string) int {
	parts := strings.Split(path, "/")
	id, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		log.Println("Не удалось получить ID:", err)
	}
	return id
}
