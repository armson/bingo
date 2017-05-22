package network

import(
    "net/http"
    "strconv"
    "time"
)

func (this *work) SetCookiePath(s string) {
    this.CookiePath = s
}
func (this *work) SetCookieDomain(s string) {
    this.CookieDomain = s
}

func (this *work) SetCookie(name, value string, sec int) {
    t := time.Now()
    str := strconv.Itoa(sec)
    d ,_:= time.ParseDuration(str+"s")
    expiration := t.Add(d)
    cookie := http.Cookie{
        Name: name, 
        Value: value, 
        Expires: expiration,
        Path:this.CookiePath,
        Domain:this.CookieDomain,
        MaxAge:sec,
        HttpOnly:true} 
    http.SetCookie(this.response, &cookie)
}

func (this *work) UnsetCookie(name string) {
    cookie := http.Cookie{
        Name: name, 
        Value: "", 
        Path:this.CookiePath,
        Domain:this.CookieDomain,
        MaxAge:-1,
        HttpOnly:true} 
    http.SetCookie(this.response, &cookie)    
}

func (this *work) Cookie(key string) string {
    cookie, _ := this.request.Cookie(key)
    if cookie == nil {
        return ""
    }
    return cookie.Value
}

func (this *work) Cookies() (m map[string]string) {
    cookies := this.request.Cookies()
    if len(cookies) == 0 {
        return
    }
    m = make(map[string]string)
    for _, v := range cookies {
        m[v.Name] = v.Value
    }
    return
}





