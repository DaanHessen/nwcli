package tui

import (
	"fmt"
	"os"

	"nwcli/pkg/news"

	tea "github.com/charmbracelet/bubbletea"
)

// LaunchTUI starts the interactive TUI for browsing articles
func LaunchTUI(articles []news.Article, title string) error {
	model, err := NewModel(articles, title)
	if err != nil {
		return fmt.Errorf("failed to create TUI model: %w", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err = p.Run()
	return err
}

// ShouldUsePager determines if we should use the pager based on flags and environment
func ShouldUsePager(noPager bool) bool {
	// Don't use pager if explicitly disabled
	if noPager {
		return false
	}

	// Don't use pager if output is redirected
	if !isTerminal() {
		return false
	}

	// Don't use pager if NWCLI_NO_PAGER is set
	if os.Getenv("NWCLI_NO_PAGER") != "" {
		return false
	}

	return true
}

// isTerminal checks if stdout is a terminal
func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}
