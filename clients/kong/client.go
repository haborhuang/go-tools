package kong

import (
	"fmt"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"net/http"

	clienttool "github.com/haborhuang/go-tools/http"
	"github.com/haborhuang/go-tools/clients/kong/types"
)

type Client struct {
	url   *url.URL
}

func NewClientOrDie(domainUrl string) *Client {
	c, err := NewClient(domainUrl)
	if nil != err {
		panic(fmt.Errorf("New client error: %v", err))
	}

	return c
}

func NewClient(domainUrl string) (*Client, error) {
	u, err := url.Parse(domainUrl)
	if nil != err {
		return nil, fmt.Errorf("Parse url error: %v", err)
	}

	return &Client{
		url:   u,
	}, nil
}

func (c *Client) newRequest() *clienttool.HttpRequest {
	return clienttool.NewHttpReq(*c.url).SetHeader("Content-Type", "application/json")
}

type response struct {
	req *clienttool.HttpRequest
}

func (c *Client) newResponse(r *clienttool.HttpRequest) *response {
	return &response{
		req: r,
	}
}

func (r *response) doRaw() (*http.Response, error) {
	resp, err := r.req.DoRaw()
	if nil != err {
		return nil, err
	}

	if err := types.ParseErr(resp); nil != err {
		return nil, err
	}

	return resp, nil
}

func (r *response) do() error {
	resp, err := r.doRaw()
	if nil != err {
		return err
	}

	resp.Body.Close()
	return nil
}

func (r *response) intoJson(expected interface{}) error {
	resp, err := r.doRaw()
	if nil != err {
		return err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err := json.Unmarshal(body, expected); nil != err {
		return fmt.Errorf("Decode response error: %v\nBody:%s", err, string(body))
	}

	return nil
}
