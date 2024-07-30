package views

import (
	"encoding/json"
	"fmt"
	"golang-excel-import/dialects"
	"golang-excel-import/models"
	"log"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm/clause"
)

type Handlers struct{}

func (H *Handlers) ProcessExcelFile(file multipart.File) {
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		log.Println("Failed to read file file:", err)
	}
	sheetname := xlsx.GetSheetName(0)
	if rows, err := xlsx.GetRows(sheetname); err != nil {
		log.Println("Failed to parse Excel file:", err)
	} else {
		var record []models.Records
		// Validate column headers
		if len(rows) < 1 || rows[0][1] != "First Name" || rows[0][2] != "Last Name" || rows[0][3] != "Gender" || rows[0][4] != "Country" || rows[0][5] != "Age" || rows[0][6] != "Date" || rows[0][7] != "Id" {
			log.Println("Invalid column headers")
			return
		}
		for _, row := range rows[1:] {
			id, err := strconv.Atoi(row[7])
			if err != nil {
				log.Println("Invalid ID format:", err)
				continue
			}
			age, err := strconv.Atoi(row[5])
			if err != nil {
				log.Println("Invalid age format:", err)
				continue
			}
			date, err := time.Parse("02/01/2006", row[6]) // Adjust the date format as per your Excel file
			if err != nil {
				log.Println("Invalid date format:", err)
				continue
			}

			reco := models.Records{
				RecordID:  id,
				FirstName: row[1],
				LastName:  row[2],
				Gender:    row[3],
				Country:   row[4],
				Age:       age,
				Date:      date,
			}
			record = append(record, reco)

		}
		go H.StoredData(record)
	}
}

func (H *Handlers) StoredData(records []models.Records) {
	if conn, err := dialects.GetConnection(); err != nil {
		log.Println("Failed to connect DB")
	} else {
		if tx := conn.Debug().Model(models.Records{}).Clauses(
			clause.OnConflict{
				Columns:   []clause.Column{{Name: "record_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "gender", "age", "date", "country"}),
			},
		).Create(&records); tx.Error != nil {
			log.Println("failed DB insert :", err)
		}
	}
	if data, err := dialects.RedisClient.Get("records"); err != nil || data == "" {
		jsonData, _ := json.Marshal(&records)
		go dialects.RedisClient.SetE("records", string(jsonData), time.Duration(5*time.Minute))
	}

}

func (H *Handlers) GetAllRecords(limit int, offset int) ([]models.Records, error) {
	var records []models.Records
	if conn, err := dialects.GetConnection(); err != nil {
		log.Println("Failed to connect DB")
		return nil, err
	} else {
		if tx := conn.Debug().Model(&models.Records{}).Order("record_id ASC").Limit(limit).Offset(offset).Find(&records); tx.Error != nil {
			log.Println("Failed to get all records :", tx.Error)
			return nil, tx.Error
		} else {
			return records, nil
		}
	}
}

func (H *Handlers) GetSingleRecords(record *models.Records) (*models.Records, error) {
	if conn, err := dialects.GetConnection(); err != nil {
		log.Println("Failed to connect DB")
		return nil, err
	} else {
		if tx := conn.Debug().Model(&models.Records{}).Where("record_id =?", record.RecordID).First(&record); tx.Error != nil {
			log.Println("Failed to get all records :", tx.Error)
			return nil, tx.Error
		} else {
			return record, nil
		}
	}
}

func (H *Handlers) Update(record *models.Records) error {
	if conn, err := dialects.GetConnection(); err != nil {
		log.Println("Failed to connect DB")
		return err
	} else {
		if tx := conn.Debug().Model(&models.Records{}).Where("record_id =?", record.RecordID).Updates(&record).First(&record); tx.Error != nil {
			log.Println("Failed to get all records :", tx.Error)
			return tx.Error
		} else {
			key := fmt.Sprintf("record:%s", strconv.Itoa(record.RecordID))
			jsonData, _ := json.Marshal(&record)
			go dialects.RedisClient.SetE(key, string(jsonData), time.Duration(2*time.Minute))
			return nil
		}
	}
}

func (H *Handlers) Delete(record *models.Records) error {
	if conn, err := dialects.GetConnection(); err != nil {
		log.Println("Failed to connect DB")
		return err
	} else {
		condition := map[string]interface{}{
			"deleted_at": time.Now(),
		}
		if tx := conn.Debug().Model(&models.Records{}).Where("record_id=?", record.RecordID).Updates(condition); tx.Error != nil {
			log.Println("Failed to get all records :", tx.Error)
			return tx.Error
		} else {
			key := fmt.Sprintf("record:%s", strconv.Itoa(record.RecordID))
			if _, err := dialects.RedisClient.Get(key); err == nil {
				go dialects.RedisClient.Delete(key)
			}
			return nil
		}
	}
}
