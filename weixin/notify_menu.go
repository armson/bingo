package weixin

import (
	"encoding/xml"
)

// 解析
func GetWxEnumClick(body []byte) (*WxEnumClick , error) {
	event := new(WxEnumClick)
	if err := xml.Unmarshal(body, event); err != nil {
		return nil, xmlDecodeInvalid
	}
	return event , nil
}


