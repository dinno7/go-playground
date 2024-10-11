package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
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

type fileTexts struct {
	line int
	text string
}

func main() {
	flag.Parse()

	isDirectTextInput := len(*inputText) > 0
	// User provide direct text to translate
	if isDirectTextInput {
		fmt.Println("💀 > ", getDirectTextMeaning(*langFrom, *langTo, *inputText))
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

	meaningChannel := make(chan fileTexts)
	wg := &sync.WaitGroup{}
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
			go translateWithChannel(line, lineNum, meaningChannel, wg)
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

	translatedFile, err := os.Create(getTranslatedFileDest(file))
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
