package lib

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"strings"
)

type ListItem struct {
	Number int
	Label string
}

func SelectFromList(question string, items []ListItem) (int, error) {
	prompt := promptui.Select{
		Label: question,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Inactive: "{{ .Number }}) {{ .Label }}",
			Selected: "Using: {{ .Label | cyan }}",
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
		fmt.Printf("Prompt failed %v\n", err)
		return -1, errors.New("prompt failed")
	}

	return i, nil
}
