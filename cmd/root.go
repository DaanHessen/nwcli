package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nwcli",
	Short: "📰 NWCLI - Your International News Command Line Interface",
	Long: `📰 NWCLI - A beautiful command line tool for reading international news

Get the latest news from popular sources worldwide, search articles, 
filter content, and enjoy a newspaper-like experience in your terminal
with beautiful markdown rendering.

Supported countries: Netherlands (nl), United States (us), United Kingdom (uk), 
Germany (de), France (fr).

Default sources include Dutch news: NOS, NU.nl, De Telegraaf, RTL Nieuws, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()
	
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("format", "f", "markdown", "output format (markdown, json, plain)")
}
