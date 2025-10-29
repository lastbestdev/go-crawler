package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func ProcessURL(input Input) (*Node, error) {
	root, err := process(input.URL, 0, input.SearchDepth)
	if err != nil {
		return nil, fmt.Errorf("error processing seed URL (%s): %v", input.URL.String(), err)
	}

	return root, nil
}

func process(url url.URL, depth, maxDepth int) (*Node, error) {
	if depth > maxDepth {
		error := fmt.Errorf("max depth reached at URL: %s. breaking search", url.String())
		return nil, error
	}

	// first check robots.txt cache
	rootUrl := url.Scheme + "://" + url.Host
	robotsTxt, found := GetRobotsTxtCache(rootUrl)

	// when crawler rules not found, fetch/parse them into the cache
	if !found {
		robotsUrl := rootUrl + "/robots.txt"
		content, err := FetchURLContent(robotsUrl)

		if err != nil {
			fmt.Printf("error fetching robots.txt: %v\n", err)
		} else {
			parsed, err := ReadRobotsTxt(content)
			if err != nil {
				fmt.Printf("error parsing robots.txt: %v\n", err)
			}

			AddRobotsTxtCache(rootUrl, *parsed)
			robotsTxt = *parsed
		}
	}

	// before proceeding to crawl, ensure abiding by robots.txt rules
	ok := CheckCrawlOk(url)
	if !ok {
		fmt.Printf("crawling URL %s is not allowed by robots.txt. skipping\n", url.String())
		return nil, nil
	}

	// observe crawl delay if specified
	if robotsTxt.CrawlDelay > 0 {
		fmt.Printf("delay crawl %d seconds for site %s\n", robotsTxt.CrawlDelay, rootUrl)
		time.Sleep(time.Duration(robotsTxt.CrawlDelay) * time.Second)
	}

	fmt.Printf("crawling URL: %s at depth %d\n", url.String(), depth)

	// allowed to crawl - fetch the URL contents
	body, err := FetchURLContent(url.String())
	if err != nil {
		return nil, err
	}

	// pages without <title> tags get default "Untitled" value
	title, err := findTitle(body)
	if err != nil {
		title = "Untitled"
	}

	node, err := MakeNode(title, url)
	if err != nil {
		fmt.Printf("error creating node: %v\n", err)
		return nil, err
	}

	links, err := findChildLinks(body)
	// when errors occur finding links, return node without children
	if err != nil {
		fmt.Printf("error finding child links for page: %v", err)
		return node, nil
	}

	// iterate child links and recursively process
	parentUrl := url.String()
	for _, link := range links {
		// handle relative links
		if strings.HasPrefix(link, "/") {
			link = parentUrl + link
		}

		link = strings.TrimSuffix(link, "/")

		// skip circular references (i.e. page links to self)
		if link == parentUrl {
			continue
		}

		childUrl, err := url.Parse(link)
		if err != nil {
			continue
		}

		child, err := process(*childUrl, depth+1, maxDepth)
		if err != nil {
			continue
		}

		node.AddChild(child)
	}

	return node, nil
}

func findTitle(body string) (string, error) {
	// regex to find title tag
	regex := regexp.MustCompile(`<title>(.*?)</title>`)
	match := regex.FindStringSubmatch(body)

	// no title found
	if len(match) < 2 {
		return "", fmt.Errorf("no title tag found in the body")
	}

	return match[1], nil
}

func findChildLinks(body string) ([]string, error) {
	var links []string

	// regex to find anchor tags
	regex := regexp.MustCompile(`<a[^>]*>(.*?)</a>`)
	matches := regex.FindAllString(string(body), -1)

	// no links found
	if matches == nil {
		return links, nil
	}

	// extract href attributes
	hrefRegex := regexp.MustCompile(`href=["'](.*?)["']`)
	for match := range matches {
		linkText := matches[match]
		hrefMatch := hrefRegex.FindStringSubmatch(linkText)
		if len(hrefMatch) > 1 {
			links = append(links, hrefMatch[1])
		}
	}

	return links, nil
}
