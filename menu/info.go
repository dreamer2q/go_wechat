package menu

type InfoNews struct {
	List []struct {
		Title      string `json:"title"`
		Author     string `json:"author"`
		Digest     string `json:"digest"`
		ShowCover  int    `json:"show_cover"`
		CoverURL   string `json:"cover_url"`
		ContentURL string `json:"content_url"`
		SourceURL  string `json:"source_url"`
	} `json:"list"`
}

type InfoBtn struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Url       string `json:"url,omitempty"`
	SubButton struct {
		List []InfoBtn `json:"list"`
	} `json:"sub_button,omitempty"`
	NewsInfo InfoNews `json:"news_info,omitempty"`
}

type Info struct {
	IsMenuOpen int `json:"is_menu_open"`
	MenuInfo   struct {
		Button []InfoBtn `json:"button"`
	} `json:"selfmenu_info"`
}

type CustomInfo struct {
	Menu struct {
		Button []InfoBtn `json:"button"`
		MenuID int64     `json:"menuid"`
	} `json:"menu"`

	CustomMenu struct {
		Button []InfoBtn `json:"button"`
		Match  MatchRule `json:"matchrule"`
		MenuID int64     `json:"menuid"`
	} `json:"conditionalmenu"`
}

type MatchInfo struct {
	Button []InfoBtn `json:"button"`
}
