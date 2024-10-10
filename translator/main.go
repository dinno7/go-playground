package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

var (
	langFrom  *string = flag.String("from", "en", "The source language")
	langTo    *string = flag.String("to", "fa", "The target language")
	inputText *string = flag.String("text", "", "Text to translate")
	filePath  *string = flag.String("path", "", "File path to translate")
)

type API interface {
	GetAPIUrl() string
}
type GoogleTranslator struct {
	from string
	to   string
	text string // The input text
}

func (gt *GoogleTranslator) GetAPIUrl() string {
	URL := url.URL{
		Scheme:   "https",
		Host:     "translate.googleapis.com",
		Path:     "/translate_a/single",
		RawQuery: "client=gtx&dt=t",
	}
	q := URL.Query()
	q.Add("sl", gt.from)
	q.Add("tl", gt.to)
	q.Add("q", gt.text)
	URL.RawQuery = q.Encode()

	return URL.String()
}

func main() {
	// flag.StringVar(text, "f", "", "The source language (alternative to -from)")
	// flag.StringVar(text, "t", "", "The target language (alternative to -to)")
	// flag.StringVar(text, "t", "", "Text to translate (alternative to -text)")
	// flag.StringVar(text, "p", "", "File path to translate (alternative to -path)")

	flag.Parse()
	isDirectTextInput := len(*inputText) > 0
	// User provide direct text to translate
	if isDirectTextInput {
		googleTranslator := GoogleTranslator{
			from: *langFrom,
			to:   *langTo,
			text: *inputText,
		}
		meaning, err := getMeaning(&googleTranslator)
		if err != nil {
			panic(err)
		}

		fmt.Println("💀 > ", meaning)
		return
	}
	// ---------------------------
	if len(*filePath) <= 0 {
		panic("Please provide file path or direct text to translate")
	}
	// User provide the file path
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	type fileTexts struct {
		line int
		text string
	}

	meaningChannel := make(chan fileTexts)
	var wg sync.WaitGroup

	fileReader := bufio.NewReader(file)
	lineNum := 0
	for {
		line, isPrefix, err := fileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}
		if isPrefix {
			fmt.Println("Line too long")
			break
		}

		// Do processes
		if len(string(line)) > 1 {
			text := string(line)
			wg.Add(1)
			go func(text string, lineNum int, ch chan fileTexts) {
				defer wg.Done()
				googleTranslator := GoogleTranslator{
					from: *langFrom,
					to:   *langTo,
					text: text,
				}
				var meaning string
				meaning, _ = getMeaning(&googleTranslator)

				ch <- fileTexts{
					line: lineNum,
					text: meaning,
				}
			}(text, lineNum, meaningChannel)
		}

		lineNum++
	}

	// Close channel after end
	go func() {
		wg.Wait()
		close(meaningChannel)
	}()

	for s := range meaningChannel {
		fmt.Println("💀 s > ", s)
	}
}

func getMeaning(api API) (string, error) {
	extractMeaning := func(data string) string {
		meaning := strings.Split(data, ",")[0]
		meaning = strings.Trim(meaning, " \"[")
		return meaning
	}

	res, err := http.Get(api.GetAPIUrl())
	if err != nil {
		return "", errors.New("💀 Something went wrong while translating from API")
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return extractMeaning(string(b)), nil
}
