package main

import (
	"context"
	"crypto/sha1"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
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
		sig := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")
		token := config.GetConf().Token
		arr := []string{token, timestamp, nonce}
		sort.Strings(arr)
		sha := sha1.Sum([]byte(strings.Join(arr, "")))
		if string(sha[:]) == sig {
			c.String(http.StatusOK, echostr)
		}
	})
	r.Run(config.GetConf().Addr)
}
func wechat(w http.ResponseWriter, r *http.Request) {
	var content, _ = ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.Println(content)
	w.Write([]byte("crazstom"))
}
