package types

import "time"

type Commit struct {
	Id          string    `json:"id"`
	ShortId     string    `json:"short_id"`
	Title       string    `json:"title"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message"`
	ParentIds   []string  `json:"parent_ids"`
}

type CommitPayload struct {
	CommitBasicInfo
	Actions []*CommitAction `json:"actions"`
}

type CommitBasicInfo struct {
	Branch        string `json:"branch"`
	CommitMessage string `json:"commit_message"`
	StartBranch   string `json:"start_branch,omitempty"`
	AuthorEmail   string `json:"author_email,omitempty"`
	AuthorName    string `json:"author_name,omitempty"`
}

const (
	CommitContentEncodingText   = "text"
	CommitContentEncodingBase64 = "base64"
)

type CommitContent struct {
	Encoding string `json:"encoding,omitempty"`
	Content  string `json:"content"`
}

const (
	CommitActionCreate = "create"
	CommitActionDelete = "delete"
	CommitActionMove   = "move"
	CommitActionUpdate = "update"
)

type CommitAction struct {
	Action       string `json:"action"`
	FilePath     string `json:"file_path"`
	PreviousPath string `json:"previous_path,omitempty"`
	CommitContent
}
