package main

import (
	"strings"

	"github.com/fatih/color"
)

type settingsData struct {
	supressedMessages   map[string]bool
	processColors       []color.Attribute
	processColorsLength int
}

type settings interface {
	ShouldSuppressMessage(message string) bool
	GetProcessColor(index int) color.Attribute
}

func getSettings() settings {
	supressedMessages := make(map[string]bool)
	supressedMessages["Terminate batch job (Y/N)?"] = true

	processColors := []color.Attribute{
		color.FgGreen,
		color.FgYellow,
		color.FgBlue,
		color.FgMagenta,
		color.FgCyan,
		color.FgHiGreen,
		color.FgHiYellow,
		color.FgHiBlue,
		color.FgHiMagenta,
		color.FgHiCyan}

	return &settingsData{
		supressedMessages:   supressedMessages,
		processColors:       processColors,
		processColorsLength: len(processColors)}
}

func (s *settingsData) ShouldSuppressMessage(message string) bool {
	test := strings.TrimRight(strings.TrimRight(strings.TrimRight(message, "\r\n"), "\n"), " ")

	_, ok := s.supressedMessages[test]

	return ok
}

func (s *settingsData) GetProcessColor(index int) color.Attribute {
	return s.processColors[index%s.processColorsLength]
}
