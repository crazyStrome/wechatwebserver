package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"wechatwebserver/config"
	"wechatwebserver/token"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	if err := config.InitConf(ctx, "conf.yaml"); err != nil {
		log.Fatal(err)
		return
	}
	if err := token.InitAccess(ctx); err != nil {
		log.Fatal(err)
		return
	}
	r := gin.Default()
	r.GET("/wechat", func(c *gin.Context) {
		echostr := c.Query("echostr")
		c.String(http.StatusOK, echostr)
	})
	r.Run(config.GetConf().Addr)
}

func wechat(w http.ResponseWriter, r *http.Request) {
	var content, _ = ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.Println(content)
	w.Write([]byte("crazstom"))
}
