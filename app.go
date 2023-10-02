package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	// 导入你的自定义包
	// 例如：
	// "your-custom-package/cypto"
	// "your-custom-package/store"
)

var store = struct {
	data string
}{}

func runServer(NEED_PASSWORD bool, PASSWORD string, port int) {
	r := gin.Default()

	r.Static("/assets", "./public/assets")
	r.StaticFile("/index", "./public/index.html")
	r.StaticFile("/clipboard.png", "./public/clipboard.png")
	r.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.Redirect(http.StatusSeeOther, "index")
	})

	// 初始化你的RSACypto和AESCypto
	MyCypto, err := NewRSACrypto() // 初始化你的RSA加密对象
	if err != nil {
		panic("can't init")
	}
	transpotCypto := NewAESCrypto() // 初始化你的AES加密对象
	PublicKey, err := MyCypto.GetPubKey()
	if err != nil {
		panic("can't init")
	}

	store.data = transpotCypto.Encrypt("test")

	r.POST("/auth", func(c *gin.Context) {
		var body map[string]string
		c.ShouldBindBodyWith(&body, binding.JSON)

		pass, exists := body["pass"]
		pubKey, exists2 := body["pubKey"]

		if !exists2 {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false})
			return
		}

		if !exists && NEED_PASSWORD {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false})
			return
		}

		if NEED_PASSWORD {

			passwordBytes, err := getPasswordFromHexString(pass)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}
			password, err := MyCypto.Decrypt(passwordBytes)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			if string(password) == PASSWORD {
				data, _ := encryptKeyForConnection(pubKey, transpotCypto.GetKey())

				c.JSON(http.StatusOK, gin.H{"ok": true, "key": data})
			} else {
				c.JSON(http.StatusOK, gin.H{"ok": false})
			}
		} else {
			data, _ := encryptKeyForConnection(pubKey, transpotCypto.GetKey())

			c.JSON(http.StatusOK, gin.H{"ok": true, "key": data})
		}
	})

	r.POST("/content", func(c *gin.Context) {
		var body map[string]string
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false})
			return
		}

		test, exists := body["test"]
		content, exists2 := body["content"]

		if !exists || !exists2 {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false})
			return
		}

		if transpotCypto.Decrypt(test) != "test" {
			c.JSON(http.StatusOK, gin.H{"ok": false})
			return
		}

		store.data = content
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	r.GET("/content", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": store.data})
	})

	r.GET("/key", func(c *gin.Context) {

		if err != nil {
			c.Abort()
		}

		c.JSON(http.StatusOK, gin.H{
			"key":      PublicKey,
			"password": NEED_PASSWORD,
		})
	})

	r.Run(fmt.Sprintf(":%d", port))
}
