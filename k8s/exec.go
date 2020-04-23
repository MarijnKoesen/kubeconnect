package k8s

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// RunCmd executes an arbitrary kubectl command
func RunCmd(arg ...string) (out string, err error) {
	cmd := exec.Command("kubectl", arg...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		err = errors.New(stderr.String())
		return
	}

	out = strings.Trim(stdout.String(), "\n")
	return
}
