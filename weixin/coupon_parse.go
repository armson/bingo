package weixin

import (
	"github.com/bitly/go-simplejson"
	"strings"
)

//解析券码详情
func (bin *Interface) getCardBaseInfo(js *simplejson.Json) (map[string]interface{} , error) {
	detail := map[string]interface{}{}
	cardType, _ := js.Get("card_type").String()
	detail["cardType"] = cardType
	detail["cardTypeName"] = GetCardTypeName(cardType)

	cardJs := js.Get(strings.ToLower(cardType))
	baseInfoJs := cardJs.Get("base_info")

	detail["id"], _ 		= baseInfoJs.Get("id").String()
	detail["logoUrl"], _ 	= baseInfoJs.Get("logo_url").String()
	detail["brandName"], _ 	= baseInfoJs.Get("brand_name").String()
	detail["title"], _ 		= baseInfoJs.Get("title").String()
	detail["colorValue"], _ = baseInfoJs.Get("color").String()
	detail["status"], _ 	= baseInfoJs.Get("status").String()
	detail["createTime"],_	= baseInfoJs.Get("create_time").Int()
	detail["updateTime"],_	= baseInfoJs.Get("update_time").Int()
	detail["color"] = GetCardColorName(detail["colorValue"].(string))

	dateInfoJs := baseInfoJs.Get("date_info")
	dateInfo := map[string]interface{}{}
	if value , err := dateInfoJs.Get("type").String(); err == nil {
		dateInfo["type"] = value
	}
	if value , err := dateInfoJs.Get("begin_timestamp").Int(); err == nil {
		dateInfo["beginTimestamp"] = value
	}
	if value , err := dateInfoJs.Get("end_timestamp").Int(); err == nil {
		dateInfo["endTimestamp"] = value
	}
	if value , err := dateInfoJs.Get("fixed_begin_term").Int(); err == nil {
		dateInfo["fixedBeginTerm"] = value
	}
	if value , err := dateInfoJs.Get("fixed_term").Int(); err == nil {
		dateInfo["fixedTerm"] = value
	}
	detail["dateInfo"] = dateInfo

	skuJs := baseInfoJs.Get("sku")
	sku := map[string]int{}
	if value , err := skuJs.Get("quantity").Int(); err == nil {
		sku["quantity"] = value
	}
	if value , err := skuJs.Get("total_quantity").Int(); err == nil {
		sku["totalQuantity"] = value
	}
	detail["sku"] = sku

	if value , err := baseInfoJs.Get("code_type").String(); err == nil {
		detail["codeType"] = value
		detail["codeTypeValue"] = GetCardCodeTypeName(value)
	}
	if value , err := baseInfoJs.Get("notice").String(); err == nil {
		detail["notice"] = value
	}
	if value , err := baseInfoJs.Get("service_phone").String(); err == nil {
		detail["servicePhone"] = value
	}
	if value , err := baseInfoJs.Get("description").String(); err == nil {
		detail["description"] = value
	}
	if value , err := baseInfoJs.Get("use_limit").Int(); err == nil {
		detail["useLimit"] = value
	}
	if value , err := baseInfoJs.Get("get_limit").Int(); err == nil {
		detail["getLimit"] = value
	}
	if value , err := baseInfoJs.Get("bind_openid").Bool(); err == nil {
		detail["bindOpenid"] = value
	}
	if value , err := baseInfoJs.Get("can_share").Bool(); err == nil {
		detail["canShare"] = value
	}
	if value , err := baseInfoJs.Get("can_give_friend").Bool(); err == nil {
		detail["canGiveFriend"] = value
	}

	//是否自定义Code码
	detail["useCustomCodeMode"] = 0
	if useCustomCode , err := baseInfoJs.Get("use_custom_code").Bool(); err == nil && useCustomCode == true {
		detail["useCustomCodeMode"] = 1
		if mode , err := baseInfoJs.Get("get_custom_code_mode").String(); err == nil && mode == "GET_CUSTOM_CODE_MODE_DEPOSIT" {
			detail["useCustomCodeMode"] = 2
		}
	}
	detail["useCustomCodeModeName"] = GetCardUseCustomCodeModeName(detail["useCustomCodeMode"].(int))

	locationIdListJs := baseInfoJs.Get("location_id_list")
	if value , err := locationIdListJs.Array(); err == nil && len(value) > 0 {
		slice := []int{}
		for k , _ := range value {
			if locationId , err := locationIdListJs.GetIndex(k).Int(); err == nil {
				slice = append(slice, locationId)
			}
		}
		detail["locationIdList"] = slice
	}
	if value , err := baseInfoJs.Get("use_all_locations").Bool(); err == nil {
		detail["useAllLocations"] = value
	}
	//使用场景入口,服务场景入口,营销场景入口
	if value , err := baseInfoJs.Get("center_title").String(); err == nil {
		detail["centerTitle"] = value
	}
	if value , err := baseInfoJs.Get("center_sub_title").String(); err == nil {
		detail["centerSubTitle"] = value
	}
	if value , err := baseInfoJs.Get("center_url").String(); err == nil {
		detail["centerUrl"] = value
	}
	if value , err := baseInfoJs.Get("custom_url_name").String(); err == nil {
		detail["customUrlName"] = value
	}
	if value , err := baseInfoJs.Get("custom_url_sub_title").String(); err == nil {
		detail["customUrlSubTitle"] = value
	}
	if value , err := baseInfoJs.Get("custom_url").String(); err == nil {
		detail["customUrl"] = value
	}
	if value , err := baseInfoJs.Get("promotion_url_name").String(); err == nil {
		detail["promotionUrlName"] = value
	}
	if value , err := baseInfoJs.Get("promotion_url_sub_title").String(); err == nil {
		detail["promotionUrlSubTitle"] = value
	}
	if value , err := baseInfoJs.Get("promotion_url").String(); err == nil {
		detail["promotionUrl"] = value
	}
	if value , err := baseInfoJs.Get("source").String(); err == nil {
		detail["source"] = value
	}

	advancedInfoJs := js.Get(strings.ToLower(cardType)).Get("advanced_info")
	useConditionJs := advancedInfoJs.Get("use_condition")
	useCondition := map[string]interface{}{}
	if value , err := useConditionJs.Get("accept_category").String(); err == nil {
		useCondition["acceptCategory"] = value
	}
	if value , err := useConditionJs.Get("reject_category").String(); err == nil {
		useCondition["rejectCategory"] = value
	}
	if value , err := useConditionJs.Get("object_use_for").String(); err == nil {
		useCondition["objectUseFor"] = value
	}
	if value , err := useConditionJs.Get("least_cost").Int(); err == nil {
		useCondition["leastCost"] = value
	}
	if value , err := useConditionJs.Get("can_use_with_other_discount").Bool(); err == nil {
		useCondition["canUseWithOtherDiscount"] = value
	}
	detail["useCondition"] = useCondition

	abstractJs := advancedInfoJs.Get("abstract")
	abstract := map[string]interface{}{}
	if value , err := abstractJs.Get("abstract").String(); err == nil {
		abstract["name"] = value
	}
	if value , err := abstractJs.Get("icon_url_list").StringArray(); err == nil {
		abstract["iconUrl"] = value[0]
	}
	detail["abstract"] = abstract

	textImageListJs := advancedInfoJs.Get("text_image_list")
	textImageList := []map[string]interface{}{}
	if items , err := textImageListJs.Array(); err == nil && len(items) > 0 {
		for i := 0 ; i < len(items); i ++ {
			itemsJs := textImageListJs.GetIndex(i)
			m := map[string]interface{}{}
			if imageUrl , err := itemsJs.Get("image_url").String(); err == nil {
				m["imageUrl"] = imageUrl
			}
			if text , err := itemsJs.Get("text").String(); err == nil {
				m["text"] = text
			}
			if len(m) > 0 {
				textImageList = append(textImageList, m)
			}
		}
	}
	detail["textImageList"] = textImageList

	businessService := []string{}
	if slice, err := advancedInfoJs.Get("business_service").StringArray(); err == nil {
		businessService = slice
	}
	detail["businessService"] = businessService


	timeLimitJs := advancedInfoJs.Get("time_limit")
	timeLimit := []map[string]interface{}{}
	if items , err := timeLimitJs.Array(); err == nil && len(items) > 0 {
		for i := 0 ; i < len(items); i ++ {
			itemsJs := timeLimitJs.GetIndex(i)
			if t , err := itemsJs.Get("type").String(); err == nil {
				m := map[string]interface{}{}
				m["type"] = t
				if value , err := itemsJs.Get("begin_hour").Int(); err == nil {
					m["beginHour"] = value
				}
				if value , err := itemsJs.Get("begin_minute").Int(); err == nil {
					m["beginMinute"] = value
				}
				if value , err := itemsJs.Get("end_hour").Int(); err == nil {
					m["endHour"] = value
				}
				if value , err := itemsJs.Get("end_minute").Int(); err == nil {
					m["endMinute"] = value
				}
				timeLimit = append(timeLimit, m)
			}
		}
	}
	detail["timeLimit"] = timeLimit

	//团购券
	if value , err := cardJs.Get("deal_detail").String(); err == nil {
		detail["dealDetail"] = value
	}
	//通用券
	if value , err := cardJs.Get("default_detail").String(); err == nil {
		detail["defaultDetail"] = value
	}
	//代金券
	if value , err := cardJs.Get("least_cost").Int(); err == nil {
		detail["leastCost"] = value
	}
	if value , err := cardJs.Get("reduce_cost").Int(); err == nil {
		detail["reduceCost"] = value
	}
	//折扣券
	if value , err := cardJs.Get("discount").Int(); err == nil {
		detail["discount"] = value
	}
	//兑换券
	if value , err := cardJs.Get("gift").String(); err == nil {
		detail["gift"] = value
	}
	return detail, nil
}