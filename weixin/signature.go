package weixin

import(
	"sort"
	"github.com/armson/bingo/utils"
	"github.com/armson/bingo/encrypt"
)

//params的参数，
// 适用微信api接入签名验证、微信服务器事件推送签名验证
func Signature(args ...string) string {
	params := []string{token()}
	params = append(params, args...)
	sort.Strings(params)
	s := utils.Slice.Join(params,"")
	return encrypt.Sha1([]byte(s))
}
//使用范例
//func checkeToken(c *bingo.Context)  {
//	signature := weixin.Signature(
//		c.GET("timestamp"),
//		c.GET("nonce"),
//	)
//	if signature == c.GET("signature") {
//		fmt.Fprintf(c.Writer , c.GET("echostr"))
//		return
//	}
//	fmt.Fprintf(c.Writer , "false")
//}



