package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	p *tea.Program
)

type model struct {
	editing bool
	err     error
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editing {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			p.CancelInput()
			m.editing = true

			editor := os.Getenv("EDITOR")
			// Note: this should be in a command, but to do that we'll most
			// likely need to build some additional bubbletea to essentially
			// pause parts of the Bubble Tea runtime.
			c := exec.Command(editor)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			m.err = c.Run()

			m.editing = false
			p.RestoreInput()
			return m, tea.HideCursor
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.editing {
		return ""
	}
	if m.err != nil {
		return "Error: " + m.err.Error()
	}
	return "Press e to open Vim. Press q to quit."
}

func main() {
	m := model{}
	p = tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
