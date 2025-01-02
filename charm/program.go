package charm

import (
	"HackerNewsCLI/api"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"runtime"
)

type model struct {
	cursor int
	items  []string
	links  []string
}

func InitialModel(titles []string, links []string) model {
	return model{
		cursor: 0,
		items:  titles,
		links:  links,
	}
}

func StartProgram() {
	stories := api.GetStories()

	if len(stories) == 0 {
		fmt.Fprintf(os.Stderr, "No stories found.\n")
		os.Exit(1)
	}

	p := tea.NewProgram(InitialModel(parseStoryTitles(stories), parseStoryLinks(stories)))
	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "r":
			refreshStories(m)
			return m, nil
		case "enter", "return":
			openBrowser(m.links[m.cursor])
		}
	}
	return m, nil
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // For Linux and other systems
		cmd = "xdg-open"
		args = []string{url}
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		return fmt.Errorf("failed to execute command %s: %v", cmd, err)
	}
	return nil
}

func refreshStories(m model) model {
	stories := api.GetStories()
	m.items = parseStoryTitles(stories)
	m.links = parseStoryLinks(stories)
	m.cursor = 0
	return m
}

func parseStoryTitles(stories []api.Story) []string {
	var titles []string
	for _, story := range stories {
		titles = append(titles, story.Title)
	}
	return titles
}

func parseStoryLinks(stories []api.Story) []string {
	var links []string
	for _, story := range stories {
		links = append(links, story.URL)
	}
	return links
}

func (m model) View() string {
	s := "Latest HN articles:\n\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, item)
	}

	s += "\n[Use ↑/↓ to navigate, r to refresh, q to quit]\n"
	return s
}
