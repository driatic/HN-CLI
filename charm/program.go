package charm

import (
	"HackerNewsCLI/api"
	"HackerNewsCLI/utils"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"os/exec"
	"runtime"
)

type model struct {
	cursor  int
	items   []string
	links   []string
	loading bool
}

func InitialModel(titles []string, links []string) model {
	return model{
		cursor:  0,
		items:   titles,
		links:   links,
		loading: false,
	}
}

func StartProgram() {
	stories := api.GetStories()

	if len(stories) == 0 {
		fmt.Fprintf(os.Stderr, "No stories found.\n")
		os.Exit(1)
	}

	p := tea.NewProgram(InitialModel(utils.ParseStoryTitles(stories), utils.ParseStoryLinks(stories)))
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
			m.loading = true
			return m, refreshStoriesCmd
		case "enter", "return":
			openBrowser(m.links[m.cursor])
		}

	case []api.Story:
		m.items = utils.ParseStoryTitles(msg)
		m.links = utils.ParseStoryLinks(msg)
		m.cursor = 0
		m.loading = false
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	if m.loading {
		return "Loading..."
	}

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

func refreshStoriesCmd() tea.Msg {
	stories := api.GetStories()
	return stories
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
