package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"xkcd/xkcd"
)

const (
	concurrency = 5 // do not abuse api too much
)

var dirPath = flag.String("dir", "data", "Directory path to cache comics records")

func main() {
	flag.Parse()
	err := os.MkdirAll(*dirPath, 0750)
	if err != nil {
		log.Fatalln(err)
	}
	lastID, err := fetchLastPublishedID()
	if err != nil {
		log.Fatalln(err)
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	for id := lastID; id >= 1; id-- {
		comicPath := filepath.Join(*dirPath, fmt.Sprintf("%d.json", id))
		if f, err := os.Open(comicPath); err == nil {
			// already cached
			f.Close()
			continue
		}
		sem <- struct{}{}
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			defer func() { <-sem }()
			err = save(id, comicPath)
			if err != nil {
				log.Println(err)
			}
		}(id)
	}
	wg.Wait()
}

const (
	urlTemplate = "https://xkcd.com/%d/info.0.json"
	urlLast     = "https://xkcd.com/info.0.json"
)

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status_code %v [url=%s]", resp.StatusCode, url)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func fetchLastPublishedID() (int, error) {
	var b []byte
	var comic *xkcd.Comic
	var err error
	if b, err = download(urlLast); err != nil {
		return 0, err
	}
	if comic, err = xkcd.Parse(b); err != nil {
		return 0, err
	}
	return comic.Num, nil
}

func save(id int, comicPath string) error {
	var b []byte
	var err error
	url := fmt.Sprintf(urlTemplate, id)
	if b, err = download(url); err != nil {
		return err
	}
	if _, err = xkcd.Parse(b); err != nil {
		// payload is invalid
		return err
	}
	if err = os.WriteFile(comicPath, b, 0640); err != nil {
		return err
	}
	return nil
}
