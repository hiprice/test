package main

import (
	"net/http"
	"log"

//	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/util"
)

const AppID = "wx96ae3fe27ad45e53"
const AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"

func main() {
	http.HandleFunc("/index", sign)
	http.ListenAndServe(":8001", nil)
}

func sign(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	log.Println(req.Form)
//	checked := checkSign(req.Form["signature"],req.Form["timestamp"],req.Form["nonce"],req.Form["token"])
	echoStr := req.Form["echoStr"]

	log.Println(echoStr)
//	if checked {
//		w.Write([]byte(echoStr))
//	}
}

func checkSign(sign,timestamp,nonce,token string) bool {

	sumSign :=  util.Sign(token,timestamp,nonce);

	//比较
	if sumSign == sign{
		return true
	}
	return false
}