package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"nwcli/pkg/news"
	"nwcli/pkg/renderer"
	"nwcli/pkg/tui"
)

// Helper functions for rendering that can be used across commands

func renderMarkdown(articles []news.Article, title string) error {
	return renderMarkdownWithPager(articles, title, false)
}

func renderMarkdownWithPager(articles []news.Article, title string, noPager bool) error {
	// Check if we should use the TUI pager
	if tui.ShouldUsePager(noPager) && len(articles) > 0 {
		return tui.LaunchTUI(articles, title)
	}

	// Fallback to regular markdown rendering
	renderer, err := renderer.NewMarkdownRenderer()
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	output, err := renderer.RenderArticles(articles, title)
	if err != nil {
		return fmt.Errorf("failed to render articles: %w", err)
	}

	fmt.Print(output)
	return nil
}

func renderJSON(articles []news.Article) error {
	data, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func renderPlain(articles []news.Article) error {
	for i, article := range articles {
		if i > 0 {
			fmt.Println("\n" + strings.Repeat("-", 50))
		}
		
		fmt.Printf("Title: %s\n", article.Title)
		fmt.Printf("Source: %s\n", article.Source)
		fmt.Printf("Published: %s\n", article.Published.Format("2006-01-02 15:04"))
		if article.Description != "" {
			fmt.Printf("Description: %s\n", article.Description)
		}
		fmt.Printf("URL: %s\n", article.Link)
	}
	return nil
}
