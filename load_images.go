package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	file, err := os.Open("images.txt")
	handleError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wg.Add(1)
		go saveImage(scanner.Text(), &wg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func saveImage(url string, wg *sync.WaitGroup) {
	response, err := http.Get(url)
	handleError(err)
	defer response.Body.Close()
	fmt.Println("Downloaded images from ", url)

	urlSlice := strings.Split(url, "/")
	os.MkdirAll("images", os.ModePerm)

	file, err := os.Create("images/" + urlSlice[len(urlSlice)-1] + ".jpg")
	handleError(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	handleError(err)

	wg.Done()
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
