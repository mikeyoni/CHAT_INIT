package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsView struct {
	currentcolor  int
	animetedcolor bool
	colorstep     int
	glitchmode    bool
	// Settings specific
	cursor int
	// Universal items
	warning   string
	active    bool
	startTime time.Time
}

func (m SettingsView) Init() tea.Cmd {
	return tick()
}

func NewSettings() SettingsView {
	s := SettingsView{
		startTime: time.Now(),
	}

	applySharedTheme(&s.currentcolor, &s.animetedcolor)
	return s
}

func (m SettingsView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.active = true
	applySharedTheme(&m.currentcolor, &m.animetedcolor)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			m.active = false
			return m, SwitchtoDash()

		case "ctrl+c":
			return m, tea.Quit

		case "q", "Q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = 2
			}

		case "down", "j":
			if m.cursor < 2 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case "i", "I", "tab":
			cycleThemeColor()
			applySharedTheme(&m.currentcolor, &m.animetedcolor)

		case "g", "G":
			toggleAnimatedColor()
			applySharedTheme(&m.currentcolor, &m.animetedcolor)

		case "y", "Y":
			m.glitchmode = !m.glitchmode
		}

	case tickMsg:

		// Convert to float first, then multiply, then back to int
		elapsed := time.Since(m.startTime).Milliseconds()

		// We use a larger multiplier if it's too slow, or check the math
		m.colorstep = int(float64(elapsed)*0.29) % 1530

		return m, tick()

	}
	return m, nil
}

func (m SettingsView) View() string {

	width := WinSize.Width

	var render string
	var warningRender string
	applySharedTheme(&m.currentcolor, &m.animetedcolor)
	themeColor := currentThemeColorHex(m.colorstep)

	Versions := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Right).
		Foreground(lipgloss.Color(themeColor))

	var boxrender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Width(width-4).Padding(0, 0).Align(lipgloss.Center)

	// Settings Content logic
	settingsTitle := lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor)).Bold(true).Render(" --- SETTINGS MENU --- ")

	menuItems := []string{"1. Profile Config", "2. Network Settings", "3. System UI"}
	var menuRender string
	for i, item := range menuItems {
		if m.cursor == i {
			menuRender += lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor)).Render("> "+item) + "\n"
		} else {
			menuRender += "  " + item + "\n"
		}
	}

	var l string
	l = currentThemeGradientText("SETTINGS", m.colorstep, m.glitchmode)

	Footther := lipgloss.NewStyle().Width(width - 10).Bold(true).
		Foreground(lipgloss.Color("rgb(0, 0, 0)"))

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	Shortcut := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Left).
		Foreground(lipgloss.Color("#ffffff9b"))

	centerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		l,
		settingsTitle,
		"\n",
		menuRender,
		render,
		warningRender, "\n",
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("'ESC' = Back 'Q' = Quit < 'I' & 'G' "), Versions.Render("v.1.02"))
	v := boxrender.Render(centerContent)

	return v
}
