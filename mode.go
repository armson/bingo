package bingo

import (
	"io"
	"os"
)

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
)
const (
	debugCode   = iota
	releaseCode   
)

var DefaultWriter io.Writer = os.Stdout
var DefaultErrorWriter io.Writer = os.Stderr

var runMode = debugCode
var modeName = DebugMode

func SetMode(value string) {
	switch value {
	case DebugMode:
		runMode = debugCode
	case ReleaseMode:
		runMode = releaseCode
	default:
		panic("run mode unknown: " + value)
	}
	modeName = value
}

func Mode() string {
	return modeName
}
