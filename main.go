package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"sync"

	"webscraper/functions"

	tea "github.com/charmbracelet/bubbletea"
)

const outputFileName = "output.txt"

var (
	visited = make(map[string]bool)
	mu      sync.Mutex
	wg      sync.WaitGroup
)

type Model struct {
	url         string
	concurrency int
	phase       int
	depth       int
	input       string
	message     string
}

// Init is called when the program starts
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key events
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyEnter:
			if m.input == "" {
				break
			}
			if m.phase == 0 {
				if isValidURL(m.input) {
					m.url = m.input
					m.phase = m.phase + 1
					m.input = ""
					break
				} else {
					m.message = "Invalid URL"
				}
			}
			if m.phase == 1 {
				if m.input == "" {
					break
				}
				num, err := strconv.Atoi(m.input)
				if err != nil {
					m.message = "Invalid number1"
					break
				}
				if num <= 0 {
					m.message = "Concurrency must be greater than 0"
					break
				}
				m.concurrency = num
				m.input = ""
				m.phase = m.phase + 1
			}
			if m.phase == 2 {
				if m.input == "" {
					break
				}
				num, err := strconv.Atoi(m.input)
				if err != nil {
					m.message = fmt.Sprintf("%v", err)
					break
				}
				if num <= 0 {
					m.message = "Depth must be greater than 0"
					break
				}
				m.depth = num
				m.input = ""
				m.phase = m.phase + 1
			}
			if m.phase == 3 {
				outputFile, err := os.Create(outputFileName)
				if err != nil {
					fmt.Println("Failed to create output file:", err)
					break
				}
				defer outputFile.Close()
				wg.Add(1)
				go scrape(m.url, m.depth, outputFile)

				wg.Wait()
				fmt.Println("Scraping completed. Output saved to", outputFileName)
				os.Exit(0)
			}
		case tea.KeyRunes:
			m.input += string(msg.Runes)
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	switch m.phase {
	case 0:
		return fmt.Sprintf("Enter URL: %s\n%s", m.input, m.message)
	case 1:
		return fmt.Sprintf("Enter Max Concurrency level: %s\n%s", m.input, m.message)
	case 2:
		return fmt.Sprintf("Enter Max Depth level: %s\n%s", m.input, m.message)
	// case 3:
	// 	return fmt.Sprintln("scrapping website")
	default:
		return fmt.Sprintf("%+v", m)
	}
}

func isValidURL(s string) bool {
	parsed, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return parsed.Scheme != "" && parsed.Host != ""
}

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

	p := tea.NewProgram(Model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
