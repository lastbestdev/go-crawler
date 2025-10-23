package main

import (
	"errors"
	"fmt"
	"net/url"
)

type Node struct {
	Title    string
	URL      *url.URL
	Children []*Node
}

// create Nodes with page title and URL
func MakeNode(title string, urlStr string) (*Node, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.New("unable to create node, invalid URL")
	}

	return &Node{
		Title:    title,
		URL:      parsedURL,
		Children: []*Node{},
	}, nil
}

// add child Nodes to parent Node
func (n *Node) AddChild(child *Node) {
	if child == nil {
		return
	}

	n.Children = append(n.Children, child)
}

// print tree to console for debugging
func (n *Node) PrintTree(level int) {
	fmt.Printf("node (title: %s, url: %s) at level %d\n", n.Title, n.URL.String(), level)

	if len(n.Children) == 0 {
		return
	}

	for _, child := range n.Children {
		child.PrintTree(level + 1)
	}
}
