package curl

import(
    "net/http"
    "crypto/tls"
    "time"
    "fmt"
    "github.com/armson/bingo"
)

func SetDefaultSettings(setting Settings) {
    mutex.Lock()
    defer mutex.Unlock()
    defaultSettings = setting
}
func (this *HttpServer) Setting(setting Settings) *HttpServer {
    this.setting = setting
    return this
}
func (this *HttpServer) SetBasicAuth(username, password string) *HttpServer {
    this.request.SetBasicAuth(username, password)
    return this
}
func (this *HttpServer) SetEnableCookie(enable bool) *HttpServer {
    this.setting.EnableCookie = enable
    return this
}
func (this *HttpServer) SetUserAgent(useragent string) *HttpServer {
    this.setting.UserAgent = useragent
    return this
}
func (this *HttpServer) Retries(times int) *HttpServer {
    this.setting.Retries = times
    return this
}
func (this *HttpServer) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *HttpServer {
    this.setting.ConnectTimeout = connectTimeout
    this.setting.ReadWriteTimeout = readWriteTimeout
    return this
}
func (this *HttpServer) SetTLSConfig(setting *tls.Config) *HttpServer {
    this.setting.TLSConfig = setting
    return this
}
func (this *HttpServer) SetHost(host string) *HttpServer {
    this.request.Host = host
    return this
}
func (this *HttpServer) Cookie(cookie *http.Cookie) *HttpServer {
    this.request.Header.Add("Cookie", cookie.String())
    return this
}
func (this *HttpServer) Cookies(cookies []*http.Cookie) *HttpServer {
    if len(cookies) < 1 { return this }
    var s []string
    for _,cookie := range cookies { 
        s = append(s, cookie.String())
    }
    this.request.Header.Add("Cookie", bingo.Slice.Join(s,";"))
    return this
}

func (this *HttpServer) SetTransport(transport http.RoundTripper) *HttpServer {
    this.setting.Transport = transport
    return this
}
func (this *HttpServer) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *HttpServer {
    this.setting.CheckRedirect = redirect
    return this
}
func (this *HttpServer) Param(key , value string) *HttpServer {
    if param, ok := this.params[key]; ok {
        this.params[key] = append(param, value)
    } else {
        this.params[key] = []string{value}
    }
    fmt.Println(this.params)
    return this
}
func (this *HttpServer) Params(params map[string]string) *HttpServer {
    if len(params) < 1 { return this }
    for key, value := range params {
        if param, ok := this.params[key]; ok {
            this.params[key] = append(param, value)
        } else {
            this.params[key] = []string{value}
        }
    }
    fmt.Println(this.params)
    return this
}



