package wx

import (
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/clbanning/mxj"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"wechat/movie"
)

type weixinQuery struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`
}

type WeixinClient struct {
	Token          string
	Query          weixinQuery
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeixinClient, error) {

	weixinClient := new(WeixinClient)

	weixinClient.Token = token
	weixinClient.Request = r
	weixinClient.ResponseWriter = w

	weixinClient.initWeixinQuery()

	if weixinClient.Query.Signature != weixinClient.signature() {
		return nil, errors.New("Invalid Signature.")
	}

	return weixinClient, nil
}

func (this *WeixinClient) initWeixinQuery() {

	var q weixinQuery

	q.Nonce = this.Request.URL.Query().Get("nonce")
	q.Echostr = this.Request.URL.Query().Get("echostr")
	q.Signature = this.Request.URL.Query().Get("signature")
	q.Timestamp = this.Request.URL.Query().Get("timestamp")
	q.EncryptType = this.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = this.Request.URL.Query().Get("msg_signature")

	this.Query = q
}

func (this *WeixinClient) signature() string {

	strs := sort.StringSlice{this.Token, this.Query.Timestamp, this.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *WeixinClient) initMessage() error {

	body, err := ioutil.ReadAll(this.Request.Body)

	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)

	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("Invalid Message.")
	}

	message, ok := m["xml"].(map[string]interface{})

	if !ok {
		return errors.New("Invalid Field `xml` Type.")
	}

	this.Message = message

	log.Println(this.Message)

	return nil
}

func (this *WeixinClient) text() {

	inMsg, ok := this.Message["Content"].(string)
	if strings.Contains(inMsg,"易达") {
		bytes := interstingReturn(this,"找帅哥干嘛!")
		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		this.ResponseWriter.Write(bytes)
	}
	logs.Info(inMsg)
	if !ok {
		return
	}
	movieInfo := movie.GetMovieInfo("1", inMsg)
	if  movieInfo == nil{
		logs.Info("movieInfo is null")
		return
	}
	var news NewsMessage
	news.InitBaseData(this,"news")
	//news.ArticleCount = 1
	news.ArticleCount = len(movieInfo)
	logs.Info(len(movieInfo))
	for _,v := range movieInfo {
		movieId := strconv.FormatInt(v.Movie_id,10)
		logs.Info(v)
		news.Articles = append(news.Articles, Article{Title:"影片名:"+v.Movie_name, Description: "豆瓣评分:"+v.Movie_grade+"分", PicURL:v.Movie_pic, URL:"https://movie.douban.com/subject/"+movieId+"/"})
	}
	//movieId := strconv.FormatInt(movieInfo.Movie_id,10)
	//news.Articles = append(news.Articles, &Article{Title:movieInfo.Movie_name, Description: movieInfo.Movie_grade+"分", PicURL:movieInfo.Movie_pic, URL:"https://movie.douban.com/subject/"+movieId+"/"})
	newsXml, err := xml.Marshal(news)
	if err != nil {
		logs.Info(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	//var reply TextMessage
	//info := movie.GetMovieInfo("1", inMsg)
	//reply.InitBaseData(this, "text")
	////reply.Content = value2CDATA(fmt.Sprintf("我收到的是：%s", inMsg))
	//log.Println(info)
	//reply.Content = value2CDATA(fmt.Sprintf("为你查到以下信息：\n%s", info))
	//replyXml, err := xml.Marshal(reply)
    //log.Println(replyXml)
	//if err != nil {
	//	log.Println(err)
	//	this.ResponseWriter.WriteHeader(403)
	//	return
	//}
	//
	//this.ResponseWriter.Header().Set("Content-Type", "news/xml")
	this.ResponseWriter.Write(newsXml)
}

func (this *WeixinClient) Run() {

	err := this.initMessage()

	if err != nil {

		logs.Info(err)
		this.ResponseWriter.WriteHeader(403)
		return
	}

	MsgType, ok := this.Message["MsgType"].(string)

	if !ok {
		this.ResponseWriter.WriteHeader(403)
		return
	}

	switch MsgType {
	case "text":
		this.text()
		break
	default:
		break
	}

	return
}

func interstingReturn(this *WeixinClient,text string) []byte {
	var reply TextMessage
	reply.InitBaseData(this, "text")
	reply.Content = value2CDATA(fmt.Sprintf(text))
	replyXml, err := xml.Marshal(reply)
	log.Println(replyXml)
	if err != nil {
		log.Println(err)
		this.ResponseWriter.WriteHeader(403)
		return nil
	}
	return replyXml
}