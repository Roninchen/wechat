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
		return
	}
	if strings.Contains(inMsg,"数据分析") {
		bytes := interstingReturn(this,"可以这样输入：\n 平均分分布情况 \n 中国大陆影片平均分分布 \n 每年上映电影数目 \n 每年上映影片平均分 \n 影片时长分布 \n 影片时长与年份 \n 影片类型数目 \n 影片分类与评分 \n 不同类型电影的平均时长 \n 影片出品数量国家统计 \n 不同国家影片的平均评分 \n 各项评分、评论等参数之间的相关图 \n")
		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		this.ResponseWriter.Write(bytes)
		return
	}
	var news NewsMessage
	news.InitBaseData(this,"news")
	if strings.Contains(inMsg,"平均分分布情况") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"平均分分布情况:", Description: "呈现正态分布，大部分电影评分在6分到8分这个区间", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/1.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/1.jpg"})
		news.Articles = append(news.Articles, Article{Title:"中国大陆影片平均分分布:", Description: "中国大陆影片评分明显低于总体电影评分", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/2.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/2.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	if strings.Contains(inMsg,"中国大陆影片平均分分布") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"中国大陆影片平均分分布:", Description: "中国大陆影片评分明显低于总体电影评分", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/2.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/2.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	if strings.Contains(inMsg,"每年上映电影数目") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"每年上映电影数目:", Description: "近年来影片产量增长迅速", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/3.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/3.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	if strings.Contains(inMsg,"每年上映影片平均分") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"每年上映影片平均分:", Description: "总有人说电影越拍越差了，通过数据分析发现还真不假，电影年度平均分越来越低。从十年前的均分在7分左右下滑到均分在6左右，近几年来的评分呈现直线下降", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/4.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/4.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	// 影片时长分布
	if strings.Contains(inMsg,"影片时长分布") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"影片时长分布:", Description: "大部分影片时长在90min左右，基本符合正态分布", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/5.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/5.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	//影片时长与年份
	if strings.Contains(inMsg,"影片时长与年份") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"影片时长与年份:", Description: "近年来影片制作商还是倾向于拍时间短的时间。特别是近代流行的拍微电影,大部分微电影时长都在60分钟以下", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/6.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/6.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	//影片类型数目
	if strings.Contains(inMsg,"影片类型数目") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"影片类型数目:", Description: "每部影片它的分类都涵盖了很多种，一部电影有多个标签属性，有统计图可知，大部分影片分类都在前面的剧情、喜剧、动作、爱情等范围内", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/7.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/7.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	//影片分类与评分
	if strings.Contains(inMsg,"影片分类与评分") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"影片分类与评分:", Description: "不认为电影类型对结果没有什么影响，从统计角度看，音乐，传记和其他类型的相对小众类型的电影评分相对较高，而恐怖片的评分比较相对差，可能是国产恐怖片拉低了总体电影评分吧", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/8.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/8.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	//不同类型电影的平均时长
	if strings.Contains(inMsg,"不同类型电影的平均时长") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"不同类型电影的平均时长:", Description: "平均时长最长的影片在历史、战争、传记这些类型中。挺合乎情理，这类影片讲述的故事线都比较长。", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/9.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/9.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	//影片出品数量国家统计
	if strings.Contains(inMsg,"影片出品数量国家统计") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"影片出品数量国家统计:", Description: "美国出品或者美国参与制作的电影数量遥遥领先，可见美国在电影行业的地位（好莱坞大片，漫威电影，DC电影，迪斯尼电影），创造了一系列全球人民都喜欢看的电影。", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/10.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/10.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	// 不同国家影片的平均评分
	if strings.Contains(inMsg,"不同国家影片的平均评分") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"不同国家影片的平均评分:", Description: "由此可见，欧洲国家整体评分都很靠前，中国香港平均评分靠后，中国大陆影片平均评分直接垫底", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/11.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/11.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	//各项评分、评论等参数之间的相关图
	if strings.Contains(inMsg,"各项评分、评论等参数之间的相关图") {
		news.ArticleCount = 1
		news.Articles = append(news.Articles, Article{Title:"各项评分、评论等参数之间的相关图:", Description: "两个变量间的皮尔逊相关系数（两个变量间协方差和标准差的商），越接近1代表正相关，越靠近-1代表越负相关，0就是代表这两个变量间增长没有任何关系，“rate”表示评分（10分满分），“stars”表示豆瓣星级（5星为满级），“1,2,3,4,5”,分别代表“一星,二星,三星,四星,五星”占比情况，“wish”表示这部电影想看的人数，“collect”表示这部电影看过的人数、“comments”、“ratings”分别代表这部电影的写了短评的人数及评价了的人数（打了分就算评价，不用写评论）。", PicURL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/12.jpg", URL:"https://raw.githubusercontent.com/Roninchen/MarkdownPhoto/master/movie/12.jpg"})

		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		newsXml, err := xml.Marshal(news)
		if err != nil {
			logs.Info(err)
			this.ResponseWriter.WriteHeader(403)
			return
		}
		this.ResponseWriter.Write(newsXml)
		return
	}
	logs.Info(inMsg)
	if !ok {
		return
	}
	movieInfo := movie.GetMovieInfo("1", inMsg)
	if  movieInfo == nil{
		logs.Info("movieInfo is null")
		bytes := interstingReturn(this,"数据暂未收录")
		this.ResponseWriter.Header().Set("Content-Type", "news/xml")
		this.ResponseWriter.Write(bytes)
		return
	}
	//news.ArticleCount = 1
	news.ArticleCount = len(movieInfo)
	logs.Info(len(movieInfo))
	for _,v := range movieInfo {
		movieId := strconv.FormatInt(v.Movie_id,10)
		logs.Info("数据：",v)
		logs.Info("影片名",v.Movie_name)
		logs.Info("豆瓣评分",v.Movie_grade)
		logs.Info("url","https://movie.douban.com/subject/"+movieId+"/")
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