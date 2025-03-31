package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HackerNewsItem struct {
	Rank        int
	Title       string
	URL         string
	Points      int
	User        string
	CommentsURL string
	CommentsNum int
}

func main() {
	url := "https://news.ycombinator.com/"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	response, err := client.Do(req)
	if err != nil {
		log.Fatal("Error fetching the URL:", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", response.StatusCode, response.Status)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body:", err)
	}

	var newsItems []HackerNewsItem

	doc.Find("tr.athing").Each(func(i int, s *goquery.Selection) {
		item := HackerNewsItem{}

		rankText := strings.TrimSpace(s.Find("span.rank").Text())
		rankText = strings.TrimSuffix(rankText, ".")
		rank, err := strconv.Atoi(rankText)
		if err == nil {
			item.Rank = rank
		}

		titleSelection := s.Find("span.titleline > a").First()
		item.Title = strings.TrimSpace(titleSelection.Text())
		item.URL, _ = titleSelection.Attr("href")

		subtext := s.Next().Find("td.subtext")
		
		pointsText := strings.TrimSpace(subtext.Find("span.score").Text())
		pointsText = strings.TrimSuffix(pointsText, " points")
		points, err := strconv.Atoi(pointsText)
		if err == nil {
			item.Points = points
		}

		item.User = strings.TrimSpace(subtext.Find("a.hnuser").Text())

		var commentLink *goquery.Selection
		subtext.Find("a").Each(func(i int, s *goquery.Selection) {
			linkText := s.Text()
			if strings.Contains(linkText, "comment") || linkText == "discuss" {
				commentLink = s
			}
		})
		
		if commentLink != nil {
			commentsHref, exists := commentLink.Attr("href")
			if exists {
				item.CommentsURL = "https://news.ycombinator.com/" + commentsHref
				
				commentsText := strings.TrimSpace(commentLink.Text())
				if commentsText == "discuss" {
					item.CommentsNum = 0
				} else {
					commentsText = strings.TrimSuffix(commentsText, " comments")
					commentsText = strings.TrimSuffix(commentsText, " comment")
					commentsNum, err := strconv.Atoi(commentsText)
					if err == nil {
						item.CommentsNum = commentsNum
					}
				}
			}
		}

		newsItems = append(newsItems, item)
	})

	for _, item := range newsItems {
		fmt.Printf("#%d: %s\n", item.Rank, item.Title)
		fmt.Printf("URL: %s\n", item.URL)
		if item.CommentsURL != "" {
			fmt.Printf("Comments URL: %s\n", item.CommentsURL)
		} else {
			fmt.Printf("Comments URL: N/A\n")
		}
		fmt.Printf("\nPoints: %d | User: %s | Comments: %d\n", item.Points, item.User, item.CommentsNum)
		fmt.Println()
	}
	
	fmt.Printf("Successfully scraped %d items from Hacker News\n", len(newsItems))
}