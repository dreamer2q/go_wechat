package menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"wechat/api/request"
)

const (
	ReqCreate       = "cgi-bin/menu/create"
	ReqGet          = "cgi-bin/get_current_selfmenu_info"
	ReqDelete       = "cgi-bin/menu/delete"
	ReqCustomCreate = "cgi-bin/menu/addconditional"
	ReqCustomDelete = "cgi-bin/menu/delconditional"
	ReqTryMatch     = "cgi-bin/menu/trymatch"
	ReqCustomGet    = "cgi-bin/menu/get"
)

type Menu struct {
	req *request.Request
}

func New(r *request.Request) *Menu {
	return &Menu{
		req: r,
	}
}

func (m *Menu) Create(menu RootMenu) error {
	menu.Type() //init type data for json marshal
	jsonData, err := json.Marshal(&menu)
	if err != nil {
		return fmt.Errorf("menu marshal: %v", err)
	}
	//start debug
	fmt.Printf("menu json: %s\n", jsonData)
	//end debug
	gorequest.New()
	_, body, err := m.req.Post(ReqCreate, nil, request.TypeJSON, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("menu create: %v", err)
	}
	ret := menuResp{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return fmt.Errorf("menu unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return fmt.Errorf("menu result: %d %s", ret.ErrCode, ret.ErrMsg)
	}
	return nil
}

func (m *Menu) Info() (*Info, error) {
	_, body, err := m.req.Get(ReqGet, nil)
	if err != nil {
		return nil, fmt.Errorf("menu info: %v", err)
	}
	info := &Info{}
	err = json.Unmarshal(body, info)
	if err != nil {
		return nil, fmt.Errorf("menu unmarshal: %v", err)
	}
	return info, nil
}

func (m *Menu) Delete() error {
	_, body, err := m.req.Get(ReqDelete, nil)
	if err != nil {
		return fmt.Errorf("menu delete: %v", err)
	}
	ret := &menuResp{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return fmt.Errorf("menu unmarshal: %v", err)
	}
	if ret.ErrCode != 0 {
		return fmt.Errorf("menu resp: %d %s", ret.ErrCode, ret.ErrMsg)
	}
	return nil
}

//TODO Support custom menu
//
//
//
