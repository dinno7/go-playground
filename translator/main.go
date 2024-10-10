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
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	langFrom       *string = flag.String("from", "en", "The source language")
	langTo         *string = flag.String("to", "fa", "The target language")
	inputText      *string = flag.String("text", "", "Text to translate")
	filePath       *string = flag.String("path", "", "File path to translate")
	isSubtitleFile *bool   = flag.Bool(
		"sub",
		false,
		"If your file is a video subtitle, provide this option as true",
	)
)

func init() {
	flag.StringVar(langFrom, "f", "en", "The source language (alternative to --from)")
	flag.StringVar(langTo, "t", "fa", "The target language (alternative to --to)")
	flag.StringVar(inputText, "x", "", "Text to translate (alternative to --text)")
	flag.StringVar(filePath, "p", "", "File path to translate (alternative to --path)")
	flag.BoolVar(
		isSubtitleFile,
		"s",
		false,
		"If your file is a video subtitle, provide this option as true (alternative to --sub)",
	)
}

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
		l, isPrefix, err := fileReader.ReadLine()
		line := strings.TrimSpace(string(l))
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
		if line != "" {
			wg.Add(1)
			go func(text string, lineNum int, ch chan fileTexts) {
				defer wg.Done()
				googleTranslator := GoogleTranslator{
					from: *langFrom,
					to:   *langTo,
					text: text,
				}
				var meaning string

				if *isSubtitleFile {
					// More comprehensive regex for various timestamp formats
					isLineTimecode := regexp.MustCompile(
						`^\d{1,2}:\d{2}:\d{2}[.,]\d{1,3}\s*-->\s*\d{1,2}:\d{2}:\d{2}[.,]\d{1,3}`,
					).MatchString(text)
					// Regex for subtitle numbers (allowing for non-sequential numbering)
					isLineSubtitleNumber := regexp.MustCompile(`^\d+$`).MatchString(text)

					if isLineTimecode || isLineSubtitleNumber {
						meaning = text
					} else {
						meaning, _ = getMeaning(&googleTranslator)
					}
				} else {
					meaning, _ = getMeaning(&googleTranslator)
				}
				ch <- fileTexts{
					line: lineNum,
					text: meaning,
				}
			}(line, lineNum, meaningChannel)
		}

		lineNum++
	}

	// Close channel after end
	go func() {
		wg.Wait()
		close(meaningChannel)
	}()

	meaningMap := make(map[int]string)
	maxLine := 0
	for s := range meaningChannel {
		meaningMap[s.line] = s.text
		if s.line > maxLine {
			maxLine = s.line
		}
	}

	translatedFileName := file.Name()
	translatedFileName = filepath.Base(translatedFileName)

	translatedFileExt := filepath.Ext(translatedFileName)

	translatedFileName = strings.TrimSuffix(translatedFileName, translatedFileExt)
	translatedFileName = fmt.Sprintf("%s_%s%s", translatedFileName, *langTo, translatedFileExt)

	translatedFile, err := os.Create(translatedFileName)
	if err != nil {
		panic(err)
	}
	defer translatedFile.Close()

	for i := 0; i <= maxLine; i++ {
		if meaningText, exist := meaningMap[i]; exist {
			_, err := translatedFile.WriteString(meaningText + "\n")
			if err != nil {
				panic(fmt.Errorf("error writing to file: %w", err))
			}
		} else {
			// Write an empty line if no text exists for this line number
			_, err := translatedFile.WriteString("\n")
			if err != nil {
				panic(fmt.Errorf("error writing to file: %w", err))
			}
		}
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
