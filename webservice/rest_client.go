package webservice

import (
	"context"
	"fmt"
	"golibs/instrument"
	"golibs/logging"
	"time"

	"github.com/liornabat/golibs/tracing"
	"gopkg.in/resty.v1"
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
	ins            *instrument.InstrumentArray
	isInstrumented bool
	tracer         *tracing.Factory
}
type ClientRequest struct {
	r         *resty.Request
	base      string
	path      string
	ins       *instrument.InstrumentArray
	reqStruct interface{}
	resStruct interface{}
	errStruct interface{}
	ctx       context.Context
	tracer    *tracing.Factory
}

type ClientResponse struct {
	StatusCode int
	Body       string
	ReqStruct  interface{}
	ResStruct  interface{}
	ErrStruct  interface{}
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

func NewRestClientWithTracer(tracer *tracing.Factory, name, host, port, user, password string) *RestClient {
	c := &RestClient{
		Name:     name,
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Routes:   make(map[string]string),
		Headers:  make(map[string]string),
		tracer:   tracer,
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
		r:      resty.R(),
		ins:    c.ins,
		ctx:    ctx,
		tracer: c.tracer,
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

func (c *RestClient) NewAuthenticatedRequest(ctx context.Context, baseName string, authToken string) *ClientRequest {
	cr := &ClientRequest{
		r:   resty.R(),
		ins: c.ins,
	}
	cr.base = baseName
	cr.r.SetHeaders(c.Headers)
	if c.User != "" {
		cr.r.SetAuthToken(authToken)
	}
	cr.r.SetContext(ctx)
	cr.path = fmt.Sprintf("%s:%s%s", c.Host, c.Port, c.Routes[baseName])
	return cr
}

func (c *RestClient) EnableInstrumenting(nameSpace, subSystem string) *RestClient {
	ins := instrument.NewInstrumentArray(nameSpace, subSystem, c.Name).
		AddCounter([]string{"type", "func", "result"}, "counters for total results of function calls").
		AddHistogram([]string{"type", "func"}, []float64{0.01, 0.5, 1, 2, 5, 10, 20, 30}, "histogram for stats results of function calls")

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
	if cr.errStruct == nil {
		cr.r.SetError(res)
		cr.errStruct = res
	}
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
	if cr.tracer != nil {
		return cr.executeWithTracing("post", cr.path)
	}
	return cr.execute("post", cr.path)
}

func (cr *ClientRequest) Get() (*ClientResponse, error) {
	if cr.tracer != nil {
		return cr.executeWithTracing("get", cr.path)
	}
	return cr.execute("get", cr.path)
}

func (cr *ClientRequest) Put() (*ClientResponse, error) {
	if cr.tracer != nil {
		return cr.executeWithTracing("put", cr.path)
	}
	return cr.execute("put", cr.path)
}

func (cr *ClientRequest) Delete() (*ClientResponse, error) {
	if cr.tracer != nil {
		return cr.executeWithTracing("delete", cr.path)
	}
	return cr.execute("delete", cr.path)
}

func (cr *ClientRequest) execute(kind, url string) (*ClientResponse, error) {
	var response *resty.Response
	var err error
	loggerHttp.Debug(fmt.Sprintf("calling %s: %s ", kind, cr.path))
	start := time.Now()
	defer func() {
		if cr.ins != nil {
			cr.ins.ObserveHistogram(float64(time.Since(start)), kind, cr.base)
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
			cr.ins.IncToCounter(kind, cr.base, "Fail")
		}
		loggerHttp.Error(err, fmt.Sprintf("calling: %s ", cr.path))
		return nil, err
	} else {
		if response.StatusCode() <= 299 {
			if cr.ins != nil {
				cr.ins.IncToCounter(kind, cr.base, "Ok")
			}
		} else {
			if cr.ins != nil {
				cr.ins.IncToCounter(kind, cr.base, "Fail")
			}
			loggerHttp.Error(nil, fmt.Sprintf("status code: %d ,calling: %s ", response.StatusCode(), cr.path))
		}

	}
	return cr.getClientResponse(response), err
}

func (cr *ClientRequest) executeWithTracing(kind, url string) (*ClientResponse, error) {
	_, span := tracing.StartSpan(cr.ctx, "rest_client/execute")
	defer span.Finish()

	span.SetHTTPUrl(url)
	span.LogKV("path", cr.path)
	var response *resty.Response
	var err error
	loggerHttp.Debug(fmt.Sprintf("calling %s: %s ", kind, cr.path))
	start := time.Now()
	defer func() {
		if cr.ins != nil {
			cr.ins.ObserveHistogram(float64(time.Since(start)), kind, cr.base)
		}
		loggerHttp.Debug(fmt.Sprintf("calling %s: %s completed", kind, cr.path))
	}()

	switch kind {
	case "get":
		span.SetHTTPMethod("GET")
		span.LogKV("parameters", cr.r.URL)
		response, err = cr.r.Get(url)

	case "post":
		span.SetHTTPMethod("POST")
		span.LogKV("body", fmt.Sprintf("%v", cr.r.Body))
		response, err = cr.r.Post(url)
	case "put":
		span.SetHTTPMethod("PUT")
		span.LogKV("body", fmt.Sprintf("%v", cr.r.Body))
		response, err = cr.r.Put(url)
	case "delete":
		span.SetHTTPMethod("DELETE")
		span.LogKV("body", fmt.Sprintf("%v", cr.r.Body))
		response, err = cr.r.Delete(url)
	}

	if err != nil {
		if cr.ins != nil {
			cr.ins.IncToCounter(kind, cr.base, "Fail")
		}
		loggerHttp.Error(err, fmt.Sprintf("calling: %s ", cr.path))
		span.SetError(err)

		return nil, err
	} else {
		span.SetHTTPStatusCode(uint16(response.StatusCode()))
		if response.StatusCode() <= 299 {
			if cr.ins != nil {
				cr.ins.IncToCounter(kind, cr.base, "Ok")
			}
		} else {
			if cr.ins != nil {
				cr.ins.IncToCounter(kind, cr.base, "Fail")
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
		ReqStruct:  cr.reqStruct,
		ResStruct:  cr.resStruct,
		ErrStruct:  cr.resStruct,
	}
	return res
}
