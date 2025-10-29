package main

import (
	"fmt"
	"net/url"
)

type Node struct {
	Title    string
	URL      url.URL
	Children []*Node
}

// create Nodes with page title and URL
func MakeNode(title string, url url.URL) (*Node, error) {
	return &Node{
		Title:    title,
		URL:      url,
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
