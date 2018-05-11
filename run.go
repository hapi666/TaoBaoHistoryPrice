package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

//通过淘口令得到对应商品的URL以及商品的对应id
func GetURL(w http.ResponseWriter, r *http.Request) {
	//http.FileServer()
	r.ParseForm()
	search := "tkl=" + "【CSDN CSDN下载  CSDN代下载  代下CSDN 代下CSDN  极速发货】，复制这条信息￥JXGJ0Itkh2F￥后打开手淘" //r.PostForm["zhikouling"] //得到前端的淘口令
	date := strings.NewReader(search)

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
	fmt.Println(res)
	//str := string(respBytes)
	//fmt.Println(str)
	u, err := url.Parse(res.Str)
	if err != nil {
		log.Fatal(err.Error())
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		log.Fatal(err.Error())
	}
	id := m["id"][0]
	//fmt.Println(id)
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	find, err := db.Query("SELECT price FROM commodity WHERE sid=?", template.HTMLEscapeString(id))
	if err != nil {
		log.Fatal(err.Error())
	}
	var result []string
	for find.Next() {
		find.Scan(&result)
	}
	if result != nil {
		// 将结果集result转换为json数据
		res, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(string(res))
		// 给前端发送数据
		w.Write(res)
	} else {
		w.Write([]byte("fail!"))
	}

}

func main() {
	http.HandleFunc("/", GetURL)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
