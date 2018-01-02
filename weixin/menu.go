package weixin

import (
	"github.com/armson/bingo/utils"
	"github.com/bitly/go-simplejson"
	"strings"
)

type Menus struct {
	*Menu
	SubButton []*Menu
}
type Menu struct {
	Type,Name, Key, Value, Url, MediaId, AppId, PagePath string
}


//创建菜单
func (bin *Interface) MenuCreate(json string) ( bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}

	//对 & 符号进行处理
	json = strings.Replace(json, "\\u0026", "%26", -1)
	json = strings.Replace(json, "&", "%26", -1)
	json = utils.String.Join("{\"button\":",json,"}")

	code ,err := utils.Json.Decode([]byte(json))
	if err != nil {
		return false , err
	}

	bin.Url = bin.cgiUrl("/cgi-bin/menu/create", accessToken)
	bin.Method = "POST"
	bin.Raw(code)
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}

	errCode, err := js.Get("errcode").Int()
	if err == nil && errCode == 0 {
		return true , nil
	}
	errMsg, _ := js.Get("errmsg").String()
	return false, &Error{errCode, errMsg}
}

//清除菜单
//例如：weixin.New(c).MenuDrop()
func (bin *Interface) MenuDrop() ( bool , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return false, err
	}
	bin.Url = bin.cgiUrl("/cgi-bin/menu/delete", accessToken)
	bin.Method = "GET"
	js , err :=  bin.Send()
	if err != nil {
		return false, requestWithoutResponse
	}
	errCode, err := js.Get("errcode").Int()
	if err != nil {
		return false, &ParseError{"MenuDelete","errcode"}
	}

	if errCode != 0 {
		errMsg, _ := js.Get("errmsg").String()
		return false, &Error{errCode, errMsg}
	}
	return true,nil
}


//获取配置菜单列表
//例如：weixin.New(c).MenuSelfList()
func (bin *Interface) MenuSelfList() ([]*Menus , error) {
	m := []*Menus{}
	accessToken, err := bin.Token()
	if err != nil {
		return m, err
	}
	bin.Url = bin.cgiUrl("/cgi-bin/get_current_selfmenu_info", accessToken)
	bin.Method = "GET"
	js , err :=  bin.Send()
	if err != nil {
		return m, requestWithoutResponse
	}
	is_menu_open, err := js.Get("is_menu_open").Int()
	if err != nil  || is_menu_open != 1 {
		return m , nil
	}

	buttonJs := js.Get("selfmenu_info").Get("button")
	buttons, err := buttonJs.Array()
	if err != nil  || len(buttons) == 0 {
		return m ,nil
	}

	//解析到Menu
	for key, _ := range buttons {
		menu := bin.menuParse(buttonJs.GetIndex(key))
		if menu == nil { continue }

		listJs := buttonJs.GetIndex(key).Get("sub_button").Get("list")
		list, err := listJs.Array()

		subs := []*Menu{}
		if err == nil  && len(list) > 0 {
			for k, _ := range list {
				if menuSub := bin.menuParse(listJs.GetIndex(k)); menuSub != nil {
					subs = append(subs, menuSub)
				}
			}
		}
		m = append(m, &Menus{menu,subs})
	}
	return m,nil
}

//获取通过接口定义菜单列表
//例如：weixin.New(c).MenuList()
func (bin *Interface) MenuList() ([]*Menus , error) {
	m := []*Menus{}
	accessToken, err := bin.Token()
	if err != nil {
		return m, err
	}
	bin.Url = bin.cgiUrl("/cgi-bin/menu/get", accessToken)
	bin.Method = "GET"
	js , err :=  bin.Send()
	if err != nil {
		return m, requestWithoutResponse
	}

	//当无配置菜单时，接口返回错误
	errCode, err := js.Get("errcode").Int()
	if err == nil || errCode > 0  {
		return m, nil
	}

	buttonJs := js.Get("menu").Get("button")
	buttons, err := buttonJs.Array()
	if err != nil  || len(buttons) == 0 {
		return m ,nil
	}

	//解析到Menu
	for key, _ := range buttons {
		menu := bin.menuParse(buttonJs.GetIndex(key))
		if menu == nil { continue }

		listJs := buttonJs.GetIndex(key).Get("sub_button")
		list, err := listJs.Array()

		subs := []*Menu{}
		if err == nil  && len(list) > 0 {
			for k, _ := range list {
				if menuSub := bin.menuParse(listJs.GetIndex(k)); menuSub != nil {
					subs = append(subs, menuSub)
				}
			}
		}
		m = append(m, &Menus{menu,subs})
	}
	return m,nil
}




func (bin *Interface) menuParse(js *simplejson.Json) (*Menu) {
	name, err := js.Get("name").String()
	if err != nil {
		return nil
	}
	t, _ := js.Get("type").String()
	key, _ := js.Get("key").String()
	value, _ := js.Get("value").String()
	url, _ := js.Get("url").String()
	mediaId, _ := js.Get("media_id").String()
	appId, _ := js.Get("appid").String()
	pagePath, _ := js.Get("pagepath").String()

	m := &Menu{
		Name     : name,
		Type     : t,
		Key      : key,
		Value    : value,
		Url      : url,
		MediaId  : mediaId,
		AppId    : appId,
		PagePath : pagePath,
	}
	return m
}

