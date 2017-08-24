package config

import (
	"os"
)

var pwd, _ = os.Getwd()

var defaultConfigs = map[string]map[string]string{
	"default": map[string]string{
		"env": 				"test",
		"runmode": 			"debug",
		"appname": 			"Bingo",
		"servername": 		"Bingo v1.1.0",
		"pid": 				pwd + "/bingo.pid",
		"enablelog": 		"true",
		"accesslog": 		pwd + "/log/access.log",
		"errorlog": 		pwd + "/log/error.log",
		"enablehttp": 		"true",
		"listentcp4":		"tcp4",
		"httpport": 		"10080",
		"httpaddr": 		"",
		"enablehttps": 		"false",
		"httpsaddr": 		"",
		"httpsport": 		"10443",
		"httpscertfile": 	pwd + "/ssl.crt",
		"httpskeyfile": 	pwd + "/ssl.key",
	},
	"cookie": map[string]string{
		"path":				"/",
		"domain":			"*",
		"secure":			"false",
		"httponly":			"false",
	},

	"mysql": map[string]string{
		"enablelog":		"true",
		"charset":			"utf8",
		"maxidleconns":		"100",
		"maxopenconns":		"500",
		"connmaxlifetime":	"3600s",
		"timeout":			"30s",
		"readtimeout":		"30s",
		"writetimeout":		"30s",
	},

	"redis": map[string]string{
		"enablelog":			"true",
		"maxretries":			"3",
		"dialtimeout":			"15s",
		"readtimeout":			"5s",
		"writetimeout":			"5s",
		"poolsize":				"1000",
		"pooltimeout":			"10s",
		"connmaxlifetime":		"45s",
		"idlecheckfrequency":	"45s",
	},
	"task": map[string]string{
		"enablelog":	"true",
		"accesslog": 	pwd + "/log/message.log",
		"errorlog": 	pwd + "/log/error.log",
	},
	"backend": map[string]string{
		"enablelog":	"true",
	},
	"response": map[string]string{
		"enablelog":	"true",
	},
}
