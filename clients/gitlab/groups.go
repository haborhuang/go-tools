package gitlab

import (
	"net/http"
	"path"
	"github.com/haborhuang/go-tools/clients/gitlab/types"
	"net/url"
	"strconv"
)

const groupsPath = "/groups"

func groupPath(gid string) string {
	return path.Join(groupsPath, gid)
}

func projectsOfGroupPath(gid string) string {
	return path.Join(groupPath(gid), projectsPath)
}

func (c *Client) GetGroup(gid string) (*types.Group, error) {
	var group *types.Group
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(groupPath(gid)),
	).intoJson(&group)

	return group, err
}

func (c *Client) ListProjectsOfGroup(gid string, paging *types.Pagination) ([]*types.Project, int, error) {
	q := make(url.Values)
	if nil != paging {
		paging.ToQuery(q)
	}

	var res []*types.Project
	respHeader := make(http.Header)
	respHeader.Set(types.RespHeaderTotalPages, "")
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(projectsOfGroupPath(gid)).Query(q),
	).extractRespHeaders(respHeader).intoJson(&res)

	pages, _ := strconv.Atoi(respHeader.Get(types.RespHeaderTotalPages))
	return res, pages, err
}