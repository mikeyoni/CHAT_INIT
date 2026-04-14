package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DirectMsgView struct {
	textInput     textinput.Model
	currentcolor  int
	animetedcolor bool
	colorstep     int
	currentlogoColor []string
	glitchmode    bool
	// Chat specific
	messages []string
	// Universal items
	warning string
	active  bool
}

func (m DirectMsgView) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tick())
}

func NewDirectMsg() DirectMsgView {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	dm := DirectMsgView{
		textInput: ti,
		messages:  []string{"System: Connection established."},
		currentlogoColor: []string{
			"Red", "Orange", "Yellow", "Green",
			"Cyan", "Blue", "Purple", "Pink",
		},
	}

	// LOAD SAVED SETTINGS HERE (Only once!)
	if Currentcolor != "" {
		if number, err := strconv.Atoi(Currentcolor); err == nil {
			dm.currentcolor = number
		}
	}

	if Animetedcolore == "true" {
		dm.animetedcolor = true
	}

	return dm
}

func (m DirectMsgView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.active = true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q", "Q":
			// Only quit if not typing in the input
			if !m.textInput.Focused() {
				return m, tea.Quit
			}

		case "enter":
			if m.textInput.Value() != "" {
				m.messages = append(m.messages, "You: "+m.textInput.Value())
				m.textInput.SetValue("")
			}

		case "i", "I", "tab":
			// If text input isn't focused or we want to cycle anyway
			if m.currentcolor >= 0 && m.currentcolor < len(m.currentlogoColor)-1 {
				m.currentcolor++
			} else {
				m.currentcolor = 0
			}
			savesettings(m.currentcolor, m.animetedcolor)

		case "g", "G":
			m.animetedcolor = !m.animetedcolor
			savesettings(m.currentcolor, m.animetedcolor)

		case "y", "Y":
			m.glitchmode = !m.glitchmode
		}

	case tickMsg:
		m.colorstep = (m.colorstep + 5) % 1530
		return m, tick()
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m DirectMsgView) View() string {

	width := WinSize.Width

	var render string
	var warningRender string
	var themeColor string

	if m.animetedcolor {
		r, g, b := getRainbowColor(m.colorstep * 2)
		themeColor = fmt.Sprintf("#%02x%02x%02x", r, g, b)
	} else {
		themeColor = "#7D56F4"
		if m.currentcolor >= 0 && m.currentcolor < len(m.currentlogoColor) {
			switch m.currentlogoColor[m.currentcolor] {
			case "Red":    themeColor = "#FF0000"
			case "Orange": themeColor = "#FF8800"
			case "Yellow": themeColor = "#FFFF00"
			case "Green":  themeColor = "#00FF00"
			case "Cyan":   themeColor = "#00FFFF"
			case "Blue":   themeColor = "#0000FF"
			case "Purple": themeColor = "#9D00FF"
			case "Pink":   themeColor = "#FF00FF"
			}
		}
	}

	Versions := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Right).
		Foreground(lipgloss.Color(themeColor))

	var boxrender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Width(width-4).Height(12).Padding(0, 1).Align(lipgloss.Left)

	var l string
	if m.animetedcolor {
		l = animetedmakeGradientText("MESSAGES", m.colorstep*2, m.glitchmode)
	} else {
		l = makeGradientText("MESSAGES", m.currentlogoColor, m.currentcolor)
	}

	// Chat log rendering
	var chatLog string
	for _, msg := range m.messages {
		chatLog += msg + "\n"
	}

	Footther := lipgloss.NewStyle().Width(width - 10).Bold(true).
		Foreground(lipgloss.Color("rgb(0, 0, 0)"))

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	Shortcut := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Left).
		Foreground(lipgloss.Color("#ffffff9b"))

	inputArea := lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor)).Render("\n > ") + m.textInput.View()

	centerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		l,
		boxrender.Render(chatLog),
		inputArea,
		render,
		warningRender, "\n",
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("'ESC' = Back 'Q' = Quit < 'I' & 'G' "), Versions.Render("v.1.02"))
	
	return lipgloss.Place(width, 0, lipgloss.Center, lipgloss.Center, centerContent)
}