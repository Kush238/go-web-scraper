package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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

	commentRegex := regexp.MustCompile(`(\d+)`)

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
		
		scoreSpan := subtext.Find("span.score")
		pointsText := strings.TrimSpace(scoreSpan.Text())
		pointsText = strings.TrimSuffix(pointsText, " points")
		points, err := strconv.Atoi(pointsText)
		if err == nil {
			item.Points = points
		}

		userLink := subtext.Find("a.hnuser")
		item.User = strings.TrimSpace(userLink.Text())

		commentLinks := subtext.Find("a")
		commentLink := commentLinks.Last()
		
		href, exists := commentLink.Attr("href")
		if exists && strings.HasPrefix(href, "item?id=") {
			item.CommentsURL = "https://news.ycombinator.com/" + href
			
			// Get the comment count
			commentText := commentLink.Text()
			if commentText == "discuss" {
				item.CommentsNum = 0
			} else if strings.Contains(commentText, "comment") {
				// Extract just the number from the text
				matches := commentRegex.FindStringSubmatch(commentText)
				if len(matches) > 0 {
					commentsNum, err := strconv.Atoi(matches[0])
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