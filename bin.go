package bingo

import (
    "sync"
    "net/http"
    "fmt"
    "os"
)

var once sync.Once
var internalEngine *Engine
var Exit func() = func(){ os.Exit(0) }


func engine() *Engine {
    once.Do(func() {
        internalEngine = Default()
    })
    return internalEngine
}

// POST is a shortcut for router.Handle("POST", path, handle)
func POST(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().POST(relativePath, handlers...)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func GET(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().GET(relativePath, handlers...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().DELETE(relativePath, handlers...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().PATCH(relativePath, handlers...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().PUT(relativePath, handlers...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().OPTIONS(relativePath, handlers...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().HEAD(relativePath, handlers...)
}

func Any(relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().Any(relativePath, handlers...)
}

func StaticFile(relativePath, filepath string) IRoutes {
    return engine().StaticFile(relativePath, filepath)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//     router.Static("/static", "/var/www")
func Static(relativePath, root string) IRoutes {
    return engine().Static(relativePath, root)
}

func StaticFS(relativePath string, fs http.FileSystem) IRoutes {
    return engine().StaticFS(relativePath, fs)
}

// Use attachs a global middleware to the router. ie. the middlewares attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func Use(middlewares ...HandlerFunc) IRoutes {
    return engine().Use(middlewares...)
}

// Run : The router is attached to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine undefinitelly unless an error happens.
func Run() (err error) {
    return engine().Run()
}

func Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
    return engine().Group(relativePath, handlers...)
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func NoRoute(handlers ...HandlerFunc) {
    engine().NoRoute(handlers...)
}

// NoMethod sets the handlers called when... TODO
func NoMethod(handlers ...HandlerFunc) {
    engine().NoMethod(handlers...)
}

func Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
    return engine().Handle(httpMethod, relativePath, handlers...)
}

func Routes() (routes RoutesInfo) {
    return engine().Routes()
}


func RunTLS(addr string, certFile string, keyFile string) (err error) {
    return engine().RunTLS(addr, certFile, keyFile)
}

func RunUnix(file string) (err error) {
    return engine().RunUnix(file)
}

func Echo(args ...interface{}){
    fmt.Println(args)
}

func Printf(s string, args ...interface{}){
    if len(args) < 1 {
        fmt.Println(s)
    } else {
        fmt.Printf(s+"\n", args...)
    }
}




