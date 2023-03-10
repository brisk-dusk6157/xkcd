package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
	"xkcd/xkcd"
)

var stopwords = map[string]bool{"": true, "a": true, "able": true, "about": true, "across": true, "after": true, "all": true, "almost": true, "also": true, "am": true, "among": true, "an": true, "and": true, "any": true, "are": true, "as": true, "at": true, "be": true, "because": true, "been": true, "but": true, "by": true, "can": true, "cannot": true, "could": true, "dear": true, "did": true, "do": true, "does": true, "either": true, "else": true, "ever": true, "every": true, "for": true, "from": true, "get": true, "got": true, "had": true, "has": true, "have": true, "he": true, "her": true, "hers": true, "him": true, "his": true, "how": true, "however": true, "i": true, "if": true, "in": true, "into": true, "is": true, "it": true, "its": true, "just": true, "least": true, "let": true, "like": true, "likely": true, "may": true, "me": true, "might": true, "most": true, "must": true, "my": true, "neither": true, "no": true, "nor": true, "not": true, "of": true, "off": true, "often": true, "on": true, "only": true, "or": true, "other": true, "our": true, "own": true, "rather": true, "said": true, "say": true, "says": true, "she": true, "should": true, "since": true, "so": true, "some": true, "than": true, "that": true, "the": true, "their": true, "them": true, "then": true, "there": true, "these": true, "they": true, "this": true, "tis": true, "to": true, "too": true, "twas": true, "us": true, "wants": true, "was": true, "we": true, "were": true, "what": true, "when": true, "where": true, "which": true, "while": true, "who": true, "whom": true, "why": true, "will": true, "with": true, "would": true, "yet": true, "you": true, "your": true}

var index = make(map[string]map[string]int)

var indexFile = flag.String("out", "index.dat", "Output file to write index data")
var dirPath = flag.String("dir", "data", "Dir containing xkcd json records to index")

func main() {
	flag.Parse()
	err := filepath.Walk(*dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return indexRecord(path)
	})
	if err != nil {
		log.Fatalln(err)
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(*indexFile, data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func indexRecord(path string) error {
	var f *os.File
	var b []byte
	var err error
	if f, err = os.Open(path); err != nil {
		return err
	}
	defer f.Close()
	if b, err = io.ReadAll(f); err != nil {
		return err
	}
	comic, _ := xkcd.Parse(b)
	scanner := bufio.NewScanner(strings.NewReader(comic.SafeTitle + " " + comic.Transcript + " " + comic.Alt))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		var wordBytes []byte
		for _, r := range scanner.Text() {
			if !unicode.IsPunct(r) {
				wordBytes = utf8.AppendRune(wordBytes, r)
			}
		}
		word := strings.ToLower(string(wordBytes))
		if !stopwords[word] {
			if _, ok := index[word]; !ok {
				index[word] = make(map[string]int)
			}
			index[word][strconv.Itoa(comic.Num)]++
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}
