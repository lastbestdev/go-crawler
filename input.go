package main

import (
	"flag"
	"fmt"
	"net/url"
)

type Input struct {
	URL         url.URL
	SearchDepth int
}

func ReadInput() (*Input, error) {
	// set depth flag
	depth := flag.Int("n", 3, "Specifies the depth links will be crawled from the input URL.")

	// read CLI input
	flag.Parse()

	// read URL argument
	arg := flag.Arg(0)
	url, err := url.Parse(arg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse input to URL: %v", err)
	}

	// validate search depth
	if *depth > 5 || *depth < 1 {
		return nil, fmt.Errorf("invalid search depth provided (%d), must be between 1 and 5", *depth)
	}

	var input = &Input{URL: *url, SearchDepth: *depth}

	return input, nil
}
