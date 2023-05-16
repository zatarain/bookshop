package controllers

import (
	"bytes"
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
	server := gin.Default()
	server.HEAD("/health", HealthCheck)
	recorder := httptest.NewRecorder()
	var capture bytes.Buffer
	log.SetOutput(&capture)

	request, exception := http.NewRequest(http.MethodHead, "/health", nil)
	assert.Nil(exception)

	// Perform the request
	server.ServeHTTP(recorder, request)

	// Check to see if the response was what you expected
	assert.Equal(http.StatusOK, recorder.Code)

	log.SetOutput(os.Stderr)
	actual := capture.String()
	expected := "OK, go!"
	assert.Contains(expected, actual)
}
