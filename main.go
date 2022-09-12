package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.POST("/distribution", func(ctx *gin.Context) {
		ctx.JSON(200, "Order distributed")
	})
	router.Run(":8081")
}
