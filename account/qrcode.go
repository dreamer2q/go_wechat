package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"wechat/request"
)

const (
	reqCreateQrCode = "cgi-bin/qrcode/create"
)

type Account struct {
	req *request.Request
}

func New(r *request.Request) *Account {
	return &Account{
		req: r,
	}
}

func (a *Account) CreateQrCode(permanent bool, ExpireIn int, sceneID int, sceneStr ...string) (*QrResult, error) {
	qrQ := qrQuery{}
	if permanent {
		qrQ.ActionName = "QR_LIMIT_SCENE"
	} else {
		qrQ.ActionName = "QR_SCENE"
		qrQ.ExpireIn = ExpireIn
	}
	if sceneStr != nil {
		qrQ.ActionInfo.Scene.SceneStr = sceneStr[0]
	} else {
		qrQ.ActionInfo.Scene.SceneID = sceneID
	}
	postBytes, err := json.Marshal(&qrQ)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}
	_, body, err := a.req.Post(reqCreateQrCode, nil, request.TypeJSON, bytes.NewReader(postBytes))
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}
	if err = request.CheckCommonError(body); err != nil {
		return nil, errors.Wrap(err, "common error")
	}
	ret := &QrResult{}
	err = json.Unmarshal(body, ret)
	return ret, errors.Wrap(err, "unmarshal")
}

func (a *Account) GetQrURL(ticket string) string {
	params := url.Values{}
	params.Add("ticket", ticket)
	return fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?%s", params.Encode())
}
