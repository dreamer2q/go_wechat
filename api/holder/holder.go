package holder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://api.weixin.qq.com/cgi-bin/"
)

type Holder struct {
	Token     string
	startTime time.Time
	expireIn  int64

	config *Config
	err    error
}

type Config struct {
	AppId     string
	AppSecret string
}

func New(c *Config) *Holder {
	return &Holder{
		config: c,
	}
}

func (h Holder) String() string {
	return h.Token
}

func (h Holder) IsExpired() bool {
	return time.Now().Unix() >= (h.startTime.Unix() + h.expireIn)
}

func (h *Holder) Refresh() error {
	param := url.Values{}
	param.Add("grant_type", "client_credential")
	param.Add("appid", h.config.AppId)
	param.Add("secret", h.config.AppSecret)

	resp, err := h.get("token", param)
	if err != nil {
		h.err = err
		return err
	}
	if resp.IsError() {
		h.err = resp
		return h.err
	}
	h.err = nil
	h.Token = resp.AccessToken
	h.expireIn = resp.ExpiresIn
	h.startTime = time.Now()
	return nil
}

func (h *Holder) GetToken() string {
	if h.IsExpired() {
		if err := h.Refresh(); err != nil {
			return ""
		}
	}
	return h.Token
}

func (h *Holder) get(endpoint string, param url.Values) (*response, error) {
	furl := fmt.Sprintf("%s%s?%s", baseURL, endpoint, param.Encode())
	resp, err := http.Get(furl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get: %s %s:%s", endpoint, resp.StatusCode, resp.Status)
	}
	ret := &response{}
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
