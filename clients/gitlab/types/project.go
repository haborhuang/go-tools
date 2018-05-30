package types

import (
	"net/url"
	"path"
)

func UrlEncodedPath(args ...string) string {
	return url.PathEscape(path.Join(args...))
}

type Project struct {
	Id                   int        `json:"id,omitempty"`
	Name                 string     `json:"name,omitempty"`
	Description          string     `json:"description,omitempty"`
	DefaultBranch        string     `json:"default_branch,omitempty"`
	Owner                *Member    `json:"owner,omitempty"`
	Public               bool       `json:"public,omitempty"`
	Path                 string     `json:"path,omitempty"`
	PathWithNamespace    string     `json:"path_with_namespace,omitempty"`
	Visibility           string     `json:"visibility,omitempty"`
	IssuesEnabled        bool       `json:"issues_enabled,omitempty"`
	MergeRequestsEnabled bool       `json:"merge_requests_enabled,omitempty"`
	WallEnabled          bool       `json:"wall_enabled,omitempty"`
	WikiEnabled          bool       `json:"wiki_enabled,omitempty"`
	CreatedAtRaw         string     `json:"created_at,omitempty"`
	Namespace            *Namespace `json:"namespace,omitempty"`
	NamespaceId          int        `json:"namespace_id,omitempty"` // Only used for create
	SshRepoUrl           string     `json:"ssh_url_to_repo"`
	HttpRepoUrl          string     `json:"http_url_to_repo"`
	WebUrl               string     `json:"web_url"`
	SharedRunners        bool       `json:"shared_runners_enabled"`
}

type Member struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	State     string `json:"state"`
	CreatedAt string `json:"created_at,omitempty"`
	// AccessLevel int
}

type Namespace struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Owner_Id    int    `json:"owner_id"`
	Created_At  string `json:"created_at"`
	Updated_At  string `json:"updated_at"`
}
