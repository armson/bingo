package weixin

//发送模板消息
//例如：
//result, err := weixin.New(m.Ctx).TemplateMessageSend(body, items);
//body:=map[string]interface{}{
// 	"touser":"OPENID",
//  "template_id":"ngqIpbwh8bUfcSsECmogfXcV14J0tQlEpBO27izEYtY",
//	"url":"http://weixin.qq.com/download",
//	"miniprogram":map[string]string{"appid":"xiaochengxuappid12345", "pagepath":"index?foo=bar"},
// }
//items := map[string][]string{
// 	"first": []string{"恭喜你购买成功！", "#173177"},
//	"keyword1": []string{"巧克力", "#173177"},
//	"keyword2": []string{"39.8元", "#173177"},
//	"keyword3": []string{"2014年9月22日", "#173177"},
//	"remark": []string{"欢迎再次购买！！", "#173177"},
// }
func (bin *Interface) TemplateMessageSend(body map[string]interface{}, items map[string][]string) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/cgi-bin/message/template/send", accessToken)
	bin.Method = "POST"
	bin.Queries["access_token"] =  accessToken

	data := map[string]map[string]string{}
	for key, item := range items {
		data[key] = map[string]string{"value":item[0], "color":item[1]}
	}
	body["data"] = data
	bin.Raw(body)
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}
	errCode, _ := js.Get("errcode").Int()
	if errCode != 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}
	return true, nil
}