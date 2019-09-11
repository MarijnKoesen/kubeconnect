package k8s

import (
	"bytes"
	"errors"
	"kubeconnect/lib"
	"os/exec"
	"strings"
)

// Pod represents a Pod in a Kubernetes Cluster
type Pod struct {
	Name string
}

// PodListItems Transforms a list of Pods to a list of ListItems to show in the Selector
func PodListItems(pods []Pod) []lib.ListItem {
	var items []lib.ListItem

	for index, pod := range pods {
		items = append(items, lib.ListItem{Number: index + 1, Label: pod.Name})
	}

	return items
}

// GetPods returns all Pods in a given Namespace and Context
func GetPods(context Context, namespace Namespace) ([]Pod, error) {
	cmd := exec.Command("kubectl", "get", "pod", "--context", context.Name, "--namespace", namespace.Name)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	lines := strings.Split(strings.Trim(stdout.String(), "\n"), "\n")

	var pods []Pod
	for _, line := range lines {
		var name = strings.Split(line, " ")[0]
		if name != "NAME" {
			pods = append(pods, Pod{Name: name})
		}
	}

	return pods, nil
}
