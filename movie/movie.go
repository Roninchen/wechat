package movie

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"wechat/conf"
)
var serviceAddress string
func init(){
	//初始化配置文件
	conf.Init("")
	serviceAddress = viper.GetString("movie.address")
}

func GetMovieInfo(method string,params string) []MovieInfo {
	conn,err:=grpc.Dial(serviceAddress,grpc.WithInsecure())
	if err != nil {
		logs.Info(err)
		return nil
	}
	defer conn.Close()
	client := NewMovieServiceClient(conn)
	movieResult, err := client.GetResult(context.Background(), &MovieRequest{Method: method, Params: params})

	if err!=nil {
		logs.Info(err)
		return nil
	}
	if movieResult.Code ==200 {
		movieInfo := make([]MovieInfo,1)
		err := json.Unmarshal(movieResult.Data, &movieInfo)
		if err != nil {
			logs.Info(err)
			return nil
		}
		return movieInfo
	}else {
		return nil
	}
}

type  MovieInfo struct {
	Id 						int64  `json:"id"`
	Movie_id 				int64  `json:"movie_id"`
	Movie_name 				string `json:"movie_name"`
	Movie_pic 				string `json:"movie_pic"`
	Movie_director 			string `json:"movie_director"`
	Movie_writer 			string `json:"movie_writer"`
	Movie_country 			string `json:"movie_country"`
	Movie_language 			string `json:"movie_language"`
	Movie_main_character 	string `json:"movie_main_character"`
	Movie_type 				string `json:"movie_type"`
	Movie_on_time  			string `json:"movie_on_time"`
	Movie_span 				string `json:"movie_span"`
	Movie_grade 			string `json:"movie_grade"`
	Remark 					string `json:"remark"`
	Movie_summary 			string `json:"movie_summary"`
	Movie_hot_comment 		string `json:"movie_hot_comment"`
	Episode 				string `json:"episode"`
	Season 					string `json:"season"`
	_Create_time 			string `json:"_create_time"`
	_Modify_time 			string `json:"_modify_time"`
	_Status 				int64  `json:"_status"`
}