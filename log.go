package bingo

import (
	"bytes"
	"runtime"
	"io"
	"os"
	"path/filepath"
	"github.com/mattn/go-isatty"
	"fmt"
	"io/ioutil"
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type BinLog interface {
	FileName() string
	String() string
	Write() io.Writer
	Rename(os.FileInfo) bool
}

func NewFileWriter(log BinLog) io.Writer {
	fileName := log.FileName()
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		fmt.Errorf("Can't create %s folder on %v",log, err)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Errorf("Can't create %s file: %v", log, err)
	}
	return file
}

func ResetFileWriter(log BinLog){
    if w, ok := log.Write().(*os.File); ok {
        file ,err := w.Stat()
        if err != nil {
            fmt.Errorf("Can't Stat %s file: %v",log, err)
        }
		log.Rename(file)
    }
}
func IsTerm(log BinLog) bool {
	isTerm := true
	if w, ok := log.Write().(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		isTerm = false
	}
	return isTerm
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

//recovery用到的方法
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}



