package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsView struct {
	currentcolor  int
	animetedcolor bool
	colorstep     int
	glitchmode    bool
	returnToLogin bool
	// Settings specific
	cursor int
	// Universal items
	warning string
	active  bool
}

func (m SettingsView) Init() tea.Cmd {
	return tick()
}

func NewSettings() SettingsView {
	s := SettingsView{}

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
			if m.returnToLogin {
				return m, SwitchtoLogin()
			}
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

		case "y", "Y":
			m.glitchmode = !m.glitchmode

		case "enter":
			switch m.cursor {
			case 0:
				cycleThemeColor()
				applySharedTheme(&m.currentcolor, &m.animetedcolor)
			case 1:
				toggleAnimatedColor()
				applySharedTheme(&m.currentcolor, &m.animetedcolor)
			case 2:
				if m.returnToLogin {
					return m, SwitchtoLogin()
				}
				return m, SwitchtoDash()
			}
		}

	case tickMsg:
		m.colorstep = currentAnimationStep()
		return m, tick()

	}
	return m, nil
}

func (m SettingsView) View() string {
	_, _, contentWidth, contentHeight := frameDimensions()
	var warningRender string
	applySharedTheme(&m.currentcolor, &m.animetedcolor)
	themeColor := currentThemeColorHex(m.colorstep)
	menuItems := []string{
		"Change Theme Color",
		"Toggle Animated Color",
		"Back",
	}
	rows := []string{}
	for i, item := range menuItems {
		rows = append(rows, menuButton(item, m.cursor == i, clamp(contentWidth-4, 20, 40), themeColor))
	}

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	bodyText := strings.Join(rows, "\n\n") + "\n\n" + statusText("use enter to apply the selected option")
	panelHeight := clamp(contentHeight-4, 10, 16)

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		screenWordmark(themeColor, m.colorstep, m.glitchmode),
		"",
		panelTitleWithBody("Settings", bodyText, contentWidth, panelHeight, themeColor),
	)

	if warningRender != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", warningRender)
	}

	body = lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		"",
		footerLine(contentWidth, "enter apply   esc back", "v1.02"),
	)

	return renderScreen(themeColor, body)
}
