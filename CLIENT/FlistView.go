package main

import (
    "fmt"
    "strconv"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type FriendlistView struct {
    currentcolor     int
    animetedcolor    bool
    colorstep        int
    currentlogoColor []string
    glitchmode       bool
    // List specific
    friends []string
    cursor  int
    // Universal items
    warning string
    active  bool
}

func (m FriendlistView) Init() tea.Cmd {
    return tick()
}

func NewFriendlist() FriendlistView {
    fl := FriendlistView{
        currentlogoColor: []string{
            "Red", "Orange", "Yellow", "Green",
            "Cyan", "Blue", "Purple", "Pink",
        },
        friends: []string{
            "Mikey [Online]",
            "Pirate_King [Away]",
            "Fedora_Pro [Online]",
            "Root_User [Busy]",
            "Gopher_01 [Offline]",
        },
    }

    // LOAD SAVED SETTINGS HERE (Only once!)
    if Currentcolor != "" {
        if number, err := strconv.Atoi(Currentcolor); err == nil {
            fl.currentcolor = number
        }
    }

    if Animetedcolore == "true" {
        fl.animetedcolor = true
    }

    return fl
}

func (m FriendlistView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.active = true

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {

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

    case tickMsg:
        m.colorstep = (m.colorstep + 5) % 1530
        return m, tick()

    }
    return m, nil
}

func (m FriendlistView) View() string {

    width := WinSize.Width

    var render string
    var warningRender string
    var themeColor string

    if m.animetedcolor {
        r, g, b := getRainbowColor(m.colorstep * 2)
        themeColor = fmt.Sprintf("#%02x%02x%02x", r, g, b)
    } else {
        themeColor = "#7D56F4"
        if m.currentcolor >= 0 && m.currentcolor < len(m.currentlogoColor) {
            switch m.currentlogoColor[m.currentcolor] {
            case "Red":    themeColor = "#FF0000"
            case "Orange": themeColor = "#FF8800"
            case "Yellow": themeColor = "#FFFF00"
            case "Green":  themeColor = "#00FF00"
            case "Cyan":   themeColor = "#00FFFF"
            case "Blue":   themeColor = "#0000FF"
            case "Purple": themeColor = "#9D00FF"
            case "Pink":   themeColor = "#FF00FF"
            }
        }
    }

    Versions := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Right).
        Foreground(lipgloss.Color(themeColor))

    var boxrender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).
        BorderForeground(lipgloss.Color(themeColor)).
        Width(width-4).Padding(0, 0).Align(lipgloss.Center)

    var l string
    if m.animetedcolor {
        l = animetedmakeGradientText("FRIENDS", m.colorstep*2, m.glitchmode)
    } else {
        l = makeGradientText("FRIENDS", m.currentlogoColor, m.currentcolor)
    }

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