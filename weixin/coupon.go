package weixin

import (
	"github.com/armson/bingo/utils"
)

//批量查询卡券列表
//例如：
//offset := bingo.Int(c.GET("offset")).Default(0).Int()
//limit := bingo.Int(c.GET("limit")).Default(10).Int()
//status := c.GET("status") //多个状态中间用,分隔，
//data, total , err := weixin.New(c).CouponList(offset, limit , status)
func (bin *Interface) CouponList(offset , count int , status string) ([]string , int, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, 0, err
	}
	bin.Url = bin.cgiUrl("/card/batchget", accessToken)
	bin.Method = "POST"

	raw := map[string]interface{}{"offset":offset, "count":count}
	status_list := CardStatusToSlice(status)
	if len(status_list) > 0 {
		raw["status_list"] = status_list
	}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return nil, 0, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return nil, 0 , &ParseError{"CardCreate","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return nil, 0 , &Error{errCode, errMsg}
	}

	total ,_ := js.Get("total_num").Int()
	list ,_  := js.Get("card_id_list").StringArray()
	return list, total, nil
}

//查看卡券详情
func (bin *Interface) CouponDetail(cardId string) (map[string]interface{} , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}
	bin.Url = bin.cgiUrl("/card/get", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]interface{}{"card_id":cardId})
	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return nil, &ParseError{"CouponDetail","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return nil, &Error{errCode, errMsg}
	}

	card := js.Get("card")
	base , err := bin.getCardBaseInfo(card)
	if err != nil {
		return nil , err
	}
	return base, nil
}

//删除优惠券
//例如：
//cardId := c.POST("cardId")
//_ , err := weixin.New(c).CouponDelete(cardId)
func (bin *Interface) CouponDelete(cardId string) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/card/delete", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]interface{}{"card_id":cardId})
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false, &ParseError{"CouponDelete","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}
	return true, nil
}

//设置买单接口
//1.设置快速买单的卡券须支持至少一家有核销员门店，否则无法设置成功；
//2.若该卡券设置了center_url（居中使用跳转链接）,须先将该设置更新为空后再设置自快速买单方可生效。
//例如：
//params := map[string]string{
//"cardId" 			: c.POST("cardId"),
//"isOpen" 			: c.POST("isOpen"),
//}
//_ , err := weixin.New(c).CouponPayCell(params)
func (bin *Interface) CouponPayCell(params map[string]string) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/card/paycell/set", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{
		"card_id"			:params["cardId"],
		"is_open"			:utils.String.Bool(params["isOpen"]),
	}
	bin.Raw(raw)
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false, &ParseError{"CardPayCell","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}
	return true, nil
}

//设置自助核销接口
//1.设置自助核销的卡券须支持至少一家门店，否则无法设置成功；
//2.若该卡券设置了center_url（居中使用跳转链接）,须先将该设置更新为空后再设置自助核销功能方可生效。
//例如：
//params := map[string]string{
//"cardId" 			: c.POST("cardId"),
//"isOpen" 			: c.POST("isOpen"),
//"needVerifyCode" 	: c.POST("needVerifyCode"),
//"needRemarkAmount" 	: c.POST("needRemarkAmount"),
//}
//_ , err := weixin.New(c).CouponSelfConsumeCell(params)

func (bin *Interface) CouponSelfConsumeCell(params  map[string]string) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/card/selfconsumecell/set", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{
		"card_id"			:params["cardId"],
		"is_open"			:utils.String.Bool(params["isOpen"]),
		"need_verify_code"	:utils.String.Bool(params["needVerifyCode"]),
		"need_remark_amount":utils.String.Bool(params["needRemarkAmount"]),
	}
	bin.Raw(raw)
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false, &ParseError{"CouponSelfConsumeCell","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}
	return true, nil
}



//创建二维码接口
//code:卡券Code码,use_custom_code字段为true的卡券必须填写，非自定义code和导入code模式的卡券不必填写。
//cardId:卡券ID
//openid:指定领取者的openid，只有该用户能领取。bind_openid字段为true的卡券必须填写，非指定openid不必填写。
//isUniqueCode:指定下发二维码，生成的二维码随机分配一个code，领取后不可再次扫描。
// 填写true或false。默认false，注意填写该字段时，卡券须通过审核且库存不为0
//outerId:领取场景值，用于领取渠道的数据统计，默认值为0，字段类型为整型，长度限制为60位数字。
// 用户领取卡券后触发的事件推送中会带上此自定义场景值。
//outerStr:uter_id字段升级版本，字符串类型，用户首次领卡时，会通过领取事件推送给商户；
// 对于会员卡的二维码，用户每次扫码打开会员卡后点击任何url，会将该值拼入url中，方便开发者定位扫码来源
//expireSeconds:指定二维码的有效时间，范围是60 ~ 1800秒。不填默认为365天有效
func (bin *Interface) CouponQrcode(coupons []map[string]interface{}, expireSeconds int) (string , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}
	bin.Url = bin.cgiUrl("/card/qrcode/create", accessToken)
	bin.Method = "POST"

	cards := []map[string]interface{}{}
	for _ , coupon := range coupons {
		card := map[string]interface{}{}
		for _ , item := range []string{"code", "cardId", "openid", "isUniqueCode", "outerId", "outerStr"} {
			if value , ok := coupon[item]; ok {
				card[utils.String.UnderLine(item)] = value
			}
		}
		cards = append(cards , card)
	}

	raw := map[string]interface{}{}
	raw["action_name"] = "QR_CARD"
	if expireSeconds <= 1800 {
		raw["expire_seconds"] = expireSeconds
	}

	if len(coupons) == 1 {
		raw["action_name"] = "QR_CARD"
		raw["action_info"] = map[string]interface{}{"card":cards[0]}
	} else {
		raw["action_name"] = "QR_MULTIPLE_CARD"
		multipleCard := map[string]interface{}{"card_list":cards}
		raw["action_info"] = map[string]interface{}{"multiple_card":multipleCard}
	}
	bin.Raw(raw)
	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}
	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return "", &ParseError{"CouponQrcode","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}

	qrcodeUrl , err := js.Get("show_qrcode_url").String()
	if err != nil {
		return "", &ParseError{"CouponQrcode","show_qrcode_url"}
	}
	return qrcodeUrl, nil
}




