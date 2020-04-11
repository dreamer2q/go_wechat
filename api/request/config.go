package request

import "time"

type Config struct {
	//basic config
	AppID     string
	AppSecret string
	AppToken  string

	//message encode
	AesEncodeKey string

	//server callback address
	Callback string
	Timeout  time.Duration
}
