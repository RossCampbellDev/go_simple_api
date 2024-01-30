package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	resp, err := http.PostForm("https://www.datastox.org",
		url.Values{"txtusername": {"ross1"}, "txtpassword": {"kaboom11"}})

	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
}
