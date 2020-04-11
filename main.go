package main

import (
	"fmt"
	"log"
	"time"
	"wechat/api"
	"wechat/api/media"
	"wechat/api/request"
	"wechat/config"
)

func main() {

	c := &request.Config{
		AppID:        config.AppID,
		AppSecret:    config.AppSecret,
		AppToken:     config.AppToken,
		AesEncodeKey: config.AesEncodeKey,
		Callback:     "wx",
		Timeout:      10 * time.Second,
	}
	r := request.New(c)
	m := media.New(r)
	res, err := m.UploadMaterial("assets/123.png", media.TypImage, true)
	if err != nil {
		log.Panicln(err)
	}
	ret, err := m.GetMaterial(res.MediaID)
	fmt.Println(ret)
	if err != nil {
		log.Panicln(err)
	}

	wx := api.New(&api.Config{
		AppID:        config.AppID,
		AppSecret:    config.AppSecret,
		AppToken:     config.AppToken,
		AesEncodeKey: config.AesEncodeKey,
		Callback:     "/wx",
	})

	wx.EventHandle = func(m *api.MessageReceive) api.MessageReply {
		log.Printf("EventHandler: user: %s: %s\n", m.FromUserName, m.Event)
		switch m.Event {
		case "subscribe":
			log.Printf("subscribe")
			return api.Text{Content: fmt.Sprintf("%s: welcome", m.FromUserName)}
		case "unsubscribe":
			log.Printf("unsubscribe")
		default:
			log.Printf("default")
		}
		return api.NoReply{}
	}

	wx.MessageHandle = func(m *api.MessageReceive) api.MessageReply {
		return api.Text{Content: "Hello"}
	}

	log.Panicln(wx.Run(":80"))
}
