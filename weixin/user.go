package weixin

import (
	"encoding/json"
	"github.com/armson/bingo/utils"
)

//获取用户基本信息（包括UnionID机制）
func (bin *Interface) UserDetail(openId string) (map[string]string, error) {
	users, err := bin.UserSearch([]string{openId})
	if err != nil {
		return nil, err
	}
	return users[openId], nil
}

//获取用户基本信息（包括UnionID机制）- 批量
func (bin *Interface) UserSearch(openIds []string) (map[string]map[string]string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}
	bin.Url = bin.cgiUrl("/cgi-bin/user/info/batchget", accessToken)
	bin.Method = "POST"
	userList := []map[string]string{}
	for _, openId := range openIds {
		userList = append(userList, map[string]string{
			"openid" 	:openId,
			"lang"		:"zh_CN",
		})
	}
	raw := map[string]interface{}{"user_list":userList}
	bin.Raw(raw)
	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	//解析返回状态
	errCode, err := js.Get("errcode").Int()
	if  err == nil || errCode > 0 {
		errMsg, _ := js.Get("errmsg").String()
		return nil, &Error{errCode, errMsg}
	}

	userInfoListJs := js.Get("user_info_list")
	userInfoList, err := userInfoListJs.Array()
	if err != nil {
		return nil, &ParseError{"UserSearch","user_info_list"}
	}
	if len(userInfoList) == 0 {
		return map[string]map[string]string{}, nil
	}

	users := map[string]map[string]string{}
	for i := 0; i < len(userInfoList); i ++ {
		u, _ := userInfoListJs.GetIndex(i).Map()
		openid := u["openid"].(string)
		subscribe, _ := u["subscribe"].(json.Number).Int64()
		user := map[string]string{
			"subscribe"	:utils.Int.String(subscribe),
			"openid"	:u["openid"].(string),
			"unionid"	:u["unionid"].(string),
		}
		if subscribe != 0 {
			sex, _ := u["sex"].(json.Number).Int64()
			user["sex"] = utils.Int.String(sex)
			user["nickname"] = u["nickname"].(string)
			user["city"] = u["city"].(string)
			user["province"] = u["province"].(string)
			user["country"] = u["country"].(string)
			user["headimgurl"] = u["headimgurl"].(string)
			user["remark"] = u["remark"].(string)
			subscribeTime, _ := u["subscribe_time"].(json.Number).Int64()
			user["subscribeTime"] = utils.Int.String(subscribeTime)
			groupId, _ := u["groupid"].(json.Number).Int64()
			user["groupId"] = utils.Int.String(groupId)


		}
		users[openid] = user
	}
	return users, nil
}