package request

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	token  *token
	client *http.Client
}

func New(c *Config) *Request {
	return &Request{
		token: &token{
			Config: c,
		},
		client: &http.Client{
			Timeout: c.Timeout,
		},
	}
}

func (r *Request) Get(endpoint string, params url.Values) (resp *http.Response, body []byte, err error) {
	var reqUrl string
	reqUrl, err = r.formatter(endpoint, params)
	if err != nil {
		return
	}
	resp, err = r.client.Get(reqUrl)
	if err != nil {
		return
	}
	return readBody(resp)
}

func (r *Request) Post(endpoint string, params url.Values, contentType string, bodyReader io.Reader) (resp *http.Response, body []byte, err error) {
	var reqUrl string
	reqUrl, err = r.formatter(endpoint, params)
	if err != nil {
		return
	}
	resp, err = r.client.Post(reqUrl, contentType, bodyReader)
	if err != nil {
		return
	}
	return readBody(resp)
}

// format url with access token
func (r *Request) formatter(endpoint string, params url.Values) (string, error) {
	access := r.token.Get()
	if access == "" {
		return "", r.token.Err
	}
	if params == nil {
		params = url.Values{}
	}
	params.Add("access_token", access)
	return fmt.Sprintf("%s%s?%s", BaseURL, endpoint, params.Encode()), nil
}

func (r *Request) Hijack() *http.Client {
	return r.client
}

func readBody(r *http.Response) (resp *http.Response, body []byte, err error) {
	resp = r
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = r.Body.Close()
	return
}
