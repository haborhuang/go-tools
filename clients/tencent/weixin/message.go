package weixin

import (
	"net/http"
	"net/url"
	"github.com/haborhuang/go-tools/clients/tencent/weixin/types"
	"errors"
)

type MPMsgClient struct {
	accessToken string
	*Client
}

func (c *Client) NewMPMsgClient(token string) *MPMsgClient {
	return &MPMsgClient{
		Client:      c,
		accessToken: token,
	}
}

func (c *MPMsgClient) GetPrivTemplates() (*types.TemplatesResp, error) {
	var res *types.TemplatesResp
	query := url.Values{}
	query.Set("access_token", c.accessToken)
	err := c.newResponse(
		c.newRequest().SubPath("template/get_all_private_template").Method(http.MethodGet).Query(query),
	).intoJson(&res)

	return res, err
}

func (c *MPMsgClient) SendTmplMsg(req *types.SendTmplMsgReq) (*types.SendTmplMsgResp, error) {
	if nil == req {
		return nil, errors.New("Empty request")
	}

	if req.ToUser == "" {
		return nil, errors.New("Missing receiver")
	}

	if req.TemplateId == "" {
		return nil, errors.New("Missing template id")
	}

	var res *types.SendTmplMsgResp
	query := url.Values{}
	query.Set("access_token", c.accessToken)
	err := c.newResponse(
		c.newRequest().SubPath("message/template/send").Method(http.MethodPost).Query(query).JsonBody(req),
	).intoJson(&res)

	return res, err
}
