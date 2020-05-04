package main

import (
	_ "Tbed/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)



func main() {

	r := gin.New()

	// middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(Auth)

	// staticFs
	r.StaticFS("/static",http.Dir("./static"))
	r.StaticFS("/img/",http.Dir(viper.GetString("img_file")))
	r.LoadHTMLGlob("views/*")

	// router
	r.GET("/", Home)
	r.POST("/upload",Upload)

	// run
	r.Run(viper.GetString("addr"))
}

func Home(c *gin.Context)  {

	_,err := c.Cookie("auth")
	password := c.Request.FormValue("p")
	if err != nil && password != viper.GetString("password") {
		c.HTML(http.StatusOK,"403.html",nil)
		return
	}
	c.SetCookie("auth", "true", 60*60*24*60, "/", "", false, true)
	c.HTML(http.StatusOK,"index.html",nil)
}

func Upload(c *gin.Context)  {

	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	fileName := fmt.Sprintf("%d-%s",time.Now().Unix(),header.Filename)
	dst := fmt.Sprintf("%s/%s",viper.Get("img_file"),fileName)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	url := fmt.Sprintf("http://%s/img/%s",c.Request.Host,fileName)
	c.JSON(http.StatusOK,gin.H{"url":url})
}

func Auth(c *gin.Context)  {

	if c.Request.Method == "GET" {
		c.Next()
		return
	}

	if _,err := c.Cookie("auth");err != nil {
		log.Println(err)
		c.JSON(http.StatusForbidden,nil)
	}
}
