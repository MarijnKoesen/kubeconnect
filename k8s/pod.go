package k8s

import (
	"fmt"
	"strings"
)

// Pod represents a Pod in a Kubernetes Cluster
type Pod struct {
	Name, Namespace, Context string
}

const podTpl = "{{range .items}}{{.metadata.name}} {{end}}"

// GetPods returns all Pods in a given Namespace and Context
func (n *Namespace) GetPods() (pods []Pod, err error) {
	out, err := runCmd("get", "pod", "--context", n.context, "--namespace", n.name, "-o", fmt.Sprintf("go-template=%s", podTpl))
	if err != nil {
		return
	}

	for _, name := range strings.Fields(out) {
		pods = append(pods, Pod{
			Name:      name,
			Namespace: n.name,
			Context:   n.context,
		})
	}

	return
}
