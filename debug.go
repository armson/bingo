package bingo

import (
    // "bytes"
    // "html/template"
    "log"
)

func init() {
    log.SetFlags(0)
}

// 待补充
func IsDebugging() bool {
    return runMode == debugCode
}

func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
    if IsDebugging() {
        nuHandlers := len(handlers)
        handlerName := nameOfFunction(handlers.Last())
        debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
    }
}

func debugPrint(format string, values ...interface{}) {
    if IsDebugging() {
        log.Printf("[debug] "+format, values...)
    }
}

func debugPrintError(err error) {
    if err != nil {
        debugPrint("[ERROR] %v\n", err)
    }
}

// func debugPrintLoadTemplate(tmpl *template.Template) {
//     if IsDebugging() {
//         var buf bytes.Buffer
//         for _, tmpl := range tmpl.Templates() {
//             buf.WriteString("\t- ")
//             buf.WriteString(tmpl.Name())
//             buf.WriteString("\n")
//         }
//         debugPrint("Loaded HTML Templates (%d): \n%s\n", len(tmpl.Templates()), buf.String())
//     }
// }


// func debugPrintWARNINGSetHTMLTemplate() {
//     debugPrint(`[WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
// at initialization. ie. before any route is registered or the router is listening in a socket:

//     router := gin.Default()
//     router.SetHTMLTemplate(template) // << good place

// `)
// }


