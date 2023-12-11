package src

import (
	"encoding/hex"
	"io"
	"log"

	"fmt"
	"net/http"

	"clipboard-go-vue/src/auth"
	"clipboard-go-vue/src/controllor"
	"clipboard-go-vue/src/lib"
	"clipboard-go-vue/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var store = struct {
	data     string
	file     []byte
	filename string
}{}

func RunServer(NEED_PASSWORD bool, PASSWORD string, port int, lang string) {
	r := gin.Default()
	logger := log.Default()

	if err := r.SetTrustedProxies([]string{
		"192.168.0.0/16",
		//"fe80::/10",
		//"127.0.0.0/8",
		"::1/128",
	}); err != nil {
		logger.Println(err.Error())
	}

	sessionPool := auth.GetInstance()

	// r.GET("/test", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{"ok": true})

	// 	ctx.JSON(200, gin.H{"ok": false})
	// })

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

	// a, _ := MyCypto.Encrypt([]byte("good moring"))
	// MyCypto.Decrypt(a)

	store.data = "{}" //clipboard quill.js data

	r.Use(middleware.CheckSize())

	// public below
	{
		r.Static("/assets", "./public/assets")
		r.StaticFile("/index.html", "./public/index.html")
		r.StaticFile("/clipboard.png", "./public/clipboard.png")
		// home page
		r.GET("/", func(ctx *gin.Context) { ctx.Redirect(http.StatusSeeOther, "/index.html") })

		r.GET("/key", func(c *gin.Context) {

			if err != nil {
				c.Abort()
			}

			c.JSON(http.StatusOK, gin.H{
				"key":      PublicKey,
				"password": NEED_PASSWORD,
				"lang":     lang,
			})
		})

		r.POST("/auth", func(c *gin.Context) {
			var body map[string]string
			c.ShouldBindBodyWith(&body, binding.JSON)

			pass, exists := body["pass"]
			//pubKey, exists2 := body["pubKey"]
			AesData, exists3 := body["aes"]

			// vaildate data
			if /*!exists2 ||*/ !exists3 {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			if !exists && NEED_PASSWORD {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			//init Crypto
			transportCrypto := lib.GetAESFromEncryptedData(AesData, MyCypto)
			if NEED_PASSWORD {

				passwordBytes, err := hex.DecodeString(pass) //get from hex
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"ok": false})
					return
				}
				passwordBytes, err = MyCypto.Decrypt(passwordBytes) //decrypt
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"ok": false})
					return
				}
				password := transportCrypto.DecryptByte(passwordBytes)

				if string(password) == PASSWORD {
					//success
					c.SetCookie("cookie-id", sessionPool.Gen(&transportCrypto), 1000, "/", "", false, true)
					c.JSON(http.StatusOK, gin.H{"ok": true})
				} else {
					c.JSON(http.StatusOK, gin.H{"ok": false})
				}
			} else {
				//success
				// 只有https可以设置secure=true
				c.SetCookie("cookie-id", sessionPool.Gen(&transportCrypto), 1000, "/", "", false, true)
				c.JSON(http.StatusOK, gin.H{"ok": true})

			}
		})
	}

	authorized := r.Group("/")

	authorized.Use(auth.NeedCookie())

	//need authorize below
	{
		authorized.POST("/content", func(c *gin.Context) {
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
			controllor.BroadcastMsg([]byte(`{"code":100}`))
			c.JSON(http.StatusOK, gin.H{"ok": true})

		})

		authorized.GET("/content", func(c *gin.Context) {
			cookieId, err := c.Cookie("cookie-id")
			transpotCypto, exists3 := sessionPool.Get(cookieId)
			if !exists3 || err != nil {
				c.JSON(http.StatusOK, gin.H{"ok": false})
				return
			}

			c.JSON(http.StatusOK, gin.H{"ok": true, "data": transpotCypto.Encrypt(store.data)})
		})

		authorized.GET("/ws", controllor.CreateWs)

		authorized.GET("/delfile", func(ctx *gin.Context) {
			store.file = []byte("")
			store.filename = ""
		})

		authorized.GET("/file", func(c *gin.Context) {
			//TODO 采用OTP、令牌机制或加密文件内容以保证通信安全
			//cookieId, err := c.Cookie("cookie-id")
			//transpotCypto, exists3 := sessionPool.Get(cookieId)
			//

			c.Writer.WriteHeader(http.StatusOK)

			c.Header("Content-Disposition", "attachment; filename="+store.filename)
			c.Header("Content-Type", "application/text/plain")
			//c.Header("Accept-Length", fmt.Sprintf("%d", len(store.file)))
			c.Writer.Write(store.file)

			// c.Header("Content-Length", fmt.Sprintf("%d", len(store.file)))
			// const bandwidth = 1024 * 512
			// for i := 0; i < len(store.file)/bandwidth-1; i++ {
			// 	c.Writer.Write(store.file[i*bandwidth : (i+1)*bandwidth])
			// 	time.Sleep(1 * time.Second)
			// }
			// c.Writer.Write(store.file[(len(store.file)/bandwidth-1)*bandwidth:])

		})
		authorized.POST("/file", func(c *gin.Context) {
			//TODO 采用OTP、令牌机制或加密文件内容以保证通信安全
			file, header, err := c.Request.FormFile("file")

			store.filename = header.Filename
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"ok": false, "msg": "文件上传失败"})
				return
			}
			check := c.Request.FormValue("key")
			cookieId, err := c.Cookie("cookie-id")
			transpotCypto, exists3 := sessionPool.Get(cookieId)
			if !exists3 || err != nil || transpotCypto.Decrypt(check) != "check" {
				c.JSON(http.StatusBadRequest, gin.H{"ok": false})
				return
			}

			var data []byte
			data, err = io.ReadAll(file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"ok": false, "msg": "文件读取失败"})
				return
			}
			file.Close()
			store.file = data
			c.JSON(http.StatusOK, gin.H{"ok": true, "msg": "上传成功"})

		})
	}

	r.Run(fmt.Sprintf(":%d", port))
	log.Println("server starting on port:" + fmt.Sprint(port))
}
