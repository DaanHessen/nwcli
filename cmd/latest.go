package cmd

import (
	"fmt"
	"strings"
	"time"

	"nwcli/pkg/news"

	"github.com/spf13/cobra"
)

var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "📈 Get the latest Dutch news",
	Long: `Fetch and display the latest news from Dutch sources including:
• NOS (Nederlandse Omroep Stichting)
• NU.nl
• De Telegraaf 
• RTL Nieuws
• AD.nl
• And more...

Articles are displayed in a beautiful newspaper-like format with images,
source information, and timestamps.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		limit, _ := cmd.Flags().GetInt("limit")
		source, _ := cmd.Flags().GetString("source")
		category, _ := cmd.Flags().GetString("category")
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")
		country, _ := cmd.Flags().GetString("country")
		fullContent, _ := cmd.Flags().GetBool("full")
		noPager, _ := cmd.Flags().GetBool("no-pager")

		if verbose {
			fmt.Printf("🔄 Fetching latest news from %s", country)
			if fullContent {
				fmt.Print(" (full articles)")
			}
			fmt.Println("...")
		}

		// Create news service with options
		newsService := news.NewNewsServiceWithOptions(country, fullContent)

		var articles []news.Article
		var err error

		if source != "" || category != "" {
			// Use filtering if source or category specified
			articles, err = newsService.FilterArticles(source, category, time.Time{}, limit)
		} else {
			// Get latest news
			articles, err = newsService.GetLatestNews(limit)
		}

		if err != nil {
			return fmt.Errorf("failed to fetch news: %w", err)
		}

		if verbose {
			fmt.Printf("✅ Found %d articles\n\n", len(articles))
		}

		// Render based on format
		switch format {
		case "json":
			return renderJSON(articles)
		case "plain":
			return renderPlain(articles)
		default: // markdown
			title := fmt.Sprintf("Latest News (%s)", strings.ToUpper(country))
			if fullContent {
				title += " - Full Articles"
			}
			return renderMarkdownWithPager(articles, title, noPager)
		}
	},
}

func init() {
	rootCmd.AddCommand(latestCmd)

	// Flags
	latestCmd.Flags().IntP("limit", "l", 20, "number of articles to show")
	latestCmd.Flags().StringP("source", "s", "", "filter by source (e.g., 'NOS', 'NU.nl')")
	latestCmd.Flags().StringP("category", "c", "", "filter by category (general, sports, tech)")
	latestCmd.Flags().StringP("country", "", "nl", "country code (nl, us, uk, de, fr)")
	latestCmd.Flags().BoolP("full", "", false, "fetch full article content instead of summaries")
	latestCmd.Flags().BoolP("no-pager", "", false, "disable interactive pager and output to stdout")
}
