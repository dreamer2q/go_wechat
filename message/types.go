package message

import "wechat/request"

type TemplateItem struct {
	TemplateID      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

type TemplateResult struct {
	request.CommonError
	List []TemplateItem `json:"template_list"`
}

type MiniProgram struct {
	AppID    string `json:"appid"`
	PagePath string `json:"pagepath"`
}

type TemplateMsg struct {
	ToUser     string       `json:"touser"` //openID
	TemplateID string       `json:"template_id"`
	URL        string       `json:"url"`
	Program    *MiniProgram `json:"miniprogram,omitempty"`
	Data       interface{}  `json:"data"`
}

/////for auto reply match rule
/////might not be useful
type autoReplyInfo struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
type autoReplyKeyword struct {
	autoReplyInfo
	MatchMode string `json:"match_mode"`
}
type autoReplyNews struct {
	autoReplyInfo
	NewsInfo struct {
		List []newsItem `json:"list"`
	} `json:"news_info"`
}
type newsItem struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Digest     string `json:"digest"`
	ShowCover  int    `json:"show_cover"`
	CoverURL   string `json:"cover_url"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url"`
}
type AutoReplySetting struct {
	request.CommonError

	AddFriendReplyOpen          int           `json:"is_add_friend_reply_open"`
	AutoReplyOpen               int           `json:"is_autoreply_open"`
	AddFriendAutoReplyInfo      autoReplyInfo `json:"add_friend_autoreply_info"`
	MessageDefaultAutoReplyInfo autoReplyInfo `json:"message_default_autoreply_info"`
	KeywordAutoReplyInfo        struct {
		List []struct {
			RuleName        string             `json:"rule_name"`
			CreateTime      int                `json:"create_time"`
			ReplyMode       string             `json:"reply_mode"`
			KeywordListInfo []autoReplyKeyword `json:"keyword_list_info"`
			ReplyListInfo   []autoReplyNews    `json:"reply_list_info"`
		} `json:"list"`
	} `json:"keyword_autoreply_info"`
}
