package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"wechat/request"
)

const (
	reqMsgPreview    = "cgi-bin/message/mass/preview"
	reqMsgSendOpenID = "cgi-bin/message/mass/send"
	reqMsgSendFilter = "cgi-bin/message/mass/sendall"
	reqVideoConvert  = "cgi-bin/media/uploadvideo"
	reqMsgDelete     = "cgi-bin/message/mass/delete"
	reqMsgStatus     = "cgi-bin/message/mass/get"
	reqSpeedSet      = "cgi-bin/message/mass/speed/set"
	reqSpeedGet      = "cgi-bin/message/mass/speed/get"
)

type MsgWrapper interface {
	Type() string
}

type ToUser []string

var _ json.Marshaler = &ToUser{}

func (t *ToUser) MarshalJSON() ([]byte, error) {
	if len(*t) == 1 {
		return []byte(fmt.Sprintf("%q", (*t)[0])), nil
	}
	return json.Marshal(*t)
}

func (t *ToUser) Add(OpenID ...string) {
	if OpenID != nil {
		*t = append(*t, OpenID...)
	}
}

type Message struct {
	ToWxName      string     `json:"towxname,omitempty"` //available only in preview
	ToUser        *ToUser    `json:"touser,omitempty"`   //send by openIDs, minimum number is 2 and maximum number is 10000
	Filter        *MsgFilter `json:"filter,omitempty"`   //send by tag, is_to_all must be false
	MsgWrapper    MsgWrapper //MpNews(article), MpText, MpVoice, MpImage, MpVideo (video should be specially handled), MpCard
	IgnoreReprint int        `json:"send_ignore_reprint"` //图文消息被判定为转载时，是否继续群发。 1为继续群发（转载），0为停止群发。 该参数默认为0。

	ClientMsgId string `json:"clientmsgid"` //使用 clientmsgid 参数，避免重复推送

	msgType string `json:"msgtype"` //unused
}

type MsgFilter struct {
	IsToAll bool `json:"is_to_all"`
	TagID   int  `json:"tag_id"`
}

type MpNews struct {
	MediaID string `json:"media_id"`
}

func (MpNews) Type() string {
	return `mpnews`
}

type MpText struct {
	Content string `json:"content"`
}

func (MpText) Type() string {
	return `text`
}

type MpVoice struct {
	MediaID string `json:"media_id"`
}

func (MpVoice) Type() string {
	return `voice`
}

type MpImage struct {
	MediaIDs           []string `json:"media_ids"`
	Recommend          string   `json:"recommend"`
	NeedOpenComment    int      `json:"need_open_comment"`
	OnlyFansCanComment int      `json:"only_fans_can_comment"`
}

func (MpImage) Type() string {
	return `image`
}

type MpVideo struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (MpVideo) Type() string {
	return `mpvideo`
}

type MpCard struct {
	CardID string `json:"card_id"`
}

func (MpCard) Type() string {
	return `wxcard`
}

type Video struct {
	VideoDescription
	MpVideo
}

func (m *Media) VideoToMessage(v *Video) (*Result, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("VideoToMessage marshal: %v", err)
	}
	_, body, err := m.req.Post(reqVideoConvert, nil, request.TypeJSON, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("VideoToMessage post: %v", err)
	}
	return checkResult(body)
}

func (m *Media) SendByFilter(msg *Message) (*Result, error) {
	return m.doSend(reqMsgSendFilter, msg)
}

func (m *Media) SendByOpenID(msg *Message) (*Result, error) {
	return m.doSend(reqMsgSendOpenID, msg)
}

func (m *Media) SendPreview(msg *Message) error {
	_, err := m.doSend(reqMsgPreview, msg)
	return err
}
func (m *Media) doSend(req string, msg *Message) (*Result, error) {
	jsonReader := msgJsonBuilder(msg)
	//debug test
	fmt.Printf("debug jsonMessage: %s\n", jsonReader)
	_, body, err := m.req.Post(req, nil, request.TypeJSON, jsonReader)
	if err != nil {
		return nil, fmt.Errorf("doSend: %s %v", req, err)
	}
	ret := &Result{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("doSend unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("doSend: %d %s", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}

func (m *Media) DeleteMsg(MsgID int, articleIndex int) error {
	postJson := fmt.Sprintf(`{"msg_id":%d,"article_idx":%d}`, MsgID, articleIndex)
	_, body, err := m.req.Post(reqMsgDelete, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return fmt.Errorf("DeleteMsg post: %v", err)
	}
	_, err = checkResult(body)
	return err
}

func (m *Media) SendStatus(MsgID string) (*Result, error) {
	postJson := fmt.Sprintf(`{"msg_id":%q}`, MsgID)
	_, body, err := m.req.Post(reqMsgStatus, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return nil, fmt.Errorf("SendStatus post: %v", err)
	}
	return checkResult(body)
}

func (m *Media) GetSendSpeed() (*Result, error) {
	_, body, err := m.req.Post(reqSpeedGet, nil, request.TypeJSON, nil)
	if err != nil {
		return nil, fmt.Errorf("GetSendSpeed post: %v", err)
	}
	return checkResult(body)
}
func (m *Media) SetSendSpeed(speedLevel int) error {
	if speedLevel < 0 || speedLevel > 4 {
		return fmt.Errorf("seppdLevel out of range(0..4)")
	}
	postJson := fmt.Sprintf(`{"speed":%d}`, speedLevel)
	_, body, err := m.req.Post(reqSpeedSet, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return fmt.Errorf("SetSendSpeed post: %v", err)
	}
	_, err = checkResult(body)
	return err
}
