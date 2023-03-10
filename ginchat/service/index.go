package service

import "github.com/gin-gonic/gin"

// GetIndex
// @Tags         首页
// @Success      200  {string}   welcome
// @Router       /index [get]
func GetIndex(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "welcome",
	})
}
