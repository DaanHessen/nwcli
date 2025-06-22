package renderer

import (
	"fmt"
	"strings"
	"time"

	"nwcli/pkg/news"

	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer handles markdown rendering with Glamour
type MarkdownRenderer struct {
	glamour *glamour.TermRenderer
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer() (*MarkdownRenderer, error) {
	// Create a glamour renderer with dark theme
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create glamour renderer: %w", err)
	}

	return &MarkdownRenderer{
		glamour: r,
	}, nil
}

// RenderArticles renders multiple articles as beautiful markdown
func (mr *MarkdownRenderer) RenderArticles(articles []news.Article, title string) (string, error) {
	if len(articles) == 0 {
		return mr.RenderMessage("📰 No articles found", "Try a different search query or check your sources.")
	}

	// Build markdown content
	var md strings.Builder

	// Header
	md.WriteString(fmt.Sprintf("# 📰 %s\n\n", title))
	md.WriteString(fmt.Sprintf("*Updated: %s*\n\n", time.Now().Format("Monday, January 2, 2006 at 15:04")))
	md.WriteString("---\n\n")

	// Articles
	for i, article := range articles {
		if i > 0 {
			md.WriteString("\n---\n\n")
		}

		// Article header with source and time
		sourceInfo := fmt.Sprintf("**%s** • %s",
			article.Source,
			formatTimeAgo(article.Published))

		md.WriteString(fmt.Sprintf("## %s\n\n", article.Title))
		md.WriteString(fmt.Sprintf("*%s*\n\n", sourceInfo))

		// Image if available
		if article.ImageURL != "" {
			md.WriteString(fmt.Sprintf("![Article Image](%s)\n\n", article.ImageURL))
		}

		// Description/Content
		if article.Description != "" {
			md.WriteString(fmt.Sprintf("%s\n\n", article.Description))
		} else if article.Content != "" {
			// Truncate content if it's very long
			content := article.Content
			if len(content) > 500 {
				content = content[:497] + "..."
			}
			md.WriteString(fmt.Sprintf("%s\n\n", content))
		}

		// Categories
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

		// Read more link
		md.WriteString(fmt.Sprintf("🔗 [Read full article](%s)\n\n", article.Link))
	}

	// Footer
	md.WriteString("---\n\n")
	md.WriteString(fmt.Sprintf("*Found %d articles • Generated with NWCLI*\n", len(articles)))

	// Render with glamour
	return mr.glamour.Render(md.String())
}

// RenderSingleArticle renders a single article in detail
func (mr *MarkdownRenderer) RenderSingleArticle(article news.Article) (string, error) {
	var md strings.Builder

	// Title
	md.WriteString(fmt.Sprintf("# %s\n\n", article.Title))

	// Metadata
	md.WriteString(fmt.Sprintf("**Source:** %s\n", article.Source))
	md.WriteString(fmt.Sprintf("**Published:** %s (%s)\n",
		article.Published.Format("Monday, January 2, 2006 at 15:04"),
		formatTimeAgo(article.Published)))
	md.WriteString(fmt.Sprintf("**URL:** %s\n\n", article.Link))

	// Categories
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

	// Image
	if article.ImageURL != "" {
		md.WriteString(fmt.Sprintf("![Article Image](%s)\n\n", article.ImageURL))
	}

	// Content
	if article.Content != "" {
		md.WriteString(fmt.Sprintf("%s\n\n", article.Content))
	} else if article.Description != "" {
		md.WriteString(fmt.Sprintf("%s\n\n", article.Description))
	}

	return mr.glamour.Render(md.String())
}

// RenderSources renders available news sources
func (mr *MarkdownRenderer) RenderSources(sources []news.Source) (string, error) {
	var md strings.Builder

	md.WriteString("# 📰 Available News Sources\n\n")
	md.WriteString("*Configure your preferred Dutch news sources*\n\n")
	md.WriteString("---\n\n")

	// Group by category
	categories := make(map[string][]news.Source)
	for _, source := range sources {
		categories[source.Category] = append(categories[source.Category], source)
	}

	for category, sources := range categories {
		md.WriteString(fmt.Sprintf("## %s\n\n", strings.Title(category)))

		for _, source := range sources {
			md.WriteString(fmt.Sprintf("### %s\n", source.Name))
			md.WriteString(fmt.Sprintf("*%s*\n\n", source.Description))
			md.WriteString(fmt.Sprintf("**URL:** %s\n", source.URL))
			md.WriteString(fmt.Sprintf("**Language:** %s\n\n", source.Language))
		}
	}

	md.WriteString("---\n\n")
	md.WriteString("*Use filters to focus on specific sources or categories*\n")

	return mr.glamour.Render(md.String())
}

// RenderMessage renders a simple message
func (mr *MarkdownRenderer) RenderMessage(title, message string) (string, error) {
	md := fmt.Sprintf("# %s\n\n%s\n", title, message)
	return mr.glamour.Render(md)
}

// RenderStats renders news statistics
func (mr *MarkdownRenderer) RenderStats(articles []news.Article) (string, error) {
	if len(articles) == 0 {
		return mr.RenderMessage("📊 No Statistics", "No articles available for analysis.")
	}

	var md strings.Builder

	md.WriteString("# 📊 News Statistics\n\n")
	md.WriteString("---\n\n")

	// Total articles
	md.WriteString(fmt.Sprintf("**Total Articles:** %d\n\n", len(articles)))

	// Sources breakdown
	sourceCount := make(map[string]int)
	categoryCount := make(map[string]int)

	for _, article := range articles {
		sourceCount[article.Source]++
		for _, cat := range article.Categories {
			categoryCount[cat]++
		}
	}

	// Top sources
	md.WriteString("## Sources\n\n")
	for source, count := range sourceCount {
		md.WriteString(fmt.Sprintf("- **%s**: %d articles\n", source, count))
	}
	md.WriteString("\n")

	// Categories
	if len(categoryCount) > 0 {
		md.WriteString("## Categories\n\n")
		for category, count := range categoryCount {
			md.WriteString(fmt.Sprintf("- **%s**: %d articles\n", category, count))
		}
		md.WriteString("\n")
	}

	// Time range
	if len(articles) > 0 {
		oldest := articles[0].Published
		newest := articles[0].Published

		for _, article := range articles {
			if article.Published.Before(oldest) {
				oldest = article.Published
			}
			if article.Published.After(newest) {
				newest = article.Published
			}
		}

		md.WriteString("## Time Range\n\n")
		md.WriteString(fmt.Sprintf("- **Newest**: %s\n", newest.Format("January 2, 2006 at 15:04")))
		md.WriteString(fmt.Sprintf("- **Oldest**: %s\n", oldest.Format("January 2, 2006 at 15:04")))
	}

	return mr.glamour.Render(md.String())
}

// formatTimeAgo formats time in a human-readable "time ago" format
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
