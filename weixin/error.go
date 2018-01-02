package weixin

import (
	"errors"
	"strconv"
)

var (
	requestWithoutResponse = errors.New("The request timeout or without a response")
	accessTokenInvalid = errors.New("access token is invalid")
	xmlDecodeInvalid =  errors.New("xml.Unmarshal is error")
)

type Error struct {
	code int
	message string
}
func (e *Error) Error() string { return "[" + strconv.Itoa(e.code) + "] " + e.message }

type ParseError struct {
	f string
	param string
}
func (e *ParseError) Error() string { return "[" + e.f + ":" + e.param + "] Do't parse valid value" }


