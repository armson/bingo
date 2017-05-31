package network

import(
    "net/http"
    "strings"
)

type work struct{
    response http.ResponseWriter
    request *http.Request
    CookiePath       string    // optional
    CookieDomain     string    // optional
    Method  string  //GET，POST，PUT，DELETE等
    Params  map[string][]string
    RemoteAddr string
    RemotePort string
}
func New(w http.ResponseWriter, r *http.Request) *work { 
    my := new(work)
    my.response = w
    my.request = r
    my.CookiePath = "/"
    my.CookieDomain = "/"
    my.Method = r.Method
    //r.Header.get("Remote_addr"),通过nginx代理获取
    addr := strings.Split(r.RemoteAddr,":")
    my.RemoteAddr = addr[0]
    my.RemotePort = addr[1]
    return my
}
func (this *work) Redirect(url string) {
    http.Redirect(this.response, this.request, url, 302)
}





