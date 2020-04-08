package holder

import "fmt"

type response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

var codeMap = map[int]string{
	-1:    "系统繁忙，此时请开发者稍候再试",
	0:     "请求成功",
	40001: "AppSecret错误或者AppSecret不属于这个公众号，请开发者确认AppSecret的正确性",
	40002: "请确保grant_type字段值为client_credential",
	40164: "调用接口的IP地址不在白名单中，请在接口IP白名单中进行设置。（小程序及小游戏调用不要求IP地址在白名单内。）",
	89503: "此IP调用需要管理员确认,请联系管理员",
	89501: "此IP正在等待管理员确认,请联系管理员",
	89506: "24小时内该IP被管理员拒绝调用两次，24小时内不可再使用该IP调用",
	89507: "1小时内该IP被管理员拒绝调用一次，1小时内不可再使用该IP调用",
}

func (r response) IsError() bool {
	return r.ErrCode != 0
}

func (r response) Error() string {
	return fmt.Sprintf("response: %d %s", r.ErrCode, r.ErrMsg)
}

func (r response) Desc() string {
	return fmt.Sprintf("response: %d %s", r.ErrCode, codeMap[r.ErrCode])
}
