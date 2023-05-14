package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMain(test *testing.T) {
	var capture bytes.Buffer
	log.SetOutput(&capture)
	main()
	log.SetOutput(os.Stderr)
	actual := capture.String()
	expected := "OK, go!"
	if !strings.Contains(actual, expected) {
		test.Errorf("Incorrect output, expected '%s' got '%s'", expected, actual)
	}
}
