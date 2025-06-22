# 📰 NWCLI - International News CLI

A beautiful command line tool for reading international news with support for multiple countries, full article content, and gorgeous markdown rendering.

## 🌟 Features

- **Multi-Country Support**: Dutch (nl), US (us), UK (uk), German (de), French (fr)
- **Full Article Content**: Choose between summaries or complete articles
- **Beautiful Rendering**: Markdown output with images using Glamour
- **Smart Caching**: Local article storage for offline reading
- **Advanced Search**: Search through titles, descriptions, and content
- **Multiple Formats**: Markdown, JSON, and plain text output
- **Source Filtering**: Filter by specific news sources or categories

## 🚀 Quick Start

```bash
# Get latest Dutch news (default)
./nwcli latest

# Get US news with full articles
./nwcli latest --country us --full --limit 10

# Search for specific topics
./nwcli search "climate change" --country uk

# Get daily digest
./nwcli digest --country nl --full

# List available sources
./nwcli sources --country de

# List supported countries
./nwcli countries
```

## 📋 Commands

### `latest` - Get Latest News
```bash
./nwcli latest [flags]

Flags:
  -l, --limit int       number of articles to show (default 20)
  -s, --source string   filter by source (e.g., 'NOS', 'CNN')
  -c, --category string filter by category (general, sports, tech)
      --country string  country code (default "nl")
      --full           fetch full article content instead of summaries
  -v, --verbose        verbose output
  -f, --format string  output format (markdown, json, plain) (default "markdown")
```

### `search` - Search Articles
```bash
./nwcli search "query" [flags]

Flags:
  -l, --limit int      number of results to show (default 20)
  -s, --source string  filter results by source
      --country string country code (default "nl")
      --full          search in full article content
```

### `digest` - Daily News Digest
```bash
./nwcli digest [flags]

Flags:
  -l, --limit int              number of articles in digest (default 15)
  -c, --categories strings     categories to include (general, sports, tech)
      --country string         country code (default "nl")
      --full                  include full article content
```

### `sources` - List News Sources
```bash
./nwcli sources [flags]

Flags:
      --country string country code (default "nl")
```

### `countries` - Supported Countries
```bash
./nwcli countries
```

Shows all supported countries with their codes and languages.

### `cache` - Cache Management
```bash
./nwcli cache stats    # Show cache statistics
./nwcli cache clear    # Clear the cache
```

## 🌍 Supported Countries

| Country | Code | Language | Sample Sources |
|---------|------|----------|----------------|
| 🇳🇱 Netherlands | `nl` | Dutch | NOS, NU.nl, De Telegraaf, RTL Nieuws |
| 🇺🇸 United States | `us` | English | CNN, BBC, Reuters, NPR |
| 🇬🇧 United Kingdom | `uk` | English | BBC UK, The Guardian, Sky News |
| 🇩🇪 Germany | `de` | German | Tagesschau, SPIEGEL, ZEIT |
| 🇫🇷 France | `fr` | French | Le Monde, France 24, Libération |

## 🎨 Output Formats

- **Markdown** (default): Beautiful newspaper-like layout with images
- **JSON**: Machine-readable format for automation
- **Plain**: Simple text output for scripting

## 💾 Caching

NWCLI automatically caches articles in `~/.nwcli/cache/` for:
- Faster subsequent access
- Offline reading capability
- Reduced network requests
- Search functionality

## 🔧 Installation

```bash
# Build from source
go build -o nwcli .

# Run directly
./nwcli --help
```

## 📖 Examples

```bash
# Morning news routine - Dutch full articles
./nwcli digest --full --limit 20

# US tech news search
./nwcli search "artificial intelligence" --country us --limit 5

# Quick German headlines
./nwcli latest --country de --limit 10 --format plain

# Sports news from multiple sources
./nwcli latest --category sports --limit 15

# Check what sources are available in France
./nwcli sources --country fr
```

## 🏗️ Architecture

- **RSS Fetching**: Reliable RSS feed parsing with `gofeed`
- **Markdown Rendering**: Beautiful terminal output with `glamour`
- **CLI Framework**: Robust command structure with `cobra`
- **Caching**: Local JSON-based article storage
- **Error Handling**: Graceful degradation when sources are unavailable

## 🤝 Contributing

Feel free to add more news sources or countries by extending the source lists in `pkg/news/service.go`.

---

*Built with ❤️ for news enthusiasts worldwide*
