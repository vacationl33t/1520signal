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
	RequestCount int       `json:"requestCount"`
	Files        []FileInfo `json:"files"`
}

type pageHandler struct {
	formFilePath string
	mu           sync.Mutex // Для синхронизации доступа к статистике
	stats        Statistics
}

func (ph *pageHandler) Web(c *gin.Context) {
	t, err := template.ParseFiles(ph.formFilePath)
	if err != nil {
		log.WithError(err).Error("Can't parse files")
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := t.Execute(c.Writer, nil); err != nil {
		log.WithError(err).Error("Can't execute template")
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

	fileHeader := c.Request.MultipartForm.File["binary_input_data"][0]

	file, err := fileHeader.Open()
	if err != nil {
		log.WithError(err).Error("Error opening file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.WithError(err).Error("Error reading file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
		return
	}

	log.Debugf("Received data: %+v", data)

	// Обновляем статистику
	ph.mu.Lock()
	ph.stats.RequestCount++
	ph.stats.Files = append(ph.stats.Files, FileInfo{
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
	})
	ph.mu.Unlock()

	result := bitcounter.Process(&bitcounter.Input{Data: data})
	c.JSON(http.StatusOK, result)
}

func (ph *pageHandler) Stat(c *gin.Context) {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	c.JSON(http.StatusOK, ph.stats)
}
