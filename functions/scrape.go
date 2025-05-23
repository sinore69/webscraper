package functions

import (
	"fmt"
	"io"
	"sync"
)

func Scrape(link string, depth int, writer io.Writer, wg *sync.WaitGroup, mu *sync.Mutex, visited *map[string]bool, sem chan struct{}) {

	defer wg.Done()

	if depth <= 0 {
		return
	}

	sem <- struct{}{}
	defer func() { <-sem }()

	mu.Lock()
	if (*visited)[link] {
		mu.Unlock()
		return
	}
	(*visited)[link] = true
	mu.Unlock()

	body, err := FetchHTML(link)
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

	body2, err := FetchHTML(link)
	if err != nil {
		return
	}
	defer body2.Close()

	links, err := ExtractLinks(link, body2)
	if err != nil {
		return
	}

	for _, l := range links {
		wg.Add(1)
		go Scrape(l, depth-1, writer, wg, mu, visited, sem)
	}
}
