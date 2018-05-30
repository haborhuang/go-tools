package gitlab

import (
	"github.com/haborhuang/go-tools/clients/gitlab/types"
	"net/http"
	"fmt"
)

func repoFilePath(projId, fpath string) string {
	return fmt.Sprintf(repoFilePathFmt, projId, fpath)
}

const repoFilePathFmt = projectPathFmt + "/repository/files/%s"

func (c *Client) CreateFile(pid, fpath string, file *types.RepoFile) (*types.SavedRepoFile, error) {
	return c.saveFile(pid, fpath, file, http.MethodPost)
}

func (c *Client) UpdateFile(pid, fpath string, file *types.RepoFile) (*types.SavedRepoFile, error) {
	return c.saveFile(pid, fpath, file, http.MethodPut)
}

func (c *Client) DeleteFile(pid, fpath string, file *types.RepoFile) error {
	return c.newResponse(
		c.newRequest().Method(http.MethodDelete).RawSubPath(repoFilePath(pid, fpath)).JsonBody(file),
	).do()
}

func (c *Client) saveFile(pid, fpath string, file *types.RepoFile, method string) (*types.SavedRepoFile, error) {
	var saved *types.SavedRepoFile
	err := c.newResponse(
		c.newRequest().Method(method).RawSubPath(repoFilePath(pid, fpath)).JsonBody(file),
	).intoJson(&saved)
	return saved, err
}