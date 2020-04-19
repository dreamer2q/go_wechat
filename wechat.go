package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"wechat/media"
	"wechat/menu"
	"wechat/message"
	"wechat/midware"
	"wechat/request"
)

type WechatAPI struct {
	EventHandle   MessageHandler
	MessageHandle MessageHandler

	Media    *media.Media
	Menu     *menu.Menu
	Template *message.Template

	config *Config
}

func New(c *Config) *WechatAPI {
	rc := &request.Config{
		AppID:        c.AppID,
		AppSecret:    c.AppSecret,
		AppToken:     c.AppToken,
		AesEncodeKey: c.AesEncodeKey,
		Callback:     c.Callback,
		Timeout:      c.Timeout,
	}
	r := request.New(rc)
	return &WechatAPI{
		MessageHandle: defaultMessageHandler,
		EventHandle:   defaultMessageHandler,
		Media:         media.New(r),
		Menu:          menu.New(r),
		Template:      message.New(r),
		config:        c,
	}
}

func (w *WechatAPI) Run(addr ...string) error {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(midware.Logger())
	r.Use(midware.Verify(w.config.AppToken))

	r.Any(w.config.Callback, w.requestHandler)
	return r.Run(addr...)
}

func (w *WechatAPI) requestHandler(c *gin.Context) {
	raw := &MessageReceive{}
	if err := c.ShouldBindXML(raw); err != nil {
		log.Printf("requestHandler: %v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Printf("request: %v\n", raw)

	f := w.MessageHandle
	if raw.IsEvent() {
		f = w.EventHandle
	}
	r := f(raw)

	//no reply
	if r == nil {
		c.String(http.StatusOK, "success")
		return
	}

	reply := &messageReply{
		messageBase: messageBase{
			ToUserName:   raw.FromUserName,
			FromUserName: raw.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      r.Type(),
		},
		MsgWrapper: messageWrapper{
			Msg: r,
		},
	}
	xmlReply, err := xml.Marshal(&reply)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Printf("xmlReply: %s\n", xmlReply)

	//use this frame provided method, to shorten code

	c.XML(http.StatusOK, &reply)
	//c.String(http.StatusOK, "success")
}

//default handler
func defaultMessageHandler(m *MessageReceive) MessageReply {
	//return  nil(noreply) to make sure that wechat server do not think we are dead
	return nil
}
