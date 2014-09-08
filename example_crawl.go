package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"skylarbpayne/crawl"
	"strings"

	"code.google.com/p/go.net/html"
)

//helper function for processing
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

type MyFetcher struct {
	m map[string]string
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

func main() {
	var fetcher MyFetcher
	fetcher.m = make(map[string]string)
	urls := []string{"http://golang.org/"}
	crawl.Crawl(urls, 2, fetcher)
}
