package main

import (
	"bottlemail/controllers"

	"github.com/gin-gonic/gin"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("./views/*.html")
	router.Static("/assets", "./assets")
	router.GET("/", controllers.Index)
	router.GET("/messages.html", controllers.Messages)
	router.POST("./send", controllers.Send)
	router.Run(":8080")
}
