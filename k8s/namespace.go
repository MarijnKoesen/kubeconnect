package k8s

import (
	"bytes"
	"errors"
	"kubeconnect/lib"
	"os/exec"
	"strings"
)

// Namespace represents a Namespace in a Kubernetes cluster
type Namespace struct {
	Name string
}

// NamespaceListItems Transforms a list of Namespaces to a list of ListItems to show in the Selector
func NamespaceListItems(namespaces []Namespace) []lib.ListItem {
	var items []lib.ListItem

	for index, namespace := range namespaces {
		items = append(items, lib.ListItem{Number: index + 1, Label: namespace.Name})
	}

	return items
}

// GetNamespaces returns all namespaces in a given Context
func GetNamespaces(context Context) ([]Namespace, error) {
	cmd := exec.Command("kubectl", "get", "ns", "--context", context.Name)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	lines := strings.Split(strings.Trim(stdout.String(), "\n"), "\n")

	var namespaces []Namespace
	for _, line := range lines {
		var name = strings.Split(line, " ")[0]
		if name != "NAME" {
			namespaces = append(namespaces, Namespace{Name: name})
		}
	}

	return namespaces, nil
}
