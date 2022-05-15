package token

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
	"wechatwebserver/config"

	"github.com/sirupsen/logrus"
)

var global atomic.Value

// InitAccess 初始化并更新 accesstoken
func InitAccess(ctx context.Context) error {
	logrus.Infof("start init accesstoken")
	token, err := getToken(ctx)
	if err != nil {
		return err
	}
	global.Store(token)
	go updateToken(ctx)
	logrus.Infof("end init accesstoken:%v", GetToken())
	return nil
}

func getToken(ctx context.Context) (string, error) {
	conf := config.GetConf()
	cli := http.Client{
		Timeout: time.Duration(conf.AccessToken.TimeoutInSec) * time.Second,
	}
	url := fmt.Sprintf("%vappid=%v&secret=%v", conf.AccessToken.URL, conf.AppID, conf.Secret)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("get token new req err:%v", err)
	}
	rsp, err := cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("get token GET err:%v", err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("get token read all err:%v", err)
	}
	type Token struct {
		AccessToken string `json:"access_token"`
		ExpireInSec int32  `json:"expires_in"`
		ErrCode     int32  `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	token := &Token{}
	if err := json.Unmarshal(body, token); err != nil {
		return "", fmt.Errorf("unmarshal token err:%v, body:%s", err, body)
	}
	logrus.Infof("access_token is:%+v", token)
	if token.ErrCode != 0 {
		return "", fmt.Errorf("token.code:%v, msg:%v", token.ErrCode, token.ErrMsg)
	}
	return token.AccessToken, nil
}

func updateToken(ctx context.Context) {
	log.Printf("start updateToken\n")
	conf := config.GetConf()
	tick := time.NewTicker(time.Duration(conf.AccessToken.DurationInSec) * time.Second)
	for {
		select {
		case <-ctx.Done():
			logrus.Infof("updateToken ctx done")
		case <-tick.C:
			for i := 0; i < int(conf.AccessToken.Retries); i++ {
				token, err := getToken(ctx)
				if err != nil {
					logrus.Errorf("updateToken get err:%v", err)
					continue
				}
				global.Store(token)
				logrus.Infof("updateToken get new token:%v", GetToken())
				break
			}
		}
	}
}

// GetToken 获取 token
func GetToken() string {
	return global.Load().(string)
}
