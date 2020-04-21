package menu

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

func (r *RootMenu) Type() {
	for _, i := range r.Menus {
		i.Type()
	}
}

//每一个实现了Item接口的结构体，如果有Key字段的话，需要设置以便用来区分具体的事件
type Item interface {
	Type()
}

type SubMenu struct {
	Name  string `json:"name"`
	Menus []Item `json:"sub_button"`
}

func (s *SubMenu) Type() {
	for _, i := range s.Menus {
		i.Type()
	}
}

type ClickMenu struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	Key  string `json:"key"`
}

func (c *ClickMenu) Type() {
	c.Typ = "click"
}

type ViewMenu struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	//url or mediaID
	Url     string `json:"url"`
	MediaID string `json:"medio_id"` //if this field is not empty, then Url will be omitted
}

func (v *ViewMenu) Type() {
	if v.MediaID != "" {
		v.Typ = "view_limited"
	}
	v.Typ = "view"
}

type MediaMenu struct {
	Typ string `json:"type"`

	Name    string `json:"name"`
	MediaID string `json:"media_id"`
}

func (i *MediaMenu) Type() {
	i.Typ = "media_id"
}

type ProgramMenu struct {
	Typ string `json:"type"`

	Name     string `json:"name"`
	Url      string `json:"url"`
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

func (p *ProgramMenu) Type() {
	p.Typ = "miniprogram"
}

type LocationMenu struct {
	Typ string `json:"type"`

	Name string `json:"name"`
	Key  string `json:"key"`
}

func (l *LocationMenu) Type() {
	l.Typ = "location_select"
}

type ScanMenu struct {
	Typ string `json:"type"`

	Name    string `json:"name"`
	Key     string `json:"key"`
	WithMsg bool   `json:"-"` //是否具有提示
}

func (s *ScanMenu) Type() {
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

type PictureMenu struct {
	Typ string `json:"type"`

	Name     string `json:"name"`
	MenuType PicType
	Key      string `json:"key"`
}

func (p *PictureMenu) Type() {
	switch p.MenuType {
	case PicSysPhoto:
		p.Typ = "pic_sysphoto"
	case PicPhotoAlbum:
		p.Typ = "pic_photo_or_album"
	case PicWx:
		p.Typ = "pic_weixin"
	}
}
