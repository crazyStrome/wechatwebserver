package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func echo(c *gin.Context) {
	logrus.Infof("echo req:%+v", c.Request)
	start := time.Now()
	defer func() {
		logrus.Infof("echo done, req:%v, cost:%v", c.Request, time.Since(start))
	}()
	echostr := c.Query("echostr")
	c.String(http.StatusOK, echostr)
}

func getJSON(val interface{}) string {
	bs, _ := json.Marshal(val)
	return string(bs)
}
