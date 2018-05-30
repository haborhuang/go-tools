package gitlab

import (
	"github.com/haborhuang/go-tools/clients/gitlab/types"
	"net/http"
	"fmt"
)

func projectPath(projId string) string {
	return fmt.Sprintf(projectPathFmt, projId)
}

const (
	projectsPath = "/projects"
	projectPathFmt = projectsPath + "/%s"
)

func (c *Client) GetProject(projId string) (*types.Project, error) {
	var p *types.Project
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(projectPath(projId)),
	).intoJson(&p)
	return p, err
}

func (c *Client) CreateProject(p *types.Project) (*types.Project, error) {
	var np *types.Project
	err := c.newResponse(
		c.newRequest().Method(http.MethodPost).SubPath(projectsPath).JsonBody(&p),
	).intoJson(&np)

	return np, err
}