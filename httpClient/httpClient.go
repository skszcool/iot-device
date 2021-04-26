package httpClient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	url2 "net/url"
	"time"
)

var restyClient *resty.Client

func init() {
	restyClient = resty.New()
	transport := http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, //连接超时时间
			KeepAlive: 30 * time.Second, //连接保持超时时间
		}).DialContext,

		//跳过证书验证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	restyClient.SetTransport(&transport)
}

type SkHttpRespBody struct {
	Code    interface{} `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type formatReqFunc func(request *resty.Request) *resty.Request

type client struct {
	request *resty.Request
}

func IsSuccessForSk(code interface{}) bool {
	return fmt.Sprintf("%v", code) == "200"
}

func IsFailedForSk(code interface{}) bool {
	return fmt.Sprintf("%v", code) != "200"
}

func C() *client {
	c := new(client)
	c.request = restyClient.R()

	return c
}

func (c *client) GetRequest(formatReq formatReqFunc) *resty.Request {
	return c.request
}

func (c *client) FormatReq(formatReq formatReqFunc) *client {
	formatReq(c.request)
	return c
}

func (c *client) SetResult(result interface{}) *client {
	c.request.SetResult(result)
	return c
}

func (c *client) SetError(result interface{}) *client {
	c.request.SetError(result)
	return c
}

func (c *client) SetHeader(header string, value string) *client {
	c.request.SetHeader(header, value)
	return c
}

func defaultApplicationJson(request *resty.Request) *resty.Request {
	if request.Header.Get("Content-Type") == "" {
		request.SetHeader("Content-Type", "application/json")
	}

	return request
}

func (c *client) Get(url string, params interface{}) (*resty.Response, error) {
	request := c.request
	switch params.(type) {
	case string:
		request.SetQueryString(params.(string))
	case map[string]string:
		request.SetQueryParams(params.(map[string]string))
	case url2.Values:
		request.SetQueryParamsFromValues(params.(url2.Values))
	default:
		return nil, errors.New("params的参数类型不合法")
	}

	return request.Get(url)
}

func (c *client) PostJSON(url string, params interface{}) (*resty.Response, error) {
	request := defaultApplicationJson(c.request)
	return request.
		SetBody(params).
		Post(url)
}

func (c *client) Put(url string, params interface{}) (*resty.Response, error) {
	request := defaultApplicationJson(c.request)
	return request.
		SetBody(params).
		Put(url)
}

func (c *client) Delete(url string, params interface{}) (*resty.Response, error) {
	request := defaultApplicationJson(c.request)

	return request.
		SetBody(params).
		Delete(url)
}

// 上传文件
func (c *client) Upload(url string, files map[string]string, params map[string]string) (*resty.Response, error) {
	request := defaultApplicationJson(c.request)

	return request.
		SetFiles(files).
		SetFormData(params).
		Post(url)
}
