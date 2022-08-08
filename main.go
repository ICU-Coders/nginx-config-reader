package main

import (
	"github.com/fatih/color"
	"log"
	"nginx-config-reader/cmd"
	"nginx-config-reader/config"
	"nginx-config-reader/output"
	"nginx-config-reader/parser"
)

func main() {

	color.New(color.FgGreen).Println(`
*********************Welcome to use***********************
*                                                        *
*    _|           _|    _|_|_|_|_|_|    _|_|_|_|_|_|     *
*    _| _|        _|    _|              _|         |     *
*    _|   _|      _|    _|              _|_|_|_|_|_|     *
*    _|     _|    _|    _|              _| _|            *
*    _|        _| _|    _|              _|    _|         *
*    _|           _|    _|_|_|_|_|_|    _|       __|     *
*                                                        *
********************nginx config reader*******************
`)

	if len(config.RootPath) == 0 {
		var err error
		config.RootPath, err = cmd.CheckNginx()
		if err != nil {
			output.ErrorFatal(1004, err)
		}
	}
	// 从根配置开始解析，如果用户输入文件，解析用户文件，否则按照nginx -t解析输出文件
	content, err := parser.ConfigParser(config.RootPath)
	if err != nil {
		log.Fatal(err)
	}

	color.New(color.FgCyan).Println("> Generate results...")
	output.Nginx(content)

}
