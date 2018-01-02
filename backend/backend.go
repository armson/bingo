package backend

import (
	"github.com/armson/bingo"
	"github.com/bitly/go-simplejson"
	"github.com/armson/bingo/utils"
	"github.com/armson/bingo/config"
	"strings"
	"io"
	"runtime"
	"net/url"
	"net/http"
	"goxapi/common"
	"bytes"
)

type Backend struct {
	Tracer 				bingo.Tracer
	Url 				string
	Queries				map[string]string
	Method 				string
	files   			map[string]string
	readers				[]Reader
	raw 				interface{}
	headers				map[string]string
	closeLog			bool
	ResponseStatus		int
}

func (b *Backend) sent() ([]byte, error) {
	method := strings.ToUpper(b.Method)

	c := Handle(method, b.Url)
	if len(b.Queries) > 0 	{c.Params(b.Queries)}
	if len(b.files) > 0 	{c.Files(b.files)}
	if len(b.readers) > 0 	{c.Readers(b.readers)}
	if b.raw != nil 		{c.Raw(b.raw)}
	if len(b.headers) > 0 	{c.Headers(b.headers)}


	// 是否启用代理
	if common.ProxyEnable == true {
		domain := utils.String.Domain(b.Url)
		if utils.Slice.In("*",common.ProxyAcceptDomain) || utils.Slice.In(domain,common.ProxyAcceptDomain) {
			c.SetProxy(func(b *http.Request) (*url.URL, error) {
				u, _ := url.ParseRequestURI(common.ProxyAgentServer)
				return u, nil
			})
		}
	}

	data, err := c.Bytes()
	if err != nil {return nil, err}

	if config.Bool("default","enableLog") && config.Bool("backend","enableLog") && b.closeLog == false {
		message := utils.String.Join("Gx:", utils.Int.String(runtime.NumGoroutine()), " Cost:",c.CostTime()," ",method," ",c.Url())
		if method != "GET" && len(b.Queries) > 0 {
			message = utils.String.Join(message, " Query: ", c.Query())
		}
		if method != "GET" && b.raw != nil {
			message = utils.String.Join(message, " Body: ", c.GetRaw())
		}
		message = utils.String.Join(message, " Response: ", string(bytes.Replace(data,[]byte("\n"),nil, -1)))
		b.Tracer.Logs("Backend",message)
	}

	//每次请求发送完成后,重置条件
	b.Queries = map[string]string{}
	b.Url = ""
	b.files = map[string]string{}
	b.readers = []Reader{}
	b.raw = nil
	b.headers = map[string]string{}
	b.ResponseStatus = c.ResponseStatus()
	return data, nil
}

func (b *Backend) Send() (*simplejson.Json, error) {
	data, err := b.sent()
	if err != nil {
		return nil, err
	}
	return simplejson.NewJson(data)
}
func (b *Backend) String() (string, error) {
	data, err := b.sent()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (b *Backend) File(formName, fileName string) *Backend {
	if b.files == nil {
		b.files = map[string]string{}
	}
	b.files[formName] = fileName
	return b
}

func (b *Backend) Files(params map[string]string) *Backend {
	if len(params) < 1 { return b }
	if b.files == nil {
		b.files = map[string]string{}
	}
	for formName, fileName := range params {
		b.files[formName] = fileName
	}
	return b
}

func (b *Backend) Reader(formField , fileName string, rc io.ReadCloser) *Backend {
	b.readers = append(b.readers , Reader{
		formField:formField,
		fileName:fileName,
		rc:rc,
	})
	return b
}

func (b *Backend) Raw(obj interface{}) *Backend {
	if obj != nil {
		b.raw = obj
	}
	return b
}

func (b *Backend) CloseLog(close bool) *Backend {
	b.closeLog = close
	return b
}
func (b *Backend) SetHeader(field, value string) *Backend {
	if b.headers == nil {
		b.headers = map[string]string{}
	}
	b.headers[field] = value
	return b
}
func (b *Backend) SetHeaders(params map[string]string) *Backend {
	if len(params) < 1 { return b }
	if b.headers == nil {
		b.headers = map[string]string{}
	}
	for field, value := range params {
		b.headers[field] = value
	}
	return b
}


