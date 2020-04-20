package menu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dreamer2q/go_wechat/request"
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

//
//自定义菜单最多包括3个一级菜单，每个一级菜单最多包含5个二级菜单。
//一级菜单最多4个汉字，二级菜单最多7个汉字，多出来的部分将会以“...”代替。
//创建自定义菜单后，菜单的刷新策略是，在用户进入公众号会话页或公众号profile页时，如果发现上一次拉取菜单的请求在5分钟以前，就会拉取一下菜单，如果菜单有更新，就会刷新客户端的菜单。测试时可以尝试取消关注公众账号后再次关注，则可以看到创建后的效果。​
//
//自定义菜单接口可实现多种类型按钮，如下：
//click：点击推事件用户点击click类型按钮后，微信服务器会通过消息接口推送消息类型为event的结构给开发者（参考消息接口指南），并且带上按钮中开发者填写的key值，开发者可以通过自定义的key值与用户进行交互；
//view：跳转URL用户点击view类型按钮后，微信客户端将会打开开发者在按钮中填写的网页URL，可与网页授权获取用户基本信息接口结合，获得用户基本信息。
//scancode_push：扫码推事件用户点击按钮后，微信客户端将调起扫一扫工具，完成扫码操作后显示扫描结果（如果是URL，将进入URL），且会将扫码的结果传给开发者，开发者可以下发消息。
//scancode_waitmsg：扫码推事件且弹出“消息接收中”提示框用户点击按钮后，微信客户端将调起扫一扫工具，完成扫码操作后，将扫码的结果传给开发者，同时收起扫一扫工具，然后弹出“消息接收中”提示框，随后可能会收到开发者下发的消息。
//pic_sysphoto：弹出系统拍照发图用户点击按钮后，微信客户端将调起系统相机，完成拍照操作后，会将拍摄的相片发送给开发者，并推送事件给开发者，同时收起系统相机，随后可能会收到开发者下发的消息。
//pic_photo_or_album：弹出拍照或者相册发图用户点击按钮后，微信客户端将弹出选择器供用户选择“拍照”或者“从手机相册选择”。用户选择后即走其他两种流程。
//pic_weixin：弹出微信相册发图器用户点击按钮后，微信客户端将调起微信相册，完成选择操作后，将选择的相片发送给开发者的服务器，并推送事件给开发者，同时收起相册，随后可能会收到开发者下发的消息。
//location_select：弹出地理位置选择器用户点击按钮后，微信客户端将调起地理位置选择工具，完成选择操作后，将选择的地理位置发送给开发者的服务器，同时收起位置选择工具，随后可能会收到开发者下发的消息。
//media_id：下发消息（除文本消息）用户点击media_id类型按钮后，微信服务器会将开发者填写的永久素材id对应的素材下发给用户，永久素材类型可以是图片、音频、视频、图文消息。请注意：永久素材id必须是在“素材管理/新增永久素材”接口上传后获得的合法id。
//view_limited：跳转图文消息URL用户点击view_limited类型按钮后，微信客户端将打开开发者在按钮中填写的永久素材id对应的图文消息URL，永久素材类型只支持图文消息。请注意：永久素材id必须是在“素材管理/新增永久素材”接口上传后获得的合法id。​
//
func (m *Menu) Create(menu RootMenu) error {
	menu.Type() //init type data for json marshal
	jsonData, err := json.Marshal(&menu)
	if err != nil {
		return fmt.Errorf("menu marshal: %v", err)
	}
	//start debug
	fmt.Printf("menu json: %s\n", jsonData)
	//end debug
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
