package httphandlers

import (
	"golang_test_task/internal/bitcounter"
	"html/template"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	_parseMaxMemory = 100 * 1024 * 1024 // 100 MB
)

type FileInfo struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

type Statistics struct {
	RequestCount int        `json:"requestCount"`
	Files        []FileInfo `json:"files"`
}

type pageHandler struct {
	formFilePath string
	mu           sync.Mutex
	stats        Statistics
}

func (ph *pageHandler) Web(c *gin.Context) {
	t, err := template.ParseFiles(ph.formFilePath)
	if err != nil {
		log.WithError(err).Error("Error parsing template file")
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := t.Execute(c.Writer, nil); err != nil {
		log.WithError(err).Error("Error executing template")
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func (ph *pageHandler) Process(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(_parseMaxMemory); err != nil {
		log.WithError(err).Error("Error parsing multipart form")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fileHeader := c.Request.MultipartForm.File["binary_input_data"]
	if len(fileHeader) == 0 {
		err := "No files found in the request"
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	file, err := fileHeader[0].Open()
	if err != nil {
		log.WithError(err).Error("Error opening file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.WithError(err).Error("Error closing file")
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		log.WithError(err).Error("Error reading file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	log.Debugf("Received data: %+v", data)

	// Обновление статистики
	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.stats.RequestCount++
	ph.stats.Files = append(ph.stats.Files, FileInfo{
		Filename: fileHeader[0].Filename,
		Size:     fileHeader[0].Size,
	})

	result := bitcounter.Process(&bitcounter.Input{Data: data})
	c.JSON(http.StatusOK, result)
}

func (ph *pageHandler) Stat(c *gin.Context) {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	c.JSON(http.StatusOK, ph.stats)
}
