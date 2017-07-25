package bingo

import(
    "path"
    "reflect"
    "runtime"
    "github.com/armson/bingo/config"
	"regexp"
	"github.com/armson/bingo/utils"
)

/*****************************/
/* 字符串                     */
/*****************************/
var (
	validLetter, _ = regexp.Compile(`^[a-zA-Z]+$`)
	validNumber, _ = regexp.Compile(`^[0-9]+$`)
	validLetterAndNumber, _ = regexp.Compile(`^[a-zA-Z0-9]+$`)
	validMobile, _ = regexp.Compile(`^1[34578]\d{9}$`)
)
type bingoString struct {
	value ,defaultValue,operate string
	regular *regexp.Regexp
	compareLen int
}
// eq:等于 ne:不等于 lt:小于 gt:大于 le:小于或等于 ge:大于或等于
func LetterAndNumber(str string) *bingoString {
	return &bingoString{
		value:		  str,
		regular:	  validLetterAndNumber,
	}
}
func Letter(str string) *bingoString {
	return &bingoString{
		value:		  str,
		regular:	  validLetter,
	}
}
func Number(str string) *bingoString {
	return &bingoString{
		value:		  str,
		regular:	  validNumber,
	}
}
func Mobile(str string) *bingoString {
	return &bingoString{
		value:		  str,
		regular:	  validMobile,
	}
}

func (bin *bingoString)EQ(limit int) *bingoString {
	bin.compareLen = limit
	bin.operate = "EQ"
	return bin
}
func (bin *bingoString)NE(limit int) *bingoString {
	bin.compareLen = limit
	bin.operate = "NE"
	return bin
}
func (bin *bingoString)LT(limit int) *bingoString {
	bin.compareLen = limit
	bin.operate = "LT"
	return bin
}
func (bin *bingoString)GT(limit int) *bingoString {
	bin.compareLen = limit
	bin.operate = "GT"
	return bin
}
func (bin *bingoString)LE(limit int) *bingoString {
	bin.compareLen = limit
	bin.operate = "LE"
	return bin
}
func (bin *bingoString)GE(limit int) *bingoString {
	bin.compareLen = limit
	bin.operate = "GE"
	return bin
}
func (bin *bingoString)Default(defaultValue string) *bingoString {
	bin.defaultValue = defaultValue
	return bin
}

func (bin *bingoString)Bool() bool {
	if !bin.regular.MatchString(bin.value) {
		return false
	}
	if bin.compareLen > 0 {
		return bin.compare()
	}
	return true
}
func (bin *bingoString)String() string {
	if bin.Bool() {
		return bin.value
	}
	return bin.defaultValue
}
func (bin *bingoString)Int() int {
	return utils.String.Int(bin.String())
}

func (bin *bingoString)compare() bool {
	strLen := len(bin.value)
	switch bin.operate {
	case "EQ":
		return strLen == bin.compareLen
	case "NE":
		return strLen != bin.compareLen
	case "LT":
		return strLen < bin.compareLen
	case "GT":
		return strLen > bin.compareLen
	case "LE":
		return strLen <= bin.compareLen
	case "GE":
		return strLen >= bin.compareLen
	}
	return strLen == bin.compareLen
}

/*****************************/
/* bingo整数或者整数字符串     */
/*****************************/
var(
	minUnixTime int64 = 368640000 //0,0,0,9,7,1981
	maxUnixTime int64 = 4102416000 //0,0,0,1,1,2100
	minUnixMSTime int64 = 368640000000 //0,0,0,9,7,1981
	maxUnixMSTime int64 = 4102416000000 //0,0,0,1,1,2100
	minInt64 int64 = -9223372036854775808
	maxInt64 int64 = 9223372036854775807
)
type bingoInt struct {
	value, compareParam, defaultValue interface{}
	operate string
	min, max int64
}
func Int(t interface{}) *bingoInt{
	return &bingoInt{
		value:		  t,
		min:		  minInt64,
		max:		  maxInt64,
		compareParam: nil,
		defaultValue: nil,
	}
}
func UnixS(t interface{}) *bingoInt{
	return &bingoInt{
		value:		  t,
		min:		  minUnixTime,
		max:		  maxUnixTime,
		compareParam: nil,
		defaultValue: nil,
	}
}
func UnixMS(t interface{}) *bingoInt{
	return &bingoInt{
		value:		  t,
		min:		  minUnixMSTime,
		max:		  maxUnixMSTime,
		compareParam: nil,
		defaultValue: nil,
	}
}

func (bin *bingoInt)EQ(param interface{}) *bingoInt {
	bin.compareParam = param
	bin.operate = "EQ"
	return bin
}
func (bin *bingoInt)NE(param interface{}) *bingoInt {
	bin.compareParam = param
	bin.operate = "NE"
	return bin
}
func (bin *bingoInt)LT(param interface{}) *bingoInt {
	bin.compareParam = param
	bin.operate = "LT"
	return bin
}
func (bin *bingoInt)GT(param interface{}) *bingoInt {
	bin.compareParam = param
	bin.operate = "GT"
	return bin
}
func (bin *bingoInt)LE(param interface{}) *bingoInt {
	bin.compareParam = param
	bin.operate = "LE"
	return bin
}
func (bin *bingoInt)GE(param interface{}) *bingoInt {
	bin.compareParam = param
	bin.operate = "GE"
	return bin
}
func (bin *bingoInt)Default(defaultValue interface{}) *bingoInt {
	bin.defaultValue = defaultValue
	return bin
}
func (bin *bingoInt)compare() bool {
	value, err := utils.Interface.Int64(bin.value)
	if err != nil { return false }

	compareParam,err := utils.Interface.Int64(bin.compareParam)
	if err != nil { return false }

	switch bin.operate {
	case "EQ":
		return value == compareParam
	case "NE":
		return value != compareParam
	case "LT":
		return value < compareParam
	case "GT":
		return value > compareParam
	case "LE":
		return value <= compareParam
	case "GE":
		return value >= compareParam
	}
	return value == compareParam
}
func (bin *bingoInt)Bool() bool {
	value,err := utils.Interface.Int64(bin.value)
	if err != nil { return false }

	if value < bin.min ||  value > bin.max {
		return false
	}
	if bin.compareParam != nil {
		return bin.compare()
	}
	return true
}
func (bin *bingoInt)Int() int {
	if bin.Bool() {
		value, _ := utils.Interface.Int(bin.value)
		return value
	}
	if bin.defaultValue != nil {
		value, _ := utils.Interface.Int(bin.defaultValue)
		return value
	}
	if bin.compareParam != nil {
		value, _ := utils.Interface.Int(bin.compareParam)
		return value
	}
	return 0
}
func (bin *bingoInt)String() string {
	if bin.Bool() {
		value, _ := utils.Interface.String(bin.value)
		return value
	}
	if bin.defaultValue != nil {
		value, _ := utils.Interface.String(bin.defaultValue)
		return value
	}
	if bin.compareParam != nil {
		value, _ := utils.Interface.String(bin.compareParam)
		return value
	}
	return ""
}



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



