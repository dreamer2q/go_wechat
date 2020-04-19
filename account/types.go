package account

type qrQuery struct {
	ExpireIn   int        `json:"expire_seconds,omitempty"`
	ActionName string     `json:"action_name"` //QR_SCENE -> tmp, QR_LIMIT_SCENE -> permanent
	ActionInfo actionInfo `json:"action_info"`
}
type actionInfo struct {
	Scene scene `json:"scene"`
}
type scene struct {
	SceneID  int    `json:"scene_id,omitempty"`
	SceneStr string `json:"scene_str,omitempty"`
}

type QrResult struct {
	Ticket   string `json:"ticket"`
	ExpireIn int    `json:"expire_seconds"`
	URL      string `json:"url"`
}
