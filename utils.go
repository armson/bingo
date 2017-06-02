package bingo

import(
    "path"
    "reflect"
    "runtime"
    "os"
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
func resolveAddress(addr []string) string {
    switch len(addr) {
    case 0:
        if port := os.Getenv("PORT"); len(port) > 0 {
            debugPrint("Environment variable PORT=\"%s\"", port)
            return ":" + port
        }
        debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
        return ":8080"
    case 1:
        return addr[0]
    default:
        panic("too much parameters")
    }
}
func filterFlags(content string) string {
    for i, char := range content {
        if char == ' ' || char == ';' {
            return content[:i]
        }
    }
    return content
}

