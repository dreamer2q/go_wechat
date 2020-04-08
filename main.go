package main

import (
	"log"
	"wepchat_subscription/api"
	"wepchat_subscription/config"
)

func main() {

	wx := api.New(&api.Config{
		AppID:        config.AppID,
		AppSecret:    config.AppSecret,
		AppToken:     config.AppToken,
		AesEncodeKey: config.AesEncodeKey,
		Callback:     "/wx",
	})

	log.Panicln(wx.Run(":80"))

}
