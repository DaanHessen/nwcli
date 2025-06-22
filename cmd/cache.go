package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"nwcli/pkg/news"
	"nwcli/pkg/renderer"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "🗂️  Manage article cache",
	Long: `Manage the local article cache used by NWCLI.

The cache stores fetched articles locally for faster access and offline reading.
Articles are automatically cached when fetching news, but you can manually
clear the cache if needed.`,
}

var cacheStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "📊 Show cache statistics",
	Long:  `Display statistics about the local article cache including number of articles, sources, and storage information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			fmt.Println("📊 Analyzing cache...")
		}

		// Create cache and get articles
		cache := news.NewArticleCache()
		articles := cache.GetCachedArticles(0) // Get all

		if len(articles) == 0 {
			fmt.Println("📭 Cache is empty")
			fmt.Println("   Run 'nwcli latest' to populate the cache")
			return nil
		}

		switch format {
		case "json":
			// TODO: Implement JSON output
			fmt.Println("JSON output not yet implemented")
			return nil
		case "plain":
			return renderCacheStatsPlain(articles, cache)
		default: // markdown
			return renderCacheStatsMarkdown(articles, cache)
		}
	},
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "🗑️  Clear the article cache",
	Long:  `Remove all cached articles from local storage. This will force fresh fetching of articles on the next command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			fmt.Println("🗑️  Clearing cache...")
		}

		cache := news.NewArticleCache()
		cache.Clear()

		fmt.Println("✅ Cache cleared successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(cacheStatsCmd)
	cacheCmd.AddCommand(cacheClearCmd)
}

func renderCacheStatsPlain(articles []news.Article, cache *news.ArticleCache) error {
	fmt.Println("Cache Statistics")
	fmt.Println("================")
	fmt.Printf("Total articles: %d\n", len(articles))
	
	if cache.IsStale() {
		fmt.Println("Status: Stale (older than 1 hour)")
	} else {
		fmt.Println("Status: Fresh")
	}
	
	// Count by source
	sourceCount := make(map[string]int)
	for _, article := range articles {
		sourceCount[article.Source]++
	}
	
	fmt.Println("\nArticles by source:")
	for source, count := range sourceCount {
		fmt.Printf("  %s: %d\n", source, count)
	}
	
	return nil
}

func renderCacheStatsMarkdown(articles []news.Article, cache *news.ArticleCache) error {
	renderer, err := renderer.NewMarkdownRenderer()
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	output, err := renderer.RenderStats(articles)
	if err != nil {
		return fmt.Errorf("failed to render stats: %w", err)
	}

	fmt.Print(output)
	return nil
}
