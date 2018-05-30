package weixin

import (
	"github.com/haborhuang/go-tools/clients/tencent/weixin/types"
	"net/http"
	"net/url"
)

type CompClient struct {
	compAppid     string
	compAppSecret string
	*Client
}

func (c *Client) NewCompClient(compAppId, compSecret string) *CompClient {
	return &CompClient{
		compAppid:     compAppId,
		compAppSecret: compSecret,
		Client:        c,
	}
}

type componentTokenReq struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

func (c *CompClient) GetCompAccessToken(ticket string) (*types.ComponentTokenRes, error) {
	var res *types.ComponentTokenRes

	err := c.newResponse(
		c.newRequest().SubPath("component/api_component_token").Method(http.MethodPost).JsonBody(componentTokenReq{
			ComponentAppid:        c.compAppid,
			ComponentAppsecret:    c.compAppSecret,
			ComponentVerifyTicket: ticket,
		}),
	).intoJson(&res)

	return res, err
}

type CertifiedCompClient struct {
	compAccessToken string
	*CompClient
}

func (c *CompClient) NewCertifiedCompClient(compAccessToken string) *CertifiedCompClient {
	return &CertifiedCompClient{
		compAccessToken: compAccessToken,
		CompClient: c,
	}
}

type preAuthCodeReq struct {
	ComponentAppid string `json:"component_appid"`
}

func (c *CertifiedCompClient) GetPreAuthCode() (*types.PreAuthCodeRes, error) {
	var res *types.PreAuthCodeRes
	query := url.Values{}
	query.Set("component_access_token", c.compAccessToken)

	err := c.newResponse(
		c.newRequest().SubPath("component/api_create_preauthcode").Method(http.MethodPost).JsonBody(preAuthCodeReq{
			ComponentAppid: c.compAppid,
		}).Query(query),
	).intoJson(&res)

	return res, err
}

type queryAuthReq struct {
	ComponentAppid    string `json:"component_appid"`
	AuthorizationCode string `json:"authorization_code"`
}

func (c *CertifiedCompClient) QueryAuth(authCode string) (*types.QueryAuthRes, error) {
	var res *types.QueryAuthRes
	query := url.Values{}
	query.Set("component_access_token", c.compAccessToken)

	err := c.newResponse(
		c.newRequest().SubPath("component/api_query_auth").Method(http.MethodPost).Query(query).JsonBody(queryAuthReq{
			ComponentAppid:    c.compAppid,
			AuthorizationCode: authCode,
		}),
	).intoJson(&res)

	return res, err
}

type authTokenReq struct {
	RefreshToken   string `json:"authorizer_refresh_token"`
	Appid          string `json:"authorizer_appid"`
	ComponentAppid string `json:"component_appid"`
}

func (c *CertifiedCompClient) GetAuthToken(appId, refreshToken string) (*types.AuthTokenRes, error) {
	var res *types.AuthTokenRes
	query := url.Values{}
	query.Set("component_access_token", c.compAccessToken)

	err := c.newResponse(
		c.newRequest().SubPath("component/api_authorizer_token").Method(http.MethodPost).Query(query).JsonBody(authTokenReq{
			ComponentAppid: c.compAppid,
			Appid:          appId,
			RefreshToken:   refreshToken,
		}),
	).intoJson(&res)

	return res, err
}
