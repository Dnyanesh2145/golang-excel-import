package main

import (
	"fmt"
	"golang-excel-import/dialects"
	"golang-excel-import/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

	if _, err := dialects.GetConnection(); err != nil {
		log.Panic(fmt.Printf("error connectin to DB: %s", err))
	} else {
		dialects.CheckModals()
	}
	redis := dialects.RedisClient
	go redis.Connect()
	routers.Endpoints(app)

	app.Run(":8080")
}
