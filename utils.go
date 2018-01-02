package bingo

import(
    "path"
    "reflect"
    "runtime"
    "github.com/armson/bingo/config"
	"regexp"
	"github.com/armson/bingo/utils"
)
type compare struct {
	logicOperator string
	mathOperator string
	arg interface{}
}
var (
	logicOperatorAndChar 	= "&&"
	logicOperatorOrChar 	= "||"
)

/*****************************/
/* 字符串                     */
/*****************************/
var (
	validLetter, _ 			= regexp.Compile(`^[a-zA-Z]+$`)
	validNumber, _ 			= regexp.Compile(`^[0-9]+$`)
	validLetterAndNumber, _ = regexp.Compile(`^[a-zA-Z0-9]+$`)
	validMobile, _ 			= regexp.Compile(`^1[34578]\d{9}$`)
	validChinese,_			= regexp.Compile(`^[\p{Han}]+$`)
	validChar,_				= regexp.Compile(`^[\p{Han}a-zA-Z0-9\-_.]+$`)
	validUrl,_				= regexp.Compile(`^(http|https)://[a-zA-Z0-9:/.|%&;?=_\-]+$`)
)
type bingoString struct {
	value ,defaultValue string
	regular *regexp.Regexp
	defaultLogicOperator string
	compares []*compare
	length int
}
// eq:等于 ne:不等于 lt:小于 gt:大于 le:小于或等于 ge:大于或等于
func LetterAndNumber(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validLetterAndNumber,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}
func Letter(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validLetter,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}
func Number(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validNumber,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}
func Mobile(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validMobile,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}
func Chinese(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validChinese,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}
func Char(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validChar,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}
func Url(str string) *bingoString {
	return &bingoString{
		value:		  			str,
		regular:	  			validUrl,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		length:					len(str),
	}
}

func (bin *bingoString)And() *bingoString {
	bin.defaultLogicOperator = logicOperatorAndChar
	return bin
}
func (bin *bingoString)Or() *bingoString {
	bin.defaultLogicOperator = logicOperatorOrChar
	return bin
}

func (bin *bingoString)EQ(limit int) *bingoString { return bin.handle(limit,"EQ") }
func (bin *bingoString)NE(limit int) *bingoString { return bin.handle(limit,"NE") }
func (bin *bingoString)LT(limit int) *bingoString { return bin.handle(limit,"LT") }
func (bin *bingoString)GT(limit int) *bingoString { return bin.handle(limit,"GT") }
func (bin *bingoString)LE(limit int) *bingoString { return bin.handle(limit,"LE") }
func (bin *bingoString)GE(limit int) *bingoString { return bin.handle(limit,"GE") }

func (bin *bingoString)Default(defaultValue string) *bingoString {
	bin.defaultValue = defaultValue
	return bin
}

func (bin *bingoString)Bool() bool {
	if !bin.regular.MatchString(bin.value) {
		return false
	}
	if len(bin.compares) > 0 {
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

func (bin *bingoString)handle(limit int, operator string) *bingoString {
	bin.compares = append(bin.compares,
		&compare{
			logicOperator:bin.defaultLogicOperator,
			mathOperator:operator,
			arg:limit,
		},
	)
	return bin
}
func (bin *bingoString)compare() bool {
	effect := true
	for _ , com := range bin.compares {
		whether := bin.whether(com)
		switch com.logicOperator {
		case logicOperatorAndChar:
			effect = effect && whether
		case logicOperatorOrChar:
			effect = effect || whether
		}
	}
	return effect
}

func (bin *bingoString)whether(com *compare) bool {
	arg , err := utils.Interface.Int(com.arg)
	if err != nil {return  false}

	switch com.mathOperator {
	case "EQ":
		return bin.length == arg
	case "NE":
		return bin.length !=  arg
	case "LT":
		return bin.length  < arg
	case "GT":
		return bin.length  > arg
	case "LE":
		return bin.length  <= arg
	case "GE":
		return bin.length  >= arg
	}
	return bin.length  == arg
}

/*****************************/
/* bingo整数或者整数字符串     */
/*****************************/
var(
	minUnixTime int64 = 0 //0,0,0,1,1,1970
	maxUnixTime int64 = 4102416000 //0,0,0,1,1,2100
	minUnixMSTime int64 = 0 //0,0,0,1,1,1970
	maxUnixMSTime int64 = 4102416000000 //0,0,0,1,1,2100
	minInt64 int64 = -9223372036854775808
	maxInt64 int64 = 9223372036854775807
)

type bingoInt struct {
	value, defaultValue interface{}
	min, max int64
	defaultLogicOperator string
	compares []*compare
}

func Int(t interface{}) *bingoInt{
	return &bingoInt{
		value:		  			t,
		min:		  			minInt64,
		max:		  			maxInt64,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		defaultValue: 			nil,
	}
}

func UnixS(t interface{}) *bingoInt{
	return &bingoInt{
		value:		  			t,
		min:		  			minUnixTime,
		max:		  			maxUnixTime,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		defaultValue: 			nil,
	}
}
func UnixMS(t interface{}) *bingoInt{
	return &bingoInt{
		value:		  			t,
		min:		  			minUnixMSTime,
		max:		  			maxUnixMSTime,
		defaultLogicOperator: 	logicOperatorAndChar,
		compares:				[]*compare{},
		defaultValue: 			nil,
	}
}
func (bin *bingoInt)And() *bingoInt {
	bin.defaultLogicOperator = logicOperatorAndChar
	return bin
}
func (bin *bingoInt)Or() *bingoInt {
	bin.defaultLogicOperator = logicOperatorOrChar
	return bin
}
func (bin *bingoInt)EQ(param interface{}) *bingoInt { return bin.handle(param,"EQ") }
func (bin *bingoInt)NE(param interface{}) *bingoInt { return bin.handle(param,"NE") }
func (bin *bingoInt)LT(param interface{}) *bingoInt { return bin.handle(param,"LT") }
func (bin *bingoInt)GT(param interface{}) *bingoInt { return bin.handle(param,"GT") }
func (bin *bingoInt)LE(param interface{}) *bingoInt { return bin.handle(param,"LE") }
func (bin *bingoInt)GE(param interface{}) *bingoInt { return bin.handle(param,"GE") }

func (bin *bingoInt)Default(defaultValue interface{}) *bingoInt {
	bin.defaultValue = defaultValue
	return bin
}
func (bin *bingoInt)Bool() bool {
	value,err := utils.Interface.Int64(bin.value)
	if err != nil { return false }

	if value < bin.min ||  value > bin.max {
		return false
	}
	if len(bin.compares) > 0 {
		return bin.compare()
	}
	return true
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
	return ""
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
	return 0
}
func (bin *bingoInt)Float() float64 {
	return float64(bin.Int())
}


func (bin *bingoInt)handle(param interface{}, operator string) *bingoInt {
	bin.compares = append(bin.compares,
		&compare{
			logicOperator:bin.defaultLogicOperator,
			mathOperator:operator,
			arg:param,
		},
	)
	return bin
}

func (bin *bingoInt)compare() bool {
	effect := true
	for _ , com := range bin.compares {
		whether := bin.whether(com)
		switch com.logicOperator {
		case logicOperatorAndChar:
			effect = effect && whether
		case logicOperatorOrChar:
			effect = effect || whether
		}
	}
	return effect
}
func (bin *bingoInt)whether(com *compare) bool {
	arg , err := utils.Interface.Int64(com.arg)
	if err != nil {return  false}

	value, err := utils.Interface.Int64(bin.value)
	if err != nil { return false }

	switch com.mathOperator {
	case "EQ":
		return value == arg
	case "NE":
		return value !=  arg
	case "LT":
		return value  < arg
	case "GT":
		return value  > arg
	case "LE":
		return value  <= arg
	case "GE":
		return value  >= arg
	}
	return value  == arg
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



