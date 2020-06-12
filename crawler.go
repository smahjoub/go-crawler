package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 1 {
		fmt.Println("Web site url is missing")
		os.Exit(1)
	}
	uri := os.Args[1]

	visited := make(map[string]bool)

	enqueue(uri, &visited)
}
