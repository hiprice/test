package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
	"github.com/chanxuehong/wechat/mp/menu"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
)

var AccessTokenServer = mp.NewDefaultAccessTokenServer(AppId, AppSecret, nil) // 一個應用只能有一個實例
var mpClient = mp.NewClient(AccessTokenServer, nil)


func main() {

	messageServeMux := mp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	//事件处理
	messageServeMux.EventHandleFunc(menu.EventTypeClick, EventMessageHandler) // 注册文本处理 Handler


	// 下面函数的几个参数设置成你自己的参数: oriId, token, appId
	mpServer := mp.NewDefaultServer("", token, "", nil, messageServeMux)

	mpServerFrontend := mp.NewServerFrontend(mpServer, mp.ErrorHandlerFunc(ErrorHandler), nil)

	// 那么可以这么注册 http.Handler
	http.Handle("/index", mpServerFrontend)
	http.ListenAndServe(":80", nil)
}

//====事件推送====
func EventMessageHandler(w http.ResponseWriter, r *mp.Request) {

	text := menu.GetClickEvent(r.MixedMsg)
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.EventKey)

	mp.WriteRawResponse(w, r, resp) // 明文模式
}

//====处理文本推送====
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}


// 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *mp.Request) {

	fmt.Println("用户事件：",r.MixedMsg.Event)
	fmt.Println(r.MixedMsg.EventKey)


	fmt.Println("==用户请求==")
	fmt.Println(string(r.RawMsgXML))

	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)

	mp.WriteRawResponse(w, r, resp) // 明文模式
	//	mp.WriteAESResponse(w, r, resp) // 安全模式
}