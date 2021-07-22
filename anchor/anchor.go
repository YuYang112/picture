package anchor

type _type int

const (
	EXTERIOR _type = iota
	INTERIOR
)

var uri = "http://v.api.lq.autohome.com.cn/Wcf/VideoService.svc/GetSevenStepClipVideosBySpecs?specIds=%d&_appid=yhcp"

var extInfo = map[string]bool{
	"车身尺寸": true,
	"轴距":   true,
	"轮胎尺寸": true,
	"发动机":  true,
	"前后灯":  true,
	"后备厢":  true,
}

var interiorInfo = map[string]bool{
	"驾驶位座椅": true,
	"仪表盘":   true,
	"方向盘":   true,
	"中控屏幕":  true,
	"空调":    true,
	"天窗":    true,
}

type Info struct {
	Tag     string `json:"tagname"`
	VideoId string `json:"videoid"`
}

type Resp struct {
	ReturnCode int `json:"return_code"`
	Result     struct {
		RespId          int64    `json:"resp_id"`
		ExteriorAnchors []Anchor `json:"exterior_anchors"`
		InteriorAnchors []Anchor `json:"interior_anchors"`
	}
}

type Anchor struct {
	Id    int
	Score float32
	Box   []float32
}
