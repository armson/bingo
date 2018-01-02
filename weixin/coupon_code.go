package weixin

import (
	"errors"
	"time"
)

//批量导入券码code ,建议最大值不要超过10000
//返回导入结果
//建议使用协程异步导入
//例如：
//skuQuantity := bingo.Int(c.POST("skuQuantity")).GT(0).LE(100000)
//go func() {
//	wxCouponModel.New(c).CouponCodeBatchSafeImport(cardId, skuQuantity.Int())
//}()
//func (m *wxCouponModel) CouponCodeBatchSafeImport(cardId string, quantity int) bool {
//	count , err := weixin.New(m.Ctx).CouponCodeBatchImportCount(cardId)  //检查库存
//	if err != nil || count >= quantity {
//		return false //是否继续
//	}
//	remain := quantity - count
//
//	sku := []*[]string{}
//	cm := couponModel.New(m.Ctx)
//	for remain > 0 {
//		maximum := remain
//		each := []string{}
//		if remain > weixin.CouponCodeEachTimeImportMaximum {
//			maximum = weixin.CouponCodeEachTimeImportMaximum
//		}
//		for i := 0; i < maximum; i++ {
//			each = append(each, cm.CreateCouponCode())
//		}
//		sku = append(sku , &each)
//		remain = remain - maximum
//	}
//	weixin.New(m.Ctx).CouponCodeBatchImport(cardId, &sku)
//	if m.CouponCodeBatchSafeImport(cardId, quantity) == false {
//		return false
//	}
//	return true
//}

func (bin *Interface) CouponCodeBatchImport(cardId string, skuPointer *[]*[]string) (map[string]int, error) {
	result := map[string]int{"success":0, "duplicate":0, "fail":0}
	accessToken, err := bin.Token()
	if err != nil {
		return result, err
	}
	bin.Url = bin.cgiUrl("/card/code/deposit", accessToken)

	//创建通道
	sku := *skuPointer
	ch := make(chan []interface{},len(sku))
	for _ , code := range sku {
		go bin.couponCodeImport(cardId, code, ch)
		time.Sleep(time.Millisecond * 50)
	}

	for i := 0 ; i < len(sku); i ++ {
		message := <-ch
		result["success"] 	= result["success"] + message[0].(int)
		result["duplicate"] = result["duplicate"] + message[1].(int)
		result["fail"] 		= result["fail"] + message[2].(int)
	}
	return result, nil
}

func (bin *Interface) couponCodeImport(cardId string, codesPointer *[]string , ch chan []interface{}) {
	codes := *codesPointer
	if len(codes) < 1 || len(codes)  > CouponCodeEachTimeImportMaximum {
		ch <- []interface{}{0,0,0,errors.New("The import amount must be greater than 0 and less than 100 for each time")}
		return
	}
	//开启新的请求，避免微信服务器屏蔽
	wx := New(bin.Tracer)
	wx.Url = bin.Url
	wx.Method = "POST"
	wx.CloseLog(true)
	raw := map[string]interface{}{"card_id":cardId, "code":codes}
	wx.Raw(raw)
	js , err :=  wx.Send()
	if err != nil {
		ch <- []interface{}{0,0,0,requestWithoutResponse}
		return
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		ch <- []interface{}{0, 0, 0, &ParseError{"CouponCodeImport","errcode"}}
		return
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		ch <- []interface{}{0,0,0,&Error{errCode, errMsg}}
		return
	}

	duplicate, _ := js.Get("duplicate_code").StringArray()
	fail, _ := js.Get("fail_code").StringArray()
	success, _ := js.Get("succ_code").StringArray()

	ch <- []interface{}{len(success), len(duplicate), len(fail) ,nil}
	return
}

//查询导入code数量接口
func (bin *Interface) CouponCodeBatchImportCount(cardId string) (int , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return 0, err
	}
	bin.Url = bin.cgiUrl("/card/code/getdepositcount", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{"card_id":cardId}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return 0, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return 0 , &ParseError{"CouponCodeBatchImportCount","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return 0 , &Error{errCode, errMsg}
	}
	count, err := js.Get("count").Int()
	if err != nil {
		return 0 , &ParseError{"CouponCodeBatchImportCount","count"}
	}
	return count , nil
}

//查询导入code数量接口
//quantity > 0 增加库存，quantity < 0 减库存
func (bin *Interface) CouponCodeStockReplace(cardId string , quantity int) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/card/modifystock", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{"card_id":cardId}
	if quantity > 0 {
		raw["increase_stock_value"] = quantity
	}
	if quantity < 0 {
		raw["reduce_stock_value"] = quantity * -1
	}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false , &ParseError{"CouponCodeStockReplace","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false , &Error{errCode, errMsg}
	}
	return true , nil
}

//查询Code详情接口
func (bin *Interface) CouponCodeDetail(cardId, code string) (map[string]interface{} , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}
	bin.Url = bin.cgiUrl("/card/code/get", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{"card_id":cardId,"code":code,"check_consume":false}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return nil , &ParseError{"CouponCodeDetail","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return nil , &Error{errCode, errMsg}
	}

	openid,_ := js.Get("openid").String()
	canConsume,_ := js.Get("can_consume").Bool()
	userCardStatus,_ := js.Get("user_card_status").String()
	card_id,_ := js.Get("card").Get("card_id").String()
	beginTime,_ := js.Get("card").Get("begin_time").Int()
	endTime,_ := js.Get("card").Get("end_time").Int()
	return map[string]interface{}{
		"openid" : openid,
		"canConsume" : canConsume,
		"userCardStatus" : userCardStatus,
		"cardId" : card_id,
		"beginTime" : beginTime,
		"endTime" : endTime,
		"userCardStatusName" : GetCardUserCardStatusName(userCardStatus),
	}, nil
}

//核销Code接口
func (bin *Interface) CouponCodeConsume(cardId, code string) (map[string]interface{} , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}
	bin.Url = bin.cgiUrl("/card/code/consume", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{"card_id":cardId,"code":code}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return nil , &ParseError{"CouponCodeConsume","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return nil , &Error{errCode, errMsg}
	}

	openid,_ := js.Get("openid").String()
	card_id,_ := js.Get("card").Get("card_id").String()
	return map[string]interface{}{
		"openid" : openid,
		"cardId" : card_id,
	}, nil
}

//code码解密
func (bin *Interface) CouponCodeDecrypt(encryptCode string) (string , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}
	bin.Url = bin.cgiUrl("/card/code/decrypt", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{"encrypt_code":encryptCode}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return "" , &ParseError{"CouponCodeConsume","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return "" , &Error{errCode, errMsg}
	}

	code, _ := js.Get("code").String()
	return code, nil
}



//设置卡券失效接口
func (bin *Interface) CouponCodeUnavailable(cardId, code ,reason string) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/card/code/unavailable", accessToken)
	bin.Method = "POST"
	raw := map[string]interface{}{"card_id":cardId,"code":code,"reason":reason}
	bin.Raw(raw)

	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false , &ParseError{"CouponCodeUnavailable","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false , &Error{errCode, errMsg}
	}
	return true, nil
}



