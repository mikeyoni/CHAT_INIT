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

		url := fmt.Sprintf("%s?user=%s&token=%s", websocketURL(), myuser, mytoken)
		conn, _, err := wsDialer.Dial(url, nil)
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
		m.colorstep = currentAnimationStep()
		return m, tea.Batch(tick(), cmd)
	}

	return m, cmd
}

func (m DirectMsgView) View() string {
	_, _, contentWidth, contentHeight := frameDimensions()
	var warningRender string
	applySharedTheme(&m.currentcolor, &m.animetedcolor)
	themeColor := currentThemeColorHex(m.colorstep)

	if m.warning != "" {
		warningRender = warnStyle.Render(m.warning)
	}

	header := compactHeader("@"+m.FriendtoMsg, "logged in as @"+myuser, contentWidth)

	messagesBox := lipgloss.NewStyle().
		Width(contentWidth).
		Height(clamp(contentHeight-8, 8, 18)).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#2b3240")).
		Background(lipgloss.Color("#151922")).
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

	inputBox := lipgloss.NewStyle().
		Width(contentWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Background(lipgloss.Color("#151922")).
		Padding(0, 1).
		Render(m.textInput.View())

	body := lipgloss.JoinVertical(
		lipgloss.Left,
		screenWordmark(themeColor, m.colorstep, m.glitchmode),
		"",
		sectionHeader("Direct Message", themeColor),
		header,
		"",
		messageHistory,
	)

	if warningRender != "" {
		body = lipgloss.JoinVertical(lipgloss.Left, body, "", warningRender)
	}

	body = lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		"",
		inputBox,
		"",
		footerLine(contentWidth, "enter send   esc back", "temporary chat"),
	)

	return renderScreen(themeColor, body)

}
