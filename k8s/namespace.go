package k8s

import (
	"fmt"
	"strings"
)

// Namespace represents a Namespace in a Kubernetes cluster
type Namespace struct {
	name, context string
}

const nsTpl = "{{range .items}}{{.metadata.name}} {{end}}"

// GetNamespaces returns all namespaces in a given Context
func (c *Context) GetNamespaces() (namespaces []Namespace, err error) {
	out, err := runCmd("get", "ns", "--context", c.name, "-o", fmt.Sprintf("go-template=%s", nsTpl))
	if err != nil {
		return
	}

	for _, name := range strings.Fields(out) {
		namespaces = append(namespaces, Namespace{
			name:    name,
			context: c.name,
		})
	}

	return
}
