# 🌐 Go CLI Web Scraper

A high-performance, concurrency-safe web scraper built in **Go**.  
This CLI tool accepts a URL, maximum recursion depth, and concurrency level, then scrapes all reachable pages—saving the content to a file and reporting the total scrape time.

---

## ✨ Features

- ✅ **CLI interface**: Simple and configurable
- 🕸️ **Recursive scraping** of nested links with user-defined depth
- ⚙️ **Concurrency control** via goroutines & semaphores (channels)
- 📁 **Saves scraped data** to a standalone file
- ⏱️ **Reports total time** taken for the scrape

---

## 🛠 Requirements

- Go 1.18 or later installed: [https://go.dev/dl/](https://go.dev/dl/)

---

## 🚀 Installation

Clone the repository:

```bash
git clone https://github.com/sinore69/webscraper.git
cd webscraper
go get
go run main.go 
```
