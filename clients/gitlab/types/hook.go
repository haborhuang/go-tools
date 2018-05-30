package types

type Hook struct {
	Id           int    `json:"id,omitempty"`
	Url          string `json:"url,omitempty"`
	CreatedAtRaw string `json:"created_at,omitempty"`
	HookFlags
}

type HookFlags struct {
	PushEvents            bool `json:"push_events"`
	IssuesEvents          bool `json:"issues_events"`
	MergeRequestsEvents   bool `json:"merge_requests_events"`
	TagPushEvents         bool `json:"tag_push_events"`
	NoteEvents            bool `json:"note_events"`
	JobEvents             bool `json:"job_events"`
	PipelineEvents        bool `json:"pipeline_events"`
	WikiEvents            bool `json:"wiki_events"`
	EnableSSLVerification bool `json:"enable_ssl_verification"`
}