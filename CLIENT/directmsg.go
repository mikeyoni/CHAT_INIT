package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DirectMsgView struct {
	FriendtoMsg   string
	textInput     textinput.Model
	currentcolor  int
	animetedcolor bool
	colorstep     int
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
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	dm := DirectMsgView{
		textInput: ti,
		messages:  []string{"System: Connection established."},
	}
	applySharedTheme(&dm.currentcolor, &dm.animetedcolor)
	return dm
}

func (m DirectMsgView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.active = true
	applySharedTheme(&m.currentcolor, &m.animetedcolor)

	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.active = false
			return m, SwitchtoDash()
		case "i", "I", "tab":
			cycleThemeColor()
			applySharedTheme(&m.currentcolor, &m.animetedcolor)
		case "g", "G":
			toggleAnimatedColor()
			applySharedTheme(&m.currentcolor, &m.animetedcolor)
		case "y", "Y":
			m.glitchmode = !m.glitchmode
		case "enter":
			// Your enter logic
		}

	case tickMsg:
		m.colorstep = (m.colorstep + 5) % 1530
		return m, tea.Batch(tick(), cmd)
	}

	return m, cmd
}

func (m DirectMsgView) View() string {

	width := WinSize.Width
	// height := WinSize.Height

	// var render string
	var warningRender string
	applySharedTheme(&m.currentcolor, &m.animetedcolor)
	themeColor := currentThemeColorHex(m.colorstep)

	// var selectedboxe = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0037")).
	// 		BorderForeground(lipgloss.Color("#ff0059")).
	// 		Border(lipgloss.RoundedBorder()).Width(30).Align(lipgloss.Center)
	// // inishializing rainbow color

	var boxrender = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Width(width - 4).
		Height(WinSize.Height - 4).
		Padding(0).
		Align(lipgloss.Center). // Keeps content centered horizontally (Left to Right)
		AlignVertical(lipgloss.Top).
		BorderTop(false).BorderBottom(false)

	// v := "\n your welcome to chat init \n"

	// 	var l string
	// 	l += "\n"
	// 	if m.animetedcolor {
	// 		l += animetedmakeGradientText(` ▄▀▀▀ █  █ █▀▀█ ▀█▀   ▀█▀ █▀▀▄ ▀█▀ ▀█▀
	//  █    █▀▀█ █▄▄█  █     █  █  █  █   █
	//  ▀▀▀ ▀  ▀ ▀  ▀  ▀    ▄█▄ ▀  ▀ ▄█▄  ▀`, m.colorstep*2, m.glitchmode)
	// 	}

	// 	if !m.animetedcolor {
	// 		l += makeGradientText(` ▄▀▀▀ █  █ █▀▀█ ▀█▀   ▀█▀ █▀▀▄ ▀█▀ ▀█▀
	//  █    █▀▀█ █▄▄█  █     █  █  █  █   █
	//  ▀▀▀ ▀  ▀ ▀  ▀  ▀    ▄█▄ ▀  ▀ ▄█▄  ▀`, m.currentlogoColor, m.currentcolor)
	// 	}

	titlebar := lipgloss.NewStyle().Background(lipgloss.Color(themeColor)).Align(lipgloss.Left).
		Width(width-4).Bold(true).Padding(0, 0).BorderBackground(lipgloss.Color(themeColor)).
		Foreground(lipgloss.Color("#00000000"))
	// usernameshotleft := lipgloss.NewStyle().Align(lipgloss.Left)
	// tomsgshow := lipgloss.NewStyle().Align(lipgloss.Right)

	// Usernameshow := usernameshotleft.Render(fmt.Sprintf("MSG TO -> @%v", m.FriendtoMsg))

	leftSide := fmt.Sprintf("  MSG TO -> @%v  ", m.FriendtoMsg)
	rightSide := fmt.Sprintf("LOGEDIN AS : @%v   ", myuser)

	totalWidth := width - 4
	occupiedWidth := lipgloss.Width(leftSide) + lipgloss.Width(rightSide)
	spacerWidth := totalWidth - occupiedWidth

	if spacerWidth < 0 {
		spacerWidth = 0
	}

	spacer := strings.Repeat(" ", spacerWidth)

	// 4. Join them and Render inside the colored bar
	title := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, spacer, rightSide)
	chastboxrender := lipgloss.NewStyle().BorderForeground(lipgloss.Color(themeColor)).Border(lipgloss.RoundedBorder()).
		Width(width - 6).Foreground(lipgloss.Color(themeColor))

	Chatbox := chastboxrender.Render(" ", REDarrowStyle, m.textInput.View())

	headerBar := titlebar.Render(title)

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	// Get actual heights
	// hH := lipgloss.Height(headerBar)
	// cH := lipgloss.Height(Chatbox)

	// Subtract an extra 2 lines just to be safe from terminal rounding errors
	spacerHeight := WinSize.Height - 6 - lipgloss.Height(headerBar) - lipgloss.Height(Chatbox)

	if spacerHeight < 0 {
		spacerHeight = 0
	}

	spacere := strings.Repeat("\n", spacerHeight)

	// Ensure the order is correct
	centerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBar,
		spacere, // Pushes everything below it to the bottom
		warningRender,
		Chatbox,
	)

	footer := lipgloss.NewStyle().
		Width(width - 6).
		Foreground(lipgloss.Color("#ffffff9b")).
		Render("'ESC' = Back  'I' = Next Color  'G' = Toggle Animation")
	centerContent = lipgloss.JoinVertical(lipgloss.Left, centerContent, footer)

	// Remove .Height() from boxrender temporarily to see if it fixes it
	// If it works without .Height(), your calculation was clipping the input
	return boxrender.Height(WinSize.Height - 4).Render(centerContent)

}
