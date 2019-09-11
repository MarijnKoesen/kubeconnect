package k8s

import (
	"bytes"
	"errors"
	"kubeconnect/lib"
	"os/exec"
	"strings"
)

// Context represents a Context in a Kubernetes Cluster
type Context struct {
	Name string
}

// ContextListItems Transform a list of Contexts to a list of ListItems to show in the Selector
func ContextListItems(contexts []Context) []lib.ListItem {
	var items []lib.ListItem

	for index, context := range contexts {
		items = append(items, lib.ListItem{Number: index + 1, Label: context.Name})
	}

	return items
}

// GetContexts returns all configured Contexts in kubectl
func GetContexts() ([]Context, error) {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o=name")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	lines := strings.Split(strings.Trim(stdout.String(), "\n"), "\n")

	var contexts []Context
	for _, name := range lines {
		if name != "NAME" {
			contexts = append(contexts, Context{Name: name})
		}
	}

	return contexts, nil
}
