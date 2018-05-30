package types

type User struct {
	Id            int    `json:"id,omitempty"`
	Username      string `json:"username,omitempty"`
	Email         string `json:"email,omitempty"`
	Name          string `json:"name,omitempty"`
	State         string `json:"state,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	Bio           string `json:"bio,omitempty"`
	Skype         string `json:"skype,omitempty"`
	LinkedIn      string `json:"linkedin,omitempty"`
	Twitter       string `json:"twitter,omitempty"`
	WebURL        string `json:"web_url"`
	ExternUid     string `json:"extern_uid,omitempty"`
	Provider      string `json:"provider,omitempty"`
	ThemeId       int    `json:"theme_id,omitempty"`
	ColorSchemeId int    `json:"color_scheme_id,omitempty"`
	AvatarUrl     string `json:"avatar_url,omitempty"`
}