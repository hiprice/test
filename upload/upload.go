package main

import (
    "github.com/chanxuehong/wechat/mp/material"
//	"github.com/chanxuehong/wechat/mp/media"
	"github.com/chanxuehong/wechat/mp"
	"fmt"
)

const (
	AppId = "wx96ae3fe27ad45e53"
	AppSecret = "ea4a0db81cc6a0017a69b0172515d5d8"
	token = "wexin"
)
//{image SQP8zwCqsiJP02ccSx2cY80w6e5q1K0FUH2QA5m8aPgQA3Ys0Xsxal8Li21sg_ia 1448586880}

func main() {

	imagePath := "../img/IMG_0129.JPG"
	accessTokenSer := mp.NewDefaultAccessTokenServer(AppId,AppSecret,nil)

//	临时素材
//	client := media.NewClient(accessTokenSer,nil)

	client := material.NewClient(accessTokenSer,nil)
	mediaInfo,url ,err := client.UploadImage(imagePath)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println("mediaInfo : ",mediaInfo)
	fmt.Println("url : ",url)

}