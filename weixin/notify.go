package weixin

import (
	"github.com/armson/bingo/utils"
	"github.com/armson/bingo/encrypt"
	"encoding/xml"
	"errors"
	"bytes"
	"encoding/binary"
	"time"
)

//消息推送入口方法
//例如：
//body, _ := c.Body();
//event , body, err := weixin.Body(body, c.GET("timestamp"), c.GET("nonce"), c.GET("msg_signature"))

func Body(body []byte, timestamp, nonce, signature string) (*WxEvent, []byte, error) {
	if len(body) == 0 {
		return nil, nil, errors.New("The content of the request body is empty")
	}
	event, err := getWxEvent(body)
	if err != nil {
		return nil, nil, err
	}
	if isEncrypt() == false {
		return event, body, nil
	}
	if Signature(timestamp, nonce, event.Encrypt) != signature {
		return nil, nil, errors.New("Request signature verification failed")
	}
	if event.Encrypt == "" {
		return nil, nil, errors.New("Encrypt in the request body content is empty")
	}
	return  decodeWxEvent(event)
}

// 解析xml获得基本数据信息
func getWxEvent(b []byte) (*WxEvent , error) {
	wxEvent := new(WxEvent)
	if err := xml.Unmarshal(b, wxEvent); err != nil {
		return nil, xmlDecodeInvalid
	}
	return wxEvent , nil
}

// 解密消息体
func decodeWxEvent(event *WxEvent) (*WxEvent,[]byte, error) {
	text, _ := decryptAES(event.Encrypt)
	t, err := getWxEvent(text)
	if err != nil {
		return nil , nil, err
	}
	return t , text , nil
}

func decryptAES(cipherText string) ([]byte, []byte) {
	plainText := encrypt.Aes.Decode(cipherText, cryptAESKey())
	le := len(plainText)
	pos := le-len(appId())
	return plainText[20:pos] , plainText[pos:]
}

func encryptAES(v interface{}, timeStamp, nonce, random string) (string, error) {
	body, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return "", err
	}

	//生成4位的body长度
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		return "", err
	}
	bodyLength := buf.Bytes()

	plainText := bytes.Join([][]byte{
		[]byte(random),
		bodyLength,
		body,
		[]byte(appId()),
	}, nil)

	cipherText := encrypt.Aes.Encode(plainText, cryptAESKey())
	return cipherText , nil
}

func output(v interface{}) (interface{}, error)  {
	if isEncrypt() == false {
		return v , nil
	}

	timeStamp := utils.Int.String(time.Now().Unix())
	nonce := utils.String.Signatures(10)
	random := utils.String.Rand(16)
	encrypt, err := encryptAES(v, timeStamp,nonce,random)
	if err != nil {
		return nil , err
	}
	msgSignature := Signature(timeStamp, nonce, encrypt)

	type xml struct {
		WxEncrypt
	}
	out := new(xml)
	out.Encrypt = encrypt
	out.MsgSignature = msgSignature
	out.Nonce = nonce
	out.TimeStamp = timeStamp
	return out , nil
}





