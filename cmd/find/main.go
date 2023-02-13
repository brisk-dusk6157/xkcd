package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"xkcd/xkcd"
)

var indexFile = flag.String("index", "index.dat", "File with index data")
var dirPath = flag.String("dir", "data", "Dir containing xkcd json records to index")

func main() {
	flag.Parse()
	searchTerm := flag.Arg(0)
	index := loadIndex(*indexFile)
	hits, ok := index[searchTerm]
	if !ok {
		fmt.Printf("The idea of \"%s\" is foreign to xkcd\n", searchTerm)
		return
	}
	for xkcdId := range hits {
		var f *os.File
		var b []byte
		var err error
		if f, err = os.Open(filepath.Join(*dirPath, xkcdId+".json")); err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()
		if b, err = io.ReadAll(f); err != nil {
			log.Println(err)
			continue
		}
		comic, _ := xkcd.Parse(b)
		fmt.Println(prettyPrint(comic))
		fmt.Println("---------------------------")
	}
}

func loadIndex(path string) (index map[string]map[string]int) {
	var err error
	var f *os.File
	var b []byte
	if f, err = os.Open(path); err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	if b, err = io.ReadAll(f); err != nil {
		log.Fatalln(err)
	}
	if err = json.Unmarshal(b, &index); err != nil {
		log.Fatalln(err)
	}
	return
}

func prettyPrint(comic *xkcd.Comic) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%s (%s) | https://xkcd.com/%d/ | %s\n", comic.SafeTitle, comic.Year, comic.Num, comic.Img)
	fmt.Fprintf(&builder, comic.Alt)
	return builder.String()
}
