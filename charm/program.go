package charm

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
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

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

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

func (m model) View() string {
	s := "Latest HN articles:\n\n"

	for i, item := range m.items {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, item)
	}

	s += "\n[Use ↑/↓ or j/k to navigate, q to quit]\n"
	return s
}
