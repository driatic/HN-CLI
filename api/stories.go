package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	baseURL     = "https://hacker-news.firebaseio.com/v0"
	newStories  = baseURL + "/newstories.json"
	itemDetails = baseURL + "/item/%d.json"
)

// todo: add paging param
func GetStories() []Story {
	storyIDs := GetStoryIDs(newStories)

	var stories []Story
	for _, id := range storyIDs[:10] {
		story := GetStoryDetails(id)
		if story.URL == "" {
			story.URL = "No URL available"
		}
		stories = append(stories, story)
	}
	return stories
}

func GetStoryIDs(url string) []int {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting story IDs: %v", err)
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		log.Fatalf("Error decoding story IDs: %v", err)
	}
	return ids
}

func GetStoryDetails(id int) Story {
	url := fmt.Sprintf(itemDetails, id)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting story details for ID %d: %v", id, err)
	}
	defer resp.Body.Close()

	var story Story
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		log.Fatalf("Error decoding story details: %v", err)
	}
	return story
}

type Story struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
