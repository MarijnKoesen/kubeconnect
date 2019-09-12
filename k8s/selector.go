package k8s

import "kubeconnect/lib"

// ContextListItems Transform a list of Context to a list of ListItems to show in the Selector
func ContextListItems(contexts []Context) (items []lib.ListItem) {
	for index, context := range contexts {
		items = append(items, lib.ListItem{Number: index + 1, Label: context.name})
	}

	return items
}

// NamespaceListItems Transform a list of Namespace to a list of ListItems to show in the Selector
func NamespaceListItems(namespaces []Namespace) (items []lib.ListItem) {
	for index, namespace := range namespaces {
		items = append(items, lib.ListItem{Number: index + 1, Label: namespace.name})
	}

	return items
}

// PodListItems Transform a list of Pod to a list of ListItems to show in the Selector
func PodListItems(pods []Pod) (items []lib.ListItem) {
	for index, pod := range pods {
		items = append(items, lib.ListItem{Number: index + 1, Label: pod.Name})
	}

	return items
}
