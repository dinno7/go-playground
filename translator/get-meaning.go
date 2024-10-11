package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type API interface {
	GetAPIUrl() string
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

func translateWithChannel(text string, lineNum int, ch chan fileTexts, wg *sync.WaitGroup) {
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
}

func getDirectTextMeaning(from, to, text string) string {
	googleTranslator := GoogleTranslator{
		from: from,
		to:   to,
		text: text,
	}
	meaning, err := getMeaning(&googleTranslator)
	if err != nil {
		panic(err)
	}

	return meaning
}

func getTranslatedFileDest(file *os.File) string {
	translatedFileName := file.Name()
	translatedFileName = filepath.Base(translatedFileName)

	absPath, err := filepath.Abs(*filePath)
	if err != nil {
		absPath = "."
	}
	parentDir := filepath.Dir(absPath)

	translatedFileExt := filepath.Ext(translatedFileName)

	translatedFileName = strings.TrimSuffix(translatedFileName, translatedFileExt)
	translatedFileName = fmt.Sprintf("%s_%s%s", translatedFileName, *langTo, translatedFileExt)
	return fmt.Sprintf("%s/%s", parentDir, translatedFileName)
}
