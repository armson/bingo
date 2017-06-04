package bingo

import (
	"fmt"
	"runtime"
)

var version = "Bingo v1.0.1"

func SetVersion(ver string) {
	version = ver
}

func GetVersion() string {
	return version
}

func PrintVersion() {
	fmt.Printf(`%s, Compiler: %s %s, Copyright (C) 2017 Armson.All Rights Reserved`,
		version,
		runtime.Compiler,
		runtime.Version())
	fmt.Println()
}

func VersionMiddleware() HandlerFunc {
	return func(c *Context) {
		c.Header("X-DRONE-VERSION", version)
		c.Next()
	}
}
