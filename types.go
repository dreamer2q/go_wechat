package wechat

import (
	"encoding/xml"
	"log"
	"reflect"
)

type MessageHandler func(m *MessageReceive) MessageReply

type Message interface {
	Type() string
}

type messageBase struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"` //sender, openID
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
}

func (m messageBase) Type() string {
	return m.MsgType
}
func (m messageBase) IsEvent() bool {
	return m.MsgType == "event"
}

type MessageReceive struct {
	//I do not know what are these for, just leave it untouched
	XMLName xml.Name `xml:"xml"`
	Text    string   `xml:",chardata"`

	//comment struct elements
	messageBase
	MsgId int64 `xml:"MsgId"`

	//MsgType: text
	Content string `xml:"Content"`
	//below types consist of it
	MediaId string `xml:"MediaId"`
	//MsgType: image
	PicUrl string `xml:"PicUrl"`

	//MsgType: voice
	Format      string `xml:"Format"`      //amr, speex
	Recognition string `xml:"Recognition"` //available when voice recognition is on

	//MsgType: video, shortvideo
	ThumbMediaId string `xml:"ThumbMediaId"`

	//MsgType: location
	LocationX string `xml:"Location_X"`
	LocationY string `xml:"Location_Y"`
	Scale     string `xml:"Scale"`
	Label     string `xml:"Label"`

	//MsgType: link
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	URL         string `xml:"Url"`

	///////////////////////////
	//MsgType: event
	Event string `xml:"Event"` //subscribe, unsubscribe

	//Event: subscript(first time to subscribe), SCAN(has subscribed)
	//Event: CLICK, VIEW
	EventKey string `xml:"EventKey"`
	Ticket   string `xml:"Ticket"`

	//Event: LOCATION(upload position event)
	Latitude  string `xml:"Latitude"`
	Longitude string `xml:"Longitude"`
	Precision string `xml:"Precision"`
}

type MessageReply interface {
	Type() string
}

type messageReply struct {
	XMLName struct{} `xml:"xml"`
	messageBase
	MsgWrapper messageWrapper
}
type messageWrapper struct {
	Msg Message
}

func (m messageWrapper) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if article, ok := m.Msg.(Articles); ok {
		err := e.EncodeElement(len(article.Items), xml.StartElement{Name: xml.Name{Local: "ArticleCount"}})
		if err != nil {
			log.Panic(err)
		}
	}
	return e.EncodeElement(m.Msg, xml.StartElement{
		Name: xml.Name{Local: reflect.TypeOf(m.Msg).Name()},
	})
}

var _ xml.Marshaler = &messageWrapper{}

type Text struct {
	Content string `xml:"Content"`
}

func (t Text) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(t.Content, xml.StartElement{Name: xml.Name{Local: "Content"}})
}

var _ xml.Marshaler = Text{}

func (Text) Type() string {
	return "text"
}

type Image struct {
	MediaId string `xml:"MediaId"`
}

func (Image) Type() string {
	return "image"
}

type Voice struct {
	MediaID string `xml:"MediaId"`
}

func (Voice) Type() string {
	return "voice"
}

type Video struct {
	MediaID     string `xml:"MediaId"`
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
}

func (Video) Type() string {
	return "video"
}

type Music struct {
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	MusicURL    string `xml:"MusicUrl"`
	HqURL       string `xml:"HQMusicUrl"`
	ThumbID     string `xml:"ThumbMediaId"`
}

func (Music) Type() string {
	return "music"
}

type ArticleItem struct {
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	PicURL      string `xml:"PicUrl"`
	URL         string `xml:"Url"`
}

type Articles struct {
	Items []ArticleItem `xml:"item"`
}

func (Articles) Type() string {
	return "news"
}
