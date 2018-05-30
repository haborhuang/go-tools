package gitlab

import (
	"github.com/haborhuang/go-tools/clients/gitlab/types"
	"net/http"
	"fmt"
	"net/url"
)

func commitsPath(projId string) string {
	return fmt.Sprintf(commitsPathFmt, projId)
}

const commitsPathFmt = projectPathFmt + "/repository/commits"

func (c *Client) CreateCommit(pid string, payload *types.CommitPayload) (*types.Commit, error) {
	var commit *types.Commit
	err := c.newResponse(
		c.newRequest().Method(http.MethodPost).RawSubPath(commitsPath(pid)).JsonBody(&payload),
	).intoJson(&commit)

	return commit, err
}

func (c *Client) ListCommits(pid, ref string, paging *types.Pagination) ([]*types.Commit, error) {
	q := make(url.Values)
	if ref != "" {
		q.Set("ref_name", ref)
	}
	if nil != paging {
		paging.ToQuery(q)
	}

	var commits []*types.Commit

	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(commitsPath(pid)).Query(q),
	).intoJson(&commits)

	return commits, err
}