package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardView struct {
	currentcolor  int
	animetedcolor bool
	colorstep     int
	glitchmode    bool
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

	Requestlist   []string
	OnlineFriends []string
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

type friendsLoadedMsg struct {
	friends []string
	online  []string
}

func fetchFriends() tea.Cmd {
	return func() tea.Msg {
		return friendsLoadedMsg{
			friends: viewflist(baseURL),
			online:  viewOnlineFriends(baseURL),
		}
	}
}

func NewDashboard() DashboardView {
	dash := DashboardView{}
	applySharedTheme(&dash.currentcolor, &dash.animetedcolor)
	return dash
}

func containsFriend(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}

	return false
}

func renderFriendStatus(friend string, onlineFriends []string) string {
	if containsFriend(onlineFriends, friend) {
		return fmt.Sprintf("@%s [online]", friend)
	}

	return fmt.Sprintf("@%s [offline]", friend)
}

func (m DashboardView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.active = true
	applySharedTheme(&m.currentcolor, &m.animetedcolor)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "q", "Q":
			return m, tea.Quit

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
			if m.friendselecting && len(m.Friendlist) > 0 {
				selectedFriend := m.Friendlist[m.friendscrolling]
				if containsFriend(m.OnlineFriends, selectedFriend) {
					m.warning = ""
					return m, SwitchtoDirectMsg(selectedFriend)
				}

				m.warning = fmt.Sprintf("@%s is offline right now.", selectedFriend)
			}

			if m.settinguseing {
				if m.settingandfriendmenu == 0 {
					return m, SwitchToSettings()
				}
				if m.settingandfriendmenu == 1 {
					return m, SwitchtoFriend()
				}
			}
		}

	case friendsLoadedMsg:
		m.Friendlist = msg.friends
		m.OnlineFriends = msg.online
		if m.friendscrolling >= len(m.Friendlist) && len(m.Friendlist) > 0 {
			m.friendscrolling = len(m.Friendlist) - 1
		}
		if len(m.Friendlist) == 0 {
			m.friendscrolling = 0
		}

		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
			return fetchFriends()()
		})

	case tickMsg:
		m.colorstep = currentAnimationStep()
		return m, tick()

	}
	return m, nil
}

func (m DashboardView) View() string {
	_, _, contentWidth, contentHeight := frameDimensions()
	var warningRender string
	applySharedTheme(&m.currentcolor, &m.animetedcolor)
	themeColor := currentThemeColorHex(m.colorstep)
	headerWidth := contentWidth
	friendPanelWidth := clamp((contentWidth*2)/3, 28, 52)
	sidebarWidth := contentWidth - friendPanelWidth - 3
	stacked := sidebarWidth < 22
	if stacked {
		friendPanelWidth = contentWidth
		sidebarWidth = contentWidth
	}
	friendPanelHeight := clamp(contentHeight-10, 10, 18)
	sidebarHeight := friendPanelHeight

	wordmark := dashboardHero(themeColor, m.colorstep, m.glitchmode)
	header := compactHeader("@"+myuser, badge(fmt.Sprintf("online %d", len(m.OnlineFriends)), true), headerWidth)

	friendRows := []string{}
	start := 0
	if len(m.Friendlist) > 0 {
		start = clamp(m.friendscrolling, 0, len(m.Friendlist)-1)
	}
	visibleRows := clamp(friendPanelHeight-4, 4, 10)
	limit := start + visibleRows
	if limit > len(m.Friendlist) {
		limit = len(m.Friendlist)
	}

	for index := start; index < limit; index++ {
		friend := m.Friendlist[index]
		isSelected := m.friendselecting && index == m.friendscrolling
		statusOn := containsFriend(m.OnlineFriends, friend)

		rowRight := badge("offline", false)
		if statusOn {
			rowRight = badge("online", true)
		}

		friendRows = append(friendRows, listRow("@"+friend, rowRight, isSelected, friendPanelWidth-6, themeColor))
	}

	friendBody := "No friends found."
	if len(friendRows) > 0 {
		friendBody = strings.Join(friendRows, "\n")
	}

	friendPanel := panelTitleWithBody("Friends Ready To Chat", friendBody, friendPanelWidth, friendPanelHeight, themeColor)

	settingsText := strings.Join([]string{
		menuButton("Settings", m.settinguseing && m.settingandfriendmenu == 0, sidebarWidth-6, themeColor),
		"",
		menuButton("Manage Friends", m.settinguseing && m.settingandfriendmenu == 1, sidebarWidth-6, themeColor),
		"",
		statusText(fmt.Sprintf("total friends  %d\nonline now    %d\n\nleft/right switch\nup/down move\nenter open", len(m.Friendlist), len(m.OnlineFriends))),
	}, "\n")

	sidebar := panelTitleWithBody("Actions", settingsText, sidebarWidth, sidebarHeight, themeColor)

	content := friendPanel
	if stacked {
		content = lipgloss.JoinVertical(lipgloss.Left, friendPanel, "", sidebar)
	} else {
		content = lipgloss.JoinHorizontal(lipgloss.Top, friendPanel, "   ", sidebar)
	}

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		wordmark,
		"",
		header,
		"",
		content,
	)
	if warningRender != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", warningRender)
	}

	body = lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		"",
		footerLine(contentWidth, "enter open chat   q quit", "v1.02"),
	)

	return renderScreen(themeColor, body)
}
