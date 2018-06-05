package sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/satori/go.uuid"
)

type Config struct {
	AccessKeyId  string
	AccessSecret string
}

type Client struct {
	signKey     string
	accessKeyId string
}

func NewClient(c Config) *Client {
	return &Client{
		accessKeyId: c.AccessKeyId,
		signKey:     c.AccessSecret + "&",
	}
}

const (
	msgUrl     = "http://dysmsapi.aliyuncs.com"
	maxNumbers = 1000
	reqMethod  = http.MethodGet
	signPrefix = reqMethod + "&%2F&"
)

type msgResp struct {
	Code      string `json:"Code"`
	BizId     string `json:"BizId"`
	RequestId string `json:"RequestId"`
	Message   string `json:"Message"`
}

func (c *Client) Send(nums []string, tmplCode, signName string, params map[string]string, extCode, outId string) (string, string, error) {
	if len(nums) > maxNumbers {
		return "", "", errors.New("Too many phone numbers")
	}

	query, err := c.GenQuery(nums, tmplCode, signName, params, extCode, outId, "JSON")
	if nil != err {
		return "", "", fmt.Errorf("generate query error: %v", err)
	}

	u, _ := url.Parse(msgUrl)
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if nil != err {
		return "", "", fmt.Errorf("Send HTTP request error: %v", err)
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		return "", "", fmt.Errorf("Response error: (%d)%s", resp.StatusCode, string(data))
	}

	var res msgResp
	if err := json.Unmarshal(data, &res); nil != err {
		return "", "", fmt.Errorf("Decode response body `%s` error: %v", string(data), err)
	}

	if res.Code != "OK" {
		return "", "", fmt.Errorf("Failed to revoke API, code: %s, message: %s.", res.Code, res.Message)
	}
	return res.BizId, res.RequestId, nil
}

func (c *Client) GenQuery(nums []string, tmplCode, signName string, params map[string]string, extCode, outId, format string) (url.Values, error) {
	query := make(url.Values)

	// System parameter
	query.Set("AccessKeyId", c.accessKeyId)
	query.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))
	query.Set("SignatureMethod", "HMAC-SHA1")
	query.Set("SignatureVersion", "1.0")
	query.Set("SignatureNonce", uuid.NewV4().String())
	if format != "" {
		query.Set("Format", format)
	}

	// Business parameter
	query.Set("Action", "SendSms")
	query.Set("Version", "2017-05-25")
	query.Set("RegionId", "cn-hangzhou")
	query.Set("PhoneNumbers", strings.Join(nums, ","))
	query.Set("SignName", signName)
	query.Set("TemplateCode", tmplCode)
	if len(params) > 0 {
		tp, err := json.Marshal(params)
		if nil != err {
			return nil, fmt.Errorf("Encode template parameters error: %v", err)
		}
		query.Set("TemplateParam", string(tp))
	}
	if outId != "" {
		query.Set("OutId", outId)
	}

	// Gen signature
	h := hmac.New(sha1.New, []byte(c.signKey))

	plain := strings.Replace(query.Encode(), "+", "%20", -1)
	//log.Println(plain)
	h.Write([]byte(signPrefix + url.QueryEscape(plain)))

	query.Set("Signature", base64.StdEncoding.EncodeToString(h.Sum(nil)))

	return query, nil
}
