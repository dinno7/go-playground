package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	var text string
	if isDirectTextInput {
		text = *inputText
	} else if len(*filePath) > 0 {
		fileText, err := os.ReadFile(*filePath)
		if err != nil {
			panic(err)
		}
		text = string(fileText)
	}

	// User provide direct text to translate
	if isDirectTextInput {
		googleTranslator := GoogleTranslator{
			from: *langFrom,
			to:   *langTo,
			text: text,
		}
		meaning, err := getMeaning(&googleTranslator)
		if err != nil {
			panic(err)
		}

		fmt.Println("💀 > ", meaning)
		return
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
