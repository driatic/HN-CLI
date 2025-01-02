package main

import (
	"HackerNewsCLI/api"
	"HackerNewsCLI/charm"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func fetchStoryTitles(stories []api.Story) []string {
	var titles []string
	for _, story := range stories {
		titles = append(titles, story.Title)
	}
	return titles
}

func fetchStoryLinks(stories []api.Story) []string {
	var links []string
	for _, story := range stories {
		links = append(links, story.URL)
	}
	return links
}

func main() {
	var stories = api.GetStories()
	//todo: add error handling if no stories

	p := tea.NewProgram(charm.InitialModel(fetchStoryTitles(stories), fetchStoryLinks(stories)))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
		os.Exit(1)
	}
}
