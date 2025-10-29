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
	fmt.Printf("seed URL: %s, search depth: %d\n", input.URL.String(), input.SearchDepth)

	root, err := ProcessURL(*input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if root == nil {
		fmt.Println("No data to display.")
		return
	}

	err = RenderTemplate(root, "templates/output.html")
	if err == nil {
		fmt.Printf("output written to /templates/output.html. run `open ./templates/output.html` to view in browser.")
	}
}
