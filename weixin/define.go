package weixin

import(
	"github.com/armson/bingo/attach"
	"strings"
)

const (
	menu_type_click = "click"
	menu_type_view = "view"
	menu_type_scancode_push = "scancode_push"
	menu_type_scancode_waitmsg = "scancode_waitmsg"
	menu_type_pic_sysphoto = "pic_sysphoto"
	menu_type_pic_photo_or_album = "pic_photo_or_album"
	menu_type_pic_weixin = "pic_weixin"
	menu_type_location_select = "location_select"
	menu_type_media_id = "media_id"
	menu_type_view_limited = "view_limited"
	menu_type_text = "text"
	menu_type_news = "news"

	material_upload_image_max_size =  1048576 // 1M
	material_upload_image_always_max_size =  2097152 // 2M

	CouponCodeEachTimeImportMaximum = 100
)
var material_upload_image_permit_ext = []string{".bmp",".png",".jpeg",".jpg",".gif"}

func MaterialUploadImageMaxSize() int {
	return material_upload_image_max_size
}
func MaterialUploadImageAlwaysMaxSize() int {
	return material_upload_image_always_max_size
}
func CheckMaterialImageMime(file *attach.Attachment) bool {
	for _, v := range material_upload_image_permit_ext {
		if m , ok := attach.Mimes[v]; ok && m == file.Mime && v == file.Ext {
			return true
		}
	}
	return false
}
func MaterialUploadImagePermitExt() []string {
	return material_upload_image_permit_ext
}

//0:非自定义Code码 1:自定义Code码 2:导入code模式
//0:use_custom_code=false
//1:use_custom_code=true
//2:use_custom_code=true get_custom_code_mode=GET_CUSTOM_CODE_MODE_DEPOSIT
var card_use_custom_code_mode = map[int]string{
	0 : "非自定义Code码",
	1 : "自定义Code码",
	2 : "导入code模式",
}
func CheckCardUseCustomCodeMode(mode int) (bool) {
	if _, ok := card_use_custom_code_mode[mode]; ok { return true }
	return false
}
func GetCardUseCustomCodeModeName(mode int) (string) {
	if name, ok := card_use_custom_code_mode[mode]; ok { return name }
	return ""
}



var card_colors = map[string]string{
	"Color010"	:"#63b359",
	"Color020"	:"#2c9f67",
	"Color030"	:"#509fc9",
	"Color040"	:"#5885cf",
	"Color050"	:"#9062c0",
	"Color060"	:"#d09a45",
	"Color070"	:"#e4b138",
	"Color080"	:"#ee903c",
	"Color081"	:"#f08500",
	"Color082"	:"#a9d92d",
	"Color090"	:"#dd6549",
	"Color100"	:"#cc463d",
	"Color101"	:"#cf3e36",
	"Color102"	:"#5E6671",
}

func GetCardColors() (map[string]string) {
	return card_colors
}
func CheckCardColor(colorName string) (bool) {
	if _, ok := card_colors[colorName]; ok {
		return true
	}
	return false
}
func GetCardColorName(colorValue string) string {
	for name, value := range card_colors {
		if value == colorValue { return name }
	}
	return ""
}
func GetCardColorValue(colorName string) string {
	if value, ok := card_colors[colorName]; ok {
		return value
	}
	return ""
}


//DISCOUNT: 折扣券,可为用户提供消费折扣
//CASH: 代金券,可为用户提供抵扣现金服务。可设置成为“满*元，减*元”
//GIFT: 兑换券,可为用户提供消费送赠品服务
//GROUPON: 团购券,可为用户提供团购套餐服务
//GENERAL_COUPON: 优惠券,即“通用券”，建议当以上四种无法满足需求时采用
const(
	CARD_TYPE_GROUPON 			= "GROUPON"
	CARD_TYPE_CASH 				= "CASH"
	CARD_TYPE_DISCOUNT 			= "DISCOUNT"
	CARD_TYPE_GIFT 				= "GIFT"
	CARD_TYPE_GENERAL_COUPON 	= "GENERAL_COUPON"
)
var card_types = map[string]string{
	CARD_TYPE_GROUPON		:"折扣券",
	CARD_TYPE_CASH			:"代金券",
	CARD_TYPE_DISCOUNT		:"兑换券",
	CARD_TYPE_GIFT			:"团购券",
	CARD_TYPE_GENERAL_COUPON:"通用券",
}
func GetCardTypes() (map[string]string) {
	return card_types
}
func CheckCardType(typeName string) (bool) {
	if _, ok := card_types[typeName]; ok { return true }
	return false
}
func GetCardTypeName(code string) (string) {
	if name, ok := card_types[code]; ok { return name }
	return ""
}

//卡券生效日期
//DATE_TYPE_FIX_TIME_RANGE:表示固定日期区间，例如：2017-09-07 00:00:00 至 2017-09-17 23:59:59
//DATE_TYPE_FIX_TERM:表示固定时长（自领取后按天算),例如：3 天生效，有效天数 30 天
const (
	DATE_TYPE_FIX_TIME_RANGE 	= "DATE_TYPE_FIX_TIME_RANGE"
	DATE_TYPE_FIX_TERM 			= "DATE_TYPE_FIX_TERM"
)
var card_date_info_types = map[string]string{
	DATE_TYPE_FIX_TIME_RANGE 	: "固定日期",
	DATE_TYPE_FIX_TERM 			: "固定时长",
}
func CheckCardDateInfoType(typeName string) (bool) {
	if _, ok := card_date_info_types[typeName]; ok {
		return true
	}
	return false
}


const(
	CARD_STATUS_NOT_VERIFY 	= "CARD_STATUS_NOT_VERIFY"
	CARD_STATUS_VERIFY_FAIL = "CARD_STATUS_VERIFY_FAIL"
	CARD_STATUS_VERIFY_OK	= "CARD_STATUS_VERIFY_OK"
	CARD_STATUS_DELETE 		= "CARD_STATUS_DELETE"
	CARD_STATUS_DISPATCH 	= "CARD_STATUS_DISPATCH"
)
var card_status = map[string]string{
	CARD_STATUS_NOT_VERIFY	:"待审核",
	CARD_STATUS_VERIFY_FAIL	:"审核失败",
	CARD_STATUS_VERIFY_OK	:"通过审核",
	CARD_STATUS_DELETE		:"卡券被商户删除",
	CARD_STATUS_DISPATCH	:"在公众平台投放过的卡券",
}
func CardStatusToSlice(s string) []string {
	if s == "" {return nil }
	r := []string{}
	for _, v := range strings.Split(s, ",") {
		v = strings.ToUpper(v)
		if _, ok := card_status[v]; ok {
			r = append(r, v)
		}
	}
	return r
}
func GetCardStatusName(statusValue string) string {
	if name, ok := card_status[statusValue]; ok {
		return name
	}
	return ""
}


const (
	CARD_CODE_TYPE_TEXT 		= "CODE_TYPE_TEXT"
	CARD_CODE_TYPE_BARCODE 		= "CODE_TYPE_BARCODE"
	CARD_CODE_TYPE_QRCODE 		= "CODE_TYPE_QRCODE"
	CARD_CODE_TYPE_ONLY_QRCODE 	= "CODE_TYPE_ONLY_QRCODE"
	CARD_CODE_TYPE_ONLY_BARCODE = "CODE_TYPE_ONLY_BARCODE"
	CARD_CODE_TYPE_NONE 		= "CODE_TYPE_NONE"
)
var card_code_types = map[string]string{
	CARD_CODE_TYPE_TEXT			:"文本",
	CARD_CODE_TYPE_BARCODE		:"一维码",
	CARD_CODE_TYPE_QRCODE		:"二维码",
	CARD_CODE_TYPE_ONLY_QRCODE	:"二维码无code显示",
	CARD_CODE_TYPE_ONLY_BARCODE	:"一维码无code显示",
	CARD_CODE_TYPE_NONE			:"不显示code和条形码类型",
}
func CheckCardCodeType(code string) (bool) {
	if _, ok := card_code_types[code]; ok {
		return true
	}
	return false
}
func GetCardCodeTypeName(code string) string {
	if name, ok := card_code_types[code]; ok {
		return name
	}
	return ""
}

var card_business_service = map[string]string{
	"BIZ_SERVICE_DELIVER"	:"外卖服务",
	"BIZ_SERVICE_FREE_PARK"	:"停车位",
	"BIZ_SERVICE_WITH_PET"	:"可带宠物",
	"BIZ_SERVICE_FREE_WIFI"	:"免费wif",
}

func CardBusinessServiceToSlice(s string) []string {
	r := []string{}
	if s == "" {return r }
	for _, v := range strings.Split(s, ",") {
		v = strings.ToUpper(v)
		if _, ok := card_business_service[v]; ok {
			r = append(r, v)
		}
	}
	return r
}


var card_time_limit_types = map[string]string{
	"MONDAY"	:"周一",
	"TUESDAY"	:"周二",
	"WEDNESDAY"	:"周三",
	"THURSDAY"	:"周四",
	"FRIDAY"	:"周五",
	"SATURDAY"	:"周六",
	"SUNDAY"	:"周日",
}

func CheckCardTimeLimitType(t string) (bool) {
	if _, ok := card_time_limit_types[t]; ok {
		return true
	}
	return false
}

var card_user_card_status = map[string]string{
	"NORMAL"		:"正常",
	"CONSUMED"		:"已核销",
	"EXPIRE"		:"已过期",
	"GIFTING"		:"转赠中",
	"GIFT_TIMEOUT"	:"转赠超时",
	"DELETE"		:"已删除",
	"UNAVAILABLE"	:"已失效",
}

func GetCardUserCardStatusName(code string) string {
	if name, ok := card_user_card_status[code]; ok {
		return name
	}
	return ""
}


type WxEvent struct {
	ToUserName		string	`xml:"ToUserName"`
	FromUserName	string	`xml:"FromUserName"`
	CreateTime		int		`xml:"CreateTime"`
	MsgType			string	`xml:"MsgType"`
	Event			string	`xml:"Event"`
	Encrypt			string	`xml:"Encrypt"`
}

type WxEnumClick struct {
	WxEvent
	EventKey		string	`xml:"EventKey"`
}
type WxCouponUserGetCard struct {
	WxEvent
	CardId					string	`xml:"CardId"`
	IsGiveByFriend			int		`xml:"IsGiveByFriend"`
	UserCardCode			string	`xml:"UserCardCode"`
	FriendUserName			string	`xml:"FriendUserName"`
	OuterId					string	`xml:"OuterId"`
	OldUserCardCode			string	`xml:"OldUserCardCode"`
	IsRestoreMemberCard		int		`xml:"IsRestoreMemberCard"`
	IsRecommendByFriend		int		`xml:"IsRecommendByFriend"`
	SourceScene				string	`xml:"SourceScene"`
}

type WxCouponUserConsumeCard struct {
	WxEvent
	CardId					string	`xml:"CardId"`
	UserCardCode			string	`xml:"UserCardCode"`
	ConsumeSource			string	`xml:"ConsumeSource"`
	LocationName			string	`xml:"LocationName"`
	StaffOpenId				string	`xml:"StaffOpenId"`
	VerifyCode				string	`xml:"VerifyCode"`
	RemarkAmount			string	`xml:"RemarkAmount"`
	OuterStr				string	`xml:"OuterStr"`
}
type WxCouponCardPassCheck struct {
	WxEvent
	CardId					string	`xml:"CardId"`
	RefuseReason			string	`xml:"RefuseReason"`
}


type WxReply struct {
	ToUserName		string	`xml:"ToUserName"`
	FromUserName	string	`xml:"FromUserName"`
	CreateTime		int64	`xml:"CreateTime"`
	MsgType			string	`xml:"MsgType"`
}

type WxEncrypt struct {
	Encrypt			string	`xml:"Encrypt"`
	MsgSignature	string	`xml:"MsgSignature"`
	TimeStamp		string	`xml:"TimeStamp"`
	Nonce			string	`xml:"Nonce"`
}

type WxMessageText struct {
	WxEvent
	Content		string	`xml:"Content"`
	MsgId		string	`xml:"MsgId"`
}







