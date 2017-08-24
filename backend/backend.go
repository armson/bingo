package backend

import (
	"github.com/armson/bingo"
	"github.com/bitly/go-simplejson"
	"github.com/armson/bingo/utils"
	"github.com/armson/bingo/config"
	"strings"
)

type Backend struct {
	Tracer bingo.Tracer
	Url string
	Queries map[string]string
	Method string
}

func (b *Backend) Send() (*simplejson.Json, error) {
	method := strings.ToUpper(b.Method)

	c := Handle(method, b.Url)
	c.Params(b.Queries)
	data, err := c.Bytes()
	if err != nil { return nil, err }

	if config.Bool("default","enableLog") && config.Bool("backend","enableLog") {
		message := utils.String.Join("Cost:",c.CostTime()," ",method," ",c.Url())
		if method != "GET" {
			message = utils.String.Join(message, " Query: ", c.Query())
		}
		message = utils.String.Join(message, " Response: ", string(data))
		b.Tracer.Logs("Backend",message)
	}
	return simplejson.NewJson(data)
}

