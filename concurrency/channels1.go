package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Anime struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Poster string `json:"poster"`
	Duration string `json:"duration"`
	Type string `json:"type"`
	Rating string `json:"rating"`
	Episodes struct {
		Sub uint16 `json:"sub"`
		Dub uint16 `json:"dub"`
	} `json:"episodes"`
}

type CategorizedAnime struct {
	Type string
	Animes []string
}

type AnimeCategoryResponse struct {
	Category string `json:"category"`
	Animes []Anime `json:"animes"`
	Genres []string `json:"genres"`
	CurrentPage uint16 `json:"currentPage"`
	TotalPages uint16 `json:"totalPages"`
	HasNextPage bool `json:"hasNextPage"`
}


func groupAnime(animeChannel <-chan Anime) {
	tvCategory := CategorizedAnime {
		Type: "TV",
	}
	specialCategory := CategorizedAnime {
		Type: "Special",
	}

	for anime := range animeChannel {
		if anime.Type == "TV" {
			tvCategory.Animes = append(tvCategory.Animes, anime.Name)
		}
		if anime.Type == "Special" {
			specialCategory.Animes = append(specialCategory.Animes, anime.Name)
		}
	}

	fmt.Println("TV Category: ", tvCategory)
	fmt.Println("\nSpecial Category: ", specialCategory)
}

func fetchAnimeSearchResult(searchQuery string, waitGroup *sync.WaitGroup) error {
	defer waitGroup.Done()

	var data AnimeCategoryResponse
	reqURL := fmt.Sprintf("https://api-aniwatch.onrender.com/anime/search?q=%s", searchQuery)

	res, err := http.Get(reqURL)
	if err != nil {
		fmt.Println("Error while fetching data")
		return nil
	}
	defer res.Body.Close()

	if parsedData := json.NewDecoder(res.Body).Decode(&data); parsedData != nil {
		fmt.Println("Error while fetching data")
		return nil
	}

	animeCh := make(chan Anime, len(data.Animes))
	for _, val := range data.Animes {
		animeCh <- val
	}
	close(animeCh)
	groupAnime(animeCh)

	return nil
}

func main() {
	var animeSearchTerm string = "tv"
	var wg sync.WaitGroup

	wg.Add(1)
	go fetchAnimeSearchResult(animeSearchTerm, &wg)
	wg.Wait()
}
