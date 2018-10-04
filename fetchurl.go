package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	mainReadURLs()
}

// TODO Reshape this to pass URLs to pool of workers,
// which will fetch, analyze, and render results.
func mainReadURLs() {
	ch := make(chan string)
	i := 0
	go func() {
		if err := readURLs(os.Stdin, ch); err != nil {
			log.Fatal("error reading URLs:", err)
		}
	}()

	for url := range ch {
		i++
		fmt.Printf("%d: %s\n", i, url)
	}
}

func readURLs(in io.Reader, urlCh chan<- string) error {
	defer close(urlCh)
	buf := bufio.NewScanner(in)
	for buf.Scan() {
		if err := buf.Err(); err != nil {
			return err
		}
		urlCh <- buf.Text()
	}
	return nil
}

func mainFetch() {
	for _, imgurl := range os.Args[1:] {
		err := fetch(imgurl, os.Stdout)
		if err != nil {
			exit(err)
		}
	}
}

func fetch(url string, out io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(out, resp.Body)
	return nil
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", path.Base(os.Args[0]), err)
	os.Exit(1)
}
