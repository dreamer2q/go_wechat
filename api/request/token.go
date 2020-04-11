package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type token struct {
	Access   string
	ExpireIn int64
	Err      error
	Start    int64

	*Config
}

func (t *token) Get() string {
	if t.IsError() || t.IsExpired() {
		t.Err = t.Refresh()
		if t.Err != nil {
			return ""
		}
	}
	return t.Access
}

func (t *token) IsError() bool {
	return t.Err != nil
}

func (t *token) IsExpired() bool {
	return (time.Now().Unix() - t.Start) > t.ExpireIn
}

func (t *token) Refresh() error {
	url := fmt.Sprintf("%s%s?grant_type=client_credential&appid=%s&secret=%s", BaseURL, "cgi-bin/token", t.AppID, t.AppSecret)
	req, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("token get: %v", err)
	}
	defer req.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("token read: %v", err)
	}
	resp := &tokenResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return fmt.Errorf("token json: %v", err)
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("token resp: %d %s", resp.ErrCode, resp.ErrMsg)
	}
	t.Access = resp.AccessToken
	t.ExpireIn = resp.ExpiresIn
	t.Start = time.Now().Unix()
	return nil
}
