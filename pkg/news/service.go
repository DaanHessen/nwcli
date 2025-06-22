package news

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Article represents a news article
type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Link        string    `json:"link"`
	Published   time.Time `json:"published"`
	Source      string    `json:"source"`
	ImageURL    string    `json:"image_url,omitempty"`
	Categories  []string  `json:"categories,omitempty"`
}

// Source represents a news source
type Source struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Category    string `json:"category"`
}

// NewsService handles news operations
type NewsService struct {
	sources     []Source
	fetcher     *RSSFetcher
	cache       *ArticleCache
	country     string
	fullContent bool
}

// NewNewsService creates a new news service
func NewNewsService() *NewsService {
	return &NewsService{
		sources:     getSourcesByCountry("nl"), // Default to Dutch
		fetcher:     NewRSSFetcher(),
		cache:       NewArticleCache(),
		country:     "nl",
		fullContent: false,
	}
}

// NewNewsServiceWithOptions creates a news service with specific options
func NewNewsServiceWithOptions(country string, fullContent bool) *NewsService {
	return &NewsService{
		sources:     getSourcesByCountry(country),
		fetcher:     NewRSSFetcher(),
		cache:       NewArticleCache(),
		country:     country,
		fullContent: fullContent,
	}
}

// GetLatestNews fetches latest news from all sources
func (ns *NewsService) GetLatestNews(limit int) ([]Article, error) {
	var allArticles []Article

	for _, source := range ns.sources {
		articles, err := ns.fetcher.FetchFromSource(source, ns.fullContent)
		if err != nil {
			fmt.Printf("Warning: Failed to fetch from %s: %v\n", source.Name, err)
			continue
		}
		allArticles = append(allArticles, articles...)
	}

	// Sort by published date (newest first)
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].Published.After(allArticles[j].Published)
	})

	// Apply limit
	if limit > 0 && len(allArticles) > limit {
		allArticles = allArticles[:limit]
	}

	// Cache articles
	ns.cache.StoreArticles(allArticles)

	return allArticles, nil
}

// SearchArticles searches articles by keywords
func (ns *NewsService) SearchArticles(query string, limit int) ([]Article, error) {
	// First try from cache
	cached := ns.cache.SearchArticles(query, limit)
	if len(cached) > 0 {
		return cached, nil
	}

	// Fetch fresh articles and search
	articles, err := ns.GetLatestNews(0) // Get all
	if err != nil {
		return nil, err
	}

	var matches []Article
	query = strings.ToLower(query)

	for _, article := range articles {
		if strings.Contains(strings.ToLower(article.Title), query) ||
			strings.Contains(strings.ToLower(article.Description), query) ||
			strings.Contains(strings.ToLower(article.Content), query) {
			matches = append(matches, article)
		}
	}

	// Apply limit
	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}

	return matches, nil
}

// FilterArticles filters articles by source and category
func (ns *NewsService) FilterArticles(sourceName, category string, since time.Time, limit int) ([]Article, error) {
	articles, err := ns.GetLatestNews(0)
	if err != nil {
		return nil, err
	}

	var filtered []Article

	for _, article := range articles {
		// Filter by source
		if sourceName != "" && !strings.EqualFold(article.Source, sourceName) {
			continue
		}

		// Filter by time
		if !since.IsZero() && article.Published.Before(since) {
			continue
		}

		// Filter by category
		if category != "" {
			found := false
			for _, cat := range article.Categories {
				if strings.EqualFold(cat, category) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, article)
	}

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	return filtered, nil
}

// GetSources returns available news sources
func (ns *NewsService) GetSources() []Source {
	return ns.sources
}

// getSourcesByCountry returns news sources for a specific country
func getSourcesByCountry(country string) []Source {
	switch strings.ToLower(country) {
	case "nl", "netherlands", "dutch":
		return getDutchSources()
	case "us", "usa", "united-states":
		return getUSSources()
	case "uk", "gb", "britain":
		return getUKSources()
	case "de", "germany", "german":
		return getGermanSources()
	case "fr", "france", "french":
		return getFrenchSources()
	default:
		// Default to Dutch sources
		return getDutchSources()
	}
}

// GetAvailableCountries returns list of supported countries
func GetAvailableCountries() []string {
	return []string{"nl", "us", "uk", "de", "fr"}
}

// getDutchSources returns Dutch news sources
func getDutchSources() []Source {
	return []Source{
		{
			Name:        "NOS",
			URL:         "https://feeds.nos.nl/nosnieuwsalgemeen",
			Description: "Nederlandse Omroep Stichting - General News",
			Language:    "nl",
			Category:    "general",
		},
		{
			Name:        "NU.nl",
			URL:         "https://www.nu.nl/rss/Algemeen",
			Description: "NU.nl - General News",
			Language:    "nl",
			Category:    "general",
		},
		{
			Name:        "De Telegraaf",
			URL:         "https://www.telegraaf.nl/rss",
			Description: "De Telegraaf - News",
			Language:    "nl",
			Category:    "general",
		},
		{
			Name:        "RTL Nieuws",
			URL:         "https://www.rtlnieuws.nl/rss.xml",
			Description: "RTL Nieuws - Latest News",
			Language:    "nl",
			Category:    "general",
		},
		{
			Name:        "AD.nl",
			URL:         "https://www.ad.nl/rss.xml",
			Description: "Algemeen Dagblad - News",
			Language:    "nl",
			Category:    "general",
		},
		{
			Name:        "NOS Sport",
			URL:         "https://feeds.nos.nl/nossport",
			Description: "NOS - Sports News",
			Language:    "nl",
			Category:    "sports",
		},
		{
			Name:        "NU.nl Tech",
			URL:         "https://www.nu.nl/rss/Tech",
			Description: "NU.nl - Technology News",
			Language:    "nl",
			Category:    "technology",
		},
	}
}

// getUSSources returns US news sources
func getUSSources() []Source {
	return []Source{
		{
			Name:        "CNN",
			URL:         "http://rss.cnn.com/rss/edition.rss",
			Description: "CNN - Breaking News",
			Language:    "en",
			Category:    "general",
		},
		{
			Name:        "BBC News",
			URL:         "http://feeds.bbci.co.uk/news/rss.xml",
			Description: "BBC News - Home",
			Language:    "en",
			Category:    "general",
		},
		{
			Name:        "Reuters",
			URL:         "https://feeds.reuters.com/reuters/topNews",
			Description: "Reuters - Top News",
			Language:    "en",
			Category:    "general",
		},
		{
			Name:        "NPR",
			URL:         "https://feeds.npr.org/1001/rss.xml",
			Description: "NPR - News",
			Language:    "en",
			Category:    "general",
		},
	}
}

// getUKSources returns UK news sources
func getUKSources() []Source {
	return []Source{
		{
			Name:        "BBC UK",
			URL:         "http://feeds.bbci.co.uk/news/uk/rss.xml",
			Description: "BBC News - UK",
			Language:    "en",
			Category:    "general",
		},
		{
			Name:        "The Guardian",
			URL:         "https://www.theguardian.com/uk/rss",
			Description: "The Guardian - UK News",
			Language:    "en",
			Category:    "general",
		},
		{
			Name:        "Sky News",
			URL:         "http://feeds.skynews.com/feeds/rss/home.xml",
			Description: "Sky News - Latest News",
			Language:    "en",
			Category:    "general",
		},
	}
}

// getGermanSources returns German news sources
func getGermanSources() []Source {
	return []Source{
		{
			Name:        "Tagesschau",
			URL:         "https://www.tagesschau.de/xml/rss2/",
			Description: "Tagesschau - Nachrichten",
			Language:    "de",
			Category:    "general",
		},
		{
			Name:        "SPIEGEL ONLINE",
			URL:         "https://www.spiegel.de/schlagzeilen/index.rss",
			Description: "SPIEGEL ONLINE - Schlagzeilen",
			Language:    "de",
			Category:    "general",
		},
		{
			Name:        "ZEIT ONLINE",
			URL:         "https://newsfeed.zeit.de/index",
			Description: "ZEIT ONLINE - Nachrichten",
			Language:    "de",
			Category:    "general",
		},
	}
}

// getFrenchSources returns French news sources
func getFrenchSources() []Source {
	return []Source{
		{
			Name:        "Le Monde",
			URL:         "https://www.lemonde.fr/rss/une.xml",
			Description: "Le Monde - À la une",
			Language:    "fr",
			Category:    "general",
		},
		{
			Name:        "France 24",
			URL:         "https://www.france24.com/fr/rss",
			Description: "France 24 - Actualités",
			Language:    "fr",
			Category:    "general",
		},
		{
			Name:        "Liberation",
			URL:         "https://www.liberation.fr/arc/outboundfeeds/rss/",
			Description: "Libération - Actualités",
			Language:    "fr",
			Category:    "general",
		},
	}
}
