package lib

import (
	"errors"
	"github.com/manifoldco/promptui"
	"strings"
)

type ListItem struct {
	Number int
	Label string
}

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
