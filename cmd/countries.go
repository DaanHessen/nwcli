package cmd

import (
	"fmt"

	"nwcli/pkg/news"

	"github.com/spf13/cobra"
)

var countriesCmd = &cobra.Command{
	Use:   "countries",
	Short: "🌍 List supported countries",
	Long: `Display all countries that NWCLI supports for news fetching.

Each country has curated news sources from major outlets in that region.
Use the country codes with the --country flag in other commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		countries := news.GetAvailableCountries()

		switch format {
		case "json":
			return renderCountriesJSON(countries)
		case "plain":
			return renderCountriesPlain(countries)
		default: // markdown
			return renderCountriesMarkdown(countries)
		}
	},
}

func init() {
	rootCmd.AddCommand(countriesCmd)
}

func renderCountriesMarkdown(countries []string) error {
	fmt.Println("# 🌍 Supported Countries\n")
	fmt.Println("NWCLI supports news sources from the following countries:\n")

	countryNames := map[string]string{
		"nl": "🇳🇱 Netherlands (Dutch)",
		"us": "🇺🇸 United States (English)",
		"uk": "🇬🇧 United Kingdom (English)",
		"de": "🇩🇪 Germany (German)",
		"fr": "🇫🇷 France (French)",
	}

	for _, code := range countries {
		if name, ok := countryNames[code]; ok {
			fmt.Printf("- **%s** - `%s`\n", name, code)
		}
	}

	fmt.Println("\n---\n")
	fmt.Println("**Usage:** Use the country code with `--country` flag")
	fmt.Println("**Example:** `nwcli latest --country us --limit 10`")

	return nil
}

func renderCountriesPlain(countries []string) error {
	fmt.Println("Supported Countries:")
	fmt.Println("==================")

	countryNames := map[string]string{
		"nl": "Netherlands (Dutch)",
		"us": "United States (English)",
		"uk": "United Kingdom (English)",
		"de": "Germany (German)",
		"fr": "France (French)",
	}

	for _, code := range countries {
		if name, ok := countryNames[code]; ok {
			fmt.Printf("  %s - %s\n", code, name)
		}
	}

	fmt.Println("\nUsage: Use the country code with --country flag")
	fmt.Println("Example: nwcli latest --country us --limit 10")

	return nil
}

func renderCountriesJSON(countries []string) error {
	fmt.Println("JSON output not yet implemented")
	return nil
}
