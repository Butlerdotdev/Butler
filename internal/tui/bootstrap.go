// Package tui provides an interactive way to bootstrap the Butler management cluster.
//
// Copyright (c) 2025, The Butler Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
)

func landingView(m model) string {
	var menu = "\nThis interactive bootstrap will walk you through selecting the configurations required to bootstrap Butler\n\n"
	return menuHelper(menu, m)
}

func processingView(p progress.Model) string {
	return sectionStyle.Render("\n🔄 Processing...\n\n" + p.View())
}

func nutanixLoginView(m model) string {
	menu := fmt.Sprintf(
		`Enter Nutanix Credentials To Bootstrap Butler:

		%s
		%s

		%s
		%s

		%s
		%s

	`,
		commandStyle.Width(64).Render("Endpoint"),
		m.inputs[0].View(),
		commandStyle.Width(64).Render("Username"),
		m.inputs[1].View(),
		commandStyle.Width(64).Render("Password"),
		m.inputs[2].View(),
	) + "\n"

	return menuHelper(menu, m)
}

func nutanixClusterSelectView(m model) string {
	return menuHelper("Select Nutanix Cluster To Bootstrap Butler:\n", m)
}

func nutanixSubnetSelectView(m model) string {
	return menuHelper("Select Nutanix Subnet To Bootstrap Butler:\n", m)
}

func bootstrapGateView(m model) string {
	return menuHelper("Selecting continue will begin the bootstrap process with your selected inputs. Note that interactive mode will exit and you will be shown the standard output of the bootstrapping process. If you do not wish to start the bootstrapping process, or would like to change your inputs, select Exit.\n", m)
}

func menuHelper(menu string, m model) string {
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
