package routers

import (
	"golang-excel-import/controllers"

	"github.com/gin-gonic/gin"
)

func Endpoints(routers *gin.Engine) {
	var controller controllers.Data

	routers.POST("/upload", controller.UploadExcel)
	routers.GET("/records", controller.GetRecords)
	routers.GET("/records/:record_id", controller.GetSingleRecord)
	routers.PUT("/records/:record_id", controller.UpdateRecord)
	routers.DELETE("/records/:record_id", controller.DeleteRecord)

}
