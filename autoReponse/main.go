package main

import (
	"log"
	"net/http"
	"fmt"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/garyburd/redigo/redis"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
	mediaId = "X9q8xoHV7-3W5E-ohorp_eiKp53vOv-2SUYipXHHA_A"
)

var AccessTokenServer = mp.NewDefaultAccessTokenServer(AppId, AppSecret, nil) // 一個應用只能有一個實例
var mpClient = mp.NewClient(AccessTokenServer, nil)


func main() {

	messageServeMux := mp.NewMessageServeMux()

	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler

	//菜单事件处理
	messageServeMux.EventHandleFunc(menu.EventTypeClick, EventMessageHandler) // 注册时间处理 Handler

	//菜单跳转链接时间处理
	messageServeMux.EventHandleFunc(menu.EventTypeView, ViewMessageHandler) // 注册时间处理 Handler

	//菜单地理位置处理
	messageServeMux.EventHandleFunc(menu.EventTypeLocationSelect, LocationEventSelectMessageHandle)

	//自动上报地理位置消息
	messageServeMux.EventHandleFunc(request.EventTypeLocation, LocationEventMessageHandle)

	// 下面函数的几个参数设置成你自己的参数: oriId, token, appId
	mpServer := mp.NewDefaultServer("", token, "", nil, messageServeMux)

	mpServerFrontend := mp.NewServerFrontend(mpServer, mp.ErrorHandlerFunc(ErrorHandler), nil)

	// 那么可以这么注册 http.Handler
	http.Handle("/index", mpServerFrontend)
	http.HandleFunc("/createMenu",CreateMenu)

	http.ListenAndServe(":80", nil)
}


func ViewMessageHandler(w http.ResponseWriter ,r *mp.Request){
	context := menu.GetViewEvent(r.MixedMsg)

	fmt.Println(context)

	content := context.MsgType

	resp := response.NewText(context.FromUserName,context.ToUserName,context.CreateTime,content)

	mp.WriteRawResponse(w,r,resp)
	fmt.Println("ViewMessageHandler")
}

func LocationEventMessageHandle(w http.ResponseWriter, r *mp.Request) {
	location := request.GetLocationEvent(r.MixedMsg)

	fmt.Println(location)
	fmt.Println("LocationEventMessageHandle")
}

//菜单处理地理位置信息
func LocationEventSelectMessageHandle(w http.ResponseWriter, r *mp.Request) {
	location := menu.GetLocationSelectEvent(r.MixedMsg)
	fmt.Println(location)

	resp := response.NewText(location.FromUserName,location.ToUserName,location.CreateTime,location.SendLocationInfo.Label)

	mp.WriteRawResponse(w,r,resp)
	fmt.Println("LocationEventSelectMessageHandle")
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	var subButtons = make([]menu.Button, 2)

	subButtons[0].SetAsViewButton("搜索", "http://www.soso.com/")
	subButtons[1].SetAsClickButton("赞一下我们", "V1001_GOOD")

	var testButtons = make([]menu.Button, 3)
	testButtons[0].SetAsClickButton("今日歌曲", "V1001_TODAY_MUSIC")
	testButtons[1].SetAsClickButton("来一张图片", "V1001_IMG")
	testButtons[2].SetAsLocationSelectButton("地理位置", "V1001_LOCATION")

	var mn menu.Menu
	mn.Buttons = make([]menu.Button, 3)
	mn.Buttons[0].SetAsSubMenuButton("test", testButtons)
	mn.Buttons[1].SetAsViewButton("视频", "http://v.qq.com/")
	mn.Buttons[2].SetAsSubMenuButton("子菜单", subButtons)

	menuClient := (*menu.Client)(mpClient)
	if err := menuClient.CreateMenu(mn); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
	w.Write([]byte("create menu:success"))
}

//====自定义事件推送====
func EventMessageHandler(w http.ResponseWriter, r *mp.Request) {

	text := menu.GetClickEvent(r.MixedMsg)
	location := menu.GetLocationSelectEvent(r.MixedMsg)

//	key := "click_count_"+text.EventKey

	var content string
	switch text.EventKey {
	case "V1001_TODAY_MUSIC":
		content = text.EventKey + "你点击了一下"
	case "V1001_GOOD":
		content = text.EventKey + "收到您的点赞，我非常高兴"
	case "V1001_IMG":	//恢复图片信息
		resp := response.NewImage(text.FromUserName,text.ToUserName,text.CreateTime,mediaId)
		mp.WriteRawResponse(w,r,resp)
	case "V1001_LOCATION":
		content = text.Event + "text - 地理位置上报成功"
		fmt.Println(content)
		fmt.Println(location.SendLocationInfo.Label)
		fmt.Println(location.SendLocationInfo.PoiName)

	default:
		content = text.EventKey + "oh ,what is wrong"
	}

	switch location.EventKey {
	case "V1001_LOCATION" :
		fmt.Println(location.EventKey)
		content = location.Event + " location - 上报地理位置成功"
		fmt.Println(location.SendLocationInfo.Label)
		fmt.Println(location.SendLocationInfo.PoiName)
	default:
		fmt.Println(location.EventKey)
		content = "location - 上报地理位置失败"
	}
//	Incr(key)

	resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, content)

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

//统计参数
func Incr(key string)(bool,error){
	var aft bool = false

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer conn.Close()

	if err != nil{
		fmt.Println("连接faild")
		return aft,err
	}

	num,err := conn.Do("INCR",key)

	if err != nil{
		return aft,err
	}
	fmt.Println("incr : num -" ,num)


	return true,nil
}