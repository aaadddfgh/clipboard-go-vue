package middleware

import "github.com/gin-gonic/gin"

func CheckSize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.ContentLength > 100*1024*1024 { //增加请求大小控制
			ctx.JSON(200, gin.H{
				"ok":   false,
				"code": 20001,
				"msg":  "too large",
			})
			ctx.Abort()
		} else {
			ctx.Next()
		}

	}
}
