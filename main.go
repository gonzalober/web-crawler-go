package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	baseurl := "https://github.com"
	resp, err := http.Get(baseurl)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body) //readAll stores evth in the memory

	fmt.Println(string(body))

	resp.Body.Close()
}
