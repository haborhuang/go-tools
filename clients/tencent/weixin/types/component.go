package types

type ComponentTokenRes struct {
	*ErrRes
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int64  `json:"expires_in"`
}

type PreAuthCodeRes struct {
	*ErrRes
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

type QueryAuthRes struct {
	*ErrRes
	AuthorizationInfo authInfo `json:"authorization_info"`
}

type authInfo struct {
	Appid string `json:"authorizer_appid"`
	AuthToken
	FuncInfo []*funcInfo `json:"func_info"`
}

type funcInfo struct {
	FuncscopeCategory struct {
		Id int64 `json:"id"`
	} `json:"funcscope_category"`
}

type AuthTokenRes struct {
	*ErrRes
	AuthToken
}

type AuthToken struct {
	AccessToken  string `json:"authorizer_access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"authorizer_refresh_token"`
}
