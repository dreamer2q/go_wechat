# Go Wechat

提供微信公众号基础能力支持

## 已实现

- [x] 自定义菜单
- 消息管理
  - [x] 接收消息
  - [x] 被动回复
  - [x] 消息加密
  - [x] 模板消息 [待改进]
  - [x] 群发消息
- [x] 素材管理
- [x] 留言管理(not tested)
- [x] 用户管理
- [x] 账号管理

## TODO

- [ ] 支持 TLS

**PS**需要更多支持，请参考[wechat](https://github.com/silenceper/wechat)

# 快速上手

## 接入指南

参考[微信开发文档](https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Access_Overview.html)

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
	wx := wc.New(&wc.Config{
		AppID:        AppID,
		AppSecret:    AppSecret,
		AppToken:     AppToken,
		AesEncodeKey: AesEncodeKey,
		Timeout:      10 * time.Second,
		Callback:     "/wx",
		Debug:        true,
	})

	wx.SetEventHandler(func(m wc.MessageReceive) wc.MessageReply {
		switch m.Event {
		case wc.EvSubscribe:
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
	})

	wx.SetMessageHandler(func(m wc.MessageReceive) wc.MessageReply {
		switch m.MsgType {
		case wc.MsgText:
			return wc.Text{Content: m.Content}
		default:
			return wc.Text{Content: "Not support type"}
		}
	})

	log.Panicln(wx.Run(":80"))
}
```

## 菜单处理

### 创建菜单

```go
    err := wx.Menu.Create(
		menu.RootMenu{
			Menus: []menu.Item{
				&menu.ClickMenu{
					Name: "点击测试",
					Key:  "click_test",
				},
				&menu.SubMenu{
					Name: "二级菜单",
					Menus: []menu.Item{
						&menu.ViewMenu{
							Name: "博客",
							Url:  "https://dreamer2q.wang",
						},
						&menu.ClickMenu{
							Name: "关于我",
							Key:  "click_about",
						},
						&menu.ClickMenu{
							Name: "点我没反应",
							Key:  "noempty",
						},
					},
				},
			},
		})
	if err != nil {
		log.Panic(err)
    }
```

### 处理菜单事件

- 方法一，订阅事件

```go
	wx.On("event.CLICK.click_test", func(msg wc.MessageReceive) wc.MessageReply {
		return wc.Text{Content: "点击菜单测试"}
	})
	wx.On("event.CLICK.click_about", func(msg wc.MessageReceive) wc.MessageReply {
		return wc.Text{Content: "关于我： 我是傻逼开发者"}
	})
	wx.On("event.CLICK.noempty", func(msg wc.MessageReceive) wc.MessageReply {
		return nil //no reply
	})
	wx.On("event.VIEW.https://dreamer2q.wang", func(msg wc.MessageReceive) wc.MessageReply {
		log.Printf("view event") //记录事件
		return nil
	})

```

- 方法二，在总事件中处理

```go
    wx.SetEventHandler(func(m wc.MessageReceive) wc.MessageReply {
		log.Printf("EventHandler: user %s: %s\n", m.FromUserName, m.Event)
		switch m.Event {
		case wc.EvClick:
			switch m.EventKey {
			case "click_test":
				return wc.Text{Content: "点击菜单测试"}
			case "click_about":
				return wc.Text{Content: "关于我，我是xxxx"}
			case "noempty":
				return nil
			}
		case wc.EvView:
			log.Printf("view event: %s", m.EventKey)
		}
		return nil
	})
```

**订阅事件优先级高，以高优先级处理结果为准。**

## 素材管理

### 上传素材

### 获取素材列表

## 用户管理

## 账号管理
