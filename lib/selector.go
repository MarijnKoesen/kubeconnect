package lib

import (
	"errors"
	"github.com/manifoldco/promptui"
	"strings"
)

// ListItem is a Item that can be displayed in the Selector as a selection item
type ListItem struct {
	Number int
	Label string
}

// SelectFromList Shows a shell based list Selector in which the user can select a ListItem
func SelectFromList(question string, selectedLabel string, items []ListItem) (int, error) {
	prompt := promptui.Select{
		Label: question,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Inactive: "{{ .Number }}) {{ .Label }}",
			Selected: selectedLabel + ": {{ .Label | cyan }}",
			Active:   "{{ .Number }}) ▸ {{ .Label | cyan }}",
		},
		Searcher: func(input string, index int) bool {
			context := items[index]
			name := strings.Replace(string(context.Number)+" "+strings.ToLower(context.Label), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input)
		},
		Size: 15,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return -1, errors.New("prompt failed")
	}

	return i, nil
}
