package test

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Secrets(c *gin.Context) {
	// 获取提交的用户名（AuthUserKey）
	user := c.MustGet(gin.AuthUserKey).(string)
	if secret, ok := secrets[user]; ok {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		return
	}
}

var secrets = gin.H{
	"misasky":    gin.H{"email": "misasky@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}