package weixin

import (
	"github.com/armson/bingo/utils"
)

//小程序用户登陆
func (bin *Interface) MiniUserLogin(jsCode string) (map[string]string , error) {
	bin.Url = utils.String.Join(miniHost(), "/sns/jscode2session")
	bin.Method = "GET"
	bin.Queries["appid"] = miniAppId()
	bin.Queries["secret"] = miniAppSecret()
	bin.Queries["js_code"] = jsCode
	bin.Queries["grant_type"] = "authorization_code"
	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}
	errCode, err := js.Get("errcode").Int()
	if err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return nil, &Error{errCode, errMsg}
	}
	mp, err := js.Map()
	if err != nil {
		return nil, &ParseError{"MiniUserLogin","js"}
	}
	return map[string]string{
		"sessionKey"	:mp["session_key"].(string),
		"openId"		:mp["openid"].(string),
		"unionId"		:mp["unionid"].(string),
	}, nil
}