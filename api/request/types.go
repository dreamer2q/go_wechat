package request

const (
	BaseURL     = "https://api.weixin.qq.com/"
	BaseURLBack = "https://api2.weixin.qq.com/"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

// Types we support.
const (
	TypeJSON       = "json"
	TypeXML        = "xml"
	TypeUrlencoded = "urlencoded"
	TypeForm       = "form"
	TypeFormData   = "form-data"
	TypeHTML       = "html"
	TypeText       = "text"
	TypeMultipart  = "multipart"
)

var types = map[string]string{
	TypeJSON:       "application/json",
	TypeXML:        "application/xml",
	TypeForm:       "application/x-www-form-urlencoded",
	TypeFormData:   "application/x-www-form-urlencoded",
	TypeUrlencoded: "application/x-www-form-urlencoded",
	TypeHTML:       "text/html",
	TypeText:       "text/plain",
	TypeMultipart:  "multipart/form-data",
}

func typeMapper(key string) string {
	if t, ok := types[key]; ok {
		return t
	}
	return key
}
