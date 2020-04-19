package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"wechat/request"
)

const (
	reqGetBlackList = "cgi-bin/tags/members/getblacklist"
	reqBlackUsers   = "cgi-bin/tags/members/batchblacklist"
	reqUnBlackUsers = "cgi-bin/tags/members/batchunblacklist"
)

func (u *User) GetBlacklist(beginOpenID ...string) (*List, error) {
	var postJson string
	if beginOpenID != nil {
		postJson = fmt.Sprintf(`{"begin_openid": %q }`, beginOpenID[0])
	} else {
		postJson = fmt.Sprintf(`{"begin_openid": "" }`)
	}
	_, body, err := u.req.Post(reqGetBlackList, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return nil, errors.Wrap(err, "post")
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

func (u *User) doBlackUsers(black bool, OpenIDs []string) error {
	postStruct := &struct {
		OpenIDList []string `json:"openid_list"`
	}{
		OpenIDList: OpenIDs,
	}
	postBytes, err := json.Marshal(postStruct)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}
	var (
		body []byte
	)
	if black {
		_, body, err = u.req.Post(reqBlackUsers, nil, request.TypeJSON, bytes.NewReader(postBytes))
	} else {
		_, body, err = u.req.Post(reqUnBlackUsers, nil, request.TypeJSON, bytes.NewReader(postBytes))
	}
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}
func (u *User) BlackUsers(OpenIDs []string) error {
	return u.doBlackUsers(true, OpenIDs)
}

func (u *User) UnBlackUsers(OpenIDs []string) error {
	return u.doBlackUsers(false, OpenIDs)
}
