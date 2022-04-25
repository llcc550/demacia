package protocol

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/levigross/grequests"
	"gitlab.u-jy.cn/xiaoyang/go-zero/core/logx"
)

const (
	ApiServiceSuccess = "0"
	ApiHttpSuccess    = http.StatusOK

	HuaweiApiServiceSuccess = "000000"

	ErrApiUrlNil              = "Illegal URL"
	ErrGoRequest              = "Request Error"
	ErrServiceRequestCode     = "HTTP Response Code Error"
	ErrServiceRequestData     = "HTTP Response Data Error"
	ErrServiceRequestJsonData = "HTTP Response Json Data Error"
	ErrServiceRequestXmlData  = "HTTP Response Xml Data Error"

	ContentTypeForm = "application/x-www-form-urlencoded; charset=UTF-8"
	ContentTypeJson = "application/json; charset=UTF-8"

	AcceptJson = "application/json"
	AcceptXml  = "application/xml"

	HeaderContentTypeKey = "Content-Type"
	HeaderAcceptKey      = "Accept"
	HeaderAuthorization  = "Authorization"
	HeaderWSSE           = "X-WSSE"
)

type DefaultProtocol struct {
	Data interface{}
}

func (p *DefaultProtocol) GetCode() string {
	return ApiServiceSuccess
}
func (p *DefaultProtocol) GetMsg() string {
	return ""
}
func (p *DefaultProtocol) GetData() interface{} {
	return p.Data
}

type Getter interface {
	SetData(ApiProtocol)
	Get() error
}

type ApiProtocol interface {
	GetCode() string
	GetMsg() string
	GetData() interface{}
}

type ApiFile grequests.FileUpload

type ApiOption struct {
	Url                string
	Method             string // net/http.MethodGet
	Data               ApiProtocol
	Param              map[string]string
	Query              map[string]string
	Header             map[string]string
	Auth               []string
	Cookie             map[string]string
	Json               interface{}
	Files              []ApiFile
	Body               io.Reader
	AcceptCodes        []int
	Client             *http.Client
	InsecureSkipVerify bool
	Host               string
}

func (opt *ApiOption) SetData(data ApiProtocol) {
	opt.Data = data
}

type ApiError struct {
	httpCode int
	code     string
	msg      string
}

func (e *ApiError) Error() string {
	return e.msg
}

func (e *ApiError) HttpCode() int {
	return e.httpCode
}

func (e *ApiError) Code() string {
	return e.code
}

func NewApiError(httpCode int, code string, msg string) *ApiError {
	return &ApiError{
		httpCode: httpCode,
		code:     code,
		msg:      msg,
	}
}

func (opt *ApiOption) checkHttpCode(code int) bool {
	for _, accept := range opt.AcceptCodes {
		if accept == code {
			return true
		}
	}
	return false
}

func (opt *ApiOption) Get() error {
	logx.Infof("api opt %+v", opt)
	gOpt := grequests.RequestOptions{
		Data:               opt.Query,
		Params:             opt.Param,
		Headers:            opt.Header,
		InsecureSkipVerify: opt.InsecureSkipVerify,
		Auth:               opt.Auth,
		Host:               opt.Host,
		JSON:               opt.Json,
		RequestBody:        opt.Body,
		HTTPClient:         opt.Client,
		Cookies: func() []*http.Cookie {
			var cookies []*http.Cookie
			if len(opt.Cookie) > 0 {
				for key, value := range opt.Cookie {
					cookies = append(cookies, &http.Cookie{
						Name:  key,
						Value: value,
					})
				}
			}
			return cookies
		}(),
		Files: func() []grequests.FileUpload {
			if opt.Files == nil || len(opt.Files) == 0 {
				return nil
			}
			var files []grequests.FileUpload
			for _, value := range opt.Files {
				files = append(files, grequests.FileUpload{
					FileName:     value.FileName,
					FileContents: value.FileContents,
					FieldName:    value.FieldName,
				})
			}
			return files
		}(),
	}
	if len(opt.AcceptCodes) == 0 {
		opt.AcceptCodes = append(opt.AcceptCodes, ApiHttpSuccess)
	}

	if len(opt.Url) == 0 {
		return NewApiError(0, ErrApiUrlNil, opt.Url)
	}

	resp, err := grequests.Req(func() string {
		if len(opt.Method) == 0 {
			opt.Method = http.MethodGet
		}
		return opt.Method
	}(), opt.Url, &gOpt)
	if err != nil {
		return NewApiError(0, ErrGoRequest, err.Error())
	}

	defer resp.Close()

	logx.Info(resp.StatusCode, opt.Url)
	if resp.Error != nil {
		return NewApiError(0, ErrGoRequest, resp.Error.Error())
	}

	if !opt.checkHttpCode(resp.StatusCode) {
		return NewApiError(resp.StatusCode, ErrServiceRequestCode, ErrServiceRequestCode)
	}
	if opt.Data != nil {
		if opt.Data.GetData() != nil {
			if v, ok := opt.Header[HeaderAcceptKey]; !ok {
				return NewApiError(resp.StatusCode, ErrServiceRequestData, resp.String())
			} else if v == AcceptJson {
				err = json.Unmarshal([]byte(resp.String()), opt.Data.GetData())
				if err != nil {
					return NewApiError(resp.StatusCode, ErrServiceRequestJsonData, resp.String())
				}
			} else if v == AcceptXml {
				err = xml.Unmarshal([]byte(resp.String()), opt.Data.GetData())
				if err != nil {
					return NewApiError(resp.StatusCode, ErrServiceRequestXmlData, resp.String())
				}
			}
		}

		if opt.Data.GetCode() != ApiServiceSuccess {
			return NewApiError(resp.StatusCode, opt.Data.GetCode(), opt.Data.GetMsg())
		}
	}
	return nil
}
