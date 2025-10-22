package main

import (
	"fmt"
)

func main() {
	input, err := ReadInput()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Received input URL: %s, search depth: %d\n", input.URL.String(), input.SearchDepth)

	root, err := ProcessURL(input.URL.String(), input.SearchDepth)
	if err != nil {
		fmt.Println(err)
		return
	}

	if root == nil {
		fmt.Println("No data to display.")
		return
	}

	root.PrintTree(1)
}
