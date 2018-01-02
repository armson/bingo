package weixin

import(
	"github.com/armson/bingo/redis"
	"github.com/armson/bingo/utils"
	"encoding/json"
	"github.com/bitly/go-simplejson"
)
func (bin *Interface) Token() (string , error) {
	return bin.grantHandle("accessToken")
}
func (bin *Interface) Ticket() (map[string]string , error) {
	s, err := bin.grantHandle("ticket")
	if err != nil {
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(s))
	if err != nil {
		return nil, err
	}
	m, err := js.Map()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"jsApiTicket"	: m["jsApiTicket"].(string),
		"apiTicket"		: m["apiTicket"].(string),
	}, nil
}
func (bin *Interface) grantHandle(name string) (string , error) {
	handleList := map[string]map[string]interface{}{
		"accessToken" : map[string]interface{}{
			"key" 		: "weixin_access_token",
			"expires" 	: 3600,
			"cgi" 		: func() (string, error) {return bin.cgiAccessToken()},
		},
		"ticket" 		: map[string]interface{}{
			"key" 		: "weixin_Ticket",
			"expires" 	: 1800,
			"cgi" 		: func() (string, error) {return bin.cgiTicket()},
		},
	}
	handle := handleList[name]
	if redis.Valid() {
		if rs, _ := redis.New(bin.Tracer).Get(handle["key"].(string)); rs != "" {
			return rs, nil
		}
	}
	cgi, _ := handle["cgi"].(func()(string, error))
	rs, err := cgi()
	if err != nil {
		return "", err
	}
	if redis.Valid() {
		redis.New(bin.Tracer).SetEx(handle["key"].(string), rs, handle["expires"].(int))
	}
	return rs, nil
}
func (bin *Interface) cgiAccessToken() (string , error) {
	bin.Url = utils.String.Join(host(), "/cgi-bin/token")
	bin.Method = "GET"
	bin.Queries["grant_type"] =  "client_credential"
	bin.Queries["appid"] =  appId()
	bin.Queries["secret"] =  appSecret()
	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}
	errCode, err := js.Get("errcode").Int()
	if err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}
	accessToken , _ := js.Get("access_token").String()
	return accessToken, nil
}
func (bin *Interface) cgiTicket() (string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}
	jsApiTicket, err := bin.doCgiTicket(accessToken, "jsapi")
	if err != nil {
		return "", err
	}
	apiTicket, err := bin.doCgiTicket(accessToken, "wx_card")
	if err != nil {
		return "", err
	}
	s, err := json.Marshal(map[string]string{"jsApiTicket": jsApiTicket, "apiTicket": apiTicket})
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (bin *Interface) doCgiTicket(accessToken, t string) (string , error) {
	bin.Url = bin.cgiUrl("/cgi-bin/ticket/getticket", accessToken)
	bin.Method = "GET"
	bin.Queries["type"] =  t
	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}
	errCode, _ := js.Get("errcode").Int()
	if errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}
	ticket, _ := js.Get("ticket").String()
	return ticket, nil
}


