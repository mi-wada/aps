package main

import (
	"fmt"
	"os"

	"github.com/mi-wada/aps/awsprofile"

	tea "github.com/charmbracelet/bubbletea"
)

// model represents the application state for the TUI
type model struct {
	profiles       []string
	cursor         int
	currentProfile string
}

func initialModel() model {
	profiles, err := awsprofile.List()
	if err != nil {
		fmt.Printf("Error getting profiles: %v\n", err)
		os.Exit(1)
	}

	return model{
		profiles:       profiles,
		cursor:         0,
		currentProfile: awsprofile.Current(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.profiles)-1 {
				m.cursor++
			}
		case "enter", " ":
			selectedProfile := m.profiles[m.cursor]
			fmt.Printf("export AWS_PROFILE=%s", selectedProfile)
			fmt.Fprintf(os.Stderr, "\nSwitched to profile: %s\n", selectedProfile)
			return m, tea.Quit
		}
	}

	return m, nil
}

// renderProfileLine formats a single profile line for display
func (m model) renderProfileLine(i int, profile string) string {
	cursor := " "
	if m.cursor == i {
		cursor = ">"
	}

	postfix := ""
	if profile == m.currentProfile {
		postfix = " [current]"
	}

	return fmt.Sprintf("%s %s%s\n", cursor, profile, postfix)
}

func (m model) View() string {
	s := "AWS Profile Switcher\n"

	for i, profile := range m.profiles {
		s += m.renderProfileLine(i, profile)
	}

	s += "\n(Use ↑/↓ or j/k to navigate, Enter to select, q to quit)"

	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
