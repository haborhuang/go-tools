package types

type TmplDataVal struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type TmplData map[string]*TmplDataVal

func (t TmplData) SetVal(key, value string, color ...string) TmplData {
	if t != nil {
		data := &TmplDataVal{
			Value: value,
		}

		if len(color) > 0 {
			data.Color = color[0]
		}

		t[key] = data
	}

	return t
}

type TemplatesResp struct {
	*ErrRes
	TmplList []*wxTmpl `json:"template_list"`
}

type wxTmpl struct {
	TmplId          string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

type MsgLink struct {
	Url         string          `json:"url,omitempty"`
	MiniProgram *MiniProgramRef `json:"miniprogram,omitempty"`
}

type MiniProgramRef struct {
	AppId    string `json:"appid"`
	PagePath string `json:"pagepath,omitempty"`
}

type SendTmplMsgReq struct {
	ToUser      string `json:"touser"`
	TemplateId  string `json:"template_id"`
	MsgLink
	Data TmplData `json:"data"`
}

type SendTmplMsgResp struct {
	*ErrRes
	MsgId int64 `json:"msg_id"`
}
