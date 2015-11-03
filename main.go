package main

import (
	"net/http"
	"log"
	"sort"
	"io"
	"crypto/sha1"
	"fmt"

//	"github.com/chanxuehong/wechat/mp"
)

const AppID = "wx96ae3fe27ad45e53"
const AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"

func main() {
	http.HandleFunc("/index", sign)
	http.ListenAndServe(":8001", nil)
}

func sign(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	checked := checkSign(req.Form["signature"],req.Form["timestamp"],req.Form["nonce"],req.Form["token"])
	echoStr := req.Form["echoStr"]

	log.Println(echoStr)
	if checked {
		w.Write([]byte(echoStr))
	}
}

func checkSign(sign,timestamp,nonce,token string) bool {
	tmp := []string{token,timestamp,nonce}
	sort.StringSlice{token,timestamp,nonce}
	//连接为字符串
	var str string
	for _,v := range tmp{
		str += v
	}
	//sha1加密
	t := sha1.New();
	io.WriteString(t,str);
	sumSign :=  fmt.Sprintf("%x",t.Sum(nil));

	//比较
	if sumSign == sign{
		return true
	}
	return false
}