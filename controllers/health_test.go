package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	var capture bytes.Buffer
	log.SetOutput(&capture)
	main()
	request, exception := http.NewRequest(http.MethodHead, "/", nil)
	assert.ErrorIs(exception, nil, "Check we don't get an error on the request.")

	// Perform the request
	server := gin.Default()
	server.ServeHTTP(recorder, request)

	// Check to see if the response was what you expected
	assert.Equal(recorder.Code, http.StatusOK, "Checking expected status.")

	log.SetOutput(os.Stderr)
	actual := capture.String()
	expected := "OK, go!"
	assert.Contains(actual, expected, fmt.Sprintf("Incorrect output, expected '%s' got '%s'", expected, actual))
}
