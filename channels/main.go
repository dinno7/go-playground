package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	links := []string{
		"http://google.com",
		"http://youtube.com",
		"https://go.dev/",
		"https://nodejs.org/",
	}

	ch := make(chan string)
	for _, link := range links {
		go checkLink(link, ch)
	}

	fmt.Println("💀 ch > ", <-ch)
	fmt.Println("💀 ch > ", <-ch)
	fmt.Println("💀 ch > ", <-ch)
	fmt.Println("💀 ch > ", <-ch)
}

func checkLink(link string, ch chan string) {
	client := http.Client{
		Timeout: time.Duration(time.Second * 3),
	}
	_, err := client.Get(link)
	if err != nil {
		// fmt.Printf("💀 The %s might be down > %s\n", link, err)
		ch <- fmt.Sprintf("💀 The %s might be down > %s\n", link, err)
		return

	}
	// fmt.Printf("💀 The %s is up\n", link)
	ch <- fmt.Sprintf("💀 The %s is up\n", link)
}
