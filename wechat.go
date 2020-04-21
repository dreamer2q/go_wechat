package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/dreamer2q/go_wechat/account"
	"github.com/dreamer2q/go_wechat/media"
	"github.com/dreamer2q/go_wechat/menu"
	"github.com/dreamer2q/go_wechat/message"
	"github.com/dreamer2q/go_wechat/midware"
	"github.com/dreamer2q/go_wechat/request"
	"github.com/dreamer2q/go_wechat/user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	Account  *account.Account

	*ev

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
		Media:    media.New(r),
		Menu:     menu.New(r),
		Template: message.New(r),
		User:     user.New(r),
		Account: account.New(r),

		//event emmiter
		ev:newEv(),

		config: c,
	}
}

func (w *API) Run(addr ...string) error {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(midware.Logger())
	r.Use(midware.Verify(w.config.AppToken))

	r.POST(w.config.Callback,w.requestHandler)

	//r.Any(w.config.Callback, w.requestHandler)
	return r.Run(addr...)
}

func (w *API) requestHandler(c *gin.Context) {
	raw := &MessageReceive{}
	if err := c.ShouldBindXML(raw); err != nil {
		log.Printf("%#v\n", errors.Wrap(err, "bindXML"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var (
		r      MessageReply
		sent   = false
		doSend = func() {
			sent = true
			reply := &xmlMsgReply{
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
			//debug only
			if w.config.debug {
			xmlReply, err := xml.Marshal(&reply)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			fmt.Printf("xmlReply: %s\n", xmlReply)
		}
	)
	for _,h := range handlers {
		r = h(*raw)
		if r != nil && !sent {
			doSend()
		}
	}
	if !sent {
		c.String(http.StatusOK, "success")
	}
}

func (w *API) SetMessageHandler(handler Handler) Unsubscribe {
	w.On(messagePrefix,handler)
	return func() {
		w.Off(messagePrefix,handler)
	}
}

func (w *API) SetEventHandler(handler Handler) Unsubscribe {
	w.On(eventPrefix,handler)
	return func() {
		w.Off(eventPrefix,handler)
	}
}
