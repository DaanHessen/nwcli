package news

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// RSSFetcher handles RSS feed fetching
type RSSFetcher struct {
	parser *gofeed.Parser
	client *http.Client
}

// NewRSSFetcher creates a new RSS fetcher
func NewRSSFetcher() *RSSFetcher {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	parser := gofeed.NewParser()
	parser.Client = client
	
	return &RSSFetcher{
		parser: parser,
		client: client,
	}
}

// FetchFromSource fetches articles from a news source
func (rf *RSSFetcher) FetchFromSource(source Source, fullContent bool) ([]Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	feed, err := rf.parser.ParseURLWithContext(source.URL, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed from %s: %w", source.Name, err)
	}

	var articles []Article
	
	for _, item := range feed.Items {
		article := Article{
			Title:       item.Title,
			Description: item.Description,
			Content:     extractContent(item, fullContent),
			Link:        item.Link,
			Source:      source.Name,
			Categories:  item.Categories,
		}
		
		// Parse published date
		if item.PublishedParsed != nil {
			article.Published = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			article.Published = *item.UpdatedParsed
		} else {
			article.Published = time.Now()
		}
		
		// Extract image URL
		if item.Image != nil && item.Image.URL != "" {
			article.ImageURL = item.Image.URL
		} else {
			// Try to find image in extensions
			article.ImageURL = extractImageFromExtensions(item)
		}
		
		// Clean up description
		article.Description = cleanDescription(article.Description)
		
		articles = append(articles, article)
	}
	
	return articles, nil
}

// extractContent extracts the best available content from feed item
func extractContent(item *gofeed.Item, fullContent bool) string {
	// If full content is requested, try to get more detailed content
	if fullContent {
		// Try content first (usually longer)
		if item.Content != "" {
			return cleanHTML(item.Content)
		}
		
		// For full content, also include extensions that might have more text
		if item.Extensions != nil {
			if contentEncoded, ok := item.Extensions["content"]; ok {
				if encoded, ok := contentEncoded["encoded"]; ok && len(encoded) > 0 {
					if encoded[0].Value != "" {
						return cleanHTML(encoded[0].Value)
					}
				}
			}
		}
	}
	
	// Default to description for summaries or fallback
	if item.Description != "" {
		content := cleanHTML(item.Description)
		// For summary mode, limit length
		if !fullContent && len(content) > 300 {
			content = content[:297] + "..."
		}
		return content
	}
	
	// Fallback to content if description is empty
	if item.Content != "" {
		content := cleanHTML(item.Content)
		if !fullContent && len(content) > 300 {
			content = content[:297] + "..."
		}
		return content
	}
	
	return ""
}

// extractImageFromExtensions tries to find image URLs in feed extensions
func extractImageFromExtensions(item *gofeed.Item) string {
	// Try media:thumbnail or media:content
	if item.Extensions != nil {
		if media, ok := item.Extensions["media"]; ok {
			if thumbnail, ok := media["thumbnail"]; ok && len(thumbnail) > 0 {
				if url, ok := thumbnail[0].Attrs["url"]; ok {
					return url
				}
			}
			if content, ok := media["content"]; ok && len(content) > 0 {
				if url, ok := content[0].Attrs["url"]; ok {
					return url
				}
			}
		}
	}
	
	// Try to extract from enclosures
	if len(item.Enclosures) > 0 {
		for _, enc := range item.Enclosures {
			if strings.HasPrefix(enc.Type, "image/") {
				return enc.URL
			}
		}
	}
	
	return ""
}

// cleanDescription removes HTML tags and cleans up text
func cleanDescription(desc string) string {
	// Remove HTML tags
	cleaned := cleanHTML(desc)
	
	// Limit length
	if len(cleaned) > 300 {
		cleaned = cleaned[:297] + "..."
	}
	
	return strings.TrimSpace(cleaned)
}

// cleanHTML removes HTML tags and entities (basic implementation)
func cleanHTML(html string) string {
	// Simple HTML tag removal
	result := html
	
	// Remove script and style tags with content
	for {
		start := strings.Index(result, "<script")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "</script>")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+9:]
	}
	
	for {
		start := strings.Index(result, "<style")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "</style>")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+8:]
	}
	
	// Remove all other HTML tags
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	
	// Decode common HTML entities
	result = strings.ReplaceAll(result, "&amp;", "&")
	result = strings.ReplaceAll(result, "&lt;", "<")
	result = strings.ReplaceAll(result, "&gt;", ">")
	result = strings.ReplaceAll(result, "&quot;", "\"")
	result = strings.ReplaceAll(result, "&#39;", "'")
	result = strings.ReplaceAll(result, "&nbsp;", " ")
	
	// Clean up whitespace
	lines := strings.Split(result, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}
	
	return strings.Join(cleanLines, "\n")
}
