package weixin

import (
	"github.com/armson/bingo/utils"
	"fmt"
	"strings"
	"errors"
)


func (bin *Interface) checkCouponBaseInfo(post *map[string]string , origin *map[string]interface{}) (map[string]interface{} , error) {
	new, old , isUpdate := *post, *origin, false
	if len(old) > 0 { isUpdate = true }

	baseInfo := map[string]interface{}{}
	for _, item := range []string{"logoUrl","brandName","color","dateInfoType","skuQuantity","title"} {
		value ,ok := new[item]
		if isUpdate && ok {
			baseInfo[utils.String.UnderLine(item)] = value
			continue
		}
		if !isUpdate {
			if !ok || value == "" {
				return nil , fmt.Errorf("Filed base_info.%s can not be empty", item)
			}
			baseInfo[utils.String.UnderLine(item)] = value
		}
	}

	//判断brandName
	if _, ok := baseInfo["brand_name"]; ok {
		if len(baseInfo["brand_name"].(string)) > 36 {
			return nil , fmt.Errorf("Filed base_info.%s, Word count limit for 12 Chinese characters", "brand_name")
		}
	}

	//判断logoUrl
	if _, ok := baseInfo["logo_url"]; ok  && !checkHttpUrl(baseInfo["logo_url"].(string)) {
		return nil , fmt.Errorf("Filed base_info.%s is not the HTTP url", "logo_url")
	}
	//判断优惠券颜色
	if _, ok := baseInfo["color"]; ok  && !CheckCardColor(baseInfo["color"].(string)) {
		return nil , fmt.Errorf("Filed base_info.%s is invalid", "color")
	}

	//判断有效期
	if _, ok := baseInfo["date_info_type"]; ok {
		dateInfo , err := bin.checkCouponDateInfo(post, origin)
		if err != nil {
			return nil ,err
		}
		baseInfo["date_info"] = dateInfo
	}
	// 判断库存
	if _, ok := baseInfo["sku_quantity"]; ok {
		sku, err := bin.checkCardSkuQuantity(baseInfo["sku_quantity"].(string))
		if err != nil {
			return nil, err
		}
		baseInfo["sku"] = sku
	}
	// 判断卡券名
	if _, ok := baseInfo["title"]; ok {
		if len(baseInfo["title"].(string)) > 27 {
			return nil , fmt.Errorf("Filed base_info.%s, Word count limit for 9 Chinese characters", "title")
		}
	}
	//判断码型
	if value, ok := new["codeType"]; ok {
		if !CheckCardCodeType(value) {
			return nil , fmt.Errorf("Filed base_info.%s is invalid", "code_type")
		}
		baseInfo["code_type"] = value
	}
	//判断卡券使用提醒
	if value, ok := new["notice"]; ok {
		if len(value) > 48 {
			return nil , fmt.Errorf("Filed base_info.%s ,byte length should be less than 48 ", "notice")
		}
		baseInfo["notice"] = value
	}
	//判断卡客服电话
	if value, ok := new["servicePhone"]; ok {
		if len(value) > 48 {
			return nil , fmt.Errorf("Filed base_info.%s ,byte length should be less than 24 ", "service_phone")
		}
		baseInfo["service_phone"] = value
	}
	//判断卡券使用说明
	if value, ok := new["description"]; ok {
		if len(value) > 3072 {
			return nil , fmt.Errorf("Filed base_info.%s ,byte length should be less than 3072 ", "description")
		}
		baseInfo["description"] = value
	}
	//每人可领券的数量限制
	if value, ok := new["useLimit"]; ok {
		i := utils.String.Int(value)
		if i < 1 || i > 100000000 {
			return nil, fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 100000000", "use_limit")
		}
		baseInfo["use_limit"] = i
	}
	//每人可核销的数量限制
	if value, ok := new["getLimit"]; ok {
		i := utils.String.Int(value)
		if i < 1 || i > 100000000 {
			return nil, fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 100000000", "get_limit")
		}
		baseInfo["get_limit"] = i
	}
	//是否指定用户领取
	if value, ok := new["bindOpenid"]; ok {
		baseInfo["bind_openid"] = utils.String.Bool(value)
	}
	//卡券领取页面是否可分享
	if value, ok := new["canShare"]; ok {
		baseInfo["can_share"] = utils.String.Bool(value)
	}
	//卡券是否可转赠
	if value, ok := new["canGiveFriend"]; ok {
		baseInfo["can_give_friend"] = utils.String.Bool(value)
	}
	//0:非自定义Code码 1:自定义Code码 2:导入code模式
	if value, ok := new["useCustomCodeMode"]; ok {
		i := utils.String.Int(value)
		if !CheckCardUseCustomCodeMode(i) {
			return nil , errors.New("use_custom_code_mode is invalid , value must be in [0,1,2]")
		}
		switch i {
		case 0:
			baseInfo["use_custom_code"] = false
			break;
		case 1:
			baseInfo["use_custom_code"] = true
			break;
		case 2:
			baseInfo["use_custom_code"] = true
			baseInfo["get_custom_code_mode"] = "GET_CUSTOM_CODE_MODE_DEPOSIT"
			if !isUpdate {
				baseInfo["sku"] = map[string]interface{}{"quantity":0}
			}
			break;
		}
	}
	//设置该卡券的适用门店
	if value, ok := new["locationIdList"]; ok && value != "" {
		baseInfo["location_id_list"] = strings.Split(value,",")
	}
	//设置来源
	if value, ok := new["source"]; ok && value != "" {
		baseInfo["source"] = value
	}

	//使用场景入口,服务场景入口,营销场景入口
	for _ , item := range []string{"centerTitle", "customUrlName", "promotionUrlName"}{
		if value, ok := new[item]; ok {
			if (isUpdate == true || (isUpdate == false && value != ""))  && len(value) <= 15 {
				baseInfo[utils.String.UnderLine(item)] = value
			}
		}
	}

	for _ , item := range []string{"centerSubTitle", "customUrlSubTitle", "promotionUrlSubTitle"}{
		if value, ok := new[item]; ok {
			if (isUpdate == true || (isUpdate == false && value != ""))  && len(value) <= 18 {
				baseInfo[utils.String.UnderLine(item)] = value
			}
		}
	}

	for _ , item := range []string{"centerUrl", "customUrl", "promotionUrl"}{
		if value, ok := new[item]; ok {
			if (isUpdate == true && (value == "" || checkHttpUrl(value)) || (isUpdate == false && checkHttpUrl(value))) {
				baseInfo[utils.String.UnderLine(item)] = value
			}
		}
	}
	//删除多余的字段
	for _, item := range []string{"date_info_type","sku_quantity"} {
		delete(baseInfo, item)
	}
	return baseInfo ,nil
}

//校验base.date_info
func (bin *Interface) checkCouponDateInfo(post *map[string]string, origin *map[string]interface{}) (map[string]interface{} , error) {
	new, old , isUpdate := *post, *origin, false
	if len(old) > 0 { isUpdate = true }

	date_info_type, ok := new["dateInfoType"]
	if !ok || !CheckCardDateInfoType(date_info_type) {
		return nil , fmt.Errorf("Filed base_info.date_info.%s is invalid", "type")
	}

	dateInfo := map[string]interface{}{}
	if value , ok := old["dateInfo"]; ok && isUpdate == true {
		dateInfo = value.(map[string]interface{})
	}
	if  isUpdate == true && dateInfo["type"] != date_info_type {
		return nil , fmt.Errorf("Filed base_info.date_info.%s,Make Sure OldDateInfoType==NewDateInfoType", "type")
	}

	if date_info_type == DATE_TYPE_FIX_TIME_RANGE {
		begin_timestamp, ok := new["dateInfoBeginTimestamp"]
		if !ok {
			return nil , fmt.Errorf("Filed base_info.date_info.%s can not be empty", "begin_timestamp")
		}
		end_timestamp, ok := new["dateInfoEndTimestamp"]
		if !ok {
			return nil , fmt.Errorf("Filed base_info.date_info.%s can not be empty", "end_timestamp")
		}
		beginTimestamp 	:= utils.String.Int(begin_timestamp)
		endTimestamp 	:= utils.String.Int(end_timestamp)
		if endTimestamp < beginTimestamp {
			return nil, fmt.Errorf("Filed base_info.date_info.%s ，End time should be greater than start time", "end_timestamp")
		}
		if isUpdate == true && (beginTimestamp > dateInfo["beginTimestamp"].(int) || endTimestamp < dateInfo["endTimestamp"].(int)) {
			return nil , errors.New("Filed base_info.date_info, Make Sure NewBeginTime<=OldBeginTime && NewEndTime>=OldEndTime")
		}
		return map[string]interface{}{
			"type"				:date_info_type,
			"begin_timestamp"	:beginTimestamp,
			"end_timestamp"		:endTimestamp,
		},nil
	}
	if date_info_type == DATE_TYPE_FIX_TERM {
		var fixedBeginTerm int
		var fixedTerm int
		fixed_begin_term, ok := new["dateInfoFixedBeginTerm"]
		if ok {
			fixedBeginTerm = utils.String.Int(fixed_begin_term)
		}
		fixed_term, ok := new["dateInfoFixedTerm"]
		if !ok {
			return nil , fmt.Errorf("Filed base_info.date_info.%s can not be empty", "fixed_term")
		}
		fixedTerm = utils.String.Int(fixed_term)

		if fixedTerm < 1 {
			return nil, fmt.Errorf("Filed base_info.date_info.%s , number of days must be greater than 0", "fixed_term")
		}
		return map[string]interface{}{
			"type"				:date_info_type,
			"fixed_begin_term"	:fixedBeginTerm,
			"fixed_term"		:fixedTerm,
		},nil
	}
	return nil , fmt.Errorf("Filed base_info.date_info.%s is invalid", "type")
}

//校验base.sku
func (bin *Interface) checkCardSkuQuantity(quantity string) (map[string]interface{} , error) {
	if  skuQuantity := utils.String.Int(quantity); skuQuantity >= 0 && skuQuantity <= 100000000 {
		return map[string]interface{}{"quantity":skuQuantity} , nil
	}
	return nil, fmt.Errorf("Filed base_info.sku.%s , must be greater than 0 and less than 100000000", "quantity")
}

//校验团购券时的，优惠券详情
func (bin *Interface) checkCouponDealDetail(post *map[string]string , isUpdate bool) (string , error) {
	info := *post
	deal_detail , ok := info["dealDetail"]
	if isUpdate == false {
		if !ok || len(deal_detail) < 1 || len(deal_detail) > 3072  {
			return "", fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 3072", "deal_detail")
		}
		return deal_detail, nil
	}
	if isUpdate == true && ok {
		if len(deal_detail) < 0 || len(deal_detail) > 3072  {
			return "", fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 3072", "deal_detail")
		}
		return deal_detail, nil
	}
	return "" , nil
}

//校验通用优惠券时的，优惠券详情
func (bin *Interface) checkCouponDefaultDetail(post *map[string]string , isUpdate bool) (string , error) {
	info := *post
	default_detail , ok := info["defaultDetail"]
	if isUpdate == false {
		if !ok || len(default_detail) < 1 || len(default_detail) > 3072  {
			return "", fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 3072", "default_detail")
		}
		return default_detail, nil
	}
	if isUpdate == true && ok {
		if len(default_detail) < 1 || len(default_detail) > 3072  {
			return "", fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 3072", "default_detail")
		}
		return default_detail, nil
	}
	return "" , nil
}

//校验代金券时的，起用金额和减免金额
func (bin *Interface) checkCouponLeastAndReduceCost(post *map[string]string , isUpdate bool) (int ,int , error) {
	info := *post
	least_cost , leastOk := info["leastCost"]
	reduce_cost ,reduceOk := info["reduceCost"]
	if isUpdate == false {
		if !leastOk || !reduceOk {
			return -1, -1, fmt.Errorf("Filed base_info.%s , must be greater than or equal to 0 and less than 1000000", "least_cost/reduce_cost")
		}
		leastCost  := utils.String.Int(least_cost)
		reduceCost := utils.String.Int(reduce_cost)
		if leastCost < 0  || leastCost  > 1000000 || reduceCost < 1  || reduceCost  > 1000000 {
			return -1, -1, fmt.Errorf("Filed base_info.%s , must be greater than or equal to 0 and less than 1000000", "least_cost/reduce_cost")
		}
		return leastCost, reduceCost, nil
	}
	if isUpdate == true {
		leastCost, reduceCost := -1 , -1
		if leastOk  { leastCost   = utils.String.Int(least_cost)  }
		if reduceOk { reduceCost  = utils.String.Int(reduce_cost) }

		if leastOk && (leastCost < 0  || leastCost  > 1000000)  {
			return -1, -1, fmt.Errorf("Filed base_info.%s , must be greater than or equal to 0 and less than 1000000", "least_cost")
		}
		if reduceOk && (reduceCost < 1  || reduceCost  > 1000000)  {
			return -1, -1, fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 1000000", "reduce_cost")
		}
		return leastCost , reduceCost , nil
	}
	return -1 ,-1 , nil
}

//校验折扣券时的，打折额度
func (bin *Interface) checkCardDiscount(post *map[string]string , isUpdate bool) (int , error) {
	info := *post
	discount , ok := info["discount"]
	if isUpdate == false {
		if !ok {
			return 0 , fmt.Errorf("Filed base_info.date_info.%s can not be empty", "discount")
		}
		discount := utils.String.Int(discount)
		if discount < 1 && discount > 99 {
			return 0, fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 100", "discount")
		}
		return discount, nil
	}
	if isUpdate == true && ok {
		discount := utils.String.Int(discount)
		if discount < 1 && discount > 99 {
			return 0, fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 100", "discount")
		}
		return discount, nil
	}
	return 0 , nil
}

//校验兑换券时的，兑换内容的名称
func (bin *Interface) checkCardGift(post *map[string]string , isUpdate bool) (string , error) {
	info := *post
	gift , ok := info["gift"]
	if isUpdate == false {
		if !ok || len(gift) < 1 || len(gift) > 3072  {
			return "", fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 3072", "gift")
		}
		return gift, nil
	}
	if isUpdate == true && ok {
		if len(gift) < 1 || len(gift) > 3072  {
			return "", fmt.Errorf("Filed base_info.%s , must be greater than 0 and less than 3072", "gift")
		}
		return gift, nil
	}
	return "" , nil
}


// 转换成微信公众号的请求格式，对参数的有效性不做校验
func (bin *Interface) formatBaseInfo(post map[string]string) map[string]interface{} {
	baseInfo := map[string]interface{}{}
	for _ , field := range []string{"brandName", "logoUrl", "codeType","title","color","notice",
		"servicePhone","description", "centerTitle","centerSubTitle","centerUrl",
		"customUrlName","customUrl","customUrlSubTitle", "promotionUrlName","promotionUrlSubTitle",
		"promotionUrl","source", "centerAppBrandUserName", "centerAppBrandPass",
		"customAppBrandUserName","customAppBrandPass","promotionAppBrandUserName","promotionAppBrandPass",
		} {
		if value, ok := post[field]; ok {
			baseInfo[utils.String.UnderLine(field)] = value
		}
	}
	for _ , field := range []string{"useLimit","getLimit"} {
		if value, ok := post[field]; ok {
			baseInfo[utils.String.UnderLine(field)] = utils.String.Int(value)
		}
	}
	for _ , field := range []string{"bindOpenid","canShare","canGiveFriend"} {
		if value, ok := post[field]; ok {
			baseInfo[utils.String.UnderLine(field)] = utils.String.Bool(value)
		}
	}
	//日期
	if dateInfoType, ok := post["dateInfoType"]; ok {
		dateInfo := map[string]interface{}{"type":dateInfoType}
		switch dateInfoType {
		case "DATE_TYPE_FIX_TIME_RANGE":
			dateInfo["begin_timestamp"] = post["dateInfoBeginTimestamp"]
			dateInfo["end_timestamp"] = post["dateInfoEndTimestamp"]
			break;
		case "DATE_TYPE_FIX_TERM":
			dateInfo["fixed_begin_term"] = post["dateInfoFixedBeginTerm"]
			dateInfo["fixed_term"] = post["dateInfoFixedTerm"]
			break;
		}
		baseInfo["date_info"] = dateInfo
	}
	//库存
	if quantity, ok := post["skuQuantity"]; ok {
		if  post["useCustomCodeMode"] == "2" {
			baseInfo["sku"] = map[string]string{"quantity":"0"}
		} else {
			baseInfo["sku"] = map[string]string{"quantity":quantity}
		}
	}
	//0:非自定义Code码 1:自定义Code码 2:导入code模式
	if value, ok := post["useCustomCodeMode"]; ok {
		switch value {
		case "0":
			baseInfo["use_custom_code"] = false
			break;
		case "1":
			baseInfo["use_custom_code"] = true
			break;
		case "2":
			baseInfo["use_custom_code"] = true
			baseInfo["get_custom_code_mode"] = "GET_CUSTOM_CODE_MODE_DEPOSIT"
			break;
		}
	}
	//设置该卡券的适用门店
	if value, ok := post["locationIdList"]; ok {
		if value != "" {
			baseInfo["location_id_list"] = strings.Split(value,",")
			baseInfo["use_all_locations"] = false
		} else {
			baseInfo["use_all_locations"] = true
		}
	}
	return baseInfo
}


