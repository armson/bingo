package weixin

import (
	"encoding/xml"
)

// 解析
func GetWxMessageText(body []byte) (*WxMessageText , error) {
	event := new(WxMessageText)
	if err := xml.Unmarshal(body, event); err != nil {
		return nil, xmlDecodeInvalid
	}
	return event , nil
}


