package bingo

import (
	"io"
	"os"
)

const (
	DebugMode   string = "debug"
	TestMode    string = "test"
	ReleaseMode string = "release"
)
const (
	debugCode   = iota
	testCode
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
	case TestMode:
		runMode = testCode
	default:
		panic("run mode unknown: " + value)
	}
	modeName = value
}

func Mode() string {
	return modeName
}
