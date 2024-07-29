package controllers

import (
	"golang-excel-import/views"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Data struct {
	handler views.Handlers
}

func (d *Data) UploadExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	// Check file extension
	if ext := filepath.Ext(file.Filename); ext != ".xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format. Please upload an .xlsx file"})
		return
	}
	go d.handler.ProcessExcelFile(f)
	c.JSON(http.StatusOK, gin.H{"message": "File is being processed"})

}

func (d *Data) GetRecords(c *gin.Context) {

	if data, err := d.handler.GetAllRecords(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	} else {
		c.JSON(http.StatusOK, data)
	}

}

// TODO write update API
// TODO write Delete API

//TODO write error logs for server,files
