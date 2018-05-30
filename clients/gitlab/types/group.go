package types

type Group struct {
	Id                        int        `json:"id,omitempty"`
	Name                      string     `json:"name,omitempty"`
	Path                      string     `json:"path,omitempty"`
	Description               string     `json:"description,omitempty"`
	Visibility                string     `json:"visibility,omitempty"`
	LfsEnabled                bool       `json:"lfs_enabled,omitempty"`
	AvatarUrl                 string     `json:"avatar_url,omitempty"`
	WebURL                    string     `json:"web_url,omitempty"`
	RequestAccessEnabled      bool       `json:"request_access_enabled,omitempty"`
	FullName                  string     `json:"full_name,omitempty"`
	FullPath                  string     `json:"full_path,omitempty"`
	ParentId                  int        `json:"parent_id,omitempty"`
	SharedRunnersMinutesLimit int        `json:"shared_runners_minutes_limit,omitempty"`
	Projects                  []*Project `json:"projects,omitempty"`
}
