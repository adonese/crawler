package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var worklist = make(chan []string)

func main() {
	var n int

	go func() {
		worklist <- os.Args[1:]
	}()
	n++

	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] && strings.HasPrefix(link, os.Args[1]) {
				seen[link] = true
				n++
				go func(link string) {
					links := crawl(link)
					fmt.Printf("Links in: %s are: %v\n\n", link, links)
					worklist <- links
				}(link)
			}
		}
	}
}

//crawl a helper function to fasciliate concurrent access to extract
func crawl(url string) []string {
	var tokens = make(chan struct{}, 20)
	tokens <- struct{}{}
	list := extract(url)
	<-tokens
	return list
}

//extract extracts links of provided URL
func extract(domain string) []string {
	var links []string

	res, err := http.Get(domain)
	if err != nil {
		return nil
	}
	defer res.Body.Close()
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					url, err := res.Request.URL.Parse(a.Val)
					if err != nil {
						continue
					}
					links = append(links, url.String())

				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}
