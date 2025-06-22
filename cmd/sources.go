package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"nwcli/pkg/news"
	"nwcli/pkg/renderer"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "📡 List available news sources",
	Long: `Display all available Dutch news sources that NWCLI can fetch from.

This includes major Dutch news outlets like NOS, NU.nl, De Telegraaf,
RTL Nieuws, and others. Each source shows its description, category,
and RSS feed URL.

Use this information to filter news by specific sources using the
--source flag in other commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		verbose, _ := cmd.Flags().GetBool("verbose")
		country, _ := cmd.Flags().GetString("country")

		if verbose {
			fmt.Printf("📡 Loading news sources for %s...\n", country)
		}

		// Create news service and get sources
		newsService := news.NewNewsServiceWithOptions(country, false)
		sources := newsService.GetSources()

		if verbose {
			fmt.Printf("✅ Found %d sources\n\n", len(sources))
		}

		// Render based on format
		switch format {
		case "json":
			// TODO: Implement JSON output
			fmt.Println("JSON output not yet implemented")
			return nil
		case "plain":
			return renderSourcesPlain(sources)
		default: // markdown
			return renderSourcesMarkdown(sources)
		}
	},
}

func init() {
	rootCmd.AddCommand(sourcesCmd)
	
	// Flags
	sourcesCmd.Flags().StringP("country", "", "nl", "country code (nl, us, uk, de, fr)")
}

func renderSourcesMarkdown(sources []news.Source) error {
	renderer, err := renderer.NewMarkdownRenderer()
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}

	output, err := renderer.RenderSources(sources)
	if err != nil {
		return fmt.Errorf("failed to render sources: %w", err)
	}

	fmt.Print(output)
	return nil
}

func renderSourcesPlain(sources []news.Source) error {
	// Determine country name from sources
	countryName := "News Sources"
	if len(sources) > 0 {
		switch sources[0].Language {
		case "nl":
			countryName = "Dutch News Sources"
		case "en":
			countryName = "English News Sources"
		case "de":
			countryName = "German News Sources"
		case "fr":
			countryName = "French News Sources"
		}
	}
	
	fmt.Printf("Available %s:\n", countryName)
	fmt.Println(strings.Repeat("=", len(countryName)+11))

	// Group by category
	categories := make(map[string][]news.Source)
	for _, source := range sources {
		categories[source.Category] = append(categories[source.Category], source)
	}

	for category, sources := range categories {
		fmt.Printf("\n%s:\n", category)
		fmt.Println(strings.Repeat("-", len(category)+1))
		
		for _, source := range sources {
			fmt.Printf("• %s\n", source.Name)
			fmt.Printf("  %s\n", source.Description)
			fmt.Printf("  URL: %s\n", source.URL)
			fmt.Printf("  Language: %s\n\n", source.Language)
		}
	}

	return nil
}
