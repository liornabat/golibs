package webservice

import (
	"context"
	"fmt"
	"golibs/instrument"
	"golibs/logging"
	"gopkg.in/resty.v1"
	"time"
)

type RestClient struct {
	Name           string
	Host           string
	Port           string
	User           string
	Password       string
	Routes         map[string]string
	Headers        map[string]string
	request        *ClientRequest
	ins            *instrument.HttpInstrument
	isInstrumented bool
}
type ClientRequest struct {
	r         *resty.Request
	base      string
	path      string
	ins       *instrument.HttpInstrument
	reqStruct interface{}
	resStruct interface{}
	errStruct interface{}
}

type ClientResponse struct {
	StatusCode int
	Body       string
	reqStruct  interface{}
	resStruct  interface{}
	errStruct  interface{}
}

var loggerHttp = logging.NewLogger("webservice/client")

func NewRestClient(name, host, port, user, password string) *RestClient {
	c := &RestClient{
		Name:     name,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Routes:   make(map[string]string),
		Headers:  make(map[string]string),
	}
	return c
}
func (c *RestClient) AddBaseRoute(name, path string) *RestClient {
	c.Routes[name] = path
	return c
}
func (c *RestClient) SetHeader(k, v string) *RestClient {
	c.Headers[k] = v
	return c
}
func (c *RestClient) NewRequest(ctx context.Context, baseName string) *ClientRequest {
	cr := &ClientRequest{
		r:   resty.R(),
		ins: c.ins,
	}
	cr.base = baseName
	cr.r.SetHeaders(c.Headers)
	if c.User != "" {
		cr.r.SetBasicAuth(c.User, c.Password)
	}
	cr.r.SetContext(ctx)
	cr.path = fmt.Sprintf("%s:%s%s", c.Host, c.Port, c.Routes[baseName])
	return cr
}

func (c *RestClient) EnableInstrumenting() *RestClient {
	ins, err := instrument.NewInstrument().SetHttpInstrument(c.Name)
	if err != nil {
		loggerHttp.Error(err, "Unable to set instrumentation")

	}
	c.ins = ins
	c.isInstrumented = true
	return c
}
func (cr *ClientRequest) SetSuffixPath(path string) *ClientRequest {
	cr.path = cr.path + path
	return cr
}
func (cr *ClientRequest) SetRequest(req interface{}) *ClientRequest {
	cr.r.SetBody(req)
	cr.reqStruct = req
	return cr
}

func (cr *ClientRequest) SetResult(res interface{}) *ClientRequest {
	cr.r.SetResult(res)
	cr.resStruct = res
	return cr
}

func (cr *ClientRequest) SetError(err interface{}) *ClientRequest {
	cr.r.SetError(err)
	cr.errStruct = err
	return cr
}
func (cr *ClientRequest) SetHeader(h, v string) *ClientRequest {
	cr.r.SetHeader(h, v)
	return cr
}

func (cr *ClientRequest) SetParams(parms map[string]string) *ClientRequest {
	cr.r.SetQueryParams(parms)
	return cr
}

func (cr *ClientRequest) Post() (*ClientResponse, error) {
	return cr.execute("post", cr.path)
}

func (cr *ClientRequest) Get() (*ClientResponse, error) {

	return cr.execute("get", cr.path)
}

func (cr *ClientRequest) Put() (*ClientResponse, error) {
	return cr.execute("put", cr.path)
}

func (cr *ClientRequest) Delete() (*ClientResponse, error) {
	return cr.execute("delete", cr.path)
}

func (cr *ClientRequest) execute(kind, url string) (*ClientResponse, error) {
	var response *resty.Response
	var err error
	loggerHttp.Debug(fmt.Sprintf("calling %s: %s ", kind, cr.path))
	start := time.Now()
	defer func() {
		if cr.ins != nil {
			cr.ins.SetObserve(kind, cr.base, float64(time.Since(start)))
		}
		loggerHttp.Debug(fmt.Sprintf("calling %s: %s completed", kind, cr.path))
	}()

	switch kind {
	case "get":
		response, err = cr.r.Get(url)
	case "post":
		response, err = cr.r.Post(url)
	case "put":
		response, err = cr.r.Put(url)
	case "delete":
		response, err = cr.r.Delete(url)
	}

	if err != nil {
		if cr.ins != nil {
			cr.ins.SetResult(kind, cr.base, false)
		}
		loggerHttp.Error(err, fmt.Sprintf("calling: %s ", cr.path))
		return nil, err
	} else {
		if response.StatusCode() <= 299 {
			if cr.ins != nil {
				cr.ins.SetResult(kind, cr.base, true)
			}
		} else {
			if cr.ins != nil {
				cr.ins.SetResult(kind, cr.base, false)
			}
			loggerHttp.Error(nil, fmt.Sprintf("status code: %d ,calling: %s ", response.StatusCode(), cr.path))
		}

	}
	return cr.getClientResponse(response), err
}

func (cr *ClientRequest) getClientResponse(response *resty.Response) *ClientResponse {
	res := &ClientResponse{
		StatusCode: response.StatusCode(),
		Body:       (string)(response.Body()),
		reqStruct:  cr.reqStruct,
		resStruct:  cr.resStruct,
		errStruct:  cr.resStruct,
	}
	return res
}
