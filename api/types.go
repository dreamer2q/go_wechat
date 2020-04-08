package api

import (
	"encoding/xml"
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
	XMLName xml.Name
	messageBase
	Msg interface{}
}

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

type NoReply struct{}

func (NoReply) Type() string {
	return "noreply"
}
