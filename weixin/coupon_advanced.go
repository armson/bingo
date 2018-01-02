package weixin

import (
	"github.com/armson/bingo/utils"
	"github.com/bitly/go-simplejson"
	"encoding/json"
	"strings"
)

func (bin *Interface) checkCouponAdvancedInfo(post *map[string]string , isUpdate bool) (map[string]interface{} , error) {
	info := *post
	advancedInfo := map[string]interface{}{}
	//使用条件
	use_condition := map[string]interface{}{}
	if value, ok := info["ucAcceptCategory"]; ok {
		use_condition["accept_category"] = value
	}
	if value, ok := info["ucRejectCategory"]; ok {
		use_condition["reject_category"] = value
	}
	if value, ok := info["ucLeastCost"]; ok {
		use_condition["least_cost"] = utils.String.Int(value)
	}
	if value, ok := info["ucObjectUseFor"]; ok {
		use_condition["object_use_for"] = value
	}
	if value, ok := info["ucCanUseWithOtherDiscount"]; ok {
		use_condition["can_use_with_other_discount"] = utils.String.Bool(value)
	}
	if len(use_condition) > 0 {
		advancedInfo["use_condition"] = use_condition
	}
	//封面
	abstract := map[string]interface{}{}
	if value, ok := info["abstractAbstract"]; ok && len(value) > 0 && len(value) <= 36  {
		abstract["abstract"] = value
	}
	if value, ok := info["abstractIconUrlList"]; ok && checkHttpUrl(value) {
		abstract["icon_url_list"] = []string{value}
	}
	if len(abstract) > 0 {
		advancedInfo["abstract"] = abstract
	}

	//图文列表，显示在详情内页
	if value, ok := info["textImageList"]; ok {
		advancedInfo["text_image_list"] = []map[string]string{{}}
		if textImageList , err := bin.checkCouponTextImageList(value); err == nil && len(textImageList) > 0 {
			advancedInfo["text_image_list"] = textImageList
		}
	}
	//商家服务类型
	if value, ok := info["businessService"]; ok {
		advancedInfo["business_service"] = CardBusinessServiceToSlice(value);
	}
	//使用时段限制
	if value, ok := info["timeLimit"]; ok {
		advancedInfo["time_limit"] = []map[string]interface{}{{}}
		if time_limit , err := bin.checkCouponTimeLimit(value); err == nil && len(time_limit) > 0 {
			advancedInfo["time_limit"] = time_limit
		}
	}
	return advancedInfo, nil
}

//解析图文列表，显示在详情内页
func (bin *Interface) checkCouponTextImageList(s string ) ([]map[string]string , error) {
	js , err := simplejson.NewJson([]byte(s))
	if err != nil {
		return nil , err
	}
	items , err := js.Array()
	if err != nil || len(items) < 1 {
		return nil , err
	}
	textImageList := []map[string]string{}
	for i := 0 ; i < len(items); i ++ {
		item := map[string]string{}
		if imageUrl , err := js.GetIndex(i).Get("imageUrl").String(); err == nil && checkHttpUrl(imageUrl) {
			item["image_url"] = imageUrl
		}
		if text , err := js.GetIndex(i).Get("text").String(); err == nil {
			item["text"] = text
		}
		if len(item) > 0 {
			textImageList = append(textImageList, item)
		}
	}
	return textImageList , nil
}

//解析timeLimit字符串
func (bin *Interface) checkCouponTimeLimit(s string ) ([]map[string]interface{} , error) {
	js , err := simplejson.NewJson([]byte(s))
	if err != nil {
		return nil , err
	}
	items , err := js.Array()
	if err != nil || len(items) < 1 {
		return nil , err
	}
	timeLimit := []map[string]interface{}{}
	for i := 0 ; i < len(items); i ++ {
		if t , err := js.GetIndex(i).Get("type").String(); err == nil && CheckCardTimeLimitType(t) {
			limit := map[string]interface{}{ "type":t }
			if value , err := js.GetIndex(i).Get("beginHour").Int(); err == nil && value > 0 && value < 24 {
				limit["begin_hour"] = value
			}
			if value , err := js.GetIndex(i).Get("beginMinute").Int(); err == nil && value > 0 && value < 60 {
				limit["begin_minute"] = value
			}
			if value , err := js.GetIndex(i).Get("endHour").Int(); err == nil && value > 0 && value < 24 {
				limit["end_hour"] = value
			}
			if value , err := js.GetIndex(i).Get("endMinute").Int(); err == nil && value > 0 && value < 60 {
				limit["end_minute"] = value
			}
			timeLimit = append(timeLimit, limit)
		}
	}
	return timeLimit,nil
}

// 转换成微信公众号的请求格式，对参数的有效性不做校验
func (bin *Interface) formatAdvancedInfo(post map[string]string) map[string]interface{} {
	advancedInfo := map[string]interface{}{}
	//使用门槛（条件）字段
	useCondition := map[string]interface{}{}
	for _ , field := range []string{"ucAcceptCategory", "ucRejectCategory","ucLeastCost","ucObjectUseFor"} {
		if value, ok := post[field]; ok {
			useCondition[utils.String.UnderLine(field)[3:]] = value
		}
	}
	if value, ok := post["ucCanUseWithOtherDiscount"]; ok {
		useCondition["can_use_with_other_discount"] = utils.String.Bool(value)
	}
	if len(useCondition) > 0 {
		advancedInfo["use_condition"] = useCondition
	}
	//abstract 封面摘要结构体名称
	abstract := map[string]interface{}{}
	for _ , field := range []string{"abstractAbstract", "abstractIconUrlList"} {
		if value, ok := post[field]; ok {
			abstract[utils.String.UnderLine(field)[9:]] = value
		}
	}
	if len(abstract) > 0 {
		advancedInfo["abstract"] = abstract
	}

	//text_image_list 图文列表，显示在详情内页
	if value, ok := post["textImageList"]; ok {
		if value != "" {
			var v []map[string]string
			if err := json.Unmarshal([]byte(value), &v); err == nil {
				for _, rows := range v {
					rows["image_url"] = rows["imageUrl"]
					delete(rows, "imageUrl")
				}
				advancedInfo["text_image_list"] = v
			}
		} else {
			advancedInfo["text_image_list"] = []map[string]string{map[string]string{}}
		}
	}
	//商家服务类型
	if value, ok := post["businessService"]; ok {
		if value != "" {
			advancedInfo["business_service"] = strings.Split(value, ",");
		} else {
			advancedInfo["business_service"] = []string{}
		}
	}
	//使用时段限制
	if value, ok := post["timeLimit"]; ok {
		if value != "" {
			var v []map[string]interface{}
			if err := json.Unmarshal([]byte(value), &v); err == nil {
				for _, rows := range v {
					rows["begin_hour"] 		= rows["beginHour"]
					rows["begin_minute"] 	= rows["beginMinute"]
					rows["end_hour"] 		= rows["endHour"]
					rows["end_minute"] 		= rows["endMinute"]
					delete(rows, "beginHour")
					delete(rows, "beginMinute")
					delete(rows, "endHour")
					delete(rows, "endMinute")
				}
				advancedInfo["time_limit"] = v
			}
		} else {
			advancedInfo["time_limit"] = []map[string]interface{}{map[string]interface{}{}}
		}
	}
	return advancedInfo
}



