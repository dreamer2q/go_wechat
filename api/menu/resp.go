package menu

type NewsItem struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	Digest     string `json:"digest"`
	ShowCover  int    `json:"show_cover"`
	CoverURL   string `json:"cover_url"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url"`
}

type NewsInfo struct {
	List []NewsItem
}

type SubButton struct {
	List []Button `json:"list"`
}

type Button struct {
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Key       string    `json:"key,omitempty"`
	Value     string    `json:"value,omitempty"`
	Url       string    `json:"url,omitempty"`
	SubButton SubButton `json:"sub_button,omitempty"`
	NewsInfo  NewsInfo  `json:"news_info,omitempty"`
}

type Info struct {
	IsMenuOpen int `json:"is_menu_open"`
	MenuInfo   struct {
		Button []Button `json:"button"`
	} `json:"selfmenu_info"`
}
