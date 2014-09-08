package crawl

import "fmt"

const num_fetches = 10

//This defines a type of object that gets urls(and also processes your data)
type Fetcher interface {
	// returns a slice of URLs found on that page. You should also perform your processing here (single loop rather than multiple loops.
	Fetch(url string) (urls []string, err error)
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

			select {
			case signal_done <- 1:
			case nd := <-signal_done:
				signal_done <- nd + 1
			}
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

	num_finished := 0

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
			num_finished += fin
		default:
			if num_running <= 0 {
				still_running = false
			}
		}
	}

	fmt.Println("Crawled ", num_finished, " pages!")
}
