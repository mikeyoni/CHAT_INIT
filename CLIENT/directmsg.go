package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DirectMsgView struct {
	FriendtoMsg      string
	textInput        textinput.Model
	currentcolor     int
	animetedcolor    bool
	colorstep        int
	currentlogoColor []string
	glitchmode       bool
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

		case "enter":

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
	// height := WinSize.Height

	// var render string
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
		

	v := "\n your welcome to chat init \n"

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
	Width(width-6).Foreground(lipgloss.Color(themeColor))

	Chatbox := chastboxrender.Render( " " , REDarrowStyle , m.textInput.View() )

	headerBar := titlebar.Render(title)

	spacerHeight := WinSize.Height - 4 - 1 - 2
	if spacerHeight < 0 {
		spacerHeight = 0
	}

	spacere := strings.Repeat("\n", spacerHeight)

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	// Join them without extra spaces
	centerContent := lipgloss.JoinVertical(
		lipgloss.Left, // Changed to Left to ensure it hugs the edge
		headerBar,
		spacere,
		warningRender,
	)

	render := lipgloss.JoinVertical(lipgloss.Left, centerContent, Chatbox)
	// Ensure boxrender has absolutely no top padding
	v = boxrender.PaddingTop(0).Render(render)

	return v
}
