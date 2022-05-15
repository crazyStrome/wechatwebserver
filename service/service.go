package service

import (
	"context"
	"wechatwebserver/config"

	"github.com/gin-gonic/gin"
)

// InitServer 初始化服务器，阻塞
func InitServer(ctx context.Context) error {
	r := gin.Default()
	r.GET("/wechat", echo)
	r.POST("/wechat", procMsg)
	return r.Run(config.GetConf().Addr)
}
