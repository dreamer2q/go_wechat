package wechat

import (
	"fmt"
	"reflect"
)

const (
	messagePrefix = "msg"
	eventPrefix   = "event"
)

type Handler func(msg MessageReceive) MessageReply

type ev struct {
	subscribers map[string][]Handler
}

func newEv() *ev {
	return &ev{
		subscribers: make(map[string][]Handler, 0),
	}
}

func (e *ev) get(event string) []Handler {
	return e.subscribers[event]
}

func (e *ev) trigger(msg *MessageReceive) []Handler {
	ret := make([]Handler, 0)
	if msg.MsgType == "" {
		return ret
	}
	if msg.IsEvent() {
		if msg.Event != "" {
			if msg.EventKey != "" {
				if es := e.get(
					fmt.Sprintf("%s.%s.%s", msg.MsgType, msg.Event, msg.EventKey),
				); es != nil {
					ret = append(ret, es...)
				}
			}
			if es := e.get(
				fmt.Sprintf("%s.%s", msg.MsgType, msg.Event),
			); es != nil {
				ret = append(ret, es...)
			}
		}
	} else {
		if es := e.get(
			fmt.Sprintf("%s.%s", messagePrefix, msg.MsgType),
		); es != nil {
			ret = append(ret, es...)
		}
		if es := e.get(messagePrefix); es != nil {
			ret = append(ret, es...)
		}
	}
	return ret
}

//添加订阅
func (e *ev) On(event string, handler Handler) Unsubscribe {
	e.subscribers[event] = append(e.subscribers[event], handler)
	return func() {
		e.Off(event, handler)
	}
}

type Unsubscribe func()

//取消订阅
func (e *ev) Off(event string, handler Handler) {
	handlers, ok := e.subscribers[event]
	if !ok {
		//do nothing
		return
	}
	newHandlers := make([]Handler, 0)
	for _, h := range handlers {
		if reflect.ValueOf(h) != reflect.ValueOf(handler) {
			newHandlers = append(newHandlers, h)
		}
	}
	e.subscribers[event] = newHandlers
}
