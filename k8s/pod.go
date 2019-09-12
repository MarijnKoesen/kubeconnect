package k8s

import (
	"fmt"
	"strings"
)

// Pod represents a Pod in a Kubernetes Cluster
type Pod struct {
	Name, Namespace, Context string
	Containers               []string
}

const podTpl = "{{range .items}}{{.metadata.name}} {{range .spec.containers}}{{ .name }} {{end}}{{\"\\n\"}}{{end}}"

// GetPods returns all Pods in a given Namespace and Context
func (n *Namespace) GetPods() (pods []Pod, err error) {
	out, err := runCmd("get", "pod", "--context", n.context, "--namespace", n.name, "-o", fmt.Sprintf("go-template=%s", podTpl))
	if err != nil {
		return
	}

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		pods = append(pods, Pod{
			Name:       fields[0],
			Containers: fields[1:],
			Namespace:  n.name,
			Context:    n.context,
		})
	}

	return
}
