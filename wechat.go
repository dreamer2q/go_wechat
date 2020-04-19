package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/dreamer2q/go_wechat/media"
	"github.com/dreamer2q/go_wechat/menu"
	"github.com/dreamer2q/go_wechat/message"
	"github.com/dreamer2q/go_wechat/midware"
	"github.com/dreamer2q/go_wechat/request"
	"github.com/dreamer2q/go_wechat/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type API struct {
	EventHandle   MessageHandler
	MessageHandle MessageHandler

	Media    *media.Media
	Menu     *menu.Menu
	Template *message.Template
	User     *user.User

	config *Config
}

func New(c *Config) *API {
	rc := &request.Config{
		AppID:        c.AppID,
		AppSecret:    c.AppSecret,
		AppToken:     c.AppToken,
		AesEncodeKey: c.AesEncodeKey,
		Callback:     c.Callback,
		Timeout:      c.Timeout,
	}
	r := request.New(rc)
	return &API{
		MessageHandle: defaultMessageHandler,
		EventHandle:   defaultMessageHandler,

		Media:    media.New(r),
		Menu:     menu.New(r),
		Template: message.New(r),
		User:     user.New(r),

		config: c,
	}
}

func (w *API) Run(addr ...string) error {

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(midware.Logger())
	r.Use(midware.Verify(w.config.AppToken))

	r.Any(w.config.Callback, w.requestHandler)
	return r.Run(addr...)
}

func (w *API) requestHandler(c *gin.Context) {
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
