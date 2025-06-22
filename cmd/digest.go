package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"nwcli/pkg/news"
)

var digestCmd = &cobra.Command{
	Use:   "digest",
	Short: "📰 Get your daily Dutch news digest",
	Long: `Generate a personalized daily news digest from Dutch sources.

This creates a newspaper-like summary with:
• Top stories from major Dutch outlets
• Categorized sections (General, Sports, Technology)
• Beautiful formatting with images
• Time stamps and source attribution

Perfect for your morning news routine!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		limit, _ := cmd.Flags().GetInt("limit")
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")
		categories, _ := cmd.Flags().GetStringSlice("categories")
		country, _ := cmd.Flags().GetString("country")
		fullContent, _ := cmd.Flags().GetBool("full")

		if verbose {
			fmt.Printf("📰 Preparing your daily %s news digest", country)
			if fullContent {
				fmt.Print(" (full articles)")
			}
			fmt.Println("...")
		}

		// Create news service with options
		newsService := news.NewNewsServiceWithOptions(country, fullContent)
		
		// Get today's articles
		today := time.Now().Truncate(24 * time.Hour)
		
		var allArticles []news.Article
		
		if len(categories) > 0 {
			// Fetch articles for specific categories
			for _, category := range categories {
				articles, err := newsService.FilterArticles("", category, today, 0)
				if err != nil {
					if verbose {
						fmt.Printf("Warning: Failed to fetch %s articles: %v\n", category, err)
					}
					continue
				}
				allArticles = append(allArticles, articles...)
			}
		} else {
			// Get all articles from today
			articles, err := newsService.FilterArticles("", "", today, 0)
			if err != nil {
				// If no articles from today, get latest
				if verbose {
					fmt.Println("No articles from today, fetching latest...")
				}
				articles, err = newsService.GetLatestNews(limit * 2) // Get more to ensure variety
				if err != nil {
					return fmt.Errorf("failed to fetch news: %w", err)
				}
			}
			allArticles = articles
		}

		// Group articles by category for better digest structure
		digestArticles := organizeDigestArticles(allArticles, limit)

		if verbose {
			fmt.Printf("✅ Prepared digest with %d articles\n\n", len(digestArticles))
		}

		// Render based on format
		switch format {
		case "json":
			return renderJSON(digestArticles)
		case "plain":
			return renderPlain(digestArticles)
		default: // markdown
			title := fmt.Sprintf("📰 Daily News Digest (%s) - %s", 
				strings.ToUpper(country),
				time.Now().Format("Monday, January 2, 2006"))
			if fullContent {
				title += " - Full Articles"
			}
			return renderMarkdown(digestArticles, title)
		}
	},
}

func init() {
	rootCmd.AddCommand(digestCmd)

	// Flags
	digestCmd.Flags().IntP("limit", "l", 15, "number of articles in digest")
	digestCmd.Flags().StringSliceP("categories", "c", []string{}, "categories to include (general, sports, tech)")
	digestCmd.Flags().StringP("country", "", "nl", "country code (nl, us, uk, de, fr)")
	digestCmd.Flags().BoolP("full", "", false, "include full article content instead of summaries")
}

// organizeDigestArticles organizes articles for a balanced digest
func organizeDigestArticles(articles []news.Article, limit int) []news.Article {
	if len(articles) <= limit {
		return articles
	}

	// Group by category and source for variety
	categoryGroups := make(map[string][]news.Article)
	sourceCount := make(map[string]int)
	
	for _, article := range articles {
		// Categorize articles
		category := "general"
		if len(article.Categories) > 0 {
			category = article.Categories[0]
		}
		
		categoryGroups[category] = append(categoryGroups[category], article)
		sourceCount[article.Source]++
	}

	// Select articles to ensure variety
	var selected []news.Article
	maxPerCategory := limit / len(categoryGroups)
	if maxPerCategory < 1 {
		maxPerCategory = 1
	}

	for _, articles := range categoryGroups {
		count := 0
		sourceUsed := make(map[string]int)
		
		for _, article := range articles {
			if count >= maxPerCategory {
				break
			}
			
			// Prefer variety in sources
			if sourceUsed[article.Source] < 2 {
				selected = append(selected, article)
				sourceUsed[article.Source]++
				count++
			}
		}
		
		if len(selected) >= limit {
			break
		}
	}

	// Fill remaining slots if needed
	if len(selected) < limit {
		used := make(map[string]bool)
		for _, article := range selected {
			used[article.Link] = true
		}
		
		for _, article := range articles {
			if len(selected) >= limit {
				break
			}
			if !used[article.Link] {
				selected = append(selected, article)
				used[article.Link] = true
			}
		}
	}

	return selected
}
