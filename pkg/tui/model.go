package tui

import (
	"fmt"
	"strings"
	"time"

	"nwcli/pkg/news"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI application state
type Model struct {
	articles      []news.Article
	currentView   ViewType
	selectedIndex int
	title         string
	renderer      *glamour.TermRenderer
	viewport      Viewport
	showHelp      bool
	windowWidth   int
	windowHeight  int
}

// ViewType represents the current view
type ViewType int

const (
	IndexView ViewType = iota
	ArticleView
)

// Viewport represents a scrollable content area
type Viewport struct {
	content      string
	scrollOffset int
	height       int
}

// NewModel creates a new TUI model
func NewModel(articles []news.Article, title string) (*Model, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	return &Model{
		articles:      articles,
		currentView:   IndexView,
		selectedIndex: 0,
		title:         title,
		renderer:      renderer,
		showHelp:      false,
		windowWidth:   80,
		windowHeight:  24,
	}, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.viewport.height = msg.Height - 4 // Leave space for header/footer

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "h", "?":
			m.showHelp = !m.showHelp

		case "esc":
			if m.currentView == ArticleView {
				m.currentView = IndexView
				m.updateViewport()
			}

		case "enter":
			if m.currentView == IndexView && len(m.articles) > 0 {
				m.currentView = ArticleView
				m.updateViewport()
			}

		case "up", "k":
			if m.currentView == IndexView {
				if m.selectedIndex > 0 {
					m.selectedIndex--
				}
			} else {
				m.scrollUp()
			}

		case "down", "j":
			if m.currentView == IndexView {
				if m.selectedIndex < len(m.articles)-1 {
					m.selectedIndex++
				}
			} else {
				m.scrollDown()
			}

		case "pgup":
			m.pageUp()

		case "pgdown":
			m.pageDown()

		case "home", "g":
			if m.currentView == IndexView {
				m.selectedIndex = 0
			} else {
				m.viewport.scrollOffset = 0
			}

		case "end", "G":
			if m.currentView == IndexView {
				m.selectedIndex = len(m.articles) - 1
			} else {
				// Scroll to bottom
				lines := strings.Count(m.viewport.content, "\n")
				maxScroll := lines - m.viewport.height + 1
				if maxScroll > 0 {
					m.viewport.scrollOffset = maxScroll
				}
			}
		}
	}

	return m, nil
}

// View renders the current view
func (m Model) View() string {
	if len(m.articles) == 0 {
		return m.renderEmpty()
	}

	switch m.currentView {
	case IndexView:
		return m.renderIndex()
	case ArticleView:
		return m.renderArticle()
	default:
		return m.renderIndex()
	}
}

// renderIndex renders the article index
func (m Model) renderIndex() string {
	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Width(m.windowWidth)

	header := fmt.Sprintf("📰 %s\n%s\n%d articles • Use ↑/↓ to navigate, Enter to read, q to quit",
		m.title,
		time.Now().Format("Monday, January 2, 2006 at 15:04"),
		len(m.articles))

	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n\n")

	// Article list
	for i, article := range m.articles {
		style := lipgloss.NewStyle().
			PaddingLeft(2)

		if i == m.selectedIndex {
			style = style.
				Bold(true).
				Foreground(lipgloss.Color("205")).
				Background(lipgloss.Color("235"))
		}

		// Article item
		sourceTime := fmt.Sprintf("[%s] %s",
			article.Source,
			formatTimeAgo(article.Published))

		item := fmt.Sprintf("• %s\n  %s",
			article.Title,
			sourceTime)

		if article.Description != "" && len(article.Description) > 0 {
			desc := article.Description
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			item += fmt.Sprintf("\n  %s", desc)
		}

		b.WriteString(style.Render(item))
		b.WriteString("\n\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		Width(m.windowWidth)

	footer := "Press 'h' for help • 'q' to quit • Enter to read article"
	if m.showHelp {
		footer = "Navigation: ↑/↓ or j/k • Enter: read article • g/G: first/last • q: quit • h: toggle help"
	}

	b.WriteString(footerStyle.Render(footer))

	return b.String()
}

// renderArticle renders the selected article
func (m Model) renderArticle() string {
	if m.selectedIndex >= len(m.articles) {
		return "Error: Invalid article selection"
	}

	article := m.articles[m.selectedIndex]

	// Generate markdown content for the article
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# %s\n\n", article.Title))
	md.WriteString(fmt.Sprintf("**Source:** %s\n", article.Source))
	md.WriteString(fmt.Sprintf("**Published:** %s (%s)\n",
		article.Published.Format("Monday, January 2, 2006 at 15:04"),
		formatTimeAgo(article.Published)))
	md.WriteString(fmt.Sprintf("**URL:** %s\n\n", article.Link))

	if len(article.Categories) > 0 {
		md.WriteString("**Categories:** ")
		for i, cat := range article.Categories {
			if i > 0 {
				md.WriteString(", ")
			}
			md.WriteString(fmt.Sprintf("`%s`", cat))
		}
		md.WriteString("\n\n")
	}

	md.WriteString("---\n\n")

	if article.ImageURL != "" {
		md.WriteString(fmt.Sprintf("![Article Image](%s)\n\n", article.ImageURL))
	}

	if article.Content != "" {
		md.WriteString(fmt.Sprintf("%s\n\n", article.Content))
	} else if article.Description != "" {
		md.WriteString(fmt.Sprintf("%s\n\n", article.Description))
	}

	// Render with glamour
	rendered, err := m.renderer.Render(md.String())
	if err != nil {
		rendered = md.String() // Fallback to plain text
	}

	m.viewport.content = rendered

	// Create scrollable view
	return m.renderScrollableContent()
}

// renderScrollableContent renders content with scrolling
func (m Model) renderScrollableContent() string {
	lines := strings.Split(m.viewport.content, "\n")

	start := m.viewport.scrollOffset
	end := start + m.viewport.height - 2 // Leave space for header

	if start < 0 {
		start = 0
	}
	if end > len(lines) {
		end = len(lines)
	}
	if start >= end {
		start = end - 1
		if start < 0 {
			start = 0
		}
	}

	var visible []string
	if start < len(lines) {
		visible = lines[start:end]
	}

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Width(m.windowWidth)

	header := fmt.Sprintf("Article %d of %d • ESC: back to index • ↑/↓: scroll • q: quit",
		m.selectedIndex+1, len(m.articles))

	// Footer with scroll indicator
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true).
		Width(m.windowWidth)

	scrollPercent := 0
	if len(lines) > m.viewport.height {
		scrollPercent = int(float64(start) / float64(len(lines)-m.viewport.height) * 100)
	}

	footer := fmt.Sprintf("Scroll: %d%% • h: help", scrollPercent)
	if m.showHelp {
		footer = "Navigation: ↑/↓ or j/k • PgUp/PgDn: page • g/G: top/bottom • ESC: back • q: quit"
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Join(visible, "\n"))
	b.WriteString("\n")
	b.WriteString(footerStyle.Render(footer))

	return b.String()
}

// renderEmpty renders empty state
func (m Model) renderEmpty() string {
	style := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("241")).
		Height(m.windowHeight).
		Width(m.windowWidth)

	return style.Render("📭 No articles found\n\nPress 'q' to quit")
}

// updateViewport updates the viewport content
func (m *Model) updateViewport() {
	if m.currentView == ArticleView {
		// Content will be generated in renderArticle
		m.viewport.scrollOffset = 0
	}
}

// Scrolling methods
func (m *Model) scrollUp() {
	if m.viewport.scrollOffset > 0 {
		m.viewport.scrollOffset--
	}
}

func (m *Model) scrollDown() {
	lines := strings.Count(m.viewport.content, "\n")
	maxScroll := lines - m.viewport.height + 1
	if maxScroll > 0 && m.viewport.scrollOffset < maxScroll {
		m.viewport.scrollOffset++
	}
}

func (m *Model) pageUp() {
	m.viewport.scrollOffset -= m.viewport.height - 2
	if m.viewport.scrollOffset < 0 {
		m.viewport.scrollOffset = 0
	}
}

func (m *Model) pageDown() {
	lines := strings.Count(m.viewport.content, "\n")
	maxScroll := lines - m.viewport.height + 1
	m.viewport.scrollOffset += m.viewport.height - 2
	if maxScroll > 0 && m.viewport.scrollOffset > maxScroll {
		m.viewport.scrollOffset = maxScroll
	}
}

// Helper function for time formatting
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("January 2, 2006")
	}
}
