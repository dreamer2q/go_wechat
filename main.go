package main

import (
	"fmt"
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

	err := wx.Token.Refresh()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("token: %s\n", wx.Token)

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
