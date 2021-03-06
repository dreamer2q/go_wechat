package media

type Base struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type Result struct {
	Base

	//upload article(news)
	Url string `json:"url"`

	Type      string `json:"type"` //values: image, voice, video, thumb
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`

	MsgID     int    `json:"msg_id"`
	MsgStatus string `json:"msg_status"`
	MsgDataID int    `json:"msg_data_id"`

	SpeedLevel int `json:"speed"`
	RealSpeed  int `json:"realspeed"`
}

type Response struct {
	Base

	//temporary video result
	VideoUrl string `json:"video_url"`
	//permanent video result
	Title       string `json:"title"`
	Description string `json:"description"`
	DownUrl     string `json:"down_url"`
	ContentType string

	//News type of material
	News []ArticleItem `json:"news_item"`

	//other types of material including image,voice,thumb
	Filename string `json:"-"`
	Data     []byte `json:"-"`
}

type VideoDescription struct {
	Title string `json:"title"`
	Intro string `json:"introduction"`
}

type ArticleItem struct {
	Title        string `json:"title"`
	ThumbMediaID string `json:"thumb_media_id"` //图文消息的封面图片素材id（必须是永久mediaID）
	ShowCoverPic int    `json:"show_cover_pic"` //是否显示封面，0为false，即不显示，1为true，即显示
	Author       string `json:"author"`
	Digest       string `json:"digest"`             //图文消息的摘要，仅有单图文消息才有摘要，多图文此处为空。如果本字段为没有填写，则默认抓取正文前64个字。
	Content      string `json:"content"`            //图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS,涉及图片url必须来源 "上传图文消息内的图片获取URL"接口获取。外部图片url将被过滤。
	URL          string `json:"url"`                //图文页的URL
	SourceURL    string `json:"content_source_url"` //图文消息的原文地址，即点击“阅读原文”后的URL

	//fields for uploading article
	Comment    int `json:"need_open_comment"`     //Uint32 是否打开评论，0不打开，1打开
	OnlyForFan int `json:"only_fans_can_comment"` //Uint32 是否粉丝才可评论，0所有人可评论，1粉丝才可评论
}

type ArticleWrapper struct {
	Articles []ArticleItem `json:"articles"`
}
type ArticleUpdateWrapper struct {
	MediaID string      `json:"media_id"`
	Index   int         `json:"index"`
	Article ArticleItem `json:"articles"`
}

type MaterialCounter struct {
	Voice   int `json:"voice_count"`
	Video   int `json:"video_count"`
	Image   int `json:"image_count"`
	Article int `json:"news_count"`
}

type MaterialList struct {
	Base

	TotalCount int `json:"total_count"`
	ItemCount  int `json:"item_count"`

	Items []struct {
		MediaID    string `json:"media_id"`
		UpdateTime int    `json:"update_time"`

		//Article
		Content struct {
			Articles []ArticleItem `json:"news_item"`
		} `json:"content"`

		//Image, Voice, Video
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"item"`
}
