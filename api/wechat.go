package api

import (
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"wechat/api/midware"
)

type WechatAPI struct {
	EventHandle   MessageHandler
	MessageHandle MessageHandler
	config        *Config
}

func New(c *Config) *WechatAPI {
	return &WechatAPI{
		MessageHandle: defaultMessageHandler,
		EventHandle:   defaultMessageHandler,
		config:        c,
	}
}

func (w *WechatAPI) Run(addr ...string) error {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(midware.Logger())
	r.Use(midware.Verify(w.config.AppToken))

	r.Any(w.config.Callback, w.requestHandler)
	return r.Run(":80")
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

	if r.Type() == "noreply" {
		c.String(http.StatusOK, "success")
		return
	}
	reply := &messageReply{
		XMLName: xml.Name{Local: "xml"},
		messageBase: messageBase{
			ToUserName:   raw.FromUserName,
			FromUserName: raw.ToUserName,
			CreateTime:   time.Now().Unix(),
			MsgType:      r.Type(),
		},
		Msg: r,
	}
	//xmlReply, err := xml.Marshal(&reply)
	//if err != nil {
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	//fmt.Printf("xmlReply: %s\n", xmlReply)

	//use this frame provided method, to shorten code
	c.XML(http.StatusOK, &reply)
	//c.String(http.StatusOK, "success")
}

//default handler
func defaultMessageHandler(m *MessageReceive) MessageReply {
	//return  noreply to make sure that wechat server do not think we are dead
	return NoReply{}
}
