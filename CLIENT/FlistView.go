package main

import (
	"time"

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
	warning   string
	active    bool
	startTime time.Time
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
		startTime: time.Now(),
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

func (m FriendlistView) View() string {

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

	var l string
	l = currentThemeGradientText("FRIENDS", m.colorstep, m.glitchmode)

	// Friend list rendering
	var listContent string
	for i, friend := range m.friends {
		if m.cursor == i {
			listContent += lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor)).Bold(true).Render("> "+friend) + "\n"
		} else {
			listContent += "  " + friend + "\n"
		}
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
		l,
		"\n",
		listContent,
		render,
		warningRender, "\n",
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("'ESC' = Back 'Q' = Quit < 'I' & 'G' "), Versions.Render("v.1.02"))
	v := boxrender.Render(centerContent)

	return v
}
