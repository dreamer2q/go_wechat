package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"wechat/request"
)

const (
	reqSetIndustry      = "cgi-bin/template/api_set_industry"
	reqGetIndustry      = "cgi-bin/template/get_industry"
	reqGetTemplateByID  = "cgi-bin/template/api_add_template"
	reqGetTemplate      = "cgi-bin/template/get_all_private_template"
	reqDelTemplate      = "cgi-bin/template/del_private_template"
	reqSendMsg          = "cgi-bin/message/template/send"
	reqGetAutoReplyRule = "cgi-bin/get_current_autoreply_info"
)

type Template struct {
	req *request.Request
}

func New(r *request.Request) *Template {
	return &Template{
		req: r,
	}
}

func (t *Template) SetIndustry(industry *IndustryType) error {
	postBytes, _ := json.Marshal(industry)

	_, body, err := t.req.Post(reqSetIndustry, nil, request.TypeJSON, bytes.NewReader(postBytes))
	if err != nil {
		return fmt.Errorf("SetIndustry post: %v", err)
	}
	return request.CheckCommonError(body)
}

func (t *Template) GetIndustry() (*IndustryResult, error) {
	_, body, err := t.req.Get(reqGetIndustry, nil)
	if err != nil {
		return nil, fmt.Errorf("GetIndustry get: %v", err)
	}
	err = request.CheckCommonError(body)
	if err != nil {
		return nil, err
	}
	ret := &IndustryResult{}
	err = json.Unmarshal(body, ret)
	return ret, nil
}

func (t *Template) GetTemplateByShortId(shortId string) (string, error) {
	reqJson := fmt.Sprintf(`{"template_id_short":%q}`, shortId)
	_, body, err := t.req.Post(reqGetTemplateByID, nil, request.TypeJSON, strings.NewReader(reqJson))
	if err != nil {
		return "", errors.Wrap(err, "post")
	}
	ret := &struct {
		request.CommonError
		TemplateID string `json:"template_id"`
	}{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return "", err
	}
	if ret.ErrCode != 0 {
		return "", errors.Errorf("%d %v", ret.ErrCode, ret.ErrMsg)
	}
	return ret.TemplateID, nil
}

func (t *Template) GetTemplates() (*TemplateResult, error) {
	_, body, err := t.req.Get(reqGetTemplate, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	ret := &TemplateResult{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if ret.ErrCode != 0 {
		return nil, errors.Errorf("%d %v", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}

func (t *Template) DeleteTemplate(templateID string) error {
	reqPost := fmt.Sprintf(`{"template_id":%q}`, templateID)
	_, body, err := t.req.Post(reqDelTemplate, nil, request.TypeJSON, strings.NewReader(reqPost))
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

//for more details information, read below
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html#0
//https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Operation_Specifications.html
//xml event will come after
func (t *Template) Send(tempMsg *TemplateMsg) (msgId int, err error) {
	jsonBytes, err := json.Marshal(tempMsg)
	if err != nil {
		return 0, errors.Wrap(err, "marshal")
	}
	_, body, err := t.req.Post(reqSendMsg, nil, request.TypeJSON, bytes.NewReader(jsonBytes))
	if err != nil {
		return 0, errors.Wrap(err, "post")
	}
	ret := &struct {
		request.CommonError
		MsgID int `json:"msgid"`
	}{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal")
	}
	if ret.ErrCode != 0 {
		return 0, errors.Errorf("%d %v", ret.ErrCode, ret.ErrMsg)
	}
	return ret.MsgID, nil
}

func (t *Template) GetAutoReplyRule() (*AutoReplySetting, error) {
	_, body, err := t.req.Get(reqGetAutoReplyRule, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	ret := &AutoReplySetting{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	if ret.ErrCode != 0 {
		return nil, errors.Errorf("%d %v", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}
