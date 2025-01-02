package utils

import "HackerNewsCLI/api"

func ParseStoryTitles(stories []api.Story) []string {
	var titles []string
	for _, story := range stories {
		titles = append(titles, story.Title)
	}
	return titles
}

func ParseStoryLinks(stories []api.Story) []string {
	var links []string
	for _, story := range stories {
		links = append(links, story.URL)
	}
	return links
}
