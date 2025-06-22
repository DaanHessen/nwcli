package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"nwcli/pkg/news"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "🔍 Search Dutch news articles",
	Long: `Search through Dutch news articles using keywords.
	
The search looks through article titles, descriptions, and content
to find relevant matches. Results are ranked by relevance and
displayed in a beautiful format.

Examples:
  nwcli search "climate change"
  nwcli search "voetbal" --limit 10
  nwcli search "politiek" --source "NOS"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		limit, _ := cmd.Flags().GetInt("limit")
		source, _ := cmd.Flags().GetString("source")
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")
		country, _ := cmd.Flags().GetString("country")
		fullContent, _ := cmd.Flags().GetBool("full")
		noPager, _ := cmd.Flags().GetBool("no-pager")

		// Combine all args into search query
		query := strings.Join(args, " ")

		if verbose {
			fmt.Printf("🔍 Searching for: '%s' in %s news", query, country)
			if fullContent {
				fmt.Print(" (full articles)")
			}
			fmt.Println()
		}

		// Create news service with options
		newsService := news.NewNewsServiceWithOptions(country, fullContent)
		
		// Search articles
		articles, err := newsService.SearchArticles(query, limit)
		if err != nil {
			return fmt.Errorf("failed to search articles: %w", err)
		}

		// Filter by source if specified
		if source != "" {
			var filtered []news.Article
			for _, article := range articles {
				if strings.EqualFold(article.Source, source) {
					filtered = append(filtered, article)
				}
			}
			articles = filtered
		}

		if verbose {
			fmt.Printf("✅ Found %d matching articles\n\n", len(articles))
		}

		// Render based on format
		switch format {
		case "json":
			return renderJSON(articles)
		case "plain":
			return renderPlain(articles)
		default: // markdown
			title := fmt.Sprintf("Search Results for '%s' (%s)", query, strings.ToUpper(country))
			if fullContent {
				title += " - Full Articles"
			}
			return renderMarkdown(articles, title)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Flags
	searchCmd.Flags().IntP("limit", "l", 20, "number of results to show")
	searchCmd.Flags().StringP("source", "s", "", "filter results by source")
	searchCmd.Flags().StringP("country", "", "nl", "country code (nl, us, uk, de, fr)")
	searchCmd.Flags().BoolP("full", "", false, "search in full article content instead of summaries")
}
