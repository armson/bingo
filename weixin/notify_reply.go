package weixin

import (
	"time"
	"errors"
)

func ReplyText(ToUserName, FromUserName, text string) (interface{}, error) {
	if ToUserName == "" || FromUserName == "" {
		return nil , errors.New("The lack of the sender or recipient")
	}
	type xml struct {
		WxReply
		Content			string
	}
	reply := new(xml)
	reply.ToUserName 	= ToUserName
	reply.FromUserName 	= FromUserName
	reply.CreateTime 	= time.Now().Unix()
	reply.MsgType 		= "text"
	reply.Content 		= text
	return output(reply)
}