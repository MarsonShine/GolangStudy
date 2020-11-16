// gin https://github.com/gin-gonic/gin
package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"gindemo/src/domain"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	r := setGinServer()
	r.Run(":8080")
}

func setGinServer() *gin.Engine {
	r := gin.Default()
	// 如果不使用默认的中间件 则调用
	// r := gin.New()
	// r.LoadHTMLGlob("./src/templates/*")
	// gin.H 是 map[string]interface{} 的一种快捷方式
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/someJSON", someJSON)

	r.GET("/html", htmlRender)
	r.GET("/jsonp", jsonpHandler)

	r.POST("/login", loginHandler)
	r.GET("json", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	r.GET("purejson", purejsonHandler)
	r.POST("/post", queryAndPostFormHandler)

	r.POST("/upload", uploadHandler)
	r.POST("/upload-multi", uploadMultipleHandler)
	// 内部调用其它api
	r.GET("/internal/invoke/otherapi", internalInvokeOtherAPI)
	return r
}

func someJSON(context *gin.Context) {
	// data := map[string]interface{}{
	// 	"lang": "go 语言",
	// 	"tag":  "<br>",
	// }
	// 输出
	// context.AsciiJSON(http.StatusOK, data)
	names := []string{"lena", "austin", "foo"}
	// 返回的响应体会加 `while(1);` 前缀，防止 json 劫持
	context.SecureJSON(200, names)
}

func htmlRender(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}

func jsonpHandler(c *gin.Context) {
	data := map[string]interface{}{
		"foo": "bar",
	}

	// /JSONP?callback=x
	// 将输出：x({\"foo\":\"bar\"})
	c.JSONP(http.StatusOK, data)
}

func loginHandler(c *gin.Context) {
	var form LoginForm
	err := c.ShouldBindWith(&form, binding.Form)
	if err != nil {
		if form.User == "user" && form.Password == "password" {
			c.JSON(200, gin.H{"status": "you are logged in"})
		} else {
			c.JSON(401, gin.H{"status": "unauthorized"})
		}
	}
}

func purejsonHandler(c *gin.Context) {
	c.PureJSON(200, gin.H{
		"html": "<b>Hello, world!</b>",
	})
}

func queryAndPostFormHandler(c *gin.Context) {
	id := c.Query("id")
	page := c.DefaultQuery("page", "0")
	name := c.PostForm("name") // form-data or x-www-form-urlencoded
	message := c.PostForm("message")
	// 用对象接收 json
	body := &domain.NameMessage{}
	if c.BindJSON(&body) == nil {
		fmt.Printf("id: %s; page: %s; json.name: %s; json.message: %s", id, page, body.Name, body.Message)
	} else {
		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	}

}

func uploadHandler(c *gin.Context) {
	// 单文件
	file, _ := c.FormFile("file")
	log.Printf(file.Filename)

	// 上传文件至指定文件
	dst := path.Join("./src/upload", file.Filename)
	// dst := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s,path: %s", err.Error(), dst))
	} else {
		c.String(200, fmt.Sprintf("'%s' uploaded! file path: %s", file.Filename, dst))
	}
}

func uploadMultipleHandler(c *gin.Context) {
	// Multipart form
	form, _ := c.MultipartForm()
	files := form.File["upload"]
	// 动态初始化数组
	msgs := make([]string, len(files))
	for i, file := range files {
		log.Println(file.Filename)
		dst := path.Join("./src/upload", file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			log.Println(fmt.Sprintf("get form err: %s,path: %s", err.Error(), dst))
		} else {
			msgs[i] = fmt.Sprintf("'%s' uploaded! file path: %s", file.Filename, dst)
		}
	}
	c.String(200, strings.Join(msgs, " "))
}

func internalInvokeOtherAPI(c *gin.Context) {
	response, err := http.Get("https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1605522765039&di=702a8d371ae8262291be62e8cdbec321&imgtype=0&src=http%3A%2F%2Fbpic.588ku.com%2Felement_origin_min_pic%2F16%2F09%2F27%2F1457ea17dfc7da9.jpg")
	if err != nil || response.StatusCode != 200 {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="gopher.png"`,
	}
	// 回写一个文件
	c.DataFromReader(200, contentLength, contentType, reader, extraHeaders)
}
