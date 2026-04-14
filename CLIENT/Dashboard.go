package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardView struct {
	currentcolor     int
	animetedcolor    bool
	colorstep        int
	currentlogoColor []string
	glitchmode       bool
	// Universal items
	warning string
	active  bool
}

func (m DashboardView) Init() tea.Cmd {
	return tick()
}

func NewDashboard() DashboardView {
	dash := DashboardView{
		currentlogoColor: []string{
			"Red", "Orange", "Yellow", "Green",
			"Cyan", "Blue", "Purple", "Pink",
		},
	}

	// LOAD SAVED SETTINGS HERE (Only once!)
	if Currentcolor != "" {
		if number, err := strconv.Atoi(Currentcolor); err == nil {
			dash.currentcolor = number
		}
	}

	if Animetedcolore == "true" {
		dash.animetedcolor = true
	}

	return dash
}

func (m DashboardView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.active = true

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q", "Q":
			return m, tea.Quit

		case "i", "I", "tab":

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
	return m, nil
}

func (m DashboardView) View() string {

	width := WinSize.Width
    // height := WinSize.Height

	var render string
	var warningRender string
	var themeColor string

	if m.animetedcolor {
		// This pulls the EXACT color from your rainbow math
		r, g, b := getRainbowColor(m.colorstep * 2)
		themeColor = fmt.Sprintf("#%02x%02x%02x", r, g, b)

	} else {
		// Standard static color logic
		themeColor = "#7D56F4" // Default
		if m.currentcolor >= 0 && m.currentcolor < len(m.currentlogoColor) {
			switch m.currentlogoColor[m.currentcolor] {
			case "Red":
				themeColor = "#FF0000"
			case "Orange":
				themeColor = "#FF8800"
			case "Yellow":
				themeColor = "#FFFF00"
			case "Green":
				themeColor = "#00FF00"
			case "Cyan":
				themeColor = "#00FFFF"
			case "Blue":
				themeColor = "#0000FF"
			case "Purple":
				themeColor = "#9D00FF"
			case "Pink":
				themeColor = "#FF00FF"
			}
		}
	}

	// if m.currentcolor >= 0 && m.currentcolor < len(m.currentlogoColor) {
	//     switch m.currentlogoColor[m.currentcolor] {
	//     case "Red":    themeColor = "#FF0000"
	//     case "Orange": themeColor = "#FF8800"
	//     case "Yellow": themeColor = "#FFFF00"
	//     case "Green":  themeColor = "#00FF00"
	//     case "Cyan":   themeColor = "#00FFFF"
	//     case "Blue":   themeColor = "#0000FF"
	//     case "Purple": themeColor = "#9D00FF"
	//     case "Pink":   themeColor = "#FF00FF"
	//     }
	// }

	// 2. Create the dynamic selection box style
	// var selectedboxe = lipgloss.NewStyle().
	// 	Bold(true).
	// 	Foreground(lipgloss.Color(themeColor)).
	// 	BorderForeground(lipgloss.Color(themeColor)).
	// 	Border(lipgloss.RoundedBorder()).
	// 	Width(50)

	// Update the version text color too!

	Versions := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Right).
		Foreground(lipgloss.Color(themeColor))

	// var selectedboxe = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0037")).
	// 		BorderForeground(lipgloss.Color("#ff0059")).
	// 		Border(lipgloss.RoundedBorder()).Width(30).Align(lipgloss.Center)
	// // inishializing rainbow color

	var boxrender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Width(width-4).Padding(0, 0).Align(lipgloss.Center)
	v := "\n your welcome to chat init \n"

	var l string

	if m.animetedcolor {
		l = animetedmakeGradientText("HEllo", m.colorstep*2, m.glitchmode)
	}

	if !m.animetedcolor {
		l = makeGradientText("hello", m.currentlogoColor, m.currentcolor)
	}

	Footther := lipgloss.NewStyle().Width(width - 10).Bold(true).
		Foreground(lipgloss.Color("rgb(0, 0, 0)"))

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	Shortcut := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Left).
		Foreground(lipgloss.Color("#ffffff9b"))

	centerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		l, render,
		warningRender, "\n",
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("'ESC' = Back 'Q' = Quit < 'I' & 'G' " ), Versions.Render("v.1.02"))
	v = boxrender.Render(centerContent)

	return v
}
