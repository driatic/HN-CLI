package charm

import (
	"HackerNewsCLI/api"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"runtime"
)

func refreshStories() tea.Msg {
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
