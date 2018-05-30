package gitlab

import (
	"github.com/haborhuang/go-tools/clients/gitlab/types"

	"net/http"
	"fmt"
)

func hooksPath(projId string) string {
	return fmt.Sprintf(hooksPathFmt, projId)
}

func hookPath(projId string, hid int) string {
	return fmt.Sprintf(hookPathFmt, projId, hid)
}

const hooksPathFmt = projectPathFmt + "/hooks"
const hookPathFmt = hooksPathFmt + "/%d"

func (c *Client) AddProjectHook(pid string, hook *types.Hook) (*types.Hook, error) {
	var nh *types.Hook
	err := c.newResponse(
		c.newRequest().Method(http.MethodPost).RawSubPath(hooksPath(pid)).JsonBody(hook),
	).intoJson(&nh)
	return nh, err
}

func (c *Client) EditProjectHook(pid string, hid int, hook *types.Hook) (*types.Hook, error) {
	var nh *types.Hook
	err := c.newResponse(
		c.newRequest().Method(http.MethodPut).RawSubPath(hookPath(pid, hid)).JsonBody(hook),
	).intoJson(&nh)
	return nh, err
}

func (c *Client) ListProjectHooks(pid string) ([]*types.Hook, error) {
	var res []*types.Hook
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).RawSubPath(hooksPath(pid)),
	).intoJson(&res)

	return res, err
}