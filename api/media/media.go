package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"net/url"
	"strings"
	"wechat/api/request"
)

const (
	reqTmpUpload  = "cgi-bin/media/upload"
	reqTmpGet     = "cgi-bin/media/get"
	reqUpload     = "cgi-bin/material/add_material"
	reqUploadNews = "cgi-bin/material/add_news"
	reqGet        = "cgi-bin/material/get_material"

	TypImage = "image"
	TypVoice = "voice"
	TypVideo = "video"
	TypThumb = "thumb"

	contentPlain = "text/plain"
)

type Media struct {
	req *request.Request
}

func New(r *request.Request) *Media {
	return &Media{
		req: r,
	}
}

//Article is the type of news in wechat development document (permanent material)
func (m *Media) UploadArticle(article *ArticleWrapper) (*Result, error) {

	//todo support markdown and automatic upload images (uploadImageForArticle) with it
	jsonBody, err := json.Marshal(article)
	fmt.Printf("debug jsonBody: %s\n", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("UploadArticle marshal: %v", err)
	}
	_, body, err := m.req.Post(reqUploadNews, nil, request.TypeJSON, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("UploadArticle post: %v", err)
	}
	ret := &Result{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("UploadArticle unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("UplaodArticle: %d %s", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}

//used to upload pictures in the Article (to satisfy wechat requirement)
func (m *Media) uploadImageForArticle() {

}

//
//1、临时素材media_id是可复用的。
//2、媒体文件在微信后台保存时间为3天，即3天后media_id失效。
//3、上传临时素材的格式、大小限制与公众平台官网一致。
//图片（image）: 2M，支持PNG\JPEG\JPG\GIF格式
//语音（voice）：2M，播放长度不超过60s，支持AMR\MP3格式
//视频（video）：10MB，支持MP4格式
//缩略图（thumb）：64KB，支持JPG格式
//
func (m *Media) UploadMaterial(filename string, typ string, permanent bool, description ...VideoDescription) (*Result, error) {
	params := url.Values{}
	params.Add("type", typ)
	postBody := &bytes.Buffer{}
	if typ == TypVideo && description == nil || typ != TypVideo && description != nil {
		return nil, fmt.Errorf("upload: video type requires description field")
	}
	post, err := newFilePost("media", filename, postBody, description...)
	if err != nil {
		return nil, fmt.Errorf("uploadtmp filepost: %v", err)
	}
	reqUrl := reqTmpUpload
	if permanent {
		reqUrl = reqUpload
	}
	_, body, err := m.req.Post(reqUrl, params, post.FormDataContentType(), postBody)
	if err != nil {
		return nil, fmt.Errorf("uploadtmp post: %v", err)
	}
	ret := &Result{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("uploadtmp unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("uploadtmp ret: %d %s", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}

func (m *Media) GetMaterial(mediaID string) (*Response, error) {
	postParam := fmt.Sprintf(`{"media_id":%q}'`, mediaID) //It might be clear not to use json.Marshal
	fmt.Printf("Debug: postParam %s", postParam)
	resp, body, err := m.req.Post(reqGet, nil, request.TypeJSON, strings.NewReader(postParam))
	if err != nil {
		return nil, fmt.Errorf("GetMaterial: %v", err)
	}
	attach := resp.Header.Get("Content-disposition") //check if other material returned directly
	ret := &Response{}
	ret.ContentType = resp.Header.Get("Content-Type")
	if attach != "" || ret.ContentType == contentPlain { //file returned directly
		_, query, err := mime.ParseMediaType(attach)
		if err != nil {
			return nil, fmt.Errorf("GetMaterial parseMediaType: %v", err)
		}
		ret.Filename = query["filename"]
		ret.Data = body
	} else {
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, fmt.Errorf("GetMaterial unmarshal: %v", err)
		}
		if ret.ErrCode != 0 {
			return nil, fmt.Errorf("GetMaterial: %d %s", ret.ErrCode, ret.ErrMsg)
		}
	}
	return ret, nil
}

func (m *Media) GetTmpMaterial(mediaID string) (*Response, error) {
	params := url.Values{}
	params.Add("media_id", mediaID)
	resp, body, err := m.req.Get(reqTmpGet, params)
	if err != nil {
		return nil, fmt.Errorf("getemp: %v", err)
	}
	attach := resp.Header.Get("Content-disposition")
	ret := &Response{}
	ret.ContentType = resp.Header.Get("Content-Type")
	if attach != "" || ret.ContentType == contentPlain { //body is the file
		_, query, err := mime.ParseMediaType(attach)
		if err != nil {
			return nil, fmt.Errorf("getTmp parseMediaType: %v", err)
		}
		ret.Filename = query["filename"]
		ret.Data = body
	} else { //need to download or error occurs
		err = json.Unmarshal(body, ret)
		if err != nil {
			return nil, fmt.Errorf("getTmp unmarshal: %v", err)
		}
		if ret.ErrCode != 0 {
			return nil, fmt.Errorf("getTmp: %d %v", ret.ErrCode, ret.ErrMsg)
		}
	}
	return ret, nil
}
