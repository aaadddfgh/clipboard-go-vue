package auth

import (
	"github.com/gin-gonic/gin"
)

func NeedCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if v, ok := ctx.Cookie("cookie-id"); ok == nil {
			if GetInstance().Exsist(v) {
				ctx.Next()
				return
			}
		}
		ctx.JSON(200, gin.H{"ok": false, "msg": "no cookies"})
		ctx.Abort()

	}
}
