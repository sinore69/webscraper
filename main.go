package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"webscraper/functions"

	tea "github.com/charmbracelet/bubbletea"
)

const outputFileName = "output.txt"

var (
	visited = make(map[string]bool)
	mu      sync.Mutex
	wg      sync.WaitGroup
)

type scrapeDoneMsg string

type Model struct {
	url         string
	concurrency int
	phase       int
	depth       int
	input       string
	message     string
	isScraping  bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func ScrapeCmd(url string, depth int, concurrency int) tea.Cmd {
	return func() tea.Msg {
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return scrapeDoneMsg("Failed to create output file: " + err.Error())
		}
		defer outputFile.Close()

		sem := make(chan struct{}, concurrency)
		defer close(sem)

		start := time.Now()
		wg.Add(1)
		go functions.Scrape(url, depth, outputFile, &wg, &mu, &visited, sem)
		wg.Wait()

		duration := time.Since(start)
		return scrapeDoneMsg(fmt.Sprintf("Scraping completed in %.2f seconds. Output saved to %s", duration.Seconds(), outputFileName))
	}
}

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
			switch m.phase {
			case 0:
				if functions.IsValidURL(m.input) {
					m.url = m.input
					m.phase++
					m.input = ""
					m.message = ""
				} else {
					m.message = "Invalid URL"
				}
			case 1:
				num, err := strconv.Atoi(m.input)
				if err != nil || num <= 0 {
					m.message = "Concurrency must be a positive number"
					break
				}
				m.concurrency = num
				m.input = ""
				m.phase++
				m.message = ""
			case 2:
				num, err := strconv.Atoi(m.input)
				if err != nil || num <= 0 {
					m.message = "Depth must be a positive number"
					break
				}
				m.depth = num
				m.input = ""
				m.phase++
				m.isScraping = true
				m.message = ""
				return m, ScrapeCmd(m.url, m.depth, m.concurrency)
			}
		case tea.KeyRunes:
			m.input += string(msg.Runes)
		}

	case scrapeDoneMsg:
		m.isScraping = false
		m.message = string(msg)
		m.phase++
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() string {
	switch m.phase {
	case 0:
		return fmt.Sprintf("Enter URL: %s\n%s", m.input, m.message)
	case 1:
		return fmt.Sprintf("Enter Max Concurrency level: %s\n%s", m.input, m.message)
	case 2:
		return fmt.Sprintf("Enter Max Depth level: %s\n%s", m.input, m.message)
	case 3:
		if m.isScraping {
			return fmt.Sprintln("Scraping...")
		}
		return fmt.Sprintln(m.message)
	case 4:
		return fmt.Sprintln(m.message)
	default:
		return fmt.Sprintf("Something went wrong %+v", m)
	}
}

func main() {
	p := tea.NewProgram(Model{})
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
