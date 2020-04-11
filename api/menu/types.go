package menu

type menuResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type RootMenu struct {
	Buttons []Item `json:"button"`
}

func (r *RootMenu) Type() {
	for _, i := range r.Buttons {
		i.Type()
	}
}

type Item interface {
	Type()
}

type SubMenu struct {
	Name  string `json:"name"`
	Items []Item `json:"sub_button"`
}

func (s *SubMenu) Type() {
	for _, i := range s.Items {
		i.Type()
	}
}

type Click struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	Key  string `json:"key"`
}

func (c *Click) Type() {
	c.Typ = "click"
}

type View struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	//url or mediaID
	Url     string `json:"url"`
	MediaID string `json:"medio_id"` //if this field is not empty, then Url will be omitted
}

func (v *View) Type() {
	if v.MediaID != "" {
		v.Typ = "view_limited"
	}
	v.Typ = "view"
}

type Image struct {
	Typ string `json:"type"`

	Name    string `json:"name"`
	MediaID string `json:"media_id"`
}

func (i *Image) Type() {
	i.Typ = "media_id"
}

type Program struct {
	Typ string `json:"type"`

	Name     string `json:"name"`
	Url      string `json:"url"`
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

func (p *Program) Type() {
	p.Typ = "miniprogram"
}

type Location struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	Key  string `json:"key"`
}

func (l *Location) Type() {
	l.Typ = "location_select"
}

type Scan struct {
	Typ string `json:"type"`

	Name    string `json:"name"`
	Key     string `json:"key"`
	WithMsg bool   `json:"-"` //是否具有提示
}

func (s *Scan) Type() {
	if s.WithMsg {
		s.Typ = "scancode_waitmsg"
	}
	s.Typ = "scancode_push"
}

type PicType int

const (
	PicSysPhoto   PicType = 1
	PicPhotoAlbum PicType = 2
	PicWx         PicType = 3
)

type Picture struct {
	Typ string `json:"type"`

	Name     string `json:"name"`
	MenuType PicType
	Key      string `json:"key"`
}

func (p *Picture) Type() {
	switch p.MenuType {
	case PicSysPhoto:
		p.Typ = "pic_sysphoto"
	case PicPhotoAlbum:
		p.Typ = "pic_photo_or_album"
	case PicWx:
		p.Typ = "pic_weixin"
	default:
		p.Typ = "error_type"
	}
}
