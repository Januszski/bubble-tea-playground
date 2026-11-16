package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	apptui "github.com/Januszski/bubble-tea-playground/internal/app-tui"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	app := apptui.NewApp()

	p := tea.NewProgram(app)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
