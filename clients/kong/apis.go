package kong

import (
	"github.com/haborhuang/go-tools/clients/kong/types"
	"path"
	"net/http"
)

func apiPath(nameOrId string) string {
	return path.Join(apisPath, nameOrId)
}

const apisPath = "/apis"

func (c *Client) GetAPI(nameOrId string) (*types.API, error) {
	var res *types.API
	err := c.newResponse(
		c.newRequest().Method(http.MethodGet).SubPath(apiPath(nameOrId)),
	).intoJson(&res)

	return res, err
}

func (c *Client) AddAPI(req *types.API) (*types.API, error) {
	var res *types.API
	err := c.newResponse(
		c.newRequest().Method(http.MethodPost).SubPath(apisPath).JsonBody(req),
	).intoJson(&res)

	return res, err
}

func (c *Client) UpdateAPI(nameOrId string, req *types.API) (*types.API, error) {
	var res *types.API
	err := c.newResponse(
		c.newRequest().Method(http.MethodPatch).SubPath(apiPath(nameOrId)).JsonBody(req),
	).intoJson(&res)

	return res, err
}

func (c *Client) DeleteAPI(nameOrId string) error {
	err := c.newResponse(
		c.newRequest().Method(http.MethodDelete).SubPath(apiPath(nameOrId)),
	).do()

	return err
}