package media

import (
	"../request"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/url"
	"strings"
)

const (
	reqTmpUpload     = "cgi-bin/media/upload"
	reqTmpGet        = "cgi-bin/media/get"
	reqUpload        = "cgi-bin/material/add_material"
	reqUploadNews    = "cgi-bin/material/add_news"
	reqGet           = "cgi-bin/material/get_material"
	reqUploadImage   = "cgi-bin/media/uploadimg"
	reqDelete        = "cgi-bin/material/del_material"
	reqArticleUpdate = "cgi-bin/material/update_news"
	reqMediaCount    = "cgi-bin/material/get_materialcount"
	reqMaterialList  = "cgi-bin/material/batchget_material"

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
func (m *Media) uploadImageForArticle(rawUri string) (retUrl string, err error) {
	var (
		bodyBytes  []byte
		bodyBuffer = &bytes.Buffer{}
	)
	if strings.HasPrefix(rawUri, "http") {
		bodyBytes, err = DownFile(rawUri)
	} else {
		bodyBytes, err = ioutil.ReadFile(rawUri)
	}
	if err != nil {
		return
	}
	fileForm, err := newFileForm("media", rawUri, bytes.NewReader(bodyBytes), bodyBuffer)
	if err != nil {
		return
	}
	//close before sending any request
	_ = fileForm.Close()
	var (
		body []byte
	)
	_, body, err = m.req.Post(reqUploadImage, nil, fileForm.FormDataContentType(), bodyBuffer)
	if err != nil {
		return
	}
	ret := &Result{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return
	}
	if ret.ErrCode != 0 {
		err = fmt.Errorf("%d %v", ret.ErrCode, ret.ErrMsg)
		return
	}
	retUrl = ret.Url
	return
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
func (m *Media) UploadMaterial(name string, in io.Reader, permanent bool, typ string, description ...VideoDescription) (*Result, error) {
	params := url.Values{}
	params.Add("type", typ)
	postBody := &bytes.Buffer{}
	post, err := newFileForm("media", name, in, postBody)
	if err != nil {
		return nil, fmt.Errorf("newFileForm: %v", err)
	}
	if typ == TypVideo {
		if description == nil {
			return nil, fmt.Errorf("missing description for type Video")
		}
		videoWriter, err := post.CreateFormField("description")
		if err != nil {
			return nil, fmt.Errorf("create description: %v", err)
		}
		//assume no error for Marshal
		videoJson, _ := json.Marshal(description[0])
		_, _ = videoWriter.Write(videoJson)
	}
	//close before sending any http request
	_ = post.Close()
	var body []byte
	if permanent {
		_, body, err = m.req.Post(reqUpload, params, post.FormDataContentType(), postBody)
	} else {
		_, body, err = m.req.Post(reqTmpUpload, params, post.FormDataContentType(), postBody)
	}
	if err != nil {
		return nil, fmt.Errorf("post: %v", err)
	}
	ret := &Result{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("response: %d %s", ret.ErrCode, ret.ErrMsg)
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
	if attach != "" || ret.ContentType != contentPlain { //body is the file
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

//can not delete temporary material
func (m *Media) Delete(mediaID string) error {
	postBody := fmt.Sprintf(`{"media":%q}`, mediaID)
	_, body, err := m.req.Post(reqDelete, nil, request.TypeJSON, strings.NewReader(postBody))
	if err != nil {
		return err
	}
	ret := &Result{}
	err = json.Unmarshal(body, ret)
	if ret.ErrCode != 0 {
		return fmt.Errorf("%d %s", ret.ErrCode, ret.ErrMsg)
	}
	return nil
}

func (m *Media) UpdateArticle(article *ArticleWrapper) error {
	jsonBytes, err := json.Marshal(article)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}
	_, body, err := m.req.Post(reqArticleUpdate, nil, request.TypeJSON, bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("post: %v", err)
	}
	ret := Result{}
	err = json.Unmarshal(body, &ret)
	if ret.ErrCode != 0 {
		return fmt.Errorf("response: %d %s", ret.ErrCode, ret.ErrMsg)
	}
	return nil
}

func (m *Media) MaterialCounter() (*MaterialCounter, error) {
	_, body, err := m.req.Get(reqMediaCount, nil)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}
	ret := &MaterialCounter{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)

	}
	return ret, nil
}

func (m *Media) GetMaterialList(typ string, offset int, count int) (*MaterialList, error) {
	postJson := fmt.Sprintf(`{"type":%q,"offset":"%d","count":"%d"}`, typ, offset, count)
	_, body, err := m.req.Post(reqMaterialList, nil, request.TypeJSON, strings.NewReader(postJson))
	if err != nil {
		return nil, fmt.Errorf("post: %v", err)
	}
	ret := &MaterialList{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return nil, fmt.Errorf("reponse: %d %v", ret.ErrCode, ret.ErrMsg)
	}
	return ret, nil
}
