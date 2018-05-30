package types

import (
	"time"
	"fmt"
	"net/url"
	"strconv"
)

type Pipeline struct {
	PipelineBrief
	BeforeSHA   string     `json:"before_sha"`
	Tag         bool       `json:"tag"`
	User        User       `json:"user"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	StartedAt   *time.Time `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`
	CommittedAt *time.Time `json:"committed_at"`
	Duration    *int64     `json:"duration"`
}

type PipelineBrief struct {
	Id     int    `json:"id"`
	SHA    string `json:"sha"`
	Ref    string `json:"ref"`
	Status string `json:"status"`
}

func (p *PipelineBrief) Finished() bool {
	return p.Status != "" && p.Status != "running" && p.Status != "pending"
}

type Job struct {
	Commit     Commit        `json:"commit"`
	CreatedAt  *time.Time    `json:"created_at"`
	FinishedAt *time.Time    `json:"finished_at"`
	StartedAt  *time.Time    `json:"started_at"`
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	Pipeline   PipelineBrief `json:"pipeline"`
	Ref        string        `json:"ref"`
	Stage      string        `json:"stage"`
	Status     string        `json:"status"`
	Tag        bool          `json:"tag"`
	User       User          `json:"user"`
}

type ListPipelinesOpts struct {
	Scope      string
	Status     string
	Ref        string
	YamlErrors bool
	Name       string
	UserName   string
	OrderBy    string
	SortAsc    bool
	Pagination
}

type Pagination struct {
	Page    int
	PerPage int
}

func (opts *ListPipelinesOpts) ToQuery() (url.Values, error) {
	if nil == opts {
		return nil, nil
	}

	if err := opts.check(); nil != err {
		return nil, err
	}

	query := make(url.Values)
	if "" != opts.Scope {
		query.Set("scopes", opts.Scope)
	}
	if "" != opts.Status {
		query.Set("status", opts.Status)
	}
	if "" != opts.Ref {
		query.Set("ref", opts.Ref)
	}
	if opts.YamlErrors {
		query.Set("yaml_errors", "true")
	}
	if "" != opts.Name {
		query.Set("name", opts.Name)
	}
	if "" != opts.UserName {
		query.Set("username", opts.UserName)
	}
	if "" != opts.OrderBy {
		query.Set("order_by", opts.OrderBy)
	}
	if opts.SortAsc {
		query.Set("sort", "asc")
	}
	opts.Pagination.ToQuery(query)
	return query, nil
}

var validPipelineScope, validPipelineStatus, validPipelineOrder map[string]bool

func init() {
	validPipelineScope = map[string]bool{
		"running":  true,
		"pending":  true,
		"finished": true,
		"branches": true,
		"tags":     true,
	}

	validPipelineStatus = map[string]bool{
		"running":  true,
		"pending":  true,
		"success":  true,
		"failed":   true,
		"canceled": true,
		"skipped":  true,
	}

	validPipelineOrder = map[string]bool{
		"id":      true,
		"status":  true,
		"ref":     true,
		"user_id": true,
	}
}

func (opts *ListPipelinesOpts) check() error {
	if nil == opts {
		return nil
	}

	if "" != opts.Scope && !validPipelineScope[opts.Scope] {
		return fmt.Errorf("Invalid scopes '%s'", opts.Scope)
	}

	if "" != opts.Status && !validPipelineStatus[opts.Status] {
		return fmt.Errorf("Invalid status '%s'", opts.Status)
	}

	if "" != opts.OrderBy && !validPipelineOrder[opts.OrderBy] {
		return fmt.Errorf("Invalid order_by '%s'", opts.OrderBy)
	}

	return opts.Pagination.check()
}

func (p Pagination) check() error {
	if p.Page < 0 {
		return fmt.Errorf("Invalid page '%d'", p.Page)
	}

	if p.PerPage < 0 || p.PerPage > 100 {
		return fmt.Errorf("Invalid per_page '%d'", p.PerPage)
	}

	return nil
}

func (p Pagination) ToQuery(vals url.Values) {
	if p.Page > 0 {
		vals.Set("page", strconv.Itoa(p.Page))
	}
	if p.PerPage > 0 {
		vals.Set("per_page", strconv.Itoa(p.PerPage))
	}
}