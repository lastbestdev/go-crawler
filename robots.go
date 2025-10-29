package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type RobotsTxt struct {
	AllowedPaths    []string
	DisallowedPaths []string
	CrawlDelay      int
}

// cache of processed robots.txt files, to avoid refetch/processing (key: root URL, value: parsed robots.txt rules)
var robotsTxtCache = make(map[string]RobotsTxt)

func AddRobotsTxtCache(url string, parsed RobotsTxt) {
	robotsTxtCache[url] = parsed
	fmt.Println("cached robots.txt rules for site ", url)
}

func GetRobotsTxtCache(url string) (RobotsTxt, bool) {
	robotsTxt, found := robotsTxtCache[url]
	return robotsTxt, found
}

func ReadRobotsTxt(content string) (*RobotsTxt, error) {
	allowedPaths := []string{}
	disallowedPaths := []string{}
	crawlDelay := 0

	userAgentRegex := regexp.MustCompile(`User-[Aa]gent: \*`)
	lines := strings.Split(content, "\n")
	parse := false
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		// begin parsing when User-agent: * rules found
		if !parse && userAgentRegex.MatchString(line) {
			parse = true
		}
		if !parse {
			continue
		}

		// stop parsing rules when line break encountered
		if line == "" {
			break
		}

		// parse Allow, Disallow, and Crawl-delay rules
		if strings.HasPrefix(line, "Allow: ") {
			allowedPaths = append(allowedPaths, strings.TrimPrefix(line, "Allow: "))
		} else if strings.HasPrefix(line, "Disallow: ") {
			disallowedPaths = append(disallowedPaths, strings.TrimPrefix(line, "Disallow: "))
		} else if strings.HasPrefix(line, "Crawl-delay: ") {
			fmt.Sscanf(line, "Crawl-delay: %d", &crawlDelay)
		}
	}

	robotsTxt := &RobotsTxt{
		AllowedPaths:    allowedPaths,
		DisallowedPaths: disallowedPaths,
		CrawlDelay:      crawlDelay,
	}

	return robotsTxt, nil
}

func CheckCrawlOk(url url.URL) bool {
	rootUrl := url.Scheme + "://" + url.Host
	robotsTxt, found := GetRobotsTxtCache(rootUrl)

	// if rules not found for site, assume free to crawl
	if !found {
		fmt.Printf("no robots.txt rules found for site %s, free to crawl\n", rootUrl)
		return true
	}

	// check if the URL is explicitly allowed
	for _, allowedPath := range robotsTxt.AllowedPaths {
		if strings.HasPrefix(url.Path, allowedPath) {
			return true
		}
	}

	// check if the URL is disallowed
	for _, disallowedPath := range robotsTxt.DisallowedPaths {
		if strings.HasPrefix(url.Path, disallowedPath) {
			return false
		}
	}

	return true
}
