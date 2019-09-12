package k8s

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func runCmd(arg ...string) (out string, err error) {
	// nolint:gosec
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
