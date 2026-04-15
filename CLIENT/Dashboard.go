package main

import (
	"fmt"
	"strconv"
	"time"

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

	// friendlist
	Friendlist             []string
	Nfriendlist            int
	IssillectedFriendtomsg bool

	Requestlist []string
	Online      bool
	Offline     bool

	startTime time.Time
}

func (m DashboardView) Init() tea.Cmd {
	return tea.Batch(tick(), fetchFriends())
}

type friendsLoadedMsg []string

func fetchFriends() tea.Cmd {
	return func() tea.Msg {
		list := viewflist(baseURL)
		return friendsLoadedMsg(list)
	}
}

func NewDashboard() DashboardView {
	dash := DashboardView{
		startTime: time.Now(),
		currentlogoColor: []string{
			"Red", "Orange", "Yellow", "Green",
			"Cyan", "Blue", "Purple", "Pink",
		},
	}
	// this loade the friend list

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

	case friendsLoadedMsg:
		m.Friendlist = msg

		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchFriends()()
		})

	case tickMsg:
		// Convert to float first, then multiply, then back to int
		elapsed := time.Since(m.startTime).Milliseconds()

		// We use a larger multiplier if it's too slow, or check the math
		m.colorstep = int(float64(elapsed)*0.29) % 1530

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
	l += "\n"
	if m.animetedcolor {
		l += animetedmakeGradientText(` ▄▀▀▀ █  █ █▀▀█ ▀█▀   ▀█▀ █▀▀▄ ▀█▀ ▀█▀
 █    █▀▀█ █▄▄█  █     █  █  █  █   █ 
 ▀▀▀ ▀  ▀ ▀  ▀  ▀    ▄█▄ ▀  ▀ ▄█▄  ▀`, m.colorstep*2, m.glitchmode)
	}

	if !m.animetedcolor {
		l += makeGradientText(` ▄▀▀▀ █  █ █▀▀█ ▀█▀   ▀█▀ █▀▀▄ ▀█▀ ▀█▀
 █    █▀▀█ █▄▄█  █     █  █  █  █   █ 
 ▀▀▀ ▀  ▀ ▀  ▀  ▀    ▄█▄ ▀  ▀ ▄█▄  ▀`, m.currentlogoColor, m.currentcolor)
	}

	Footther := lipgloss.NewStyle().Width(width - 10).Bold(true).
		Foreground(lipgloss.Color("rgb(0, 0, 0)"))

	titlebar := lipgloss.NewStyle().Background(lipgloss.Color(themeColor)).Align(lipgloss.Left).
		Width(width - 4).Bold(true).Foreground(lipgloss.Color("#00000000")).Render(fmt.Sprintf(" LOGEDIN AS : @%v", myuser))

	render += "\n"
	render += titlebar
	render += "\n"

	Dashe := lipgloss.NewStyle().Width(width-4).Padding(1, 1).Align(lipgloss.Center)
	Fv := lipgloss.NewStyle().Margin(0).Align(lipgloss.Center).PaddingBottom(1)
	Flistbox := lipgloss.NewStyle().Width(40).Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).Margin(0, 4)
	title := yellotext.Render("	FRIENDS	")
	title += "\n"
	F := ""
	F += "\n"
	for I, _ := range m.Friendlist {

		F += Fv.Render(m.Friendlist[I])
		F += "\n"

		if I > 2 {
			break
		}

	}

	flist := Flistbox.Render(title, F)
	// in here we gonna also add the list of the friend print them in there

	Settingbtn := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).Bold(true)

	ManageFriend := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).Bold(true)

	statuse := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.ThickBorder()).Bold(true).Foreground(lipgloss.Color(themeColor)).Blink(true).
		BorderForeground(lipgloss.Color(themeColor)).Padding(0, 0)

	status := statuse.Render(fmt.Sprintf("Total Friends : 10\nOnline : 4"))
	settings := Settingbtn.Render(" SETTING ")
	managefriend := ManageFriend.Render(" MANAGE FRIEND ")

	button := lipgloss.JoinVertical(lipgloss.Left, settings, managefriend, "", status)

	Dash := Dashe.Render(lipgloss.JoinHorizontal(lipgloss.Top, flist, button))

	render += Dash

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	Shortcut := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Left).
		Foreground(lipgloss.Color("#ffffff9b"))

	centerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		l, render,
		warningRender,
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("Quit = 'Q' < 'I' & 'G'"), Versions.Render("v.1.02"))
	v = boxrender.Render(centerContent)

	return v
}
