package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

//通过淘口令得到对应商品的URL
func GetURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//search := "tkl=" + "【@港仔文艺男 夏季韩版潮流宽松休闲裤男士街头纯色直筒裤九分裤】，复制这条信息€7JHd0ENgkIa€后打开👉淘宝👈" //r.PostForm["zhikouling"] //得到前端的淘口令
	//date := strings.NewReader(search)
	sssss1 := r.FormValue("name")
	sssss:="tkl="+sssss1
	date := strings.NewReader(sssss)
	urll := "http://api.chaozhi.hk/tb/tklParse"
	request, err := http.NewRequest("POST", urll, date)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Encoding", "gzip, deflate")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	request.Header.Set("Content-Length", "320")

	request.Header.Set("Host", "api.chaozhi.hk")
	request.Header.Set("Origin", "http://tool.chaozhi.hk")
	request.Header.Set("Proxy-Connection", "keep-alive")
	request.Header.Set("Referer", "http://tool.chaozhi.hk/")

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	re := string(respBytes)
	fmt.Println(re)
	res := gjson.Get(re, "data.url")
	//fmt.Println(res)

	//Date := strings.NewReader(res.Str)
	str := url.QueryEscape(res.Str)
	urlll := "http://tool.manmanbuy.com/m/history.aspx?DA=1&action=gethistory&url=" + str + "&token=jb8n37e966ca1a60164089724f0b00ffd84865vxq8z6"
	//fmt.Println(urlll)

	req, err := http.NewRequest("GET", urlll, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept", "application/json, text/javascript, */*")

	Client := http.Client{}
	response, err := Client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}
	defer response.Body.Close()
	respBytes1, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}
	result := gjson.Get(string(respBytes1), "jsData")
	fmt.Println(result.Str)
	//temp:=make([]string,1000)
	////temp=append(temp,result.Str)
	//temp=strings.Split(result.Str,",")
	//for _,v:=range temp {
	//	fmt.Println(v)
	//}
	bb:=[]byte(result.Str)
	//fmt.Println(result.Type)
	w.Write(bb)
}

func main() {
	http.HandleFunc("/", GetURL)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
