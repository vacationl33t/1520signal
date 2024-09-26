package httphandlers

import (
	"golang_test_task/internal/bitcounter"
	"html/template"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	_parseMaxMemory = 100 * 1024 * 1024
)

type pageHandler struct {
	formFilePath string
}

func (ph pageHandler) Web(c *gin.Context) {
	t, err := template.ParseFiles(ph.formFilePath)
	if err != nil {
		log.WithError(err).Fatal("Can`t parse files")
	}

	if err := t.Execute(c.Writer, nil); err != nil {
		log.WithError(err).Fatal("Can`t execute template")
	}
}

func (ph pageHandler) Process(c *gin.Context) {
	c.Request.ParseMultipartForm(_parseMaxMemory)

	fileHeader := c.Request.MultipartForm.File["binary_input_data"][0]

	file, _ := fileHeader.Open()
	defer file.Close()

	data, _ := io.ReadAll(file)

	log.Debugf("Received data: %+v", data)

	result := bitcounter.Process(&bitcounter.Input{Data: data})

	c.JSON(http.StatusOK, result)
}
