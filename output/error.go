package output

import (
	"github.com/fatih/color"
	"os"
)

var ErrorCode = map[int]string{
	1001: "Reading nginx root file error",
	1002: "Reading nginx root file error",
	1003: "Reading nginx root file error",
	1004: "Reading nginx root file error",
}
var errorColor = color.New(color.FgRed)

func ErrorFatal(code int, msg interface{}) {
	errorColor.Println(code, ErrorCode[code], msg)
	os.Exit(1)
}
