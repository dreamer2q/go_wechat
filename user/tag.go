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
	reqTagAdd         = "cgi-bin/tags/create"
	reqTagGet         = "cgi-bin/tags/get"
	reqTagUpdate      = "cgi-bin/tags/update"
	reqTagDelete      = "cgi-bin/tags/delete"
	reqGetByTagID     = "cgi-bin/user/tag/get"
	reqTagToUsers     = "cgi-bin/tags/members/batchtagging"
	reqUntagToUsers   = "cgi-bin/tags/members/batchuntagging"
	reqGetTagFromUser = "cgi-bin/tags/getidlist"
)

func (u *User) AddTag(tagName string) (tagID int, err error) {
	posJson := fmt.Sprintf(`{"tag":{"name": %q }}`, tagName)
	var body []byte
	_, body, err = u.req.Post(reqTagAdd, nil, request.TypeJSON, strings.NewReader(posJson))
	if err != nil {
		return 0, errors.Wrap(err, "post")
	}
	if err = request.CheckCommonError(body); err != nil {
		return
	}
	ret := &struct {
		Tag struct {
			ID int `json:"id"`
		} `json:"tag"`
	}{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal")
	}
	return ret.Tag.ID, nil
}

func (u *User) GetTags() (*Tags, error) {
	_, body, err := u.req.Get(reqTagGet, nil)
	if err != nil {
		return nil, errors.Wrap(err, "get")
	}
	ret := &Tags{}
	if err = request.CheckCommonError(body); err != nil {
		return nil, errors.Wrap(err, "comment error")
	}
	err = json.Unmarshal(body, err)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	return ret, nil
}

func (u *User) UpdateTag(tagID int, tagName string) error {
	postJson := fmt.Sprintf(`{"tag":{"id": %d, "name": %q }}`, tagID, tagName)
	_, body, err := u.req.Post(reqTagUpdate, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (u *User) DeleteTag(tagID int) error {
	postJson := fmt.Sprintf(`{"tag":{"id": %d }}`, tagID)
	_, body, err := u.req.Post(reqTagDelete, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (u *User) GetUsersByTag(tagID int, nextOpenID ...string) (*List, error) {
	var postJson string
	if nextOpenID != nil {
		postJson = fmt.Sprintf(`{"tagid": %d, "next_openid":%q }`, tagID, nextOpenID[0])
	} else {
		postJson = fmt.Sprintf(`{"tagid": %d, "next_openid":"" }`, tagID)
	}
	_, body, err := u.req.Post(reqGetByTagID, nil, request.TypeJSON, strings.NewReader(postJson))
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

func (u *User) AddTagToUsers(tg *TagUsers) error {
	return u.doTagBatch(true, tg)
}
func (u *User) DelTagToUsers(tg *TagUsers) error {
	return u.doTagBatch(false, tg)
}
func (u *User) doTagBatch(addTag bool, tg *TagUsers) error {
	var (
		postBytes []byte
		body      []byte
		err       error
	)
	postBytes, err = json.Marshal(tg)
	if err != nil {
		return errors.Wrap(err, "marshal")
	}
	if addTag {
		_, body, err = u.req.Post(reqTagToUsers, nil, request.TypeJSON, bytes.NewReader(postBytes))
	} else {
		_, body, err = u.req.Post(reqUntagToUsers, nil, request.TypeJSON, bytes.NewReader(postBytes))
	}
	if err != nil {
		return errors.Wrap(err, "post")
	}
	return request.CheckCommonError(body)
}

func (u *User) GetTagsFromUser(OpenID string) (tagList []int, err error) {
	postJson := fmt.Sprintf(`{"openid": %q }`, OpenID)
	var body []byte
	_, body, err = u.req.Post(reqGetTagFromUser, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}
	if err = request.CheckCommonError(body); err != nil {
		return nil, errors.Wrap(err, "common error")
	}
	ret := &struct {
		TagList []int `json:"tagid_list"`
	}{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	return ret.TagList, nil
}


