package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var target string
var wg sync.WaitGroup

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please input URL")
		return
	}

	target = os.Args[1]

	queue := make(chan string)
	filteredQueue := make(chan string)
	wg.Add(1)
	go func() { queue <- target }()
	go func() {
		var urls = make(map[string]bool)
		for uri := range queue {
			if !urls[uri] {
				urls[uri] = true
				filteredQueue <- uri
			} else {
				wg.Done()
			}
		}
	}()

	for i := 0; i < 10; i++ {
		go func() {
			for nextLink := range filteredQueue {
				RunCrawler(nextLink, queue)
				wg.Done()
			}
		}()
	}

	wg.Wait()

}

func RunCrawler(uri string, queue chan string) {
	if uri == "" {
		return
	}

	fmt.Println("Fetching", uri)

	response, err := http.Get(uri)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	doc.Find("a").Each(func(i int, q *goquery.Selection) {
		attr, exists := q.Attr("href")
		if exists {
			nextLink := TrimUrl(attr)
			wg.Add(1)
			go func() { queue <- nextLink }()
		}
	})
}

func TrimUrl(uri string) string {
	uri = strings.TrimSuffix(uri, "/")
	validUrl, err := url.Parse(uri)
	if err != nil {
		return ""
	}

	targetUrl, _ := url.Parse(target)

	if strings.Contains(validUrl.String(), targetUrl.Host) {
		return uri
	}

	return ""

}
