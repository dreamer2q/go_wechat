package api

const (
	BaseURL = "https://api.weixin.qq.com/cgi-bin/"
)

type Config struct {
	//basic config
	AppID        string
	AppSecret    string
	AppToken     string

	//message encode
	AesEncodeKey string

	//server callback address
	Callback     string
}
