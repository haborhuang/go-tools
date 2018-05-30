package gitlab

import (
	"net/url"
	"fmt"

	clienttool "github.com/haborhuang/go-tools/http"
	"github.com/haborhuang/go-tools/clients/gitlab/types"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Url    string
	APIVer string
	Token string
}

type Client struct {
	url *url.URL
	token string
}

func NewClientOrDie(conf Config) *Client {
	c, err := NewClient(conf)
	if nil != err {
		panic(fmt.Errorf("New client error: %v", err))
	}

	return c
}

func NewClient(conf Config) (*Client, error) {
	u, err := url.Parse(conf.Url)
	if nil != err {
		return nil, fmt.Errorf("Parse url error: %v", err)
	}

	if "" == conf.APIVer {
		conf.APIVer = "v4"
	}

	u.Path = "/api/" + conf.APIVer

	return &Client{
		url: u,
		token: conf.Token,
	}, nil
}

func (c *Client) newRequest() *clienttool.HttpRequest {
	return clienttool.NewHttpReq(*c.url).SetHeader("PRIVATE-TOKEN", c.token)
}

type response struct {
	req *clienttool.HttpRequest
	respHeader http.Header
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

	if r.respHeader != nil {
		for k := range r.respHeader {
			r.respHeader.Set(k, resp.Header.Get(k))
		}
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

func (r *response) extractRespHeaders(header http.Header) *response {
	r.respHeader = header
	return r
}