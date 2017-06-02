package network

import(
    "strings"
    "strconv"
    "github.com/armson/bingo/utils"
)

var requestMethod []string = []string{"GET","POST","PUT","DELETE","HEAD","OPTIONS"}

func (this *work) Inputs() (map[string][]string) {
    if len(this.Params) > 0 {
        return this.Params
    }
    m := map[string][]string{}
    for _, v := range requestMethod {
        if v == this.Method {
            this.request.ParseForm()
            values := this.request.Form
            for k, val := range values {
                m[k] = val
            }
            this.Params = m
            return m
        }
    }
    return m
}

func (this *work) GetString(args ...string) string {
    m := this.Inputs()
    if v, ok := m[args[0]]; ok {
        return strings.Join(v, ",")
    }
    if len(args) > 1 {
        return args[1]
    }
    return ""
}

// string到int ：int,err:=strconv.Atoi(string)  
// string到int64  ： int64, err := strconv.ParseInt(string, 10, 64)  
// int到string  ：string:=strconv.Itoa(int)  
// nt64到string  ： string:=strconv.FormatInt(int64,10)

func (this *work) GetInt(args ...interface{}) int64 {
    m := this.Inputs()
    k := args[0].(string)
    if v, ok := m[k]; ok {
        s := strings.Join(v, ",")
        i, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return 0
        } 
        return i
    }
    if len(args) > 1 {
        utils.Interface.Int(args[1]) 
    }
    return 0
}





