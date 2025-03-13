package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
)

// Helper to render menu options
func renderMenu(options []string, cursor int) string {
	var menu = "\nThis interactive bootstrap will walk you through selecting the configurations required to bootstrap Butler\n\n"
	maxWidth := 46 // Set max text width (excluding emoji)

	for i, option := range options {
		cursorMarker := "  " // Default for unselected

		// Ensure consistent width & padding
		optionStyle := descriptionStyle.Width(maxWidth)

		if cursor == i {
			cursorMarker = cursorStyle.Render(">")
			optionStyle = selectedStyle.Width(maxWidth) // Highlight selection
		}

		// Apply padding and ensure it renders correctly
		menu += fmt.Sprintf("%s %s\n", cursorMarker, optionStyle.Render(option))
	}

	// Explicitly set section width to prevent shifting
	return sectionStyle.Width(110).Render(menu)
}

// Processing View
func processingView(p progress.Model) string {
	return sectionStyle.Render("\n🔄 Processing...\n\n" + p.View())
}

func nutanixLogin(m model) string {
	menu := fmt.Sprintf(
		`Enter Nutanix Credentials To Bootstrap Butler:

		%s
		%s

		%s
		%s

	`,
		commandStyle.Width(64).Render("Username"),
		m.inputs[0].View(),
		commandStyle.Width(64).Render("Password"),
		m.inputs[1].View(),
	) + "\n"

	maxWidth := 64 // Set max text width (excluding emoji)

	for i, option := range m.options {
		cursorMarker := "  " // Default for unselected

		// Ensure consistent width & padding
		optionStyle := descriptionStyle.Width(maxWidth)

		if m.cursor == i+m.inputsCount {
			cursorMarker = cursorStyle.Render(">")
			optionStyle = selectedStyle.Width(maxWidth) // Highlight selection
		}

		// Apply padding and ensure it renders correctly
		menu += fmt.Sprintf("%s %s\n", cursorMarker, optionStyle.Render(option))
	}

	// Explicitly set section width to prevent shifting
	return sectionStyle.Width(110).Render(menu)
}
