package counting_semaphore

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

//limit 20 goroutine and file descriptor
// using buffer channel
var token = make(chan struct{}, 20)

func crawl2(url string) []string {
	fmt.Println("URL::", url)
	token <- struct{}{} // if the channel'queue is 20 => lock the token
	// heavy computation
	time.Sleep(time.Second)
	list, err := extract(url)
	<-token //other goroutines done crawl => release the token
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error::", err)
	}
	return list
}

func UseCrawl2() {
	worklist := make(chan []string)
	var n int // number of pending send to worklist
	n++
	go func() {
		worklist <- os.Args[1:]
	}()
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl2(link)
				}(link)
			}
		}
	}
}

// hepler function -- no need to read

func extract(url string) ([]string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s\n", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v\n", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
