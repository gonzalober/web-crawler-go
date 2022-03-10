package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var (
	//channel and concurrency to fix endless loop
	queue       = make(chan UrlRank)
	visitedLink = make(map[string]int)
)

type UrlRank struct {
	url  string
	rank int
}

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println("Missing URL input")
		os.Exit(1)
	}
	// fmt.Println(arguments)
	u, _ := time.ParseDuration("60s")
	ctx, cancel := context.WithTimeout(context.Background(), u)
	defer cancel()

	go func() {
		queue <- UrlRank{url: arguments[0], rank: 0}
	}()

	intChan := make(chan int)
	//concurrency asyn to not exhaust all the resources
	// go func() {

	// queue <- arguments[0]
	go crawl(queue, ctx, intChan)

	// }()

	<-intChan
	// for href := range queue {
	// if !visitedLink[href] {
	// crawl(queue)
	// }

	// }

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isInDenyList(extension string) bool {
	denyList := []string{"css", "svg", "xml", "png", "json"}
	for _, denyListElem := range denyList {
		if denyListElem == extension {
			return true
		}
	}
	return false
}

func crawl(urlChan chan UrlRank, ctx context.Context, intChan chan int) {

	fmt.Println("\n", visitedLink)

	for href := range urlChan {

		select {
		case <-ctx.Done():
			fmt.Println("done")
			intChan <- 1
		default:
			fmt.Println("NOT DONE")
		}

		rank, ok := visitedLink[href.url]
		fmt.Println(rank, ok)
		fmt.Printf("=======> %v \n", href)
		// "https://github.com/gonzalober/about"
		baseurl := href.url
		resp, err := http.Get(baseurl) //get the html element
		checkError(err)

		body, err := ioutil.ReadAll(resp.Body) //readAll stores evthing in  memory
		checkError(err)

		r := regexp.MustCompile("href=\"([^\"]+)")

		//extract anchor tags
		result := r.FindAllStringSubmatch(string(body), -1)

		for _, num := range result {
			str := num[1]

			// if isInDenyList(str[len(str)-3:]) {
			// 	continue
			// }
			// fmt.Println("----->>>> %v \n" + num[1])

			properUrl := addHostToPath(str, href.url)
			//concurrency asyn to not exhaust all the resources
			go func() { urlChan <- UrlRank{properUrl, href.rank + 1} }()
			// crawl(addHostToPath(num[1], href))

		}
		resp.Body.Close()

	}

}

func addHostToPath(path, baseUrl string) string {
	uri, err := url.Parse(path)
	if err != nil {
		return ""
	}
	// fmt.Println("URI HOST", uri)

	base, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}
	// fmt.Println("BASE HOST", base)

	//takes the host from base, and path from uri
	//if it has own host do nothing else base.Host+ uri.Path
	addHostUri := base.ResolveReference(uri)

	return addHostUri.String()
}
