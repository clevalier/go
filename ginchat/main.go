package main

import "ginchat/router"

func main() {
	//官方示例
	//r := gin.Default()
	//r.GET("/ping", func(context *gin.Context) {
	//	context.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	//r.Run()

	r := router.Router()
	r.Run(":8081")

}
