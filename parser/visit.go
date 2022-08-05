package parser

var visitMap = make(map[string]bool)

func isVisited(path string) bool {
	if has := visitMap[path]; has {
		return true
	}
	visitMap[path] = true
	return false
}
