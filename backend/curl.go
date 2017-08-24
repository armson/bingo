package backend

import(
    "bytes"
    "net/http"
    "net/url"
    "strings"
    "io/ioutil"
    "compress/gzip"
    "github.com/armson/bingo/utils"
    "time"
    "net"
)

func (this *HttpServer) Request() *http.Request {
    return this.request
}
func (this *HttpServer) Response() (*http.Response, error) {
    if this.response.StatusCode != 0 { return this.response, nil }
    response, err := this.ci()
    if err != nil { return nil, err }
    this.response = response
    return response, nil
}
func (this *HttpServer) Header(key, value string) *HttpServer {
    this.request.Header.Set(key, value)
    return this
}
func (this *HttpServer) Url() string {
    return this.request.URL.String()
}
func (this *HttpServer) Query() string {
    return utils.Map.HttpBuildQuery(this.params)
}
func (this *HttpServer) CostTime() string {
	return this.cost.String()
}
func (this *HttpServer) Bytes() ([]byte, error) {
    if this.body != nil { return this.body, nil }
    response, err := this.Response()
    if err != nil { return nil, err }
    if response.Body == nil { return nil, nil }
    defer response.Body.Close()

    if this.setting.Gzip && response.Header.Get("Content-Encoding") == "gzip" {
        reader, err := gzip.NewReader(response.Body)
        if err != nil { return nil, err }
        this.body, err = ioutil.ReadAll(reader)
        return this.body, err
    }
    this.body, err = ioutil.ReadAll(response.Body)
    return this.body, err
}
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        conn, err := net.DialTimeout(netw, addr, cTimeout)
        if err != nil { return nil, err }
        err = conn.SetDeadline(time.Now().Add(rwTimeout))
        return conn, err
    }
}
func (this *HttpServer) ci() (response *http.Response, err error) {
    query := utils.Map.HttpBuildQuery(this.params)

    if this.request.Method == "GET" && len(query) > 0 {
        if strings.Contains(this.url, "?") {
            this.url = utils.String.Join(this.url,"&",query)
        } else {
            this.url = utils.String.Join(this.url,"?",query)
        }
    }

    if (this.request.Method == "POST" || this.request.Method == "PUT" || this.request.Method == "DELETE") && this.request.Body == nil {
        if len(query) > 0 {
            this.Header("Content-Type", "application/x-www-form-urlencoded")
            bf := bytes.NewBufferString(query)
            this.request.Body = ioutil.NopCloser(bf)
            this.request.ContentLength = int64(len(query))
        }
    }

    url, err := url.Parse(this.url)
    if err != nil { return nil, err }
    this.request.URL = url

    trans := this.setting.Transport
    if trans == nil {
        trans = &http.Transport{
            TLSClientConfig:     this.setting.TLSConfig,
            Dial:                TimeoutDialer(this.setting.ConnectTimeout, this.setting.ReadWriteTimeout),
            MaxIdleConnsPerHost: -1,
        }
    } else {
        if t, ok := trans.(*http.Transport); ok {
            if t.TLSClientConfig == nil { t.TLSClientConfig = this.setting.TLSConfig }
            if t.Dial == nil { t.Dial = TimeoutDialer(this.setting.ConnectTimeout, this.setting.ReadWriteTimeout) }
        }
    }

    var jar http.CookieJar
    if this.setting.EnableCookie {
        if cookieJar == nil { createDefaultCookie() }
        jar = cookieJar
    }

    client := &http.Client{
        Transport: trans,
        Jar:       jar,
    }

    if this.setting.UserAgent != "" && this.request.Header.Get("User-Agent") == "" {
        this.request.Header.Set("User-Agent", this.setting.UserAgent)
    }

    if this.setting.CheckRedirect != nil {
        client.CheckRedirect = this.setting.CheckRedirect
    }

    start := time.Now()
    for i := 0; this.setting.Retries == -1 || i <= this.setting.Retries; i++ {
        response, err = client.Do(this.request)
        if err == nil { break }
    }
	this.cost = time.Since(start)
    return response, err
}


