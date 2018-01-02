package weixin

import (
	"github.com/armson/bingo/backend"
	"github.com/armson/bingo"
	"github.com/armson/bingo/config"
	"github.com/armson/bingo/utils"
	"strings"
	"github.com/armson/bingo/encrypt"
	"regexp"
)

type Interface struct {
	backend.Backend
}

func New(trance bingo.Tracer) *Interface {
	n := &Interface{}
	n.Tracer = trance
	n.Queries = map[string]string{}
	return n
}

func (bin *Interface) Params(args map[string]string) *Interface {
	for arg, value := range args {
		bin.Queries[arg] =  value
	}
	return bin
}

//uri , token string
func (bin *Interface) cgiUrl(args ...string) string {
	if len(args) == 1 {
		return utils.String.Join(host(), args[0])
	}

	if strings.Index(args[0], "?") > -1 {
		return utils.String.Join(host(), args[0],"&access_token=",args[1])
	}
	return utils.String.Join(host(), args[0],"?access_token=",args[1])
}

func appId() string { return  config.String("weixin","appId") }
func appSecret() string { return  config.String("weixin","appSecret") }
func token() string { return  config.String("weixin","token") }
func host() string { return  config.String("weixin","host") }
func isEncrypt() bool { return  config.Bool("weixin","isEncrypt") }
func cryptAESKey() []byte {
	var encodeKey string
	if encodeKey = config.String("weixin","encodingAESKey"); encodeKey == "" {
		return []byte{}
	}
	return  encrypt.Base64.Decode(encodeKey + "=")
}
func checkHttpUrl(url string) bool {
	validUrl,_	:= regexp.Compile(`^(http|https)://[a-zA-Z0-9:/.|%&;?=_]+$`)
	return validUrl.MatchString(url) && len(url) < 200
}
func miniAppId() string { return  config.String("weixin","miniAppId") }
func miniAppSecret() string { return  config.String("weixin","miniAppSecret") }
func miniHost() string { return  config.String("weixin","miniHost") }
func mapStringString(v map[string]interface{}) map[string]string {
	mp := map[string]string{}
	if len(v) == 0 {return mp}
	for key, value := range v {
		if val, ok := value.(string); ok {
			mp[key] = val
		} else {
			mp[key] = ""
		}
	}
	return mp
}
