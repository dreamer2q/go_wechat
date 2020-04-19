# Go Wechat

API Implementation for wechat subscription account in Golang

# 快速上手

**API 很可能会有很大的变动**

## 基本的消息与事件处理

```go
package main

import (
	"fmt"
	wc "github.com/dreamer2q/go_wechat"
	"github.com/pkg/errors"
	"log"
	"time"
)

func main() {

    //初始化
	wx := wc.New(&wc.Config{
		AppID:        AppID,
		AppSecret:    AppSecret,
		AppToken:     AppToken,
		AesEncodeKey: AesEncodeKey,
		Timeout:      10 * time.Second,
		Callback:     "/wx",
	})

    //设置事件处理函数，空为不处理
	wx.EventHandle = func(m *wc.MessageReceive) wc.MessageReply {
		log.Printf("EventHandler: user: %s: %s\n", m.FromUserName, m.Event)
		switch m.Event {
        case wc.EvSubscribe:
            //根据openID获取用户信息
            userInfo, err := wx.User.GetUserInfo(m.FromUserName)
			if err != nil {
				log.Printf("error: %#v", errors.Wrap(err, "getUserInfo"))
				return wc.Text{Content: "something went wrong"}
			}
			return wc.Text{Content: fmt.Sprintf("%s welcome you", userInfo.Nickname)}
		case wc.EvUnsubscribe:
			log.Printf("unsubscribe event")
		}
		//nil means no reply
		return nil
	}

    //设置消息处理函数，空为不处理
	wx.MessageHandle = func(m *wc.MessageReceive) wc.MessageReply {
		switch m.MsgType {
		case wc.MsgText:
			return wc.Text{Content: m.Content}
		default:
			return wc.Text{Content: "Not support type"}
		}
	}

	log.Panicln(wx.Run(":80"))
}

```

## 菜单处理

> 下面的格式和可能会大改

```go
err := wx.Menu.Create(
		menu.RootMenu{
			Buttons: []menu.Item{
				&menu.SubMenu{
					Name: "开始",
					Items: []menu.Item{
						&menu.View{
							Name: "博客",
							Url:  "https://dreamer2q.wang",
						},
						&menu.Click{
							Name: "点击测试",
							Key:  "click_test",
						},
					},
				},
				&menu.SubMenu{
					Name: "关于",
					Items: []menu.Item{
						&menu.Click{
							Name: "关于我",
							Key:  "click_aboutme",
						},
					},
				},
			},
		})
	if err != nil {
		log.Panic(err)
	}
```

菜单事件的消息处理

```go
wx.EventHandle = func(m *wc.MessageReceive) wc.MessageReply {
		log.Printf("EventHandler: user %s: %s\n", m.FromUserName, m.Event)
		switch m.Event {
		case wc.EvClick:
			switch m.EventKey {
			case "click_test":
				return wc.Text{Content: "点击测试"}
			case "click_aboutme":
				return wc.Text{Content: "关于我： 我是 dreamer2q"}
			}
		case wc.EvView:
			fmt.Printf("Event: view: %#v", m)
		}
		//nil means no reply
		return nil
	}
```

## 更多用例

```go
//TODO
```
