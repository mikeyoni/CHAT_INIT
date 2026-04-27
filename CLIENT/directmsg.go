package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type DirectMsgView struct {
	FriendtoMsg   string
	textInput     textinput.Model
	currentcolor  int
	animetedcolor bool
	colorstep     int
	glitchmode    bool
	conn          *websocket.Conn
	readerStarted bool
	// Chat specific
	my_msg       []string
	Incoming_msg []string
	messages     []string
	// Universal items
	warning string
	active  bool
}

type directMsgConnectedMsg struct {
	conn *websocket.Conn
}

type directMsgReceivedMsg struct {
	text string
}

type directMsgErrorMsg struct {
	text string
}

type directMsgConnectionClosedMsg struct {
	text string
}

var directMessageEvents = make(chan tea.Msg, 32)

func (m DirectMsgView) Init() tea.Cmd {
	cmds := []tea.Cmd{textinput.Blink, tick()}

	if m.conn == nil && m.FriendtoMsg != "" {
		cmds = append(cmds, connectDirectSocket())
	}

	if m.readerStarted {
		cmds = append(cmds, waitForDirectEvent())
	}

	return tea.Batch(cmds...)
}

func NewDirectMsg(friend string) DirectMsgView {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.Prompt = ""
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = WinSize.Width - 9

	dm := DirectMsgView{
		FriendtoMsg: friend,
		textInput:   ti,
		messages: []string{
			"System: Connection established.",
		},
		my_msg:       []string{},
		Incoming_msg: []string{},
	}
	applySharedTheme(&dm.currentcolor, &dm.animetedcolor)
	return dm
}

func connectDirectSocket() tea.Cmd {
	return func() tea.Msg {
		if myuser == "" || mytoken == "" {
			return directMsgErrorMsg{text: "Missing user token."}
		}

		url := fmt.Sprintf("ws://localhost:4040/chat?user=%s&token=%s", myuser, mytoken)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			return directMsgErrorMsg{text: fmt.Sprintf("Failed to connect chat socket: %v", err)}
		}

		return directMsgConnectedMsg{conn: conn}
	}
}

func startDirectReader(conn *websocket.Conn) {
	go func() {
		for {
			_, payload, err := conn.ReadMessage()
			if err != nil {
				directMessageEvents <- directMsgConnectionClosedMsg{text: fmt.Sprintf("Chat socket closed: %v", err)}
				return
			}

			directMessageEvents <- directMsgReceivedMsg{text: string(payload)}
		}
	}()
}

func waitForDirectEvent() tea.Cmd {
	return func() tea.Msg {
		return <-directMessageEvents
	}
}

func senderFromWireMessage(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "[") {
		return ""
	}

	endIdx := strings.Index(trimmed, "]")
	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(trimmed[1:endIdx])
}

func (m DirectMsgView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.active = true
	applySharedTheme(&m.currentcolor, &m.animetedcolor)

	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case directMsgConnectedMsg:
		m.conn = msg.conn
		m.warning = ""
		if !m.readerStarted {
			startDirectReader(m.conn)
			m.readerStarted = true
		}
		return m, waitForDirectEvent()

	case directMsgReceivedMsg:
		sender := senderFromWireMessage(msg.text)
		if sender == "" || sender == m.FriendtoMsg {
			m.Incoming_msg = append(m.Incoming_msg, msg.text)
			m.messages = append(m.messages, msg.text)
		} else {
			m.warning = fmt.Sprintf("New message from @%s", sender)
		}

		return m, waitForDirectEvent()

	case directMsgErrorMsg:
		m.warning = msg.text
		return m, nil

	case directMsgConnectionClosedMsg:
		m.warning = msg.text
		m.conn = nil
		m.readerStarted = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.active = false
			return m, SwitchtoDash()
		case "ctrl+i", "ctrl+I", "tab":
			cycleThemeColor()
			applySharedTheme(&m.currentcolor, &m.animetedcolor)
		case "ctrl+3", "ctrl+G":
			toggleAnimatedColor()
			applySharedTheme(&m.currentcolor, &m.animetedcolor)
		case "ctrl+y", "ctrl+Y":
			m.glitchmode = !m.glitchmode
		case "enter":
			text := strings.TrimSpace(m.textInput.Value())
			if text == "" {
				return m, cmd
			}

			if m.FriendtoMsg == "" {
				m.warning = "No friend selected."
				return m, cmd
			}

			if m.conn == nil {
				m.warning = "Chat socket is offline."
				return m, cmd
			}

			wireMsg := fmt.Sprintf("tusr:%s:user:%s", m.FriendtoMsg, text)
			if err := m.conn.WriteMessage(websocket.TextMessage, []byte(wireMsg)); err != nil {
				m.warning = fmt.Sprintf("Failed to send message: %v", err)
				m.conn = nil
				m.readerStarted = false
				return m, cmd
			}

			rendered := fmt.Sprintf("[ me ] : %s", text)
			m.my_msg = append(m.my_msg, rendered)
			m.messages = append(m.messages, rendered)
			m.textInput.SetValue("")
			m.warning = ""
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

	messagesBox := lipgloss.NewStyle().
		Width(width-6).
		Height(WinSize.Height-10).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Padding(1, 2)

	var chatBody string
	if len(m.messages) == 0 {
		chatBody = "No messages yet."
	} else {
		start := 0
		maxVisible := WinSize.Height - 14
		if maxVisible < 1 {
			maxVisible = 1
		}
		if len(m.messages) > maxVisible {
			start = len(m.messages) - maxVisible
		}
		chatBody = strings.Join(m.messages[start:], "\n")
	}

	messageHistory := messagesBox.Render(chatBody)

	// Get actual heights
	// hH := lipgloss.Height(headerBar)
	// cH := lipgloss.Height(Chatbox)

	// Subtract an extra 2 lines just to be safe from terminal rounding errors
	spacerHeight := WinSize.Height - 8 - lipgloss.Height(headerBar) - lipgloss.Height(Chatbox) - lipgloss.Height(messageHistory)

	if spacerHeight < 0 {
		spacerHeight = 0
	}

	spacere := strings.Repeat("\n", spacerHeight)

	// Ensure the order is correct
	centerContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerBar,
		messageHistory,
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
