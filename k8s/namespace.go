package k8s

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type Namespace struct {
	Name   string
}

func GetNamespaces(context Context) ([]Namespace, error) {
	cmd := exec.Command("kubectl", "get", "ns", "--context", context.Name);

	var stdout,stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil,errors.New(stderr.String())
	}

	lines := strings.Split(strings.Trim(stdout.String(), "\n"), "\n")

	var namespaces []Namespace
	for _, line := range lines {
		var name = strings.Split(line, " ")[0]
		if name != "NAME" {
			namespaces = append(namespaces, Namespace{Name: name})
		}
	}

	return namespaces,nil
}
