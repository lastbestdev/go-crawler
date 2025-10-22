package main

import (
	"fmt"
	"regexp"
	"strings"
)

func ProcessURL(url string, searchDepth int) (*Node, error) {
	root, err := process(url, 0, searchDepth)
	if err != nil {
		return nil, fmt.Errorf("error processing URL (%s): %v", url, err)
	}

	return root, nil
}

func process(url string, depth, maxDepth int) (*Node, error) {
	if depth > maxDepth {
		error := fmt.Errorf("max depth reached at URL: %s. breaking search", url)
		return nil, error
	}
	urlRegex := regexp.MustCompile(`^https?://`)
	if !urlRegex.MatchString(url) {
		error := fmt.Errorf("invalid URL format (%s). skipping url", url)
		return nil, error
	}

	body, err := FetchURLContent(url)
	if err != nil {
		fmt.Printf("Error fetching URL (%s): %v\n", url, err)
		return nil, err
	}

	title, err := findTitle(body)
	// pages without <title> tags found get default "Untitled" value
	if err != nil {
		fmt.Printf("Error finding title: %v\n", err)
		title = "Untitled"
	}

	node, err := MakeNode(title, url)
	if err != nil {
		fmt.Printf("Error creating node: %v\n", err)
		return nil, err
	}

	links, err := findChildLinks(body)
	// when errors occur finding links, return node without children
	if err != nil {
		fmt.Printf("Error processing content: %v\n", err)
		return node, nil
	}

	for _, link := range links {
		// handle relative links
		if strings.HasPrefix(link, "/") {
			link = url + link
		}

		link = strings.TrimSuffix(link, "/")

		// skip circular references (i.e. page links to self)
		if link == url {
			continue
		}

		child, err := process(link, depth+1, maxDepth)
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
