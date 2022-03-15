package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

var (
	//channel and concurrency to fix endless loop
	queue       = make(chan Node)
	visitedLink = make(map[string]bool)
)

type Node struct {
	Depth int
	Url   string
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

func crawl(url string, depth int) {
	fmt.Println(">>>", visitedLink, "---\n")

	visitedLink[url] = true

	fmt.Printf("=======> %v \n", url)

	baseurl := url
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

		node := Node{Url: addHostToPath(str, url), Depth: depth + 1}
		//concurrency asyn to not exhaust all the resources
		go func() { queue <- node }()
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

	limit, err := strconv.Atoi(arguments[1])
	checkError(err)
	n := Node{Url: arguments[0], Depth: 1}
	//concurrency is async to not exhaust all the resources
	go func() {
		queue <- n
	}()
	fmt.Printf("--------%v \n", queue)

	for node := range queue {

		if node.Depth > limit {
			fmt.Println("The depth limit has been reached")
			os.Exit(0)
		}

		if !visitedLink[node.Url] {
			fmt.Println("~~~~~~~VISITED", visitedLink)
			crawl(node.Url, node.Depth)
		}

	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
