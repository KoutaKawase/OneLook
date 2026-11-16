package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const queryPlaceholder = "{query}"

type SearchLink struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	URLTemplate   string `json:"url_template"`
	QueryEncoding string `json:"query_encoding"`
}

type linkConfig struct {
	Links []SearchLink `json:"links"`
}

type DisplayLink struct {
	Name string
	URL  string
}

type LinkGroup struct {
	Category string
	Links    []DisplayLink
}

// loadLinks reads and validates the link configuration. Invalid entries are
// returned as warnings so one bad entry does not disable the other links.
func loadLinks(path string) ([]SearchLink, []error, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read link configuration %q: %w", path, err)
	}

	var config linkConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, nil, fmt.Errorf("parse link configuration %q: %w", path, err)
	}
	if len(config.Links) == 0 {
		return nil, nil, fmt.Errorf("link configuration %q contains no links", path)
	}

	validLinks := make([]SearchLink, 0, len(config.Links))
	var warnings []error
	for i, link := range config.Links {
		if err := validateLink(link); err != nil {
			warnings = append(warnings, fmt.Errorf("link %d (%q): %w", i+1, link.Name, err))
			continue
		}
		validLinks = append(validLinks, link)
	}
	if len(validLinks) == 0 {
		return nil, warnings, fmt.Errorf("link configuration %q contains no valid links", path)
	}

	return validLinks, warnings, nil
}

func validateLink(link SearchLink) error {
	if strings.TrimSpace(link.Name) == "" {
		return fmt.Errorf("name must not be empty")
	}
	if strings.TrimSpace(link.Category) == "" {
		return fmt.Errorf("category must not be empty")
	}
	if strings.Count(link.URLTemplate, queryPlaceholder) != 1 {
		return fmt.Errorf("url_template must contain %q exactly once", queryPlaceholder)
	}
	if link.QueryEncoding != "path" && link.QueryEncoding != "query" {
		return fmt.Errorf("unsupported query_encoding %q; use %q or %q", link.QueryEncoding, "path", "query")
	}

	probeURL := strings.Replace(link.URLTemplate, queryPlaceholder, "search-term", 1)
	parsed, err := url.Parse(probeURL)
	if err != nil {
		return fmt.Errorf("invalid url_template: %w", err)
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Hostname() == "" {
		return fmt.Errorf("url_template must be an absolute HTTP or HTTPS URL")
	}
	if parsed.User != nil {
		return fmt.Errorf("url_template must not contain user information")
	}

	schemeEnd := strings.Index(link.URLTemplate, "://")
	authorityStart := schemeEnd + len("://")
	authorityEnd := strings.IndexAny(link.URLTemplate[authorityStart:], "/?#")
	if authorityEnd == -1 {
		authorityEnd = len(link.URLTemplate) - authorityStart
	}
	if strings.Contains(link.URLTemplate[authorityStart:authorityStart+authorityEnd], queryPlaceholder) {
		return fmt.Errorf("url_template must not place %q in the URL host", queryPlaceholder)
	}

	return nil
}

// buildURL applies the encoding method declared by the link configuration.
// Link configuration is validated by loadLinks before this function is used.
func buildURL(link SearchLink, word string) string {
	var encoded string
	switch link.QueryEncoding {
	case "path":
		encoded = url.PathEscape(word)
	case "query":
		encoded = url.QueryEscape(word)
	default:
		return ""
	}
	return strings.Replace(link.URLTemplate, queryPlaceholder, encoded, 1)
}

func groupLinks(links []SearchLink, word string) []LinkGroup {
	groups := make([]LinkGroup, 0)
	groupIndexes := make(map[string]int)

	for _, link := range links {
		index, exists := groupIndexes[link.Category]
		if !exists {
			index = len(groups)
			groupIndexes[link.Category] = index
			groups = append(groups, LinkGroup{Category: link.Category})
		}

		groups[index].Links = append(groups[index].Links, DisplayLink{
			Name: link.Name,
			URL:  buildURL(link, word),
		})
	}

	return groups
}
