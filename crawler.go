package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"code.google.com/p/go.net/html"
)

const num_fetches = 100

func get_url_base(url string) string {
	count := 0
	for i, char := range url {
		if char == '/' {
			count++
		}

		if count == 3 {
			return url[:i]
		}
	}
	return url
}

//This defines a type of object that gets urls(and also processes your data)
type Fetcher interface {
	// returns a slice of URLs found on that page. You should also perform your processing here (single loop rather than multiple loops.
	Fetch(url string) (urls []string, err error)
}

type MyFetcher struct {
	m map[string]string
}

func Crawl(urls []string, depth int, fetcher Fetcher) {
	if depth <= 0 {
		return
	}

	type Search struct {
		pages []string
		depth int
	}

	signal_done := make(chan int, num_fetches*depth)
	new_pages := make(chan Search, num_fetches*depth)
	buffer := make(chan int, num_fetches)
	visited_urls := make(map[string]bool)
	num_errors := 0

	var crawler func(url string, depth int, fetcher Fetcher)
	crawler = func(url string, depth int, fetcher Fetcher) {
		defer func() {
			fmt.Println("Done")
			signal_done <- 1
		}()

		if depth < 0 {
			return
		}

		fmt.Println("found: ", url)

		buffer <- 1
		urls, err := fetcher.Fetch(url)
		<-buffer
		if err != nil {
			num_errors += 1
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Waiting to send new urls")

		new_pages <- Search{urls, depth - 1}
	}

	new_pages <- Search{urls, depth}

	still_running := true
	num_running := 0
	for still_running {
		select {
		case search := <-new_pages:
			for _, page := range search.pages {
				if !visited_urls[page] {
					visited_urls[page] = true
					num_running += 1
					fmt.Println(num_running)
					go crawler(page, search.depth, fetcher)
				}
			}
		case fin := <-signal_done:
			num_running -= fin
		default:
			if num_running <= 0 {
				still_running = false
			}
		}
	}
}

func main() {
	var fetcher MyFetcher
	fetcher.m = make(map[string]string)
	urls := []string{"http://golang.org/"}
	Crawl(urls, 2, fetcher)
}

func (fetcher MyFetcher) Fetch(url string) (urls []string, err error) {
	//GET HTTP request.
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println()
		fmt.Println(url)
		fmt.Println(err)
		fmt.Println()
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println()
		fmt.Println(url)
		fmt.Println(err)
		fmt.Println()
		return
	}

	body := string(data)
	fetcher.m[url] = body
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		fmt.Println()
		fmt.Println(url)
		fmt.Println(err)
		fmt.Println()
		return
	}

	url_base := get_url_base(url)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					u := a.Val
					if u == ".." || len(u) == 0 {
						continue
					} else if len(u) < 4 || u[0:4] != "http" {
						if u[0] == '/' {
							u = url_base + u
						} else {
							u = url + u
						}
					} else if len(u) > 15 && u[7:16] == "localhost" {
						continue
					}
					urls = append(urls, u)
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return
}
