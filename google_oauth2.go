package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)



var admobReadonlyScope = "https://www.googleapis.com/auth/admob.readonly"

var admobReportScope = "https://www.googleapis.com/auth/admob.report"

//获取授权码流程
//1, auth 访问auth接口 浏览器跳转Google授权页面，选择账号授权
//授权完成之后，重定向到redirect_uris，并且会携带code,用此code交换刷新token，和访问token，刷新token长久有效，可保存


//下载的客户端凭证
//redirect_uris 重定向地址写回环地址即可，私有IP授权会失败
var config = `{"installed":{"client_id":"","redirect_uris":["http://127.0.0.1:8080/call"]}}`

func main() {
  //不同api收钱scope不同，需要参考Google文档，此处以admob为示例
	conf, err := google.ConfigFromJSON([]byte(config), admobReadonlyScope, admobReportScope)
	if err != nil {
		panic(err)
	}
	//cli := http.DefaultClient
	app := gin.Default()
	//上面config配置的回调地址，会传递一个code进来
	app.GET("/call", func(c *gin.Context) {
		code := c.Request.URL.Query().Get("code")
		if code == "" {
			c.JSON(200, "fuck !")
			return
		}
		//拿这个code，换一个刷新token, 访问token，刷新token长久有效，如果失效，跑一下这个程序
		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			c.JSON(200, err)
			return
		}

		c.JSON(200, tok)
		return
	})
	//访问这个接口，转到授权登录页面
	app.GET("/auth", func(c *gin.Context) {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		c.Redirect(http.StatusFound, url)
	})
	app.Run(":8080")
}
