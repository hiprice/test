package main

import (
	"log"
	"net/http"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

// 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *mp.Request) {
	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)
	log.Println(text)
	//mp.WriteRawResponse(w, r, resp) // 明文模式
	mp.WriteAESResponse(w, r, resp) // 安全模式
}

func main() {
//	aesKey, err := util.AESKeyDecode("encodedAESKey") // 这里 encodedAESKey 改成你自己的参数
//	if err != nil {
//		panic(err)
//	}

	messageServeMux := mp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	// 下面函数的几个参数设置成你自己的参数: oriId, token, appId
	mpServer := mp.NewDefaultServer("", token, AppId, nil, messageServeMux)

	mpServerFrontend := mp.NewServerFrontend(mpServer, mp.ErrorHandlerFunc(ErrorHandler), nil)

	// 如果你在微信后台设置的回调地址是
	//   http://xxx.yyy.zzz/wechat
	// 那么可以这么注册 http.Handler
	http.Handle("/index", mpServerFrontend)
	http.ListenAndServe(":80", nil)
}