package k8s

import (
	"strings"
)

// Context represents a Context in a Kubernetes Cluster
type Context struct {
	name string
}

// GetContexts returns all configured Contexts in kubectl
func GetContexts() (contexts []Context, err error) {
	out, err := runCmd("config", "get-contexts", "-o", "name")
	if err != nil {
		return
	}

	for _, name := range strings.Split(out, "\n") {
		contexts = append(contexts, Context{
			name: name,
		})
	}

	return
}
