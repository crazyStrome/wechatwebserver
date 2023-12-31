package main

import (
	"context"
	"flag"
	"wechatwebserver/client"
	"wechatwebserver/config"
	"wechatwebserver/service"
	"wechatwebserver/token"

	"github.com/sirupsen/logrus"
)


func main() {
	initLogrus()
	ctx := context.Background()
	if err := config.InitConf(ctx, "conf.yaml"); err != nil {
		logrus.Fatal(err)
		return
	}
	if err := token.InitAccess(ctx); err != nil {
		logrus.Fatal(err)
		return
	}
	if err := client.Init(); err != nil {
		logrus.Fatal(err)
		return
	}
	client.Test(ctx)
	if err := service.InitServer(ctx); err != nil {
		logrus.Fatal(err)
		return
	}

}

var logLevel string

func initLogrus() {
    flag.StringVar(&logLevel, "log_level", "debug", "default is debug,other is info error")
    flag.Parse()

	if logLevel == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if logLevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
    } else {
        logrus.SetLevel(logrus.DebugLevel)
    }
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
}
