package curl

// 以POST举例，其他方法雷同
// c := curl.Post("http://open.100.com/client/v1/store/city?clientType=3")
// c.Param("debug","1")
// c.Params(map[string]string{"keyword":"海","cityId":"110100"})
// cookie ,_ := this.Ctx.Request.Cookie("anditiesAuther")
// c.Cookie(cookie)
// c.Cookies(this.Ctx.Request.Cookies())
// r , _ := c.Map()

import(
    "net/http"
    "net/http/cookiejar"
    "sync"
    "net/url"
    "log"
    "time"
    "crypto/tls"
    "encoding/json"
)

var cookieJar http.CookieJar
var mutex sync.Mutex

type HttpServer struct {
    request     *http.Request
    response    *http.Response
    url         string
    params      map[string][]string
    body        []byte
    setting     Settings
}
type Settings struct {
    UserAgent        string
    ConnectTimeout   time.Duration
    ReadWriteTimeout time.Duration
    TLSConfig        *tls.Config
    Transport        http.RoundTripper
    CheckRedirect    func(req *http.Request, via []*http.Request) error
    EnableCookie     bool
    Gzip             bool
    Retries          int
}
var defaultSettings = Settings{
    UserAgent:        "Bingo Server",
    ConnectTimeout:   60 * time.Second,
    ReadWriteTimeout: 60 * time.Second,
    EnableCookie:     true,
    Gzip:             true,
    Retries:          3,
}

func createDefaultCookie() {
    mutex.Lock()
    defer mutex.Unlock()
    cookieJar, _ = cookiejar.New(nil)
}

func createHttpServer(method ,rawurl string) *HttpServer {
    u, err := url.Parse(rawurl)
    if err != nil { log.Fatalln("Bingo HttpServer:", err) }
    request := http.Request{
        Method:     method,
        URL:        u,
        Header:     make(http.Header),
        Proto:      "HTTP/1.1",
        ProtoMajor: 1,
        ProtoMinor: 1,
    }
    return &HttpServer{
        request:     &request,
        response:    &http.Response{},
        url:         rawurl,
        params:      map[string][]string{},
        setting:     defaultSettings,
    }
}
func Get(rawurl string) *HttpServer {
    return createHttpServer("GET", rawurl)
}
func Post(rawurl string) *HttpServer {
    return createHttpServer("POST", rawurl)
}
func Put(rawurl string) *HttpServer {
    return createHttpServer("PUT", rawurl)
}
func Delete(rawurl string) *HttpServer {
    return createHttpServer("DELETE", rawurl)
}
func Head(rawurl string) *HttpServer {
    return createHttpServer("HEAD", rawurl)
}
func (this *HttpServer) String() (string, error) {
    data, err := this.Bytes()
    if err != nil { return "", err }
    return string(data), nil
}
func (this *HttpServer) Map() (m map[string]interface{}, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return m, err}
    return m, json.Unmarshal(data, &m)
}
func (this *HttpServer) Int() (i int, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return i, err}
    return i,json.Unmarshal(data, &i)
}
func (this *HttpServer) Float() (f float64, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return f, err}
    return f,json.Unmarshal(data, &f)
}
func (this *HttpServer) Bool() (b bool, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return b, err}
    return b,json.Unmarshal(data, &b)
}
func (this *HttpServer) Strings() (s []string, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return s, err}
    return s, json.Unmarshal(data, &s)
}
func (this *HttpServer) Ints() (i []int, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return i, err}
    return i,json.Unmarshal(data, &i)
}
func (this *HttpServer) Floats() (f []float64, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return f, err}
    return f,json.Unmarshal(data, &f)
}
func (this *HttpServer) Bools() (b []bool, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return b, err}
    return b,json.Unmarshal(data, &b)
}
func (this *HttpServer) Interfaces() (i []interface{}, err error) {
    var data []byte
    data, err = this.Bytes()
    if err != nil { return i, err}
    return i,json.Unmarshal(data, &i)
}

