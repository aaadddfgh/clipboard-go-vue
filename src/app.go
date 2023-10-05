package src

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"clipboard-go-vue/src/auth"
	"clipboard-go-vue/src/lib"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var store = struct {
	data string
}{}

func RunServer(NEED_PASSWORD bool, PASSWORD string, port int) {
	r := gin.Default()

	sessionPool := auth.GetInstance()

	r.Static("/assets", "./public/assets")
	r.StaticFile("/index", "./public/index.html")
	r.StaticFile("/clipboard.png", "./public/clipboard.png")
	r.NoRoute(func(c *gin.Context) {
		// 实现内部重定向
		c.Redirect(http.StatusSeeOther, "index")
	})

	// 初始化你的RSACypto和AESCypto
	MyCypto, err := lib.NewRSACrypto() // 初始化你的RSA加密对象
	if err != nil {
		panic("can't init")
	}
	//transpotCypto := lib.NewAESCrypto() // 初始化你的AES加密对象
	PublicKey, err := MyCypto.GetPubKey()
	if err != nil {
		panic("can't init")
	}

	a, _ := MyCypto.Encrypt([]byte("good moring"))
	MyCypto.Decrypt(a)

	store.data = "{}"

	r.GET("/key", func(c *gin.Context) {

		if err != nil {
			c.Abort()
		}

		c.JSON(http.StatusOK, gin.H{
			"key":      PublicKey,
			"password": NEED_PASSWORD,
		})
	})

	r.POST("/auth", func(c *gin.Context) {
		var body map[string]string
		c.ShouldBindBodyWith(&body, binding.JSON)

		pass, exists := body["pass"]
		//pubKey, exists2 := body["pubKey"]
		AesData, exists3 := body["aes"]

		if /*!exists2 ||*/ !exists3 {
			c.JSON(http.StatusOK, gin.H{"ok": false})
			return
		}

		if !exists && NEED_PASSWORD {
			c.JSON(http.StatusOK, gin.H{"ok": false})
			return
		}
		transportCrypto := lib.GetAESFromEncryptedData(AesData, MyCypto)
		if NEED_PASSWORD {
			//DAES( DRSA(DHEX(byte) ) )
			passwordBytes, err := hex.DecodeString(pass)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}
			passwordBytes, err = MyCypto.Decrypt(passwordBytes)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}
			password := transportCrypto.DecryptByte(passwordBytes)
			if string(password) == PASSWORD {

				c.SetCookie("cookie-id", sessionPool.Gen(&transportCrypto), 1000, "/", "", true, true)
				c.JSON(http.StatusOK, gin.H{"ok": true})
			} else {
				c.JSON(http.StatusOK, gin.H{"ok": false})
			}
		} else {

			c.SetCookie("cookie-id", sessionPool.Gen(&transportCrypto), 1000, "/", "", true, true)
			c.JSON(http.StatusOK, gin.H{"ok": true})

		}
	})

	authorized := r.Group("/")

	authorized.Use(auth.NeedCookie())
	{
		r.POST("/content", func(c *gin.Context) {
			var body map[string]string
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			test, exists := body["test"]
			content, exists2 := body["content"]
			cookieId, err := c.Cookie("cookie-id")
			transpotCypto, exists3 := sessionPool.Get(cookieId)

			if !exists || !exists2 || !exists3 || err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			if transpotCypto.Decrypt(test) != "test" {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			store.data = transpotCypto.Decrypt(content)
			c.JSON(http.StatusOK, gin.H{"ok": true})
		})

		r.GET("/content", func(c *gin.Context) {
			cookieId, err := c.Cookie("cookie-id")
			transpotCypto, exists3 := sessionPool.Get(cookieId)
			if !exists3 || err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": transpotCypto.Encrypt(store.data)})
		})
	}

	r.Run(fmt.Sprintf(":%d", port))
}
