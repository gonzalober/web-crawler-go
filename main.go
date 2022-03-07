package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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

func crawl( ) {
	baseurl := "https://github.com/gonzalober"
	resp, err := http.Get(baseurl)
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body) //readAll stores evth in the memory
	checkError(err)

	r := regexp.MustCompile("href=\"([^\"]+)")

	result := r.FindAllStringSubmatch(string(body), -1)

	for _, num := range result {
		str := num[1]

		if isInDenyList(str[len(str)-3:]) {
			continue
		}
		fmt.Println(num[1])
	}

	resp.Body.Close()

}

func main() {
	crawl()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
