package main

import (
	"fmt"

	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/mass/mass2all"
	"github.com/chanxuehong/wechat/mp/message/mass"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
)

var AccessTokenServer = mp.NewDefaultAccessTokenServer(AppId, AppSecret, nil) // 一個應用只能有一個實例
var mpClient = mp.NewClient(AccessTokenServer, nil)
func main() {
	msg := mass2all.NewText("send message to every body !")

	maClient := mass.NewClient(AccessTokenServer,nil)

	if re,err := maClient.MassToAll(msg);err != nil{
		fmt.Println(re)
		fmt.Println(err)
	}

	fmt.Println("ok")
}