package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"wechat/api/request"
)

//TODO comment related APIs

const (
	reqCommentOpen        = "cgi-bin/comment/open"
	reqCommentClose       = "cgi-bin/comment/close"
	reqCommentList        = "cgi-bin/comment/list"
	reqCommentMarkElect   = "cgi-bin/comment/markelect"
	reqCommentUnmarkElect = "cgi-bin/comment/unmarkelect"
	reqCommentDel         = "cgi-bin/comment/delete"
	reqCommentDelReply    = "cgi-bin/comment/reply/delete"
	reqCommentAddReply    = "cgi-bin/comment/reply/add"
)

type GetComment struct {
	MsgID string      `json:"msg_data_id"`
	Index uint32      `json:"index"`
	Begin uint32      `json:"begin"`
	Count uint32      `json:"count"`
	Type  CommentType `json:"type"`
}
type CommentType uint32

const (
	CTypeBoth     CommentType = 0
	CTypeCommon   CommentType = 1
	CTypeSelected CommentType = 2
)

type CommentResult struct {
	request.CommonError
	Total    uint32 `json:"total"`
	Comments []struct {
		UserCommentID string `json:"user_comment_id"`
		OpenID        string `json:"openid"`
		CreateTime    int64  `json:"create_time"`
		Content       string `json:"content"`
		Reply         *struct {
			Content    string `json:"content"`
			CreateTime int64  `json:"create_time"`
		} `json:"reply,omitempty"`
	} `json:"comment"`
}

func (m *Media) CommentList(req *GetComment) (*CommentResult, error) {
	postBytes, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}
	_, body, err := m.req.Post(reqCommentList, nil, request.TypeJSON, bytes.NewReader(postBytes))
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}
	ret := &CommentResult{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("%d %v", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}

func (m *Media) CommentMark(MsgID uint32, UserID uint32, index ...uint32) error {
	return m.doCommentSelect(true, MsgID, UserID, index...)
}
func (m *Media) CommentUnmark(MsgID uint32, UserID uint32, index ...uint32) error {
	return m.doCommentSelect(false, MsgID, UserID, index...)
}
func (m *Media) doCommentSelect(mark bool, MsgID uint32, UserID uint32, index ...uint32) error {
	var (
		i    uint32
		body []byte
		err  error
	)
	if index != nil {
		i = index[0]
	}
	postJson := fmt.Sprintf(`{"msg_data_id":%d ,"index": %d, "user_comment_id": %d}`, MsgID, i, UserID)
	if mark {
		_, body, err = m.req.Post(reqCommentMarkElect, nil, request.TypeJSON, strings.NewReader(postJson))
	} else {
		_, body, err = m.req.Post(reqCommentUnmarkElect, nil, request.TypeJSON, strings.NewReader(postJson))
	}
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (m *Media) CommentOpen(mediaID string, index ...int) error {
	return m.doCommentSwitch(true, mediaID, index...)
}
func (m *Media) CommentClose(mediaID string, index ...int) error {
	return m.doCommentSwitch(false, mediaID, index...)
}
func (m *Media) doCommentSwitch(open bool, mediaID string, index ...int) error {
	var (
		postJson string
		body     []byte
		err      error
	)
	if index != nil {
		postJson = fmt.Sprintf(`{"msg_data_id":%q, "index":%q }`, mediaID, index[0])
	} else {
		postJson = fmt.Sprintf(`{"msg_data_id":%q}`, mediaID)
	}
	if open {
		_, body, err = m.req.Post(reqCommentOpen, nil, request.TypeJSON, strings.NewReader(postJson))
	} else {
		_, body, err = m.req.Post(reqCommentClose, nil, request.TypeJSON, strings.NewReader(postJson))
	}
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (m *Media) CommentDe(isReply bool, MsgID uint32, CommentID uint32, index ...uint32) error {
	var (
		postJson string
		body     []byte
		err      error
	)
	if index != nil {
		postJson = fmt.Sprintf(`{"msg_data_id":%d ,"index":%d ,"user_comment_id":%d }`, MsgID, index[0], CommentID)
	} else {
		postJson = fmt.Sprintf(`{"msg_data_id":%d ,"user_comment_id":%d }`, MsgID, CommentID)
	}
	if isReply {
		_, body, err = m.req.Post(reqCommentDelReply, nil, request.TypeJSON, strings.NewReader(postJson))
	} else {
		_, body, err = m.req.Post(reqCommentDel, nil, request.TypeJSON, strings.NewReader(postJson))
	}
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (m *Media) CommentReply(msgID uint32, commentID uint32, content string, index ...uint32) error {
	var (
		postJson string
	)
	if index != nil {
		postJson = fmt.Sprintf(`{"msg_data_id": %d, "index":%d, "user_comment_id": %d, "content":%q }`, msgID, index[0], commentID, content)
	} else {
		postJson = fmt.Sprintf(`{"msg_data_id": %d, "user_comment_id": %d, "content":%q }`, msgID, commentID, content)
	}
	_, body, err := m.req.Post(reqCommentAddReply, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}
