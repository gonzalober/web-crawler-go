package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var (
	//channel and concurrency to fix endless loop
	queue = make(chan string)
)

func isInDenyList(extension string) bool {
	denyList := []string{"css", "svg", "xml", "png", "json"}
	for _, denyListElem := range denyList {
		if denyListElem == extension {
			return true
		}
	}
	return false
}

func crawl(href string) {
	fmt.Printf("=======> %v \n", href)
	// "https://github.com/gonzalober/about"
	baseurl := href
	resp, err := http.Get(baseurl)
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body) //readAll stores evth in the memory
	checkError(err)

	r := regexp.MustCompile("href=\"([^\"]+)")

	result := r.FindAllStringSubmatch(string(body), -1)

	for _, num := range result {
		str := num[1]

		fmt.Printf("----->>>> %v \n", str[len(str)-3:])
		if isInDenyList(str[len(str)-3:]) {
			continue
		}
		// fmt.Println("----->>>> %v \n" + num[1])

		properUrl := addHostToPath(num[1], href)
		go func() { queue <- properUrl }()
		// crawl(addHostToPath(num[1], href))
	}

	resp.Body.Close()

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

func main() {
	arguments := os.Args[1:]

	if len(arguments) == 0 {
		fmt.Println("Missing URL input")
		os.Exit(1)
	}
	// fmt.Println(arguments)

	//concurrency asyn to not exhaust all the resources
	go func() {
		queue <- arguments[0]
	}()
	fmt.Printf("--------%v", queue)

	for href := range queue {
		crawl(href)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
