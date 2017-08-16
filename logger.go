package bingo

import (
    "fmt"
    "io"
    "os"
    "time"
    "github.com/mattn/go-isatty"
    "github.com/armson/bingo/config"
    "github.com/armson/bingo/utils"
    "path/filepath"
    "sync"
	"strings"
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

var defaultLoggerWriter io.Writer
var onlyOne sync.Once

func SetLoggerWriter() {
    onlyOne.Do(func() {
        if runMode == debugCode {
            defaultLoggerWriter = os.Stdout
        }
        if runMode == releaseCode {
            defaultLoggerWriter = NewFileWriter()
        }
    })
}

func NewFileWriter() io.Writer {
    fileName := config.String("accessLog")
    if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
        fmt.Errorf("Can't create accessLog folder on %v", err)
    }
    file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
    if err != nil {
        fmt.Errorf("Can't create accessLog file: %v", err)
    }
    return file
}

func ResetLoggerWriter(){
    if w, ok := defaultLoggerWriter.(*os.File); ok {
        file ,err := w.Stat()
        if err != nil {
            fmt.Errorf("Can't Stat accessLog file: %v", err)
        }
        modTime := file.ModTime()
        nowTime := time.Now()
        if modTime.Year() != nowTime.Year() || modTime.YearDay() != nowTime.YearDay() || modTime.Hour() != nowTime.Hour() {
            format := "2006010215"
            fileName := config.String("accessLog")
            err := os.Rename(fileName, fileName+"."+modTime.Format(format))
            if err != nil {
                fmt.Errorf("Can't Rename [%s] file: %v", fileName, err)
            }
            defaultLoggerWriter = NewFileWriter()
        }
    }
}

func Logger() HandlerFunc {
    SetLoggerWriter()
    isTerm := true
    if w, ok := defaultLoggerWriter.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
        isTerm = false
    }
    return func(c *Context) {
        start := time.Now()
        // Process request
        c.Next()
        ResetLoggerWriter()

        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()
        var statusColor, methodColor string
        if isTerm {
            statusColor = colorForStatus(statusCode)
            methodColor = colorForMethod(method)
        }
		hd := c.RequestHeaderString()
        comment := c.Errors.ByType(ErrorTypePrivate).String()

		end := time.Now()
		latency := end.Sub(start)
        fmt.Fprintf(defaultLoggerWriter, "[Bingo] %v |%s %3d %s| %13v | %s |%s %s %s| %s %s | %s \n%s",
            end.Format("2006/01/02 - 15:04:05"),
            statusColor, statusCode, reset,
            latency,
            clientIP,
            methodColor,  method, reset, 
            c.Request.URL, hd,
            strings.Join(c.logs," "),
            comment,
        )

    }
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

/***********************************/
/* cron日志 */
/***********************************/
var defaultCronLoggerWriter io.Writer
var CronOnlyOne sync.Once

func SetCronLoggerWriter() {
	CronOnlyOne.Do(func() {
		if runMode == debugCode {
			defaultCronLoggerWriter = os.Stdout
		}
		if runMode == releaseCode {
			defaultCronLoggerWriter = NewCronFileWriter()
		}
	})
}

func NewCronFileWriter() io.Writer {
	fileName := config.String("task","accessLog")
	if err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm); err != nil {
		fmt.Errorf("Can't create task accessLog folder on %v", err)
	}
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Errorf("Can't create task accessLog file: %v", err)
	}
	return file
}

type CronLogger struct {
	Begin time.Time
	Messages    []string
}

func NewCronLogger() *CronLogger{
	return &CronLogger{
		Begin: time.Now(),
		Messages: []string{},
	}
}
func (log *CronLogger) Record(args ...string) {
	if config.Bool("task","enableLog") == false || config.Bool("default","enableLog") == false {
		return
	}

	var message string
	if len(args) < 1 {
		panic("Func Task Logs params is shorter")
	}
	pose := utils.Int.String(len(log.Messages))
	if len(args) == 1 {
		message = utils.String.Join("[Node##",pose,"] ",args[0])
	}
	if len(args) > 1 {
		message = utils.String.Join("[",args[0],"##",pose,"] ",args[1])
	}
	log.Messages = append(log.Messages , message)
}
func (log *CronLogger) Save(script string, err ...bool) {
	if config.Bool("task","enableLog") == false || config.Bool("default","enableLog") == false {
		return
	}
	isTerm := true
	if w, ok := defaultCronLoggerWriter.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		isTerm = false
	}
	log.ResetCronLoggerWriter()

	var resultColor string
	result := "ERROR"
	if isTerm {
		resultColor = red
		if len(err) == 0 || err[0] == true {
			resultColor = green
		}
	}
	if len(err) == 0 || err[0] == true {
		result = "SUCCESS"
	}

	end := time.Now()
	latency := end.Sub(log.Begin)

	fmt.Fprintf(defaultCronLoggerWriter, "[Bingo] %v |%s %s %s| %13v | %s  | %s \n",
		end.Format("2006/01/02 - 15:04:05"),
		resultColor, result, reset,
		latency,
		script,
		strings.Join(log.Messages," "),
	)
}
func (log *CronLogger) ResetCronLoggerWriter(){
	if w, ok := defaultCronLoggerWriter.(*os.File); ok {
		file ,err := w.Stat()
		if err != nil {
			fmt.Errorf("Can't Stat accessLog file: %v", err)
		}
		modTime := file.ModTime()
		nowTime := time.Now()
		if modTime.Year() != nowTime.Year() || modTime.YearDay() != nowTime.YearDay() || modTime.Hour() != nowTime.Hour() {
			format := "2006010215"
			fileName := config.String("task","accessLog")
			err := os.Rename(fileName, fileName+"."+modTime.Format(format))
			if err != nil {
				fmt.Errorf("Can't Rename [%s] file: %v", fileName, err)
			}
			defaultCronLoggerWriter = NewCronFileWriter()
		}
	}
}

func (log *CronLogger) Logs(args ...string) {
	log.Record(args...)
}



