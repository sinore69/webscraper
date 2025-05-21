package main

import (
	"fmt"
	"io"
	"os"
	"sync"

	"webscraper/functions"
)

const outputFileName = "output.txt"


var (
	visited = make(map[string]bool)
	mu      sync.Mutex
	wg      sync.WaitGroup
)

func scrape(link string, depth int, writer io.Writer) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	mu.Lock()
	if visited[link] {
		mu.Unlock()
		return
	}
	visited[link] = true
	mu.Unlock()

	body, err := functions.FetchHTML(link)
	if err != nil {
		mu.Lock()
		fmt.Fprintf(writer, "\n--- URL: %s ---\nError: %v\n", link, err)
		mu.Unlock()
		return
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		mu.Lock()
		fmt.Fprintf(writer, "\n--- URL: %s ---\nError reading content: %v\n", link, err)
		mu.Unlock()
		return
	}

	mu.Lock()
	fmt.Fprintf(writer, "\n--- URL: %s ---\n%s\n", link, content)
	mu.Unlock()

	body2, err := functions.FetchHTML(link)
	if err != nil {
		return
	}
	defer body2.Close()

	links, err := functions.ExtractLinks(link, body2)
	if err != nil {
		return
	}

	for _, l := range links {
		wg.Add(1)
		go scrape(l, depth-1, writer)
	}
}


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url> [depth]")
		return
	}

	startURL := os.Args[1]
	depth := 3 
	if len(os.Args) >= 3 {
		fmt.Sscanf(os.Args[2], "%d", &depth)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		return
	}
	defer outputFile.Close()

	wg.Add(1)
	go scrape(startURL, depth, outputFile)

	wg.Wait()
	fmt.Println("Scraping completed. Output saved to", outputFileName)
}
