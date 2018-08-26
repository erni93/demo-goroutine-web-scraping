package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/iconv.v1"
)

type dictionary struct {
	Words []word `xml:"words>word"`
}

type word struct {
	Value    string    `xml:"value"`
	Examples []example `xml:"examples>example"`
}

type example struct {
	English string `xml:"english"`
	Spanish string `xml:"spanish"`
}

const (
	wordsFilename   = "words_alpha.txt"
	xmlSaveFilename = "./dictionary.xml"
	baseURL         = "https://www.linguee.es/ingles-espanol/search?query="
)

var (
	urlRegex = regexp.MustCompile(`[^ ]+\.[^ ]+`)
)

func main() {
	var dictionary dictionary
	var sem = make(chan struct{}, runtime.NumCPU())
	in := make(chan string)
	out := make(chan word)
	var wg sync.WaitGroup
	fileWords := getFileWords()
	// fileWords := []string{"hello", "car", "chair", "anoisgabnoiwga"} // Uncomment to test with 4 words
	for _, fileWord := range fileWords {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			getXMLWord(in, out)
			<-sem
		}()
		go appendToDictionary(&dictionary, out)
		in <- fileWord
	}
	wg.Wait()
	close(in)
	close(out)
	saveXML(dictionary)
}

func getFileWords() []string {
	file, err := os.Open(wordsFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}

func getXMLWord(in <-chan string, out chan<- word) {
	start := time.Now()

	fileWord := <-in
	res, err := http.Get(baseURL + fileWord)
	if err != nil {
		fmt.Printf("%s error: %v\n", fileWord, err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("%s status code error: %v\n", fileWord, res.StatusCode)
		return
	}

	cd, err := iconv.Open("utf-8", "ISO-8859-15")
	if err != nil {
		fmt.Printf("%s error: %v\n", fileWord, err)
		return
	}
	defer cd.Close()

	utfBody := iconv.NewReader(cd, res.Body, 0)
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		fmt.Printf("%s error: %v\n", fileWord, err)
		return
	}
	xmlWord := word{Value: fileWord}
	doc.Find("tbody.examples>tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		xmlWord.Examples = append(xmlWord.Examples, example{English: formatText(cells.First().Text()), Spanish: formatText(cells.Last().Text())})
	})
	fmt.Printf("Word %s finish in %.2f \n", fileWord, time.Since(start).Seconds())
	out <- xmlWord
}

func appendToDictionary(dictionary *dictionary, out <-chan word) {
	dictionary.Words = append(dictionary.Words, <-out)
}

func formatText(cellText string) string {
	text := strings.TrimSpace(cellText)
	text = strings.Replace(text, "  ", " ", -1)
	text = strings.Replace(text, "[...]", "", -1)
	text = strings.Replace(text, "\n", "", -1)
	return urlRegex.ReplaceAllString(text, "")
}

func saveXML(dictionary dictionary) bool {
	output, err := xml.MarshalIndent(dictionary, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return false
	}
	err = ioutil.WriteFile(xmlSaveFilename, output, 1024)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return false
	}
	return true
}
