package wx

import (
	"encoding/xml"
	"strconv"
	"time"
)

type Base struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   CDATAText
}

func (b *Base) InitBaseData(w *WeixinClient, msgtype string) {

	b.FromUserName = value2CDATA(w.Message["ToUserName"].(string))
	b.ToUserName = value2CDATA(w.Message["FromUserName"].(string))
	b.CreateTime = value2CDATA(strconv.FormatInt(time.Now().Unix(), 10))
	b.MsgType = value2CDATA(msgtype)
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type TextMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}

type NewsMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	ArticleCount int    `xml:"ArticleCount"`
	Articles []Article `xml:"Articles>item,omitempty"`
}
type Item struct {
	Title CDATAText
	Description CDATAText
	PicUrl CDATAText
	Url CDATAText
}

//Article 单篇文章
type Article struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicURL      string `xml:"PicUrl,omitempty"`
	URL         string `xml:"Url,omitempty"`
}