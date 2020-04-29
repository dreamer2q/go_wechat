package wechat

import (
	"github.com/dreamer2q/go_wechat/account"
	"github.com/dreamer2q/go_wechat/media"
	"github.com/dreamer2q/go_wechat/menu"
	"github.com/dreamer2q/go_wechat/message"
	"github.com/dreamer2q/go_wechat/request"
	"github.com/dreamer2q/go_wechat/user"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

type API struct {
	Media    *media.Media
	Menu     *menu.Menu
	Template *message.Template
	User     *user.User
	Account  *account.Account

	*Ev

	Handlers []gin.HandlerFunc

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
	ev := newEv()
	api := &API{
		Media:    media.New(r),
		Template: message.New(r),
		User:     user.New(r),
		Account:  account.New(r),
		Menu:     menu.New(r),

		//event emmiter
		Ev: ev,

		config: c,
	}

	api.Handlers = []gin.HandlerFunc{
		api.logger(),       //logger 需要记录完整的事件，需要处于第一个
		api.verifier(),     //微信请求认证
		api.debugger(),     //输出请求和发送的body
		api.encryptor(),    //消息加密的透明代理
		api.requestHandler, //消息处理
	}

	return api
}

//attach route to a gin
func (w *API) AttachToGin(r *gin.Engine) {
	r.GET(w.config.Callback, w.Handlers...)  //fired when setting callback addr
	r.POST(w.config.Callback, w.Handlers...) //fired when message/event comes in
}

//default instance
func (w *API) Run(addr ...string) error {
	r := gin.New()
	r.Use(gin.Recovery())

	w.AttachToGin(r)
	return r.Run(addr...)
}

//main handle
func (w *API) requestHandler(c *gin.Context) {
	raw := &MessageReceive{}
	if err := c.ShouldBindXML(raw); err != nil {
		log.Printf("%#v", errors.Wrap(err, "bindXML"))
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

			//use framework provided method
			c.XML(http.StatusOK, reply)
		}
	)
	handlers := w.trigger(raw)
	for _, h := range handlers {
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
	w.On(messagePrefix, handler)
	return func() {
		w.Off(messagePrefix, handler)
	}
}

func (w *API) SetEventHandler(handler Handler) Unsubscribe {
	w.On(eventPrefix, handler)
	return func() {
		w.Off(eventPrefix, handler)
	}
}
