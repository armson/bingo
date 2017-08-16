package bingo

import(
    "context"
    "time"
    "net"
    "net/http"
    "math"
    "strings"
    "fmt"
    "io"
    "net/url"
    "encoding/json"
    "github.com/armson/bingo/utils"
	"github.com/armson/bingo/config"
	"bytes"
)

const (
    MIMEJSON              = "application/json"
    MIMEHTML              = "text/html"
    MIMEXML               = "application/xml"
    MIMEXML2              = "text/xml"
    MIMEPlain             = "text/plain"
    MIMEPOSTForm          = "application/x-www-form-urlencoded"
    MIMEMultipartPOSTForm = "multipart/form-data"
    MIMEPROTOBUF          = "application/x-protobuf"
)

var jsonContentType = []string{"application/json; charset=utf-8"}
var plainContentType = []string{"text/plain; charset=utf-8"}
var defaultCookieSetting = map[string]string{
    "path":"/",
    "domain":"/",
    "secure":"false",
    "httpOnly":"false",
}



const abortIndex int8 = math.MaxInt8 / 2

type Context struct {
    writermem responseWriter
    Request   *http.Request
    Writer    ResponseWriter

    Params   Params
    handlers HandlersChain
    index    int8

    engine   *Engine
    Keys     map[string]interface{}
    Errors   errorMsgs
    Accepted []string

    cookieSetting   map[string]string
    logs    []string
}
var _ context.Context = &Context{}

/************************************************/
/* 通过Get、MustGet、Set，可以在上下文中，传递变量     */
/************************************************/
func (c *Context) Get(key string) (value interface{}, exists bool) {
    if c.Keys != nil {
        value, exists = c.Keys[key]
    }
    return
}
func (c *Context) MustGet(key string) interface{} {
    if value, exists := c.Get(key); exists {
        return value
    }
    panic("Key \"" + key + "\" does not exist")
}
func (c *Context) Set(key string, value interface{}) {
    if c.Keys == nil {
        c.Keys = make(map[string]interface{})
    }
    c.Keys[key] = value
}
/**************************************/
/* 通过路由传递的变量                    */
/**************************************/
func (c *Context) Param(key string) string {
    return c.Params.ByName(key)
}
/**************************************/
/* 通过url获取变量，类似PHP中$_GET的方法   */
/**************************************/
func (c *Context) GET(key string) string {
    value, _ := c.GetQuery(key)
    return value
}
func (c *Context) GETs(key string) []string {
    values, _ := c.GetQueryArray(key)
    return values
}
func (c *Context) GetQuery(key string) (string, bool) {
    if values, ok := c.GetQueryArray(key); ok {
        return values[0], ok
    }
    return "", false
}
func (c *Context) DefaultQuery(key, defaultValue string) string {
    if value, ok := c.GetQuery(key); ok {
        return value
    }
    return defaultValue
}
func (c *Context) GetQueryArray(key string) ([]string, bool) {
    req := c.Request
    if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
        return values, true
    }
    return []string{}, false
}
/***********************************/
/* Post的数据，相当于PHP中的$_POST    */
/***********************************/
func (c *Context) POST(key string) string {
    value, _ := c.GetPostForm(key)
    return value
}
func (c *Context) POSTs(key string) []string {
    values, _ := c.GetPostFormArray(key)
    return values
}
func (c *Context) DefaultPostForm(key, defaultValue string) string {
    if value, ok := c.GetPostForm(key); ok {
        return value
    }
    return defaultValue
}
func (c *Context) GetPostForm(key string) (string, bool) {
    if values, ok := c.GetPostFormArray(key); ok {
        return values[0], ok
    }
    return "", false
}
func (c *Context) GetPostFormArray(key string) ([]string, bool) {
    req := c.Request
    req.ParseForm()
    req.ParseMultipartForm(32 << 20) // 32 MB
    if values := req.PostForm[key]; len(values) > 0 {
        return values, true
    }
    if req.MultipartForm != nil && req.MultipartForm.File != nil {
        if values := req.MultipartForm.Value[key]; len(values) > 0 {
            return values, true
        }
    }
    return []string{}, false
}
func (c *Context) PostFormQuery() (string) {
    req := c.Request
    req.ParseForm()
    return req.PostForm.Encode()
}

/***********************************/
/* cookie的相关读取、设置方法          */
/***********************************/
func (c *Context) SetCookiePath(path string) {
    if path == "" { path = "/" }
    c.cookieSetting["path"] = path
}
func (c *Context) SetCookieDomain(domain string) {
    c.cookieSetting["domain"] = domain
}
func (c *Context) SetCookie(name string,value string, maxAge int) {
    http.SetCookie(c.Writer, &http.Cookie{
        Name:     name,
        Value:    url.QueryEscape(value),
        MaxAge:   maxAge,
        Path:     c.cookieSetting["path"],
        Domain:   c.cookieSetting["domain"],
        Secure:   utils.String.Bool(c.cookieSetting["secure"]),
        HttpOnly: utils.String.Bool(c.cookieSetting["httpOnly"]),
    })
}
func (c *Context) UnsetCookie(name string) {
    c.SetCookie(name, "", -1)
}

func (c *Context) Cookie(name string) string {
    cookie, err := c.Request.Cookie(name)
    if err != nil {
        return ""
    }
    val, _ := url.QueryUnescape(cookie.Value)
    return val
}

func (c *Context) Cookies() (m map[string]string) {
    cookies := c.Request.Cookies()
    if len(cookies) == 0 {
        return
    }
    m = make(map[string]string)
    for _, v := range cookies {
        m[v.Name] = v.Value
    }
    return
}

/***********************************/
/* 响应输出方法 String和JSON          */
/***********************************/
func (c *Context) String(code int, format string, values ...interface{}) {
    c.Status(code)
	c.Header("Access-Control-Allow-Origin","*")
	c.Header("Access-Control-Allow-Methods","POST, GET, OPTIONS, PUT, DELETE, HEAD")
    writeContentType(c.Writer, plainContentType)
    if len(values) > 0 {
        fmt.Fprintf(c.Writer, format, values...)
    } else {
        io.WriteString(c.Writer, format)
    }
}
func (c *Context) JSON(code int, obj interface{}) {
    c.Status(code)
    c.Header("Access-Control-Allow-Origin","*")
    c.Header("Access-Control-Allow-Methods","POST, GET, OPTIONS, PUT, DELETE, HEAD")
    writeContentType(c.Writer, jsonContentType)
    if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
        panic(err)
    }
	if config.Bool("default","enableLog") && config.Bool("response","enableLog") {
		s ,_ := json.Marshal(obj)
		c.Logs("Response", string(s))
	}
}
func (c *Context) StringOK(format string, values ...interface{}) {
    c.String(http.StatusOK,format,values...)
}
func (c *Context) JsonOK(obj interface{}) {
    c.JSON(http.StatusOK, obj)
}

/***********************************/
/* 获取客户端请求头部信息以及IP        */
/***********************************/
func (c *Context) ContentType() string {
    return filterFlags(c.RequestHeader("Content-Type"))
}
func (c *Context) Header(key, value string) {
    if len(value) == 0 {
        c.Writer.Header().Del(key)
    } else {
        c.Writer.Header().Set(key, value)
    }
}
func (c *Context) ClientIP() string {
    if c.engine.ForwardedByClientIP {
        clientIP := strings.TrimSpace(c.RequestHeader("X-Real-Ip"))
        if len(clientIP) > 0 {
            return clientIP
        }
        clientIP = c.RequestHeader("X-Forwarded-For")
        if index := strings.IndexByte(clientIP, ','); index >= 0 {
            clientIP = clientIP[0:index]
        }
        clientIP = strings.TrimSpace(clientIP)
        if len(clientIP) > 0 {
            return clientIP
        }
    }
    if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
        return ip
    }
    return ""
}

/******************/
/* 设置响应信息     */
/******************/
func (c *Context) Status(code int) {
    c.writermem.WriteHeader(code)
}
func (c *Context) SetAccepted(formats ...string) {
    c.Accepted = formats
}

/************/
/* 跳转     */
/************/

func (c *Context) Redirect(location string) {
    http.Redirect(c.Writer, c.Request, location, 302)
}
/************/
/* 设置日志  */
/************/
func (c *Context) Logs(args ...string) {
	var message string
	if len(args) < 1 {
		panic("Func (c *Context) Logs params is shorter")
	}
	pose := utils.Int.String(len(c.logs))
	if len(args) == 1 {
		message = utils.String.Join("[Node##",pose,"] ",args[0])
	}
	if len(args) > 1 {
		message = utils.String.Join("[",args[0],"##",pose,"] ",args[1])
	}
	c.logs = append(c.logs , message)
}

func (c *Context) File(filepath string) {
    http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) Next() {
    c.index++
    s := int8(len(c.handlers))
    for ; c.index < s; c.index++ {
        c.handlers[c.index](c)
    }
}

func (c *Context) RequestHeader(key string) string {
    if values, _ := c.Request.Header[key]; len(values) > 0 {
        return values[0]
    }
    return ""
}

func (c *Context) RequestHeaderString() string {
	buf := bytes.Buffer{}
	for k, header := range c.Request.Header {
		if k != "Content-Type" {
			buf.WriteString(k)
			buf.WriteString(":")
			buf.WriteString(header[0])
			buf.WriteString(" ")
		}
	}
	buf.WriteString("Post-Query:")
	buf.WriteString(c.PostFormQuery())
	return buf.String()
}

func (c *Context) AbortWithStatus(code int) {
    c.Status(code)
    c.Writer.WriteHeaderNow()
    c.Abort()
}

func (c *Context) Abort() {
    c.index = abortIndex
}

func writeContentType(w http.ResponseWriter, value []string) {
    header := w.Header()
    if val := header["Content-Type"]; len(val) == 0 {
        header["Content-Type"] = value
    }
}

func (c *Context) reset() {
    c.Writer = &c.writermem
    c.Params = c.Params[0:0]
    c.handlers = nil
    c.index = -1
    c.Keys = nil
    c.Errors = c.Errors[0:0]
    c.Accepted = nil
    c.cookieSetting = defaultCookieSetting
    c.logs = []string{}
}
func (c *Context) HandlerName() string {
    return nameOfFunction(c.handlers.Last())
}

func (c *Context) IsAborted() bool {
    return c.index >= abortIndex
}
func (c *Context) AbortWithError(code int, err error) *Error {
    c.AbortWithStatus(code)
    return c.Error(err)
}
func (c *Context) Error(err error) *Error {
    var parsedError *Error
    switch err.(type) {
    case *Error:
        parsedError = err.(*Error)
    default:
        parsedError = &Error{
            Err:  err,
            Type: ErrorTypePrivate,
        }
    }
    c.Errors = append(c.Errors, parsedError)
    return parsedError
}

func (c *Context) Copy() *Context {
    var cp = *c
    cp.writermem.ResponseWriter = nil
    cp.Writer = &cp.writermem
    cp.index = abortIndex
    cp.handlers = nil
    return &cp
}

/**************************************/
/* Context的interface的实现的四个方法    */
/**************************************/
func (c *Context) Deadline() (deadline time.Time, ok bool) {
    return
}

func (c *Context) Done() <-chan struct{} {
    return nil
}

func (c *Context) Err() error {
    return nil
}

func (c *Context) Value(key interface{}) interface{} {
    if key == 0 {
        return c.Request
    }
    if keyAsString, ok := key.(string); ok {
        val, _ := c.Get(keyAsString)
        return val
    }
    return nil
}







