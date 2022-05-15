package service

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CDATA 包裹
type CDATA struct {
	Test string `xml:",cdata"`
}

// Msg 微信发过来的消息格式
type Msg struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName CDATA    `xml:"ToUserName"`
	CreateTime uint64   `xml:"CreateTime"`
	MsgType    CDATA    `xml:"MsgType"`
	MsgId      uint64   `xml:"MsgId"`
	// 文本消息 text
	Content CDATA `xml:"Content"`
	// 图片消息 image
	PicUrl  CDATA `xml:"PicUrl"`
	MediaId CDATA `xml:"MediaId"`
	// 语音消息 voice
	Format      CDATA `xml:"Format"`
	Recognition CDATA `xml:"Recognition"`
	// 视频消息 video，小视频 shortvideo
	ThumbMediaId CDATA `xml:"ThumbMediaId"`
	// 地理位置 location
	Location_X float64 `xml:"Location_X"`
	Location_Y float64 `xml:"Location_Y"`
	Scale      int32   `xml:"Scale"`
	Label      CDATA   `xml:"Label"`
}

func procMsg(c *gin.Context) {
	logrus.Infof("procMsg req:%v", c.Request)
	start := time.Now()
	defer func() {
		logrus.Infof("procMsg done, req:%v, cost:%v", c.Request, time.Since(start))
	}()
	if err := handleMsg(c); err != nil {
		logrus.Errorf("handleMsg err:%v", err)
	}
	c.String(http.StatusOK, "")
}

func handleMsg(c *gin.Context) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("read body err:%v", err)
	}
	msg := &Msg{}
	if err := xml.Unmarshal(body, msg); err != nil {
		return fmt.Errorf("unmarshal err:%v, body:%s", err, body)
	}
	logrus.Infof("handMsg:%v", getJSON(msg))
	return nil
}
