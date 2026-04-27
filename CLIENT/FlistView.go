package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FriendlistView struct {
	currentcolor  int
	animetedcolor bool
	colorstep     int
	glitchmode    bool
	// List specific
	friends []string
	cursor  int
	// Universal items
	warning string
	active  bool
}

func (m FriendlistView) Init() tea.Cmd {
	// Only start the rainbow tick
	cmds := []tea.Cmd{tick()}

	// Only add fetchFriends if it's NOT already running
	if !isFriendLoopRunning {
		isFriendLoopRunning = true
		cmds = append(cmds, fetchFriends())
	}
	return tea.Batch(cmds...)

}

func NewFriendlist() FriendlistView {
	fl := FriendlistView{
		friends: []string{
			"Mikey [Online]",
			"Pirate_King [Away]",
			"Fedora_Pro [Online]",
			"Root_User [Busy]",
			"Gopher_01 [Offline]",
		},
	}
	applySharedTheme(&fl.currentcolor, &fl.animetedcolor)
	return fl
}

func (m FriendlistView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.cursor = len(m.friends) - 1
			}

		case "down", "j":
			if m.cursor < len(m.friends)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}

		case "y", "Y":
			m.glitchmode = !m.glitchmode

		}

	case tickMsg:
		m.colorstep = currentAnimationStep()
		return m, tick()

	}
	return m, nil
}

func (m FriendlistView) View() string {
	_, _, contentWidth, contentHeight := frameDimensions()
	var warningRender string
	applySharedTheme(&m.currentcolor, &m.animetedcolor)
	themeColor := currentThemeColorHex(m.colorstep)
	rows := []string{}
	visibleRows := clamp(contentHeight-4, 5, 10)
	start := 0
	if len(m.friends) > visibleRows {
		start = clamp(m.cursor-(visibleRows/2), 0, len(m.friends)-visibleRows)
	}
	end := clamp(start+visibleRows, 0, len(m.friends))

	for i, friend := range m.friends {
		if i < start || i >= end {
			continue
		}
		rows = append(rows, listRow(friend, "", m.cursor == i, contentWidth-4, themeColor))
	}

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	listContent := "No items."
	if len(rows) > 0 {
		listContent = strings.Join(rows, "\n")
	}

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		screenWordmark(themeColor, m.colorstep, m.glitchmode),
		"",
		panelTitleWithBody("Friend Manager", listContent, contentWidth, clamp(contentHeight-4, 10, 18), themeColor),
	)

	if warningRender != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", warningRender)
	}

	body = lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		"",
		footerLine(contentWidth, "esc back   q quit", "v1.02"),
	)

	return renderScreen(themeColor, body)
}
