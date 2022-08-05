package config

import "flag"

var RootPath string
var ShowLog bool
var SortType string // port/host/uri/dir
var Match string

func init() {
	flag.StringVar(&RootPath, "i", "", "请输入待解析文件路径(Please enter the path of the file to be parsed)")
	flag.StringVar(&SortType, "s", "port", "输出结果排序依据(Output results are sorted based on)")
	flag.BoolVar(&ShowLog, "l", false, "是否显示日志路径(Whether to display the log path)")
	flag.StringVar(&Match, "m", "", "高亮匹配到的字符(Matching character filter)")
	flag.Parse()
}
