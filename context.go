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
}
var _ context.Context = &Context{}

/************************************/

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


/************************************/
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


func (c *Context) Status(code int) {
    c.writermem.WriteHeader(code)
}

func (c *Context) ClientIP() string {
    if c.engine.ForwardedByClientIP {
        clientIP := strings.TrimSpace(c.requestHeader("X-Real-Ip"))
        if len(clientIP) > 0 {
            return clientIP
        }
        clientIP = c.requestHeader("X-Forwarded-For")
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

func (c *Context) requestHeader(key string) string {
    if values, _ := c.Request.Header[key]; len(values) > 0 {
        return values[0]
    }
    return ""
}

func (c *Context) AbortWithStatus(code int) {
    c.Status(code)
    c.Writer.WriteHeaderNow()
    c.Abort()
}

func (c *Context) Abort() {
    c.index = abortIndex
}


func (c *Context) String(code int, format string, values ...interface{}) {
    c.Status(code)
    writeContentType(c.Writer, plainContentType)
    if len(values) > 0 {
        fmt.Fprintf(c.Writer, format, values...)
    } else {
        io.WriteString(c.Writer, format)
    }
}


func (c *Context) JSON(code int, obj interface{}) {
    c.Status(code)
    writeContentType(c.Writer, jsonContentType)

    if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
        panic(err)
    }
}


func writeContentType(w http.ResponseWriter, value []string) {
    header := w.Header()
    if val := header["Content-Type"]; len(val) == 0 {
        header["Content-Type"] = value
    }
}

func (c *Context) Param(key string) string {
    return c.Params.ByName(key)
}

// It is shortcut for `c.Request.URL.Query().Get(key)`
//      GET /path?id=1234&name=Manu&value=
//      c.Query("id") == "1234"
//      c.Query("name") == "Manu"
//      c.Query("value") == ""
//      c.Query("wtf") == ""
func (c *Context) Query(key string) string {
    value, _ := c.GetQuery(key)
    return value
}
// GetQuery is like Query(), it returns the keyed url query value
// if it exists `(value, true)` (even when the value is an empty string),
// othewise it returns `("", false)`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//      GET /?name=Manu&lastname=
//      ("Manu", true) == c.GetQuery("name")
//      ("", false) == c.GetQuery("id")
//      ("", true) == c.GetQuery("lastname")
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

// GetQueryArray returns a slice of strings for a given query key, plus
// a boolean value whether at least one value exists for the given key.
func (c *Context) GetQueryArray(key string) ([]string, bool) {
    req := c.Request
    if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
        return values, true
    }
    return []string{}, false
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
}

func (c *Context) QueryArray(key string) []string {
    values, _ := c.GetQueryArray(key)
    return values
}



// PostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns an empty string `("")`.
func (c *Context) PostForm(key string) string {
    value, _ := c.GetPostForm(key)
    return value
}

// DefaultPostForm returns the specified key from a POST urlencoded form or multipart form
// when it exists, otherwise it returns the specified defaultValue string.
// See: PostForm() and GetPostForm() for further information.
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

// PostFormArray returns a slice of strings for a given form key.
// The length of the slice depends on the number of params with the given key.
func (c *Context) PostFormArray(key string) []string {
    values, _ := c.GetPostFormArray(key)
    return values
}

// GetPostFormArray returns a slice of strings for a given form key, plus
// a boolean value whether at least one value exists for the given key.
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

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
    return filterFlags(c.requestHeader("Content-Type"))
}

func (c *Context) SetAccepted(formats ...string) {
    c.Accepted = formats
}

func (c *Context) Header(key, value string) {
    if len(value) == 0 {
        c.Writer.Header().Del(key)
    } else {
        c.Writer.Header().Set(key, value)
    }
}

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

func (c *Context) Redirect(location string) {
    http.Redirect(c.Writer, c.Request, location, 302)
}


func (c *Context) Copy() *Context {
    var cp = *c
    cp.writermem.ResponseWriter = nil
    cp.Writer = &cp.writermem
    cp.index = abortIndex
    cp.handlers = nil
    return &cp
}





// // HTML renders the HTTP template specified by its file name.
// // It also updates the HTTP code and sets the Content-Type as "text/html".
// // See http://golang.org/doc/articles/wiki/
// func (c *Context) HTML(code int, name string, obj interface{}) {
//     instance := c.engine.HTMLRender.Instance(name, obj)
//     c.Render(code, instance)
// }

// // IndentedJSON serializes the given struct as pretty JSON (indented + endlines) into the response body.
// // It also sets the Content-Type as "application/json".
// // WARNING: we recommend to use this only for development propuses since printing pretty JSON is
// // more CPU and bandwidth consuming. Use Context.JSON() instead.
// func (c *Context) IndentedJSON(code int, obj interface{}) {
//     c.Render(code, render.IndentedJSON{Data: obj})
// }

// // JSON serializes the given struct as JSON into the response body.
// // It also sets the Content-Type as "application/json".


// // XML serializes the given struct as XML into the response body.
// // It also sets the Content-Type as "application/xml".
// func (c *Context) XML(code int, obj interface{}) {
//     c.Render(code, render.XML{Data: obj})
// }

// // YAML serializes the given struct as YAML into the response body.
// func (c *Context) YAML(code int, obj interface{}) {
//     c.Render(code, render.YAML{Data: obj})
// }

// // String writes the given string into the response body.


// // Redirect returns a HTTP redirect to the specific location.


// // Data writes some data into the body stream and updates the HTTP code.
// func (c *Context) Data(code int, contentType string, data []byte) {
//     c.Render(code, render.Data{
//         ContentType: contentType,
//         Data:        data,
//     })
// }

// // File writes the specified file into the body stream in a efficient way.


// // SSEvent writes a Server-Sent Event into the body stream.
// func (c *Context) SSEvent(name string, message interface{}) {
//     c.Render(-1, sse.Event{
//         Event: name,
//         Data:  message,
//     })
// }

// func (c *Context) Stream(step func(w io.Writer) bool) {
//     w := c.Writer
//     clientGone := w.CloseNotify()
//     for {
//         select {
//         case <-clientGone:
//             return
//         default:
//             keepOpen := step(w)
//             w.Flush()
//             if !keepOpen {
//                 return
//             }
//         }
//     }
// }

// /************************************/
// /******** CONTENT NEGOTIATION *******/
// /************************************/

// type Negotiate struct {
//     Offered  []string
//     HTMLName string
//     HTMLData interface{}
//     JSONData interface{}
//     XMLData  interface{}
//     Data     interface{}
// }

// func (c *Context) Negotiate(code int, config Negotiate) {
//     switch c.NegotiateFormat(config.Offered...) {
//     case binding.MIMEJSON:
//         data := chooseData(config.JSONData, config.Data)
//         c.JSON(code, data)

//     case binding.MIMEHTML:
//         data := chooseData(config.HTMLData, config.Data)
//         c.HTML(code, config.HTMLName, data)

//     case binding.MIMEXML:
//         data := chooseData(config.XMLData, config.Data)
//         c.XML(code, data)

//     default:
//         c.AbortWithError(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server"))
//     }
// }

// func (c *Context) NegotiateFormat(offered ...string) string {
//     assert1(len(offered) > 0, "you must provide at least one offer")

//     if c.Accepted == nil {
//         c.Accepted = parseAccept(c.requestHeader("Accept"))
//     }
//     if len(c.Accepted) == 0 {
//         return offered[0]
//     }
//     for _, accepted := range c.Accepted {
//         for _, offert := range offered {
//             if accepted == offert {
//                 return offert
//             }
//         }
//     }
//     return ""
// }



