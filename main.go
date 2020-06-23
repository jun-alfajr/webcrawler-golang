package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var target string
var urls = make(map[string]bool)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please input URL")
		return
	}

	target = os.Args[1]
}

func RunCrawler(uri string) {
	if url == "" {
		return
	}

	if !urls[uri] {
		urls[uri] = true
	} else {
		return
	}

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
		attr, exists := q.attr("href")
		if exists {
			nextLink := TrimUrl(attr)
		}
	})
}

func TrimUrl(uri string) {
	uri = strings.TrimSuffix(uri, "/")
	validUrl, err := url.Parse(uri)
	if err != nil {
		return ""
	}

}
