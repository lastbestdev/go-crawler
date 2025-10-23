package main

import (
	"fmt"
	"html/template"
	"os"
)

// render the search tree with HTML template, write to specified output file
func RenderTemplate(root *Node, outputPath string) error {
	tmpl, err := template.New("tree").ParseFiles("templates/tree.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return err
	}
	defer file.Close()

	err = tmpl.ExecuteTemplate(file, "tree.html", root)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return err
	}

	return nil
}
