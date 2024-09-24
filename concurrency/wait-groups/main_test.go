package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	wg.Add(1)

	updateMessage("alpha")

	wg.Wait()

	if msg != "alpha" {
		t.Error("Expected to find alpha, but it is not there")
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	msg = "alpha"
	printMessage()

	_ = w.Close()

	res, _ := io.ReadAll(r)
	output := string(res)

	os.Stdout = stdOut

	if !strings.Contains(output, "alpha") {
		t.Error("Expected to find alpha, but it not there")
	}
}
