package types

type RepoFile struct {
	CommitBasicInfo
	CommitContent
	LastCommitId string `json:"last_commit_id,omitempty"`
}

type SavedRepoFile struct {
	FileName string `json:"file_name"`
	Branch   string `json:"branch"`
}

const (
	TreeObjTypeTree = "tree"
	TreeObjTypeBlob = "blob"
)

type RepoTreeObj struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Path string `json:"path"`
	Mode string `json:"mode"`
}
