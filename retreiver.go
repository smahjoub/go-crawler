package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type link struct {
	source string
	target string
}

func enqueue(uri string, visited *map[string]bool) {

	links, _ := retrieve(uri, visited)

	for _, l := range links {
		if !(*visited)[l.target] {
			enqueue(l.target, visited)
		}
	}
}

func retrieve(uri string, visited *map[string]bool) ([]link, error) {
	resp, err := http.Get(uri)
	(*visited)[uri] = true
	fmt.Println("Fetching:", uri)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	u, _ := url.Parse(uri)
	host := u.Host

	links := []link{}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		lu, err := url.Parse(href)
		if err != nil {
			return
		}
		if isInternalURL(host, lu) {
			links = append(links, link{
				source: uri,
				target: u.ResolveReference(lu).String(),
			})
		}

	})

	return unique(links), nil
}

func isInternalURL(host string, lu *url.URL) bool {

	if lu.IsAbs() {
		return strings.EqualFold(host, lu.Host)
	}
	return len(lu.Host) == 0
}

// insures that there is no repetition
func unique(s []link) []link {
	keys := make(map[link]bool)
	list := []link{}
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
