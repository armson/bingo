package backend

import (
	"github.com/armson/bingo/utils"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// 以POST举例，其他方法雷同
// c := backend.Post("http://open.100.com/client/v1/store/city?clientType=3")
// c.Params(map[string]string{"keyword":"海","cityId":"110100"})
// cookie ,_ := this.Ctx.Request.Cookie("anditiesAuther")
// c.Cookie(cookie)
// c.Cookies(this.Ctx.Request.Cookies())
// r , _ := c.Map()

var defaultHttpSettings = HttpSettings{
	UserAgent:        "Bingo Server",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	EnableCookie:     true,
	Gzip:             true,
	Retries:          3,
}

var defaultCookieJar http.CookieJar
var settingMutex sync.Mutex

func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

func SetDefaultSetting(setting HttpSettings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultHttpSettings = setting
}

func newHttpServer(method ,rawurl string) *HttpServer {
	u, err := url.Parse(rawurl)
	if err != nil { log.Fatalln("Bingo HttpServer:", err) }

	request := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	return &HttpServer{
		url:     	rawurl,
		request:    &request,
		params:  	map[string][]string{},
		files:   	map[string]string{},
		readers: 	[]Reader{},
		setting: 	defaultHttpSettings,
		response:   &http.Response{},
	}
}

func Get(rawurl string) *HttpServer 			{	return newHttpServer("GET", rawurl)	}
func Post(rawurl string) *HttpServer		 	{	return newHttpServer("POST", rawurl) }
func Put(rawurl string) *HttpServer 			{	return newHttpServer("PUT", rawurl) }
func Delete(rawurl string) *HttpServer 			{	return newHttpServer("DELETE", rawurl) }
func Head(rawurl string) *HttpServer 			{	return newHttpServer("HEAD", rawurl) }
func Handle(method,rawurl string) *HttpServer 	{	return newHttpServer(method, rawurl) }

// BeegoHTTPSettings is the http.Client setting
type HttpSettings struct {
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	EnableCookie     bool
	Gzip             bool
	Retries          int // if set to -1 means will retry forever
}

type HttpServer struct {
	request     *http.Request
	response    *http.Response
	url     	string
	body    	[]byte

	params  	map[string][]string
	files   	map[string]string
	readers		[]Reader
	raw 		interface{}

	setting		HttpSettings
	cost		time.Duration
}

type Reader struct {
	formField string
	fileName string
	rc io.ReadCloser
}

func (b *HttpServer) Request() *http.Request {	return b.request }
func (b *HttpServer) Setting(setting HttpSettings) *HttpServer {
	b.setting = setting
	return b
}

func (b *HttpServer) SetBasicAuth(username, password string) *HttpServer {
	b.request.SetBasicAuth(username, password)
	return b
}
func (b *HttpServer) SetEnableCookie(enable bool) *HttpServer {
	b.setting.EnableCookie = enable
	return b
}

func (b *HttpServer) SetUserAgent(useragent string) *HttpServer {
	b.setting.UserAgent = useragent
	return b
}

func (b *HttpServer) Retries(times int) *HttpServer {
	b.setting.Retries = times
	return b
}

func (b *HttpServer) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *HttpServer {
	b.setting.ConnectTimeout = connectTimeout
	b.setting.ReadWriteTimeout = readWriteTimeout
	return b
}

func (b *HttpServer) SetTLSClientConfig(setting *tls.Config) *HttpServer {
	b.setting.TLSClientConfig = setting
	return b
}
func (b *HttpServer) SetHost(host string) *HttpServer {
	b.request.Host = host
	return b
}

func (b *HttpServer) Header(key, value string) *HttpServer {
	b.request.Header.Set(key, value)
	return b
}
func (b *HttpServer) Headers(params map[string]string) *HttpServer {
	for key, value := range params {
		b.request.Header.Set(key, value)
	}
	return b
}

func (b *HttpServer) SetProtocolVersion(vers string) *HttpServer {
	if len(vers) == 0 {
		vers = "HTTP/1.1"
	}
	if major, minor, ok := http.ParseHTTPVersion(vers); ok {
		b.request.Proto = vers
		b.request.ProtoMajor = major
		b.request.ProtoMinor = minor
	}
	return b
}

func (b *HttpServer) SetCookie(cookie *http.Cookie) *HttpServer {
	b.request.Header.Add("Cookie", cookie.String())
	return b
}

func (b *HttpServer) SetCookies(cookies []*http.Cookie) *HttpServer {
	if len(cookies) < 1 { return b }
	var s []string
	for _,cookie := range cookies {
		s = append(s, cookie.String())
	}
	b.request.Header.Add("Cookie", strings.Join(s,";"))
	return b
}

func (b *HttpServer) SetTransport(transport http.RoundTripper) *HttpServer {
	b.setting.Transport = transport
	return b
}

//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (b *HttpServer) SetProxy(proxy func(*http.Request) (*url.URL, error)) *HttpServer {
	b.setting.Proxy = proxy
	return b
}

func (b *HttpServer) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *HttpServer {
	b.setting.CheckRedirect = redirect
	return b
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (b *HttpServer) Param(key, value string) *HttpServer {
	if param, ok := b.params[key]; ok {
		b.params[key] = append(param, value)
	} else {
		b.params[key] = []string{value}
	}
	return b
}

func (b *HttpServer) Params(params map[string]string) *HttpServer {
	if len(params) < 1 { return b }
	for key, value := range params {
		if param, ok := b.params[key]; ok {
			b.params[key] = append(param, value)
		} else {
			b.params[key] = []string{value}
		}
	}
	return b
}

func (b *HttpServer) File(formName, fileName string) *HttpServer {
	b.files[formName] = fileName
	return b
}

func (b *HttpServer) Files(params map[string]string) *HttpServer {
	if len(params) < 1 { return b }
	for formName, fileName := range params {
		b.files[formName] = fileName
	}
	return b
}

func (b *HttpServer) Reader(formField , fileName string, rc io.ReadCloser) *HttpServer {
	b.readers = append(b.readers , Reader{
		formField:formField,
		fileName:fileName,
		rc:rc,
	})
	return b
}

func (b *HttpServer) Readers(r []Reader) *HttpServer {
	b.readers = r
	return b
}

func (b *HttpServer) Raw(obj interface{}) *HttpServer {
	if obj != nil {
		b.raw = obj
	}
	return b
}

func (b *HttpServer) GetRaw() string {
	bs, err := json.Marshal(b.raw)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (b *HttpServer) Url() string {
	return b.request.URL.String()
}

func (b *HttpServer) Query() string {
	return utils.Map.HttpBuildQuery(b.params)
}

func (b *HttpServer) CostTime() string {
	return b.cost.String()
}

func (b *HttpServer) Response() (*http.Response, error) {
	if b.response.StatusCode != 0 {
		return b.response, nil
	}
	response, err := b.ci()
	if err != nil {
		return nil, err
	}
	b.response = response
	return response, nil
}
func (b *HttpServer) ResponseStatus() int {
	return b.response.StatusCode
}

func (b *HttpServer) XML(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

func (b *HttpServer) Bytes() ([]byte, error) {
	if b.body != nil {
		return b.body, nil
	}

	response, err := b.Response()
	if err != nil {
		return nil, err
	}

	if response.Body == nil {
		return nil, nil
	}
	defer response.Body.Close()

	if b.setting.Gzip && response.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		b.body, err = ioutil.ReadAll(reader)
		return b.body, err
	}
	b.body, err = ioutil.ReadAll(response.Body)
	return b.body, err
}

func (b *HttpServer) String() (string, error) {
	data, err := b.Bytes()
	if err != nil { return "", err }
	return string(data), nil
}
func (b *HttpServer) Map() (m map[string]interface{}, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return m, err}
	return m, json.Unmarshal(data, &m)
}
func (b *HttpServer) Int() (i int, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return i, err}
	return i,json.Unmarshal(data, &i)
}
func (b *HttpServer) Float() (f float64, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return f, err}
	return f,json.Unmarshal(data, &f)
}
func (b *HttpServer) Bool() (t bool, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return t, err}
	return t,json.Unmarshal(data, &t)
}
func (b *HttpServer) Strings() (s []string, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return s, err}
	return s, json.Unmarshal(data, &s)
}
func (b *HttpServer) Ints() (i []int, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return i, err}
	return i,json.Unmarshal(data, &i)
}
func (b *HttpServer) Floats() (f []float64, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return f, err}
	return f,json.Unmarshal(data, &f)
}
func (b *HttpServer) Bools() (t []bool, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return t, err}
	return t,json.Unmarshal(data, &t)
}
func (b *HttpServer) Interfaces() (i []interface{}, err error) {
	var data []byte
	data, err = b.Bytes()
	if err != nil { return i, err}
	return i,json.Unmarshal(data, &i)
}

// 将响应的内容写入文件，一般用户用于处理附件
// ToFile不于Bytes()同时使用。
func (b *HttpServer) ToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	response, err := b.Response()
	if err != nil {
		return err
	}
	if response.Body == nil {
		return nil
	}
	defer response.Body.Close()

	_, err = io.Copy(f, response.Body)
	return err
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil { return nil, err }
		err = conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, err
	}
}