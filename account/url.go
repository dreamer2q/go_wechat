package account

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"wechat/request"
)

const (
	reqShortenUrl = "cgi-bin/shorturl"
)

func (a *Account) ShortenUrl(longUrl string) (shortUrl string, err error) {
	postJson := fmt.Sprintf(`{"action":"long2short","long_url": %q }`, longUrl)
	_, body, err := a.req.Post(reqShortenUrl, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return "", errors.Wrap(err, "post")
	}
	if err = request.CheckCommonError(body); err != nil {
		return "", errors.Wrap(err, "common error")
	}
	ret := &struct {
		ShortUrl string `json:"short_url"`
	}{}
	err = json.Unmarshal(body, ret)
	return ret.ShortUrl, err
}
