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
	friendselecting        bool
	friendscrolling        int

	settinguseing        bool
	settingandfriendmenu int

	Requestlist []string
	Online      bool
	Offline     bool

	startTime time.Time
}

var isFriendLoopRunning = false

func (m DashboardView) Init() tea.Cmd {
	// Only start the rainbow tick
	cmds := []tea.Cmd{tick()}

	// Only add fetchFriends if it's NOT already running
	if !isFriendLoopRunning {
		isFriendLoopRunning = true
		cmds = append(cmds, fetchFriends())
	}
	return tea.Batch(cmds...)
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
			
		case "left":
			m.friendselecting = true
			m.settinguseing = false

		case "right":
			m.friendselecting = false
			m.settinguseing = true

		case "down":
			if m.friendselecting && m.friendscrolling < len(m.Friendlist)-1 && !m.settinguseing {
				m.friendscrolling++
			}
			if m.settinguseing && !m.friendselecting && m.settingandfriendmenu < 1 {
				m.settingandfriendmenu++
			}

		case "up":
			if m.friendselecting && m.friendscrolling > 0 {
				m.friendscrolling--
			}
			if m.settinguseing && m.settingandfriendmenu > 0 {
				m.settingandfriendmenu--
			}

		case "y", "Y":
			m.glitchmode = !m.glitchmode
		
		case "enter":
			if m.settinguseing {
				if m.settingandfriendmenu == 0 {
					return m , SwitchToSettings()
				}
				if m.settingandfriendmenu == 1 {
					return  m , SwitchtoFriend()
				}
			}
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


var SelectedFriend = lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor)).
	Width(25).Align(lipgloss.Center).
	Border(lipgloss.ThickBorder()).Bold(true).
	BorderTop(true).
    BorderLeft(false).
    BorderRight(false).
    BorderBottom(true).
	BorderForeground(lipgloss.Color(themeColor))


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
	SelectedFlistbox := lipgloss.NewStyle().Width(40).Align(lipgloss.Center).
		Border(lipgloss.ThickBorder()).Margin(0, 4).
		BorderForeground(lipgloss.Color(themeColor))

	
		

	title := yellotext.Render(" -- FRIENDS TO MSG -- ")
	title += "\n"
	F := ""
	F += "\n"

	for I := range m.Friendlist {

		if !m.friendselecting {
			F += Fv.Render(fmt.Sprintf("@%v",m.Friendlist[I]))
			F += "\n"
		} else if m.friendselecting {

			if I == 0 {

				F += SelectedFriend.Render(REDarrowStyle, greentext.Render(fmt.Sprintf( "@%v", m.Friendlist[m.friendscrolling])))

			} else {
				if I+m.friendscrolling > len(m.Friendlist)-1 {
					F += Fv.Render(" ")
				} else {

					F += Fv.Render(fmt.Sprintf("@%s", m.Friendlist[I+m.friendscrolling]))
				}
			}
			F += "\n"
		}

		if I > 2 {
			break
		}

	}

	var flist string

	if !m.friendselecting {

		flist = Flistbox.Render(title, F)

	} else {

		flist = SelectedFlistbox.Render(title, F)

	}

	// in here we gonna also add the list of the friend print them in there

	Settingbtn := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).Bold(true)

	Selectedsettingbtn := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.ThickBorder()).Bold(true).BorderForeground(lipgloss.Color(themeColor)).
		Foreground(lipgloss.Color("#f2ff00"))

	ManageFriend := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).Bold(true)

	SelectedManageFriend := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.ThickBorder()).Bold(true).BorderForeground(lipgloss.Color(themeColor)).
		Foreground(lipgloss.Color("#f2ff00"))

	statuse := lipgloss.NewStyle().Width(20).Align(lipgloss.Center).
		Border(lipgloss.ThickBorder()).Bold(true).Foreground(lipgloss.Color(themeColor)).
		Padding(0, 0).BorderTop(false).BorderBottom(false).BorderLeft(false).BorderRight(false).
		MarginLeft(1)

	status := statuse.Render(fmt.Sprintf("Total Friends : %v\nOnline : 4 \n\n USE : <- -> ^ v ", len(m.Friendlist)))

	settings := ""
	managefriend := ""

	if m.settinguseing && m.settingandfriendmenu == 0 {

		settings = Selectedsettingbtn.Render(REDarrowStyle, Redtext.Render(" SETTING "))
		managefriend = ManageFriend.Render(" MANAGE FRIEND ")

	} else if m.settinguseing && m.settingandfriendmenu == 1 {

		settings = Settingbtn.Render(" SETTING ")
		managefriend = SelectedManageFriend.Render(REDarrowStyle, Redtext.Render(" MANAGE FRIEND "))

	} else {
		settings = Settingbtn.Render(" SETTING ")
		managefriend = ManageFriend.Render(" MANAGE FRIEND ")
	}

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
