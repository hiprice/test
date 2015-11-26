package main

import (
	"github.com/chanxuehong/wechat/mp/media"
	"github.com/chanxuehong/wechat/mp"
	"fmt"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
)

func main() {

	imagePath := "../IMG_0129.JPG"
	accessTokenSer := mp.NewDefaultAccessTokenServer(AppId,AppSecret,nil)

	client,err := media.NewClient(accessTokenSer,nil)
	if err != nil{
		fmt.Println(err)
	}

	mediaInfo ,err := client.UploadImage(imagePath)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(mediaInfo)
}