package user

type Tags struct {
	Tags []TagItem `json:"tags"`
}
type TagItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TagUsers struct {
	OpenIDList []string `json:"openid_list"`
	TagID      int      `json:"tagid"`
}

type Language string

const (
	LangZh = "zh_CN"
	LangZW = "zh_TW"
	LangEn = "en"
)

type GetInfo struct {
	OpenID string   `json:"openid"`
	Lang   Language `json:"lang,omitempty"`
}

type Info struct {
	Subscribe      int    `json:"subscribe"`
	Openid         string `json:"openid"`
	Nickname       string `json:"nickname"`
	Sex            int    `json:"sex"`
	Language       string `json:"language"`
	City           string `json:"city"`
	Province       string `json:"province"`
	Country        string `json:"country"`
	HeadImgUrl     string `json:"headimgurl"`
	SubscribeTime  int64  `json:"subscribe_time"`
	UnionID        string `json:"unionid"`
	Remark         string `json:"remark"`
	GroupID        int    `json:"groupid"`
	TagList        []int  `json:"tagid_list"`
	SubscribeScene string `json:"subscribe_scene"`
	QrScene        int    `json:"qr_scene"`
	QrSceneStr     string `json:"qr_scene_str"`
}
