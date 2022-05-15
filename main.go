package main

import (
	"context"
	"wechatwebserver/config"
	"wechatwebserver/service"
	"wechatwebserver/token"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	if err := config.InitConf(ctx, "conf.yaml"); err != nil {
		logrus.Fatal(err)
		return
	}
	if err := token.InitAccess(ctx); err != nil {
		logrus.Fatal(err)
		return
	}
	if err := service.InitServer(ctx); err != nil {
		logrus.Fatal(err)
		return
	}
}

func initLogrus() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
