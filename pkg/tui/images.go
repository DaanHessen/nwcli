package tui

import (
	"fmt"
	"os"
	"strings"
)

// TerminalImageSupport detects available image protocols
type TerminalImageSupport struct {
	SupportsKitty bool
	SupportsSixel bool
	SupportsITerm bool
}

// DetectImageSupport checks what image protocols the terminal supports
func DetectImageSupport() TerminalImageSupport {
	support := TerminalImageSupport{}
	
	// Check for Kitty terminal
	if strings.Contains(os.Getenv("TERM"), "kitty") || 
	   os.Getenv("KITTY_WINDOW_ID") != "" {
		support.SupportsKitty = true
	}
	
	// Check for iTerm2 
	if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		support.SupportsITerm = true
	}
	
	// Check for Sixel support (basic check)
	termFeatures := os.Getenv("TERM_FEATURES")
	if strings.Contains(termFeatures, "sixel") {
		support.SupportsSixel = true
	}
	
	return support
}

// RenderImage attempts to render an image in the terminal if supported
func (support TerminalImageSupport) RenderImage(imageURL string) string {
	if imageURL == "" {
		return ""
	}
	
	// For now, we'll just provide a placeholder with the URL
	// Full implementation would require downloading and encoding images
	if support.SupportsKitty || support.SupportsSixel || support.SupportsITerm {
		return fmt.Sprintf("\n🖼️  Image: %s\n(Terminal image rendering is detected but not yet implemented)\n\n", imageURL)
	}
	
	// Fallback for terminals without image support
	return fmt.Sprintf("\n🔗 Image available: %s\n\n", imageURL)
}

// GetImagePlaceholder returns a styled placeholder for images
func GetImagePlaceholder(imageURL string) string {
	if imageURL == "" {
		return ""
	}
	
	support := DetectImageSupport()
	return support.RenderImage(imageURL)
}
