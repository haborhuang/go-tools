package gitlab

import (
	"io"
	"net/url"
	"fmt"
	"mime"
	"github.com/haborhuang/go-tools/clients/gitlab/types"
	"net/http"
	"strconv"
)

func repoArchPath(projId string) string {
	return fmt.Sprintf(repoArchPathFmt, projId)
}

const repoArchPathFmt = projectPathFmt + "/repository/archive"

func (c *Client) RepoArchive(projId, sha string) (string, io.ReadCloser, error) {
	r := c.newRequest().
		RawSubPath(repoArchPath(projId))
	if sha != "" {
		q := make(url.Values)
		q.Set("sha", sha)
		r.Query(q)
	}


	resp, err := c.newResponse(r).doRaw()
	if nil != err {
		return "", nil, err
	}

	_, params, err := mime.ParseMediaType(resp.Header.Get("content-disposition"))
	if nil != err {
		resp.Body.Close()
		return "", nil, fmt.Errorf("Parse archive name error: %v", err)
	}

	return params["filename"], resp.Body, nil
}

func repoTreePath(projId string) string {
	return fmt.Sprintf(repoTreePathFmt, projId)
}

const repoTreePathFmt = projectPathFmt + "/repository/tree"

func (c *Client) RepoTree(pid, path, ref string, recursive bool, paging *types.Pagination) ([]*types.RepoTreeObj, int, error) {
	q := make(url.Values)
	if path != "" {
		q.Set("path", path)
	}
	if ref != "" {
		q.Set("ref", ref)
	}
	if recursive {
		q.Set("recursive", "true")
	}
	if nil != paging {
		paging.ToQuery(q)
	}

	var tree []*types.RepoTreeObj
	respHeader := make(http.Header)
	respHeader.Set(types.RespHeaderTotalPages, "")
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(repoTreePath(pid)).Query(q),
	).extractRespHeaders(respHeader).intoJson(&tree)

	pages, _ := strconv.Atoi(respHeader.Get(types.RespHeaderTotalPages))
	return tree, pages, err
}