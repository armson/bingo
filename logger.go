package bingo

import (
    "fmt"
    "io"
    "os"
    "time"
    "github.com/armson/bingo/config"
    "github.com/armson/bingo/utils"
    "sync"
	"strings"
	"net/http/httputil"
	"bytes"
)

/***********************************/
/* 上下文access日志 */
/***********************************/
func Logger() HandlerFunc {
	log := new(accessLog)
	log.SetWriter()
	isTerm := IsTerm(log)
	return func(c *Context) {
		start := time.Now()
		// Process request
		c.Next()
		ResetFileWriter(log)

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
		fmt.Fprintf(log.w, "[Bingo] %v |%s %3d %s| %13v | %s |%s %s %s| %s %s | %s \n%s",
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

type accessLog struct {
	w io.Writer
}
var defaultAccessWrite io.Writer
var accessOnce sync.Once

func (log *accessLog) SetWriter()  {
	accessOnce.Do(func() {
		if runMode == debugCode {
			defaultAccessWrite = os.Stdout
		}
		if runMode == releaseCode {
			defaultAccessWrite = NewFileWriter(log)
		}
	})
	log.w = defaultAccessWrite
}

func (*accessLog) FileName() string {
	return config.String("accessLog")
}

func (*accessLog) String() string {
	return "accessLog"
}

func (log *accessLog) Write() io.Writer {
	return log.w
}

func (*accessLog) Rename(f os.FileInfo) bool {
	modTime := f.ModTime()
	nowTime := time.Now()
	if modTime.Year() != nowTime.Year() || modTime.YearDay() != nowTime.YearDay() || modTime.Hour() != nowTime.Hour() {
		format := "2006010215"
		fileName := config.String("accessLog")
		err := os.Rename(fileName, fileName+"."+modTime.Format(format))
		if err != nil {
			fmt.Errorf("Can't Rename [%s] file: %v", fileName, err)
		}
		return true
	}
	return false
}

/***********************************/
/* 上下文panic日志 */
/***********************************/
type errorLog struct {
	w io.Writer
}
var defaultErrorWrite io.Writer
var errorOnce sync.Once

func (log *errorLog) SetWriter()  {
	errorOnce.Do(func() {
		if runMode == debugCode {
			defaultErrorWrite = os.Stdout
		}
		if runMode == releaseCode {
			defaultErrorWrite = NewFileWriter(log)
		}
	})
	log.w = defaultErrorWrite
}

func (*errorLog) FileName() string {
	return config.String("errorLog")
}

func (*errorLog) String() string {
	return "errorLog"
}

func (log *errorLog) Write() io.Writer {
	return log.w
}

func (*errorLog) Rename(f os.FileInfo) bool {
	modTime := f.ModTime()
	nowTime := time.Now()
	if modTime.Year() != nowTime.Year() || modTime.YearDay() != nowTime.YearDay() {
		format := "20060102"
		fileName := config.String("errorLog")
		err := os.Rename(fileName, fileName+"."+modTime.Format(format))
		if err != nil {
			fmt.Errorf("Can't Rename [%s] file: %v", fileName, err)
		}
		return true
	}
	return false
}

func Recovery() HandlerFunc {
	log := new(errorLog)
	log.SetWriter()
	isTerm := IsTerm(log)
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				httprequest = bytes.Trim(httprequest,"\r\n")
				var color string
				if isTerm {
					color = "\x1b[31m"
				}

				fmt.Fprintf(log.w, "%s[Bingo] %v \n[Recovery] panic recovered:\n%s\n[Error] %s\n%s%s",
					color,
					time.Now().Format("2006/01/02 - 15:04:05"),
					string(httprequest),
					err,
					stack,
					reset)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

/***********************************/
/* cronAccess日志 */
/***********************************/
type CronAccessLog struct {
	Begin time.Time
	Messages    []string
	w io.Writer
}
var defaultCronAccessWrite io.Writer
var cronAccessOnce sync.Once

func NewCronAccessLog() *CronAccessLog{
	cron := &CronAccessLog{
		Begin: time.Now(),
		Messages: []string{},
	}
	cron.SetWriter()
	return cron
}

func (log *CronAccessLog) SetWriter()  {
	cronAccessOnce.Do(func() {
		if runMode == debugCode {
			defaultCronAccessWrite = os.Stdout
		}
		if runMode == releaseCode {
			defaultCronAccessWrite = NewFileWriter(log)
		}
	})
	log.w = defaultCronAccessWrite
}

func (log *CronAccessLog) Record(args ...string) {
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

func (log *CronAccessLog) Save(script string, err ...bool) {
	if config.Bool("task","enableLog") == false || config.Bool("default","enableLog") == false {
		return
	}
	isTerm := IsTerm(log)
	ResetFileWriter(log)

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
	fmt.Fprintf(log.w, "[Bingo] %v |%s %s %s| %13v | %s  | %s \n",
		end.Format("2006/01/02 - 15:04:05"),
		resultColor, result, reset,
		latency,
		script,
		strings.Join(log.Messages," "),
	)
}

func (log *CronAccessLog) Logs(args ...string) {
	log.Record(args...)
}

func (*CronAccessLog) FileName() string {
	return config.String("task","accessLog")
}

func (*CronAccessLog) String() string {
	return "taskAccessLog"
}

func (log *CronAccessLog) Write() io.Writer {
	return log.w
}

func (*CronAccessLog) Rename(f os.FileInfo) bool {
	modTime := f.ModTime()
	nowTime := time.Now()
	if modTime.Year() != nowTime.Year() || modTime.YearDay() != nowTime.YearDay() || modTime.Hour() != nowTime.Hour() {
		format := "2006010215"
		fileName := config.String("task","accessLog")
		err := os.Rename(fileName, fileName+"."+modTime.Format(format))
		if err != nil {
			fmt.Errorf("Can't Rename [%s] file: %v", fileName, err)
		}
		return true
	}
	return false
}




