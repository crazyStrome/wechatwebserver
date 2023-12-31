package test

import (
	"encoding/xml"
	"testing"
	"wechatwebserver/service"
)

func TestUnmarshalMsg(t *testing.T) {
	msgStr := `<xml>
  <ToUserName><![CDATA[toUser]]></ToUserName>
  <FromUserName><![CDATA[fromUser]]></FromUserName>
  <CreateTime>1348831860</CreateTime>
  <MsgType><![CDATA[text]]></MsgType>
  <Content><![CDATA[this is a test]]></Content>
  <MsgId>1234567890123456</MsgId>
  <PicUrl><![CDATA[this is a url]]></PicUrl>
  <MediaId><![CDATA[media_id]]></MediaId>
  <Format><![CDATA[Format]]></Format>
  <Recognition><![CDATA[腾讯微信团队]]></Recognition>
  <ThumbMediaId><![CDATA[thumb_media_id]]></ThumbMediaId>
  <Location_X>23.134521</Location_X>
  <Location_Y>113.358803</Location_Y>
  <Scale>20</Scale>
  <Label><![CDATA[位置信息]]></Label>
</xml>
`
	msg := &service.Request{}
	if err := xml.Unmarshal([]byte(msgStr), msg); err != nil {
		t.Error(err)
	}
	t.Logf("%+v", msg)
}
