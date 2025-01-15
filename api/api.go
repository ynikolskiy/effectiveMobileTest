package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func GetSongDetail(group string, song string) (*SongDetail, error) {
	url := fmt.Sprintf("https://api.com/info?group=%s&song=%s", group, song)

	client := http.Client{
		Timeout: 10 * time.Second, // Таймаут запроса
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Не получилось сделать запрос к API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Ошибка от API: %d", resp.StatusCode)
		return nil, fmt.Errorf("ошибка от API: %d", resp.StatusCode)
	}

	var songDetail SongDetail
	err = json.NewDecoder(resp.Body).Decode(&songDetail)
	if err != nil {
		log.Printf("Не получилось десериализовать: %v", err)
		return nil, err
	}

	return &songDetail, nil
}
