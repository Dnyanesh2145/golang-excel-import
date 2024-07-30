package controllers

import (
	"encoding/json"
	"fmt"
	"golang-excel-import/dialects"
	"golang-excel-import/models"
	"golang-excel-import/views"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

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
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	if data, err := d.handler.GetAllRecords(limit, offset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch record"})
		return
	} else {
		c.JSON(http.StatusOK, data)
	}

}

func (d *Data) GetSingleRecord(c *gin.Context) {
	var record models.Records
	recordid, err := strconv.Atoi(c.Param("record_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record id"})
	}
	record.RecordID = recordid

	key := fmt.Sprintf("record:%s", strconv.Itoa(record.RecordID))

	if data, err := dialects.RedisClient.Get(key); err == nil && data != "" {
		if err := json.Unmarshal([]byte(data), &record); err == nil {
			c.JSON(http.StatusOK, record)
			return
		}
	} else {
		if data, err := d.handler.GetSingleRecords(&record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
			return
		} else {
			jsonData, _ := json.Marshal(&record)
			go dialects.RedisClient.SetE(key, string(jsonData), time.Duration(2*time.Minute))
			c.JSON(http.StatusOK, data)
		}
	}

}

func (d *Data) UpdateRecord(c *gin.Context) {
	var record models.Records

	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := d.handler.Update(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record in MySQL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record updated successfully"})

}

func (d *Data) DeleteRecord(c *gin.Context) {
	var record models.Records
	recordid, err := strconv.Atoi(c.Param("record_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record id"})
	}
	record.RecordID = recordid

	if err := d.handler.Delete(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Delete record in MySQL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})

}
