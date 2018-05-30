package gitlab

import (
	"net/http"
	"net/url"
	"fmt"

	"github.com/haborhuang/go-tools/clients/gitlab/types"
)

const (
	createPipelinePathFmt = projectPathFmt + "/pipeline"
	pipelinesPathFmt = projectPathFmt + "/pipelines"
	pipelinePathFmt = pipelinesPathFmt + "/%d"
)

func createPipelinePath(projId string) string {
	return fmt.Sprintf(createPipelinePathFmt, projId)
}

func pipelinesPath(projId string) string {
	return fmt.Sprintf(pipelinesPathFmt, projId)
}

func pipelinePath(projId string, pipelineId int) string {
	return fmt.Sprintf(pipelinePathFmt, projId, pipelineId)
}

func (c *Client) CreatePipeline(pid, ref string) (*types.Pipeline, error) {
	q := url.Values{}
	q.Set("ref", ref)

	var p *types.Pipeline
	err := c.newResponse(
		c.newRequest().Method(http.MethodPost).RawSubPath(createPipelinePath(pid)).Query(q),
	).intoJson(&p)

	return p, err
}

func (c *Client) ListPipelines(pid string, opts *types.ListPipelinesOpts) ([]*types.PipelineBrief, error) {
	q, err := opts.ToQuery()
	if nil != err {
		return nil, fmt.Errorf("Check list pipelines parameters error: %v", err)
	}

	var ps []*types.PipelineBrief
	err = c.newResponse(
		c.newRequest().Debug().Method(http.MethodGet).RawSubPath(pipelinesPath(pid)).Query(q),
	).intoJson(&ps)

	return ps, err
}

func (c *Client) GetPipeline(projId string, pipelineId int) (*types.Pipeline, error) {
	var p *types.Pipeline
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(pipelinePath(projId, pipelineId)),
	).intoJson(&p)

	return p, err
}