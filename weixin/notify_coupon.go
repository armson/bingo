package weixin

import (
	"encoding/xml"
)

// 解析领券
func GetWxCouponUserGetCard(body []byte) (*WxCouponUserGetCard , error) {
	event := new(WxCouponUserGetCard)
	if err := xml.Unmarshal(body, event); err != nil {
		return nil, xmlDecodeInvalid
	}
	return event , nil
}

// 解析核销券
func GetWxCouponUserConsumeCard(body []byte) (*WxCouponUserConsumeCard , error) {
	event := new(WxCouponUserConsumeCard)
	if err := xml.Unmarshal(body, event); err != nil {
		return nil, xmlDecodeInvalid
	}
	return event , nil
}

// 解析审核事件推送
func GetWxCouponCardPassCheck(body []byte) (*WxCouponCardPassCheck , error) {
	event := new(WxCouponCardPassCheck)
	if err := xml.Unmarshal(body, event); err != nil {
		return nil, xmlDecodeInvalid
	}
	return event , nil
}




