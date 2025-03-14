package tui

import (
	"butler/internal/adapters/providers/nutanix"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Model definition
type model struct {
	cursor      int
	options     []string
	currentView string
	rootCmd     *cobra.Command
	message     string
	inputs      []textinput.Model
	inputsCount int
}

// Messages for async operations
type tickMsg struct{}
type commandCompleteMsg struct {
	err    error
	output string
}

var logger *zap.Logger

// UI Styles
var (
	outerBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("#50FA7B")).
				Padding(1, 2).
				Width(120).
				Align(lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Background(lipgloss.Color("#282A36")).
			Padding(1, 2).
			Align(lipgloss.Center).
			Width(100).
			Bold(true)

	sectionStyle = lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#8BE9FD")).
			Padding(1, 3).
			Width(110).
			MarginTop(1).
			Align(lipgloss.Center)

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Bold(true)

	descriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Bold(false)

	footerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#6272A4")).
			Background(lipgloss.Color("#44475A")).
			Align(lipgloss.Center).
			Padding(1, 2).
			Width(110)
)

func renderHeader() string {
	// Title (Gold) - No Background
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")). // Gold (Readable in dark & light mode)
		Bold(true).
		Width(100).
		Align(lipgloss.Center).
		Render("🚀 Welcome to Butler")

	// Subtitle (Green)
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50FA7B")). // Green (Good contrast)
		Width(100).
		Align(lipgloss.Center).
		Render("🤖 \"Efficiency through automation.\"")

	// Decorative Separator (Cyan)
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8BE9FD")). // Cyan (Matches theme)
		Width(100).
		Align(lipgloss.Center).
		Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Border Style (NO Background!)
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#8BE9FD")). // Cyan
		Padding(1, 2).
		Width(110) // Matches overall layout

	// Assemble the header
	return borderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, separator),
	)
}

func renderFooter() string {
	// Define colors for key bindings & descriptions
	keyColor := lipgloss.Color("#FF79C6")  // Magenta (Good contrast in dark & light mode)
	descColor := lipgloss.Color("#8BE9FD") // Cyan (Good contrast in dark & light mode)

	// Styles for key bindings & descriptions (no background)
	keyStyle := lipgloss.NewStyle().Foreground(keyColor).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(descColor).Bold(false)

	// Properly spaced & aligned footer
	footerContent := lipgloss.JoinHorizontal(lipgloss.Center,
		keyStyle.Render("↑")+" "+descStyle.Render("Move Up")+"   ",
		keyStyle.Render("↓")+" "+descStyle.Render("Move Down")+"   ",
		keyStyle.Render("↵/enter")+" "+descStyle.Render("Select")+"   ",
		keyStyle.Render("ctrl+c")+" "+descStyle.Render("Quit"),
	)

	// Remove background to prevent UI corruption
	return footerStyle.
		UnsetBackground(). // Ensure no background box issues
		Width(110).
		Render(footerContent)
}

func initialModel(rootCmd *cobra.Command) model {
	availableCommands := []string{}
	availableCommands = append(availableCommands, "📟   Continue")
	availableCommands = append(availableCommands, "🚪   Exit")
	inputs := make([]textinput.Model, 3)

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 64
		t.Width = 64

		switch i {
		case 0:
			t.Placeholder = "https://nutanix.example.com"
		case 1:
			t.Placeholder = "Username"
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}

	return model{
		options:     availableCommands,
		currentView: "main",
		rootCmd:     rootCmd,
		inputs:      inputs,
		inputsCount: 0,
	}
}

// Init initializes the TUI program
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Run Cobra Command
func runCobraCommand(rootCmd *cobra.Command, args ...string) tea.Cmd {
	return func() tea.Msg {
		var stdout, stderr bytes.Buffer

		cmd, _, err := rootCmd.Find(args)
		if err != nil {
			return commandCompleteMsg{err: fmt.Errorf("command not found: %s", args)}
		}

		cmd.SetOut(&stdout)
		cmd.SetErr(&stderr)
		rootCmd.SetArgs(args)

		time.Sleep(2 * time.Second)
		_, err = cmd.ExecuteC()
		output := stdout.String()

		_ = stderr.String()
		return commandCompleteMsg{err: err, output: output}
	}
}

// Tick function
func tick() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

// StartTUI
func StartTUI(rootCmd *cobra.Command, log *zap.Logger) {
	logger = log
	p := tea.NewProgram(initialModel(rootCmd), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		logger.Fatal("Failed to start TUI", zap.Error(err))
		os.Exit(1)
	}
}

// View renders the TUI
func (m model) View() string {
	header := renderHeader()
	footer := renderFooter()

	var body string
	switch m.currentView {
	case "main":
		body = renderMenu(m.options, m.cursor)
	case "nutanix_login":
		body = nutanixLogin(m)
	}

	return outerBorderStyle.Render(
		lipgloss.JoinVertical(lipgloss.Top, header, body, descriptionStyle.Render(m.message), footer),
	)
}

// Update handles key inputs and updates the model state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case commandCompleteMsg:
		if msg.err != nil {
			m.message = fmt.Sprintf("❌ Error: %s", msg.err)
		} else {
			m.message = fmt.Sprintf("✅ Command Completed Successfully!\n\n%s", msg.output)
		}
		return m, nil

	case tickMsg:
		if m.currentView == "nutanix_login" {
			nutanixClient := nutanix.NewNutanixClient(nil, m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value(), logger)
			nutanixAdapter := nutanix.NewNutanixAdapter(nutanixClient, logger)
			uuids, err := nutanixAdapter.GetClusterUuids()
			return m, func() tea.Msg { return commandCompleteMsg{err: err, output: fmt.Sprintf("%v", uuids)} }
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < m.inputsCount+len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			switch m.currentView {
			case "main":
				switch m.options[m.cursor] {
				case "📟   Continue":
					m.currentView = "nutanix_login"
					m.options = []string{"📟   Next"}
					m.cursor = 0
					m.inputsCount = 3
					m.inputs[m.cursor].Focus()
				case "🚪   Exit":
					return m, tea.Quit
				}
			case "nutanix_login":
				if m.cursor >= len(m.inputs) {
					m.options = append(m.options, "🔄   Processing...")
					m.cursor = 0
					m.inputsCount = 0
					return m, tick()
				}
			}
		}
	}

	// Handle text input
	if m.currentView == "nutanix_login" {
		// update cursor to match the current input or option
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		if m.cursor < len(m.inputs) {
			m.inputs[m.cursor].Focus()
		}
		var cmd tea.Cmd
		m.inputs[0], cmd = m.inputs[0].Update(msg)
		m.inputs[1], cmd = m.inputs[1].Update(msg)
		m.inputs[2], cmd = m.inputs[2].Update(msg)
		return m, cmd
	}

	return m, nil
}
