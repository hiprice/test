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

	/* 返回结果
	//临时素材
	{image SQP8zwCqsiJP02ccSx2cY80w6e5q1K0FUH2QA5m8aPgQA3Ys0Xsxal8Li21sg_ia 1448586880}

	//永久素材
	mediaInfo :  X9q8xoHV7-3W5E-ohorp_eiKp53vOv-2SUYipXHHA_A
	url :  https://mmbiz.qlogo.cn/mmbiz/u3RgicO3YCJa2aljOxNXWHKOU8tvenDZWVTCqptOn71bo2OCKr5A6xTSJDxwaZFoHPewvSPC3Xbw9ibX9DWvMsmQ/0?wx_fmt=jpeg

	*/
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println("mediaInfo : ",mediaInfo)
	fmt.Println("url : ",url)

}