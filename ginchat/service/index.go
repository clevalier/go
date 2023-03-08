package service

import "github.com/gin-gonic/gin"

func GetIndex(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "welcome",
	})
}
