package main

//Recurisve Web Crawling
import (
	"flag"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/jackdanger/collectlinks"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	visitedPage   = make(map[string]bool)
	visitedLinks  = make(map[string]bool)
	BrokenPage    = make(map[string]string)
	UrlCrawlCount = 0
	Linkcount     = 0
	brokenLinks   = 0
	baseURL       = ""
)

func main() {
	flag.Parse()
	args := flag.Args()
	//fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}
	baseURL = args[0]

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = fmt.Sprintf("crawling %s, please wait ", baseURL)
	s.Start()
	enqueue(baseURL)

	tmpCount := 0
	s.Stop()
	fmt.Println("================================================================================================")
	fmt.Println("================================================================================================")
	fmt.Println("Broken Links:", brokenLinks, "Ok Links:", Linkcount, "Web Pages Crawled:", UrlCrawlCount)
	for key, value := range BrokenPage {
		tmpCount++
		fmt.Println(fmt.Sprintf("[%v] \n broken  : %s \n source: %s", tmpCount, key, value))
	}
}

// fixUrl converts all relative links to absolute links
func fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func enqueue(uri string) {
	//fmt.Println("fetching", uri)
	visitedPage[uri] = true
	resp, err := http.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		absolute := fixUrl(link, uri)
		if !strings.Contains(absolute, baseURL) && !visitedLinks[absolute] {
			visitedLinks[absolute] = true
			checkWebStatus(absolute, uri)
		}
		if strings.Contains(absolute, baseURL) && !visitedPage[absolute] {
			UrlCrawlCount++
			//fmt.Println(absolute)
			checkWebStatus(absolute, uri)
			enqueue(absolute)
		}
	}
}

// checkWebStatus checks all given links if they are invalid
func checkWebStatus(urlParams string, baseline string) {
	resp, _ := http.Get(urlParams)
	if resp != nil && resp.StatusCode == 200 {
		Linkcount++
	} else {
		brokenLinks++
		BrokenPage[urlParams] = baseline
	}
}

//TODO - Redis
