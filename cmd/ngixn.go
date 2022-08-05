package cmd

import (
	output2 "nginx-config-reader/output"
	"os/exec"
	"regexp"
)

func CheckNginx() (string, error) {
	command := exec.Command("nginx", "-t")
	output, err := command.CombinedOutput()
	if err != nil {
		output2.ErrorFatal(1001, err)
	}
	result := string(output)
	compile, err := regexp.Compile("the configuration file(.*?)syntax is ok")
	if err != nil {
		output2.ErrorFatal(1002, err)
	}

	submatch := compile.FindStringSubmatch(result)
	if len(submatch) < 2 {
		output2.ErrorFatal(1003, result)
	}
	return submatch[1], nil
}
