package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Initialize Colly collector
	c := colly.NewCollector()

	// Find and print article titles
	c.OnHTML("a.storylink", func(e *colly.HTMLElement) {
		fmt.Println("Title:", e.Text)
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with error:", err)
	})

	// Start scraping
	err := c.Visit("https://news.ycombinator.com/")
	if err != nil {
		log.Fatal(err)
	}
}
