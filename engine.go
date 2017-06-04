package bingo

import(
    "sync"
    "net/http"
    "os"
    "net"
)

var default404Body = []byte("404 page not found")
var default405Body = []byte("405 method not allowed")
var mimePlain = []string{MIMEPlain}


type HandlerFunc func(*Context)
type HandlersChain []HandlerFunc

type Engine struct {
        RouterGroup
        // //HTMLRender  render.HTMLRender
        allNoRoute  HandlersChain
        allNoMethod HandlersChain
        noRoute     HandlersChain
        noMethod    HandlersChain
        pool        sync.Pool
        trees       methodTrees

        RedirectTrailingSlash bool

        RedirectFixedPath bool

        HandleMethodNotAllowed bool
        ForwardedByClientIP    bool
}

var _ IRouter = &Engine{}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
    assert(path[0] == '/', "path must begin with '/'")
    assert(len(method) > 0, "HTTP method can not be empty")
    assert(len(handlers) > 0, "there must be at least one handler")

    debugPrintRoute(method, path, handlers)
    root := engine.trees.get(method)
    if root == nil {
        root = new(node)
        engine.trees = append(engine.trees, methodTree{method: method, root: root})
    }
    root.addRoute(path, handlers)
}

func (c HandlersChain) Last() HandlerFunc {
    length := len(c)
    if length > 0 {
        return c[length-1]
    }
    return nil
}
// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
func New() *Engine {
    engine := &Engine{
        RouterGroup: RouterGroup{
            Handlers: nil,
            basePath: "/",
            root:     true,
        },
        RedirectTrailingSlash:  true,
        RedirectFixedPath:      false,
        HandleMethodNotAllowed: false,
        ForwardedByClientIP:    true,
        trees:                  make(methodTrees, 0, 9),
    }
    engine.RouterGroup.engine = engine
    engine.pool.New = func() interface{} {
        return engine.allocateContext()
    }
    return engine
}

func Default() *Engine {
    engine := New()
    engine.Use(Logger(), Recovery())
    return engine
}

func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
    engine.RouterGroup.Use(middleware...)
    engine.rebuild404Handlers()
    engine.rebuild405Handlers()
    return engine
}

func (engine *Engine) rebuild404Handlers() {
    engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
    engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

func redirectTrailingSlash(c *Context) {
    req := c.Request
    path := req.URL.Path
    code := 301 // Permanent redirect, request with GET method
    if req.Method != "GET" {
        code = 307
    }

    if len(path) > 1 && path[len(path)-1] == '/' {
        req.URL.Path = path[:len(path)-1]
    } else {
        req.URL.Path = path + "/"
    }
    debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
    http.Redirect(c.Writer, req, req.URL.String(), code)
    c.writermem.WriteHeaderNow()
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
    req := c.Request
    path := req.URL.Path

    fixedPath, found := root.findCaseInsensitivePath(
        cleanPath(path),
        trailingSlash,
    )
    if found {
        code := 301 // Permanent redirect, request with GET method
        if req.Method != "GET" {
            code = 307
        }
        req.URL.Path = string(fixedPath)
        debugPrint("redirecting request %d: %s --> %s", code, path, req.URL.String())
        http.Redirect(c.Writer, req, req.URL.String(), code)
        c.writermem.WriteHeaderNow()
        return true
    }
    return false
}
func (engine *Engine) allocateContext() *Context {
    return &Context{engine: engine}
}
func (engine *Engine) Run() (err error) {
    defer func() { debugPrintError(err) }()

    address := resolveAddress()
    debugPrint("Listening and serving HTTP on %s\n", address)
    err = http.ListenAndServe(address, engine)
    return
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := engine.pool.Get().(*Context)
    c.writermem.reset(w)
    c.Request = req
    c.reset()

    engine.handleHTTPRequest(c)
    engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(context *Context) {
    httpMethod := context.Request.Method
    path := context.Request.URL.Path

    // Find root of the tree for the given HTTP method
    t := engine.trees
    for i, tl := 0, len(t); i < tl; i++ {
        if t[i].method == httpMethod {
            root := t[i].root
            // Find route in tree
            handlers, params, tsr := root.getValue(path, context.Params)
            if handlers != nil {
                context.handlers = handlers
                context.Params = params
                context.Next()
                context.writermem.WriteHeaderNow()
                return

            } else if httpMethod != "CONNECT" && path != "/" {
                if tsr && engine.RedirectTrailingSlash {
                    redirectTrailingSlash(context)
                    return
                }
                if engine.RedirectFixedPath && redirectFixedPath(context, root, engine.RedirectFixedPath) {
                    return
                }
            }
            break
        }
    }

    // TODO: unit test
    if engine.HandleMethodNotAllowed {
        for _, tree := range engine.trees {
            if tree.method != httpMethod {
                if handlers, _, _ := tree.root.getValue(path, nil); handlers != nil {
                    context.handlers = engine.allNoMethod
                    serveError(context, 405, default405Body)
                    return
                }
            }
        }
    }
    context.handlers = engine.allNoRoute
    serveError(context, 404, default404Body)
}

func serveError(c *Context, code int, defaultMessage []byte) {
    c.writermem.status = code
    c.Next()
    if !c.writermem.Written() {
        if c.writermem.Status() == code {
            c.writermem.Header()["Content-Type"] = mimePlain
            c.Writer.Write(defaultMessage)
        } else {
            c.writermem.WriteHeaderNow()
        }
    }
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
    engine.noRoute = handlers
    engine.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
    engine.noMethod = handlers
    engine.rebuild405Handlers()
}

type RoutesInfo []RouteInfo
type RouteInfo  struct {
    Method  string
    Path    string
    Handler string
}

func (engine *Engine) Routes() (routes RoutesInfo) {
    for _, tree := range engine.trees {
        routes = iterate("", tree.method, routes, tree.root)
    }
    return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
    path += root.path
    if len(root.handlers) > 0 {
        routes = append(routes, RouteInfo{
            Method:  method,
            Path:    path,
            Handler: nameOfFunction(root.handlers.Last()),
        })
    }
    for _, child := range root.children {
        routes = iterate(path, method, routes, child)
    }
    return routes
}

func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) (err error) {
    debugPrint("Listening and serving HTTPS on %s\n", addr)
    defer func() { debugPrintError(err) }()

    err = http.ListenAndServeTLS(addr, certFile, keyFile, engine)
    return
}

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) {
    debugPrint("Listening and serving HTTP on unix:/%s", file)
    defer func() { debugPrintError(err) }()

    os.Remove(file)
    listener, err := net.Listen("unix", file)
    if err != nil {
        return
    }
    defer listener.Close()
    err = http.Serve(listener, engine)
    return
}



// func (engine *Engine) LoadHTMLGlob(pattern string) {
//     if IsDebugging() {
//         debugPrintLoadTemplate(template.Must(template.ParseGlob(pattern)))
//         engine.HTMLRender = render.HTMLDebug{Glob: pattern}
//     } else {
//         templ := template.Must(template.ParseGlob(pattern))
//         engine.SetHTMLTemplate(templ)
//     }
// }

// func (engine *Engine) LoadHTMLFiles(files ...string) {
//     if IsDebugging() {
//         engine.HTMLRender = render.HTMLDebug{Files: files}
//     } else {
//         templ := template.Must(template.ParseFiles(files...))
//         engine.SetHTMLTemplate(templ)
//     }
// }

// func (engine *Engine) SetHTMLTemplate(templ *template.Template) {
//     if len(engine.trees) > 0 {
//         debugPrintWARNINGSetHTMLTemplate()
//     }
//     engine.HTMLRender = render.HTMLProduction{Template: templ}
// }











// // Conforms to the http.Handler interface.








