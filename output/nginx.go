package output

import (
	"fmt"
	"github.com/ICU-Coders/table"
	"nginx-config-reader/config"
	"nginx-config-reader/parser"
	"sort"
	"strconv"
	"strings"
)

type Location struct {
	URI        string
	Dir        string
	Listen     string
	ServerName string
	Method     string
}

type NginxLogPath struct {
	Type string
	Path string
}

type NginxLog struct {
	Module   string
	LogPaths []NginxLogPath
}

func Nginx(nginx *parser.ObjectExpression) {
	var paths []*Location
	var nginxLogs []NginxLog
	// https -> servers -> locations
	nginxLogs = append(nginxLogs, NginxLog{
		Module:   "default",
		LogPaths: makeNginxLog(nginx),
	})
	https := nginx.GetMustObject("http")
	for _, http := range https {

		nginxLogs = append(nginxLogs, NginxLog{
			Module:   "http",
			LogPaths: makeNginxLog(http),
		})
		servers := http.GetMustObject("server")
		for _, server := range servers {
			listen := server.GetFirst("listen").(*parser.MultiValue).GetValueWithSep(" ")
			defaultRoot := server.GetFirstMustString("root")
			locations := server.GetMustObject("location")
			serverName := server.GetFirstMustString("server_name")

			nginxLogs = append(nginxLogs, NginxLog{
				Module:   serverName + ":" + listen,
				LogPaths: makeNginxLog(http),
			})

			for _, location := range locations {
				lo := new(Location)
				uri := location.GetFirst("path").(*parser.MultiValue).GetValue()
				body := location.GetFirstMustObject("body")
				dir := body.GetFirstMustString("root")
				lo.Method = "root"
				alias := body.GetFirstMustString("alias")
				if len(dir) == 0 {
					dir = alias
					lo.Method = "alias"
				}
				proxyPass := body.GetFirstMustString("proxy_pass")
				if len(dir) == 0 {
					dir = proxyPass
					lo.Method = "proxy"
				}
				if len(dir) == 0 {
					dir = defaultRoot
					lo.Method = "default"
				}
				lo.URI = uri
				lo.Dir = dir
				lo.Listen = listen
				lo.ServerName = serverName
				paths = append(paths, lo)
			}
		}
	}

	if config.ShowLog {
		fmt.Println()
		fmt.Println()

		fmt.Println("Log Table:")
		var content [][]string
		for _, nginxLog := range nginxLogs {
			for _, path := range nginxLog.LogPaths {
				var sub []string
				sub = append(sub, nginxLog.Module, path.Type, path.Path)
				content = append(content, sub)
			}

		}
		table.Show([]string{"Module", "Type", "Path"}, content)

		fmt.Println()
		fmt.Println()
	}

	switch strings.ToLower(config.SortType) {
	case "host", "server", "server_name":
		sort.Slice(paths, func(i, j int) bool {
			return paths[i].ServerName < paths[j].ServerName
		})
	case "uri":
		sort.Slice(paths, func(i, j int) bool {
			return paths[i].URI < paths[j].URI
		})
	case "dir":
		sort.Slice(paths, func(i, j int) bool {
			return paths[i].Dir < paths[j].Dir
		})
	default:
		sort.Slice(paths, func(i, j int) bool {
			ii, err := strconv.ParseInt(paths[i].Listen, 10, 64)
			if err != nil {
				return paths[i].Listen < paths[j].Listen
			}
			jj, err := strconv.ParseInt(paths[j].Listen, 10, 64)
			if err != nil {
				return paths[i].Listen < paths[j].Listen
			}
			return ii < jj
		})
	}

	var content [][]string
	for _, path := range paths {
		var sub []string
		if len(config.Match) > 0 {
			if strings.Contains(path.ServerName, config.Match) ||
				strings.Contains(path.Listen, config.Match) ||
				strings.Contains(path.URI, config.Match) ||
				strings.Contains(path.Method, config.Match) ||
				strings.Contains(path.Dir, config.Match) {
				sub = append(sub, path.ServerName, path.Listen, path.URI, path.Method, path.Dir)
			}
		} else {
			sub = append(sub, path.ServerName, path.Listen, path.URI, path.Method, path.Dir)
		}
		content = append(content, sub)
	}

	fmt.Println("Index Table:")

	table.Show([]string{"ServerName", "Listen", "URI", "Method", "Dir"}, content)

}

func makeNginxLog(nginx *parser.ObjectExpression) []NginxLogPath {
	if !config.ShowLog {
		return nil
	}
	var logs []NginxLogPath
	errorLogs := nginx.GetMustString("error_log")
	for _, errorLog := range errorLogs {
		logs = append(logs, NginxLogPath{
			Type: "error",
			Path: errorLog,
		})
	}
	accessLogs := nginx.GetMustString("access_log")
	for _, accessLog := range accessLogs {
		logs = append(logs, NginxLogPath{
			Type: "access",
			Path: accessLog,
		})
	}
	return logs
}
