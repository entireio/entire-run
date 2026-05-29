package ui

import (
	"errors"
	"fmt"

	"charm.land/huh/v2"
)

var ErrSelectCanceled = errors.New("select canceled")

func Select(title, hint string, labels []string) (int, error) {
	if len(labels) == 0 {
		return 0, nil
	}

	options := make([]huh.Option[int], 0, len(labels))
	for i, label := range labels {
		options = append(options, huh.NewOption(label, i))
	}

	selected := 0
	height := len(options) + 2
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(title).
				Description(hint).
				Options(options...).
				Height(height).
				Value(&selected),
		),
	).WithTheme(huh.ThemeFunc(huh.ThemeDracula))

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) || errors.Is(err, huh.ErrTimeout) {
			return 0, ErrSelectCanceled
		}
		return 0, fmt.Errorf("select prompt: %w", err)
	}
	return selected, nil
}
