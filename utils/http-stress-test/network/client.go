package network

import (
	"github.com/valyala/fasthttp"
)

type HttpClient struct {
	fastHttpClient *fasthttp.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		fastHttpClient: &fasthttp.Client{},
	}
}

func (c *HttpClient) SendRequest(target string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(target)

	err := c.fastHttpClient.Do(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
