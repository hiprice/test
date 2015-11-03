package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/menu"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
)

func main() {
	log.Println("Wechat Service: Start!")
	http.HandleFunc("/", procRequest)
	http.HandleFunc("/menu", createMenu())
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Wechat Service: ListenAndServe failed, ", err)
	}
	log.Println("Wechat Service: Stop!")
}

func makeSignature(timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signatureGen := makeSignature(timestamp, nonce)

	signatureIn := strings.Join(r.Form["signature"], "")
	if signatureGen != signatureIn {
		return false
	}
	echostr := strings.Join(r.Form["echostr"], "")
	fmt.Fprintf(w, echostr)
	return true
}

func procRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !validateUrl(w, r) {
		log.Println("Wechat Service: this http request is not from Wechat platform!")
		return
	}
	log.Println("Wechat Service: validateUrl Ok!")
}

func createMenu() {

	var AccessTokenServer = mp.NewDefaultAccessTokenServer(AppId, AppSecret, nil) // 一個應用只能有一個實例
	var mpClient = mp.NewClient(AccessTokenServer, nil)

	var subButtons = make([]menu.Button, 2)

	subButtons[0].SetAsViewButton("搜索", "http://www.soso.com/")
	subButtons[1].SetAsClickButton("赞一下我们", "V1001_GOOD")

	var mn menu.Menu
	mn.Buttons = make([]menu.Button, 3)
	mn.Buttons[0].SetAsClickButton("今日歌曲", "V1001_TODAY_MUSIC")
	mn.Buttons[1].SetAsViewButton("视频", "http://v.qq.com/")
	mn.Buttons[2].SetAsSubMenuButton("子菜单", subButtons)

	menuClient := (*menu.Client)(mpClient)
	if err := menuClient.CreateMenu(mn); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ok")
}

func GetToken() {

}