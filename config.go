package wechat

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

	//request timeout
	Timeout time.Duration

	debug bool
}

func (c *Config) SetDebug(d bool) {
	c.debug = d
}
