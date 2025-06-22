package news

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// ArticleCache handles caching of articles
type ArticleCache struct {
	articles   []Article
	lastUpdate time.Time
	mu         sync.RWMutex
	cacheDir   string
}

// NewArticleCache creates a new article cache
func NewArticleCache() *ArticleCache {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".nwcli", "cache")
	
	// Create cache directory if it doesn't exist
	os.MkdirAll(cacheDir, 0755)
	
	cache := &ArticleCache{
		cacheDir: cacheDir,
	}
	
	// Load cached articles
	cache.loadFromDisk()
	
	return cache
}

// StoreArticles stores articles in cache
func (ac *ArticleCache) StoreArticles(articles []Article) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	
	// Merge with existing articles, avoiding duplicates
	existingMap := make(map[string]bool)
	for _, article := range ac.articles {
		existingMap[article.Link] = true
	}
	
	for _, article := range articles {
		if !existingMap[article.Link] {
			ac.articles = append(ac.articles, article)
			existingMap[article.Link] = true
		}
	}
	
	ac.lastUpdate = time.Now()
	
	// Save to disk
	ac.saveToDisk()
}

// SearchArticles searches cached articles
func (ac *ArticleCache) SearchArticles(query string, limit int) []Article {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	
	var matches []Article
	query = strings.ToLower(query)
	
	for _, article := range ac.articles {
		if strings.Contains(strings.ToLower(article.Title), query) ||
		   strings.Contains(strings.ToLower(article.Description), query) ||
		   strings.Contains(strings.ToLower(article.Content), query) {
			matches = append(matches, article)
		}
		
		if limit > 0 && len(matches) >= limit {
			break
		}
	}
	
	return matches
}

// GetCachedArticles returns all cached articles
func (ac *ArticleCache) GetCachedArticles(limit int) []Article {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	
	if limit > 0 && len(ac.articles) > limit {
		return ac.articles[:limit]
	}
	
	return ac.articles
}

// IsStale checks if cache is stale (older than 1 hour)
func (ac *ArticleCache) IsStale() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	
	return time.Since(ac.lastUpdate) > time.Hour
}

// Clear clears the cache
func (ac *ArticleCache) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	
	ac.articles = nil
	ac.lastUpdate = time.Time{}
	
	// Remove cache file
	os.Remove(filepath.Join(ac.cacheDir, "articles.json"))
}

// saveToDisk saves articles to disk
func (ac *ArticleCache) saveToDisk() {
	cacheFile := filepath.Join(ac.cacheDir, "articles.json")
	
	data := map[string]interface{}{
		"articles":    ac.articles,
		"last_update": ac.lastUpdate,
	}
	
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Warning: Failed to marshal cache data: %v\n", err)
		return
	}
	
	err = os.WriteFile(cacheFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Warning: Failed to save cache to disk: %v\n", err)
	}
}

// loadFromDisk loads articles from disk
func (ac *ArticleCache) loadFromDisk() {
	cacheFile := filepath.Join(ac.cacheDir, "articles.json")
	
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		// Cache file doesn't exist or can't be read
		return
	}
	
	var cacheData map[string]interface{}
	err = json.Unmarshal(data, &cacheData)
	if err != nil {
		fmt.Printf("Warning: Failed to unmarshal cache data: %v\n", err)
		return
	}
	
	// Parse articles
	if articlesData, ok := cacheData["articles"].([]interface{}); ok {
		for _, articleData := range articlesData {
			if articleMap, ok := articleData.(map[string]interface{}); ok {
				article := parseArticleFromMap(articleMap)
				ac.articles = append(ac.articles, article)
			}
		}
	}
	
	// Parse last update time
	if lastUpdateStr, ok := cacheData["last_update"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, lastUpdateStr); err == nil {
			ac.lastUpdate = parsed
		}
	}
}

// parseArticleFromMap parses an article from a map
func parseArticleFromMap(data map[string]interface{}) Article {
	article := Article{}
	
	if title, ok := data["title"].(string); ok {
		article.Title = title
	}
	if desc, ok := data["description"].(string); ok {
		article.Description = desc
	}
	if content, ok := data["content"].(string); ok {
		article.Content = content
	}
	if link, ok := data["link"].(string); ok {
		article.Link = link
	}
	if source, ok := data["source"].(string); ok {
		article.Source = source
	}
	if imageURL, ok := data["image_url"].(string); ok {
		article.ImageURL = imageURL
	}
	if publishedStr, ok := data["published"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, publishedStr); err == nil {
			article.Published = parsed
		}
	}
	if categoriesData, ok := data["categories"].([]interface{}); ok {
		for _, cat := range categoriesData {
			if catStr, ok := cat.(string); ok {
				article.Categories = append(article.Categories, catStr)
			}
		}
	}
	
	return article
}
