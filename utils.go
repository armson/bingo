package bingo

import(
    "path"
    "reflect"
    "runtime"
    "github.com/armson/bingo/config"
)

func joinPaths(absolutePath, relativePath string) string {
    if len(relativePath) == 0 { 
        return absolutePath 
    }
    finalPath := path.Join(absolutePath, relativePath)
    appendSlash := lastChar(relativePath) == '/' && lastChar(finalPath) != '/'
    if appendSlash { 
        return finalPath + "/" 
    }
    return finalPath
}

func lastChar(str string) uint8 {
    size := len(str)
    if size == 0 { 
        panic("The length of the string can't be 0") 
    }
    return str[size-1]
}

func assert(guard bool, text string) {
    if !guard {
        panic(text)
    }
}
func nameOfFunction(f interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
func resolveAddress() string {
    var port string
    p  := config.String("httpPort")
    if p == "" {
        port = ":8080"
    } else {
        port = ":" + p
    }

    h := config.String("httpAddr")
    if h == "" {
        return h + port
    }
    return port
}
func filterFlags(content string) string {
    for i, char := range content {
        if char == ' ' || char == ';' {
            return content[:i]
        }
    }
    return content
}



