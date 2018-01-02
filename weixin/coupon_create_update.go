package weixin

import (
	"errors"
	"strings"
)

//创建卡券接口
//cardType:	*必填,卡券类型，参考define.go中的card_types定义
//logoUrl		:*必填,卡券的商户logo，建议像素为300*300。
//brandName	:*必填,商户名字,字数上限为12个汉字。[无法更新]
//color			:*必填,卡券颜色
//dateInfoType	：*必填,可选值DATE_TYPE_FIX_TIME_RANGE：表示固定日期区间 DATE_TYPE_FIX_TERM：表示固定时长（自领取后按天算）
// [无法更新]，如需求修改dateInfoBeginTimestamp、dateInfoEndTimestamp等时，dateInfoType：必须和已有类型保持一致
//dateInfoBeginTimestamp		: *必填, type为DATE_TYPE_FIX_TIME_RANGE时专用，表示起用时间，[更新时：必须小于等于已有开始日期]
//dateInfoEndTimestamp		: *必填,type为DATE_TYPE_FIX_TIME_RANGE时专用，表示结束时间，[更新时：必须大于等于已有开始日期]
//dateInfoFixedBeginTerm	: *必填,type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天开始生效，领取后当天生效填写0
//dateInfoFixedTerm	: *必填,type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天内有效，不支持填写0
//skuQuantity		: *必填,卡券库存的数量，上限为100000000。在创建时，use_custom_code_mode=2,sku_quantity必须等于0, [无法更新]
//title	: *必填,卡券名，字数上限为9个汉字。(建议涵盖卡券属性、服务及金额)。
//dealDetail :*必填,团购券专用，优惠券详情(指优惠说明),限制1-3072个字符之间，[无法更新]
//defaultDetail：*必填,通用券专用，优惠券详情(指优惠说明),限制1-3072个字符之间，[无法更新]
//leastCost:*必填,代金券专用,表示起用金额（单位为分）,如果无起用门槛则填0。此处的优先级低于uc_least_cost, [无法更新]
//reduceCost:*必填,代金券专用,代金券专用，表示减免金额。（单位为分）, [无法更新]
//discount:*必填,折扣券专用,表示打折额度（百分比）。填30就是七折。[无法更新]
//gift:*必填,兑换券专用,填写兑换内容的名称。[无法更新]

//**************以下为选填项*******************//
//codeType：码型，参考define.go中的card_types定义
//notice：卡券使用提醒(使用设置-操作提示)，字数上限为16个汉字。
//servicePhone：客服电话（限制24字节）
//description：卡券使用说明，字数上限为1024个汉字。
//useLimit：每人可核销的数量限制,不填写默认为50。[无法更新]
//getLimit：每人可领券的数量限制,不填写默认为50。
//useCustomCodeMode: 0:非自定义Code码 1:自定义Code码 2:导入code模式,参考define.go中的card_use_custom_code_mode定义。[无法更新]
//bindOpenid：是否指定用户领取，填写true或false。默认为false。通常指定特殊用户群体投放卡券或防止刷券时选择指定用户领取。
//canShare：卡券领取页面是否可分享。
//canGiveFriend：卡券是否可转赠。
//locationIdList:设置该卡券的适用门店,[]int{}，默认无指定门店，传空值表示全部门店适用

//centerTitle: 使用场景入口，仅卡券被用户领取且处于有效状态时显示（未到有效期、转赠中、核销后不显示）。
//centerSubTitle：
//centerUrl：
//示例：立即使用
//例如： "centerTitle": "顶部居中按钮", "centerSubTitle": "按钮下方的wording", "centerUrl": "www.qq.com",

//centerAppBrandUserName:卡券跳转的小程序的user_name，仅可跳转该公众号绑定的小程序,例如：gh_86a091e50ad4@app
//centerAppBrandPass:卡券跳转的小程序的path，例如：API/cardPage
//customAppBrandUserName:同上
//customAppBrandPass:同上
//promotionAppBrandUserName:同上
//promotionAppBrandPass:同上

//customUrlName: 服务场景入口, 仅卡券被用户领取且处于有效状态时显示（转赠中、核销后不显示）
//customUrl:
//customUrlSubTitle:
//示例：在线商城
//例如：  "customUrlName": "立即使用", "customUrl": "http://www.qq.com", "customUrlSubTitle": "6个汉字tips",

//promotionUrlName:营销场景入口，卡券处于正常状态、转赠中、核销后等异常状态均显示该入口。
//promotionUrlSubTitle:
//promotionUrl:
//示例：再次购买

//source: "大众点评" 。[无法更新]
//ucAcceptCategory:指定可用的商品类目，仅用于代金券类型，填入后将在券面拼写适用于xxx 。[无法更新]
//ucRejectCategory:指定不可用的商品类目，仅用于代金券类型，填入后将在券面拼写不适用于xxxx 。[无法更新]
//ucLeastCost:满减门槛字段，可用于兑换券和代金券，填入后将在全面拼写消费满xx元可用。[无法更新]
//ucObjectUseFor:购买xx可用类型门槛，仅用于兑换，填入后自动拼写购买xxx可用。[无法更新]
//ucCanUseWithOtherDiscount:不可以与其他类型共享门槛，填写false时系统将在使用须知里,拼写“不可与其他优惠共享”，[无法更新]
// 填写true时系统将在使用须知里,拼写“可与其他优惠共享”，默认为true。[无法更新]

//abstractAbstract:封面摘要简介
//abstractIconUrlList:封面URL

//textImageList:图文列表，显示在详情内页,Json结构[]map[string]string,更新时，使用空将清空图文列表
//textImageList:imageUrl
//textImageList:text


//businessService:商家服务类型，多个服务中间使用,分隔，可选值 BIZ_SERVICE_DELIVER 外卖服务； BIZ_SERVICE_FREE_PARK 停车位；
// BIZ_SERVICE_WITH_PET 可带宠物；BIZ_SERVICE_FREE_WIFI 免费wifi，注意：更新时，微信接口不支持清空该数据

//timeLimit :使用时段限制，Json结构[]map[string]interface{}，更新时，使用空或者[],那么将清空该时间限制
//timeLimit :type 此处只控制显示，不控制实际使用逻辑，不填默认不显示, 限制类型枚举值：支持填入
//MONDAY 周一 TUESDAY 周二 WEDNESDAY 周三 THURSDAY 周四 FRIDAY 周五 SATURDAY 周六 SUNDAY 周日
//timeLimit :beginHour
//timeLimit :beginMinute
//timeLimit :endHour
//timeLimit :endMinute
//例如timeLimit=[{"type":"MONDAY","beginMinute":10,"beginHour":10,"endHour":11,"endMinute":11},
// {"type":"MONDAY","beginMinute":8,"beginHour":8,"endHour":9,"endMinute":9},
// {"type":"SUNDAY"},{"type":"FRIDAY","beginMinute":9,"beginHour":9,"endHour":22,"endMinute":22}]


func (bin *Interface) CouponCreate(post *map[string]string) ( string , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}

	info := *post
	cardType, ok := info["cardType"]
	if !ok || !CheckCardType(cardType) {
		return "" , errors.New("Incorrect cardType")
	}

	//整理数据成json格式
	item := map[string]interface{}{}

	baseInfo ,err := bin.checkCouponBaseInfo(post, &map[string]interface{}{})
	if err != nil {
		return "", err
	}
	item["base_info"] = baseInfo

	//不同类型的优惠券的专有信息
	//团购券详情
	if cardType == CARD_TYPE_GROUPON {
		var deal_detail string
		if deal_detail ,err = bin.checkCouponDealDetail(post,false); err != nil {
			return "", err
		}
		item["deal_detail"] = deal_detail
	}
	//通用券详情
	if cardType == CARD_TYPE_GENERAL_COUPON {
		var default_detail string
		if default_detail ,err = bin.checkCouponDefaultDetail(post,false); err != nil {
			return "", err
		}
		item["default_detail"] = default_detail
	}

	//代金券
	if cardType == CARD_TYPE_CASH {
		var least_cost ,reduce_cost int
		if least_cost, reduce_cost, err = bin.checkCouponLeastAndReduceCost(post,false); err != nil {
			return "", err
		}
		item["least_cost"] = least_cost
		item["reduce_cost"] = reduce_cost
	}

	//折扣券
	if cardType == CARD_TYPE_DISCOUNT {
		var discount int
		if discount ,err = bin.checkCardDiscount(post,false); err != nil {
			return "", err
		}
		item["discount"] = discount
	}

	//兑换券
	if cardType == CARD_TYPE_GIFT {
		var gift string
		if gift ,err = bin.checkCardGift(post,false); err != nil {
			return "", err
		}
		item["gift"] = gift
	}

	advancedInfo ,err := bin.checkCouponAdvancedInfo(post, false)
	if err != nil {
		return "", err
	}
	item["advanced_info"] = advancedInfo

	//整理数据成json格式
	card := map[string]interface{}{}
	card["card_type"] = strings.ToUpper(cardType)
	card[strings.ToLower(cardType)] = item

	bin.Url = bin.cgiUrl("/card/create", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]interface{}{"card":card})
	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return "", &ParseError{"CouponCreate","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}

	cardId, err := js.Get("card_id").String()
	if err != nil {
		return "", &ParseError{"CouponCreate","card_id"}
	}
	return cardId, nil
}


//优惠券更新
func (bin *Interface) CouponUpdate(cardId string, post *map[string]string, origin *map[string]interface{}) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}

	new ,old := *post , *origin
	//删除接口不支持更新的字段
	for _ , field := range []string{"brandName","skuQuantity","defaultDetail","ucAcceptCategory","ucRejectCategory",
		"ucLeastCost","ucObjectUseFor","ucCanUseWithOtherDiscount","source" } {
		delete(new, field)
	}

	if len(new) < 1 {
		return false , errors.New("There is no need to update the data")
	}
	cardType := old["cardType"].(string)
	if !CheckCardType(cardType) {
		return false , errors.New("Incorrect cardType")
	}

	//整理数据成json格式
	item := map[string]interface{}{}

	baseInfo ,err := bin.checkCouponBaseInfo(post, origin)
	if err != nil {
		return false, err
	}
	item["base_info"] = baseInfo

	advancedInfo ,err := bin.checkCouponAdvancedInfo(post, true)
	if err != nil {
		return false, err
	}
	item["advanced_info"] = advancedInfo

	card := map[string]interface{}{}
	card["card_id"] = cardId
	card[strings.ToLower(cardType)] = item

	bin.Url = bin.cgiUrl("/card/update", accessToken)
	bin.Method = "POST"
	bin.Raw(card)
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false, &ParseError{"CouponUpdate","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}


	sendCheck, err := js.Get("send_check").Bool()
	if err != nil {
		return false, &ParseError{"CouponUpdate","send_check"}
	}
	return sendCheck, nil
}

//优惠券快速创建接口，
// 注：该接口不验证参数，在调用前请先验证参数的准确性
func (bin *Interface) CouponCreateQuick(post map[string]string) (string , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}
	cardType, ok := post["cardType"]
	if !ok || !CheckCardType(cardType) {
		return "" , errors.New("Incorrect cardType")
	}
	//整理数据成json格式
	item := map[string]interface{}{
		"base_info"		:bin.formatBaseInfo(post),
		"advanced_info"	:bin.formatAdvancedInfo(post),
	}
	//不同类型的优惠券的专有信息
	switch cardType {
	case CARD_TYPE_GROUPON:
		item["deal_detail"] = post["dealDetail"]
		break;
	case CARD_TYPE_GENERAL_COUPON:
		item["default_detail"] = post["defaultDetail"]
		break;
	case CARD_TYPE_CASH:
		item["least_cost"] = post["leastCost"]
		item["reduce_cost"] = post["reduceCost"]
		break;
	case CARD_TYPE_DISCOUNT:
		item["discount"] = post["discount"]
		break;
	case CARD_TYPE_GIFT:
		item["gift"] = post["gift"]
		break;
	}
	//整理数据成json格式
	card := map[string]interface{}{}
	card["card_type"] = strings.ToUpper(cardType)
	card[strings.ToLower(cardType)] = item

	bin.Url = bin.cgiUrl("/card/create", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]interface{}{"card":card})
	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}
	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return "", &ParseError{"CouponCreate","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}
	cardId, err := js.Get("card_id").String()
	if err != nil {
		return "", &ParseError{"CouponCreate","card_id"}
	}
	return cardId, nil
}
//优惠券快速更新接口
// 注：该接口不验证参数，在调用前请先验证参数的准确性
func (bin *Interface) CouponUpdateQuick(cardId, cardType string, post map[string]string) (bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	//删除接口不支持更新的字段
	for _ , field := range []string{"brandName","skuQuantity", "useCustomCodeMode", "source",
		"ucAcceptCategory", "ucRejectCategory","ucLeastCost","ucObjectUseFor", "ucCanUseWithOtherDiscount",
		"cardType", "dealDetail","defaultDetail","leastCost","reduceCost","discount","gift",
	} {
		delete(post, field)
	}
	if len(post) < 1 {
		return true , nil
	}
	//整理数据成json格式
	card := map[string]interface{}{}
	card["card_id"] = cardId
	card[strings.ToLower(cardType)] = map[string]interface{}{
		"base_info"		:bin.formatBaseInfo(post),
		"advanced_info"	:bin.formatAdvancedInfo(post),
	}

	bin.Url = bin.cgiUrl("/card/update", accessToken)
	bin.Method = "POST"
	bin.Raw(card)
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false, &ParseError{"CouponUpdate","errcode"}
	}
	if  errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}
	sendCheck, err := js.Get("send_check").Bool()
	if err != nil {
		return false, &ParseError{"CouponUpdate","send_check"}
	}
	return sendCheck, nil
}








