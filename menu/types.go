package menu

import (
	"fmt"
	"github.com/dreamer2q/go_wechat"
)

type Callback = wechat.Handler
type ev = *wechat.Ev

type RootMenu struct {
	Menus []Item     `json:"button"`
	Match *MatchRule `json:"matchrule,omitempty"`
}

type MatchRule struct {
	//there is at least one not null
	TagID              string `json:"tag_id,omitempty"`
	Sex                string `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

func (r *RootMenu) Type(e ev) {
	for _, i := range r.Menus {
		i.Type(e)
	}
}

type Item interface {
	Type(ev)
}

type SubMenu struct {
	Name  string `json:"name"`
	Menus []Item `json:"sub_button"`
}

func (s *SubMenu) Type(e ev) {
	for _, i := range s.Menus {
		i.Type(e)
	}
}

type ClickMenu struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	Key  string `json:"key"`

	On Callback `json:"-"`
}

func (c *ClickMenu) Type(e ev) {
	c.Typ = "click"
	if c.On != nil {
		e.On("event.CLICK.addr_"+c, c.On)
	}
}

type ViewMenu struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	//url or mediaID
	Url     string `json:"url"`
	MediaID string `json:"medio_id"` //if this field is not empty, then Url will be omitted

	On Callback `json:"-"`
}

func (v *ViewMenu) Type(e ev) {
	if v.MediaID != "" {
		v.Typ = "view_limited"
	}
	v.Typ = "view"
	if v.On != nil {
		e.On("event.VIEW."+v.Url, v.On)
	}
}

type MediaMenu struct {
	Typ string `json:"type"`

	Name    string `json:"name"`
	MediaID string `json:"media_id"`

	//nocallback
}

func (i *MediaMenu) Type(e ev) {
	i.Typ = "media_id"
}

type ProgramMenu struct {
	Typ string `json:"type"`

	Name     string `json:"name"`
	Url      string `json:"url"`
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`

	On Callback `json:"-"`
}

func (p *ProgramMenu) Type(e ev) {
	p.Typ = "miniprogram"
	if p.On != nil {
		e.On("event.view_miniprogram."+p.PagePath, p.On)
	}
}

type LocationMenu struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	Key  string `json:"key"`

	On Callback `json:"-"`
}

func (l *LocationMenu) Type(e ev) {
	l.Typ = "location_select"
	if l.On != nil {
		e.On("event.location_select.addr_"+l, l.On)
	}
}

type ScanMenu struct {
	Typ string `json:"type"`

	Name    string `json:"name"`
	Key     string `json:"key"`
	WithMsg bool   `json:"-"` //是否具有提示

	On Callback `json:"-"`
}

func (s *ScanMenu) Type(e ev) {
	if s.WithMsg {
		s.Typ = "scancode_waitmsg"
	}
	s.Typ = "scancode_push"
	if s.On != nil {
		e.On(fmt.Sprintf("event.%s.addr_"+s, s.Typ), s.On)
	}
}

type PicType int

const (
	PicSysPhoto   PicType = 1
	PicPhotoAlbum PicType = 2
	PicWx         PicType = 3
)

type PictureMenu struct {
	Typ string `json:"type"`

	Name     string `json:"name"`
	MenuType PicType
	Key      string `json:"key"`

	On Callback `json:"-"`
}

func (p *PictureMenu) Type(e ev) {
	switch p.MenuType {
	case PicSysPhoto:
		p.Typ = "pic_sysphoto"
	case PicPhotoAlbum:
		p.Typ = "pic_photo_or_album"
	case PicWx:
		p.Typ = "pic_weixin"
	}
	if p.On != nil {
		e.On(fmt.Sprintf("event.%s.addr_"+p, p.Typ), p.On)
	}
}
