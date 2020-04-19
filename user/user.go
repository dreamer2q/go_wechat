package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dreamer2q/go_wechat/request"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

const (
	reqGetUsers     = "cgi-bin/user/get"
	reqSetRemark    = "cgi-bin/user/info/updateremark"
	reqGetUserInfo  = "cgi-bin/user/info"
	reqGetUsersInfo = "cgi-bin/user/info/batchget"
)

type User struct {
	req *request.Request
}

func New(r *request.Request) *User {
	return &User{
		req: r,
	}
}

type List struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenIDs []string `json:"openid"`
	}
	Next string `json:"next_openid"`
}

func (u *User) GetUsers(nextOpenID ...string) (*List, error) {
	params := url.Values{}
	if nextOpenID != nil {
		params.Add("next_openid", nextOpenID[0])
	}
	_, body, err := u.req.Get(reqGetUsers, params)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	if err = request.CheckCommonError(body); err != nil {
		return nil, errors.Wrap(err, "common error")
	}
	ret := &List{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	return ret, nil
}

func (u *User) SetRemark(openID string, remark string) error {
	postJson := fmt.Sprintf(`{"openid": %q ,"remark": %q }`, openID, remark)
	_, body, err := u.req.Post(reqSetRemark, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (u *User) GetUserInfo(openID string, lang ...Language) (*Info, error) {
	params := url.Values{}
	params.Add("openid", openID)
	if lang != nil {
		params.Add("lang", string(lang[0]))
	}
	_, body, err := u.req.Get(reqGetUserInfo, params)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	if err = request.CheckCommonError(body); err != nil {
		return nil, errors.Wrap(err, "common error")
	}
	ret := &Info{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	return ret, nil
}

func (u *User) GetUsersInfo(list *GetInfo) ([]Info, error) {
	postBytes, err := json.Marshal(list)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}
	_, body, err := u.req.Post(reqGetUsersInfo, nil, request.TypeJSON, bytes.NewReader(postBytes))
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}
	if err = request.CheckCommonError(body); err != nil {
		return nil, errors.Wrap(err, "common error")
	}
	ret := &struct {
		UserList []Info `json:"user_list"`
	}{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	return ret.UserList, nil
}
