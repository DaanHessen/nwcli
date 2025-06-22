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

	// Ensure we have at least one article to select
	selectedIndex := 0
	if len(articles) == 0 {
		selectedIndex = -1 // No articles to select
	}

	return &Model{
		articles:      articles,
		currentView:   IndexView,
		selectedIndex: selectedIndex,
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

	case tea.MouseMsg:
		switch msg.Action {
		case tea.MouseActionPress:
			if msg.Button == tea.MouseButtonWheelUp {
				if m.currentView == IndexView {
					if m.selectedIndex > 0 {
						m.selectedIndex--
					}
				} else {
					m.scrollUp()
				}
			} else if msg.Button == tea.MouseButtonWheelDown {
				if m.currentView == IndexView {
					if m.selectedIndex < len(m.articles)-1 {
						m.selectedIndex++
					}
				} else {
					m.scrollDown()
				}
			} else if msg.Button == tea.MouseButtonLeft {
				if m.currentView == IndexView && len(m.articles) > 0 {
					// Calculate which article was clicked based on mouse position
					// Each article takes about 6 lines with new styling
					headerHeight := 6 // Approximate header height
					if msg.Y >= headerHeight {
						clickedIndex := (msg.Y - headerHeight) / 6
						if clickedIndex >= 0 && clickedIndex < len(m.articles) {
							m.selectedIndex = clickedIndex
							m.currentView = ArticleView
							m.updateViewport()
						}
					}
				}
			}
		}

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
			if m.currentView == IndexView && len(m.articles) > 0 && m.selectedIndex >= 0 {
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

	// Header with enhanced styling
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5F87D7")).
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(m.windowWidth - 4).
		Align(lipgloss.Center)

	header := fmt.Sprintf("📰 %s", m.title)
	subHeader := fmt.Sprintf("%s • %d articles available",
		time.Now().Format("Monday, January 2, 2006 at 15:04"),
		len(m.articles))

	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n")
	
	// Sub-header
	subHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8888AA")).
		Align(lipgloss.Center).
		Width(m.windowWidth).
		MarginBottom(1)
	
	b.WriteString(subHeaderStyle.Render(subHeader))
	b.WriteString("\n\n")

	// Check if we have articles to display
	if len(m.articles) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			Width(m.windowWidth).
			Padding(2)
		
		b.WriteString(emptyStyle.Render("📭 No articles found\n\nPress 'q' to quit"))
		return b.String()
	}

	// Article list with enhanced styling
	for i, article := range m.articles {
		isSelected := i == m.selectedIndex
		
		// Base article container
		containerStyle := lipgloss.NewStyle().
			Width(m.windowWidth - 4).
			MarginBottom(1).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder())

		if isSelected {
			containerStyle = containerStyle.
				BorderForeground(lipgloss.Color("#FF6B9D")).
				Background(lipgloss.Color("#2D1B4E")).
				Bold(true)
		} else {
			containerStyle = containerStyle.
				BorderForeground(lipgloss.Color("#3C3C3C")).
				Background(lipgloss.Color("#1A1A1A"))
		}

		// Title styling
		titleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Bold(true)
		
		if isSelected {
			titleStyle = titleStyle.Foreground(lipgloss.Color("#FFB6D9"))
		}

		// Source and time styling
		metaStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			MarginTop(1)
		
		if isSelected {
			metaStyle = metaStyle.Foreground(lipgloss.Color("#BBBBBB"))
		}

		// Description styling  
		descStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			MarginTop(1)
		
		if isSelected {
			descStyle = descStyle.Foreground(lipgloss.Color("#EEEEEE"))
		}

		// Build article content
		var articleContent strings.Builder
		
		// Selection indicator
		indicator := "  "
		if isSelected {
			indicator = "▶ "
		}
		
		articleContent.WriteString(indicator + titleStyle.Render(article.Title))
		
		// Source and time
		sourceTime := fmt.Sprintf("📡 %s • 🕒 %s",
			article.Source,
			formatTimeAgo(article.Published))
		articleContent.WriteString("\n" + metaStyle.Render(sourceTime))

		// Description
		if article.Description != "" {
			desc := article.Description
			if len(desc) > 120 {
				desc = desc[:117] + "..."
			}
			articleContent.WriteString("\n" + descStyle.Render("💬 " + desc))
		}

		b.WriteString(containerStyle.Render(articleContent.String()))
		b.WriteString("\n")
	}

	// Enhanced footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Background(lipgloss.Color("#1A1A1A")).
		Padding(1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Width(m.windowWidth - 4).
		Align(lipgloss.Center).
		MarginTop(1)

	footer := "🖱️  Mouse & scroll wheel supported • ⏎ Enter to read • ↑/↓ or j/k to navigate • g/G first/last • h help • q quit"
	if m.showHelp {
		footer = "📖 Navigation: ↑/↓ or j/k or mouse wheel • ⏎ Enter: read article • 🖱️ Click: select & read • g/G: first/last • ESC: back • q: quit • h: toggle help"
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
		// Use terminal image support if available
		imagePlaceholder := GetImagePlaceholder(article.ImageURL)
		md.WriteString(imagePlaceholder)
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
	end := start + m.viewport.height - 3 // Leave space for header and footer

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

	// Enhanced header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5F87D7")).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		Width(m.windowWidth - 4).
		Align(lipgloss.Center)

	article := m.articles[m.selectedIndex]
	header := fmt.Sprintf("📖 Article %d of %d • %s",
		m.selectedIndex+1, len(m.articles), article.Source)

	// Calculate scroll percentage
	scrollPercent := 0
	if len(lines) > m.viewport.height {
		scrollPercent = int(float64(start) / float64(len(lines)-m.viewport.height) * 100)
	}

	// Enhanced footer with scroll indicator
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Background(lipgloss.Color("#1A1A1A")).
		Padding(0, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Width(m.windowWidth - 4).
		Align(lipgloss.Center)

	footer := fmt.Sprintf("📊 %d%% • 🖱️ Mouse wheel supported • ⬅ ESC back • ↑/↓ scroll • h help", scrollPercent)
	if m.showHelp {
		footer = "🖱️ Mouse wheel or ↑/↓ j/k: scroll • PgUp/PgDn: page • g/G: top/bottom • ⬅ ESC: back to index • q: quit • h: toggle help"
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render(header))
	b.WriteString("\n\n")
	b.WriteString(strings.Join(visible, "\n"))
	b.WriteString("\n\n")
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
		// Always start at the top when entering article view
		m.viewport.scrollOffset = 0
		m.viewport.height = m.windowHeight - 4 // Account for header/footer
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
