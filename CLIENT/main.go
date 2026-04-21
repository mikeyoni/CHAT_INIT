package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func exit() {

	cmd := exec.Command("exit")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

var baseURL string
var twidth int
var (
	userpromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Bold(true).
			Padding(0, 1).
			MarginLeft(2).
			Render("👤 ENTER USERNAME ")
	passpromptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9500ff")).
			Bold(true).
			Padding(0, 1).
			MarginLeft(2).
			Render("🔒 ENTER PASSWORD ")
	EmailpromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0088")).
				Bold(true).
				Padding(0, 1).
				MarginLeft(2).
				Render("✉️ ENTER EMAIL	 ")

	// 2. Create an "Arrow" style
	arrowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true).
			Render("\n   ❯ ")

	// 2. Create an "Arrow" style
	REDarrowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("rgb(205, 0, 0)")).
			Bold(true).
			Render("❯")

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")). // Purple
			Border(lipgloss.RoundedBorder()).      // Nice rounded box
			Padding(0, 3).                         // Space inside the box
			MarginLeft(2)                          // Space from the left edge

	texts = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ff9d")). // Purple
		Border(lipgloss.RoundedBorder()).      // Nice rounded box
		Padding(0, 3).                         // Space inside the box
		MarginLeft(2)                          // Space from the left edge

	boxe = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff")).
		Border(lipgloss.RoundedBorder()).Width(30).Align(lipgloss.Center)
)

var (
	// Primary Colors
	Red    = color.RGBA{255, 0, 0, 255}
	Green  = color.RGBA{0, 255, 0, 255}
	Blue   = color.RGBA{0, 0, 255, 255}
	Yellow = color.RGBA{255, 255, 0, 255}
	Cyan   = color.RGBA{0, 255, 255, 255}
	Pink   = color.RGBA{255, 0, 255, 255} // Also known as Magenta
	Purple = color.RGBA{157, 0, 255, 255} // Neon Purple
	Orange = color.RGBA{255, 136, 0, 255}
)

var yellotext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd900")).Bold(true)
var Redtext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
var warnStyle = lipgloss.NewStyle().Width(50).Align(lipgloss.Center)
var greentext = lipgloss.NewStyle().Foreground(lipgloss.Color("#3cff00")).Bold(true)
var purpultext = lipgloss.NewStyle().Foreground(lipgloss.Color("rgb(255, 0, 0)")).Bold(true)
var cynetext = lipgloss.NewStyle().Foreground(lipgloss.Color("#dcbaff"))
var wboldtext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)
var fwboldtext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true).Align(lipgloss.Right)

// this is the neon colur chaning mathod

func getRainbowColor(step int) (int, int, int) {
	step = step % 1530 // Keeps it in the 0-1530 range

	switch {
	case step < 255: // Red to Yellow
		return 255, step, 0
	case step < 510: // Yellow to Green
		return 510 - step, 255, 0
	case step < 765: // Green to Cyan
		return 0, 255, step - 510
	case step < 1020: // Cyan to Blue
		return 0, 1020 - step, 255
	case step < 1275: // Blue to Magenta
		return step - 1020, 0, 255
	default: // Magenta to Red
		return 255, 0, 1530 - step
	}
}

var (
	mytoken        string
	myuser         string
	Currentcolor   string
	Animetedcolore string
)

// this is login function it do post login info in a url in json form package

func login(url string, username string, password string) (bool, string) {

	url = fmt.Sprintf("%s/login", string(url))

	data := map[string]string{

		"username": username,
		"password": password,
	}

	jsondata, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		return false, "Network Problem Faild To Canect To Server !"
	}

	defer resp.Body.Close()

	bodybytes, _ := io.ReadAll(resp.Body)
	massage := string(bodybytes)

	parts := strings.Split(massage, ":")

	// ... after strings.Split ...
	if len(parts) >= 3 && parts[0] == "success" {
		savecradenshial(username, parts[2])
		mytoken = parts[2]
		myuser = username
		return true, "Login Successful"
	}

	if len(parts) >= 2 && parts[0] == "not" {
		return false, parts[1] // Return the actual error message
	}

	return false, "Unknown Server Error"
}

// this is register function this do post register info to the srever in json form

func emailcheck(url string, email string, username string, password string) (bool, string) {
	var waring string
	urle := fmt.Sprintf("%s/signup", url)

	data := map[string]string{

		"email":    email,
		"username": username,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {

		waring = fmt.Sprintf("\n Faild to Marshel json data : %v ", err)
		return false, waring
	}

	resp, err := http.Post(urle, "application/json-data", bytes.NewBuffer(jsondata))

	if err != nil {

		waring = fmt.Sprintf("\n Faild to post register data : %v ", err)

		return false, waring
	} else {

		msg := fmt.Sprintf("\n Sucessfully Posted Register data : %v ", resp.Status)
		fmt.Print(greentext.Render(msg))

	}

	defer resp.Body.Close()

	newbytes, _ := io.ReadAll(resp.Body)

	massage := string(newbytes)

	if massage == "success" {
		waring = "OK"
		return true, waring
	}

	if massage == "nosuccess" {
		waring = "NO"

	}

	return false, waring
}

var userinput string

func register(url string, email string, username string, password string) bool {

	url = fmt.Sprintf("%s/confarmregister", url)

	data := map[string]string{

		"email":    email,
		"username": username,
		"password": password,
		"otp":      userinput,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {

		msg := fmt.Sprintf("\n Faild to Marshel json data : %v ", err)
		fmt.Print(Redtext.Render(msg))

		return false
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {

		msg := fmt.Sprintf("\n Faild to Post Register info : %v ", err)
		fmt.Print(Redtext.Render(msg))

		return false
	} else {

		msg := fmt.Sprintf("\n Successfully Posted the Register data : %v ", resp.Status)
		fmt.Print(greentext.Render(msg))

	}

	defer resp.Body.Close()

	newbystes, _ := io.ReadAll(resp.Body)
	servermessage := string(newbystes)

	partsmessge := strings.Split(servermessage, ":")

	if len(partsmessge) > 0 && partsmessge[0] == "success" {
		fmt.Print(greentext.Render("\n login successfully \n"))
		savecradenshial(username, partsmessge[2])
		mytoken = partsmessge[2]
		return true
	}

	return false
}

// forget password

func sentforgetpassreq(email string) bool {
	url := fmt.Sprintf("%s/forgetpass", baseURL)
	data := map[string]string{"email": email}
	jsondata, _ := json.Marshal(data)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	parts := strings.Split(string(body), ":")

	if parts[0] == "done" {
		return true
	}

	return false
}

func (m LoginPageView) forgetpass(newpass string) bool {

	otpInput := m.otp
	newPassInput := newpass

	resetURL := fmt.Sprintf("%s/forgetpass?otp=%s&user=%s&new=%s", baseURL, otpInput, m.email, newPassInput)

	resp2, err := http.Post(resetURL, "application/json", nil)
	if err != nil {
		return false
	}

	defer resp2.Body.Close()

	finalBody, _ := io.ReadAll(resp2.Body)
	parts := strings.Split(string(finalBody), ":")

	if parts[0] == "success" {

		return true
	}

	return false

}

func savecradenshial(username string, tokeen string) {

	contents := fmt.Sprintf("user=%s\ntoken=%s", username, tokeen)

	err := os.WriteFile(".env", []byte(contents), 0644)

	if err != nil {
		msg := fmt.Sprintf("\n Faild to save the cradential : %v ", err)
		fmt.Print(Redtext.Render(msg))

	}

}

func savesettings(color int, animetedcolor bool) {
	yes := "false"
	if animetedcolor {
		yes = "true"
	}
	settings := fmt.Sprintf("\ncurrentcolor=%v\nanimetedcolor=%v", color, yes)
	err := os.WriteFile(".setings", []byte(settings), 0644)
	if err != nil {
		fmt.Printf(" \n \n \n \n faild to save setting's: %v \n", err)
	}
}

func chate(tusr string, token string, user string) {

	url := fmt.Sprintf("ws://localhost:4040/chat?user=%s&token=%s", user, token)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		msg := fmt.Sprintf("\n handshake faild : %v ", err)
		fmt.Print(Redtext.Render(msg))
		return
	}

	defer conn.Close()

	fmt.Print(cynetext.Render(" We are in Cannected to the server \n  ( Prss enter to Start Chat ) \n "))

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {

			_, p, err := conn.ReadMessage()

			if err != nil {
				msg := fmt.Sprintf("\n Faild to resive message : %v ", err)
				fmt.Print(Redtext.Render(msg))
				return
			}

			fmt.Println(yellotext.Render(" \n " + string(p)))

		}

	}()
	// hmm
	scanner := bufio.NewScanner(os.Stdin)

	for {

		if scanner.Scan() {

			this := fmt.Sprintf("\n TO -- > %s  ", tusr)
			fmt.Print(cynetext.Render(this) + Redtext.Render("[ !back to exit ]") + arrowStyle)
			text := scanner.Text()

			if text == "" {
				continue
			}

			if text == "!back" {
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				conn.Close()
				cls()
				DM(baseURL)
				return
			}

			msg := fmt.Sprintf("tusr:%s:user:%s", tusr, text)
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg))

			if err != nil {
				fmt.Printf("\n %v ", err)
				break
			}

		}
	}

	conn.Close()
	<-done

}

// ACtions SRQ for sentfriend req RFQ for refect friend req AFQ for accept firend req
// DLF for delate form friend

func todo(url, token string, user string, action string, targetuser string) {

	var act string

	switch action {

	case "SRQ":
		act = "sentfreq"
	case "RFQ":
		act = "rejectfreq"
	case "AFQ":
		act = "acceptfreq"
	case "DLF":
		act = "delatfre"

	}

	Url := fmt.Sprintf("%s/do?user=%s&token=%s&act=%s&tar=%s", url, user, token, act, targetuser)

	resp, err := http.Post(Url, "application/json", nil)

	if err != nil {
		fmt.Printf("\n Failed to send friend request: %v", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("\n Successfully action sented to %s!", targetuser)
	} else {
		fmt.Printf("\n Server returned error: %s", resp.Status)
	}

	rebytes, _ := io.ReadAll(resp.Body)

	message := string(rebytes)

	fmt.Printf(" \n %v  \n", message)
}

func getEmail() string {

	fmt.Print(EmailpromptStyle + REDarrowStyle)

	// 3. Take the input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

var flist []string

func DM(url string) {
	if mytoken != "" && myuser != "" {
		for {
			cls()

			flist = viewflist(url)

			title := "           Direct Message        \n"
			footer := fmt.Sprintf("\n          Chouse Numbers to msg them <3      \n %v ", purpultext.Render("  0. Exit   111. Friend Manage "))

			var list string

			if len(flist) > 0 {

				for i := range flist {

					list += fmt.Sprintf(" %v • %v\n", i+1, flist[i])

				}

			} else {

				fmt.Println(texts.Render(title + " No friend in this New World  "))
			}

			fmt.Println(texts.Render(title + wboldtext.Render(list) + footer))

			scanner := bufio.NewScanner(os.Stdin)

			fmt.Printf("%v", arrowStyle)

			if scanner.Scan() {

				cleanText := strings.TrimSpace(scanner.Text())
				choice, err := strconv.Atoi(cleanText)

				if err != nil {
					fmt.Printf("\n ❌ That's not a number, %v \n", myuser)
					time.Sleep(1 * time.Second)
					continue
				}

				if choice == 0 {
					return
				}
				if choice == 111 {
					// freindsetting()
					return
				}

				if choice > 0 && choice <= len(flist) {

					targeteruser := flist[choice-1]

					cls()
					chate(targeteruser, mytoken, myuser)

					break

				} else {
					fmt.Println("❌ Invalid friend number!")
					time.Sleep(1 * time.Second)
					continue

				}

			}

		}
	} else {

		// manue(baseURL)
		return
	}

}

func tokenchekcing(url string) bool {

	if mytoken == "" || myuser == "" {
		return false
	}

	// var resp *http.Response

	url = fmt.Sprintf("%v/checking?token=%s", url, mytoken)

	resp, err := http.Post(url, "application/checking", nil)

	if err != nil {
		return false
	}

	replaybyte, err := io.ReadAll(resp.Body)

	if err != nil {
		return false
	}

	replay := ""
	replay = string(replaybyte)

	if replay == "ok" {
		return true
	} else if replay == "no" {
		return false
	}

	defer resp.Body.Close()
	return false

}

func viewflist(url string) []string {

	url = fmt.Sprintf("%v/viewflist?token=%v&user=%v", url, mytoken, myuser)

	resp, err := http.Post(url, "aplication/flistreq", nil)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	serverbytes, _ := io.ReadAll(resp.Body)

	var Flist []string

	err = json.Unmarshal(serverbytes, &Flist)

	if err != nil {
		return nil
	}

	return Flist
}

func viewReqlist(url string) []string {

	url = fmt.Sprintf("%v/viewReqlist?token=%v&user=%v", url, mytoken, myuser)

	resp, err := http.Post(url, "aplication/Rlistreq", nil)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	serverbytes, _ := io.ReadAll(resp.Body)

	var Rlist []string

	err = json.Unmarshal(serverbytes, &Rlist)

	if err != nil {
		return nil
	}

	return Rlist
}

var Reqlist []string

type seasionState int

const (
	LoginState      seasionState = iota // 0
	DashState                           // 1
	SettingState                        // 2
	FriendlistState                     // 3
	DirectMsgState                      // 4 <--- Add this
)

type rootModel struct {
	islogedin    bool
	loginAttampt int
	state        seasionState
	login        LoginPageView
	dash         DashboardView
	friendlist   FriendlistView
	settings     SettingsView
	directmsg    DirectMsgView // <--- Add this
}

func (m rootModel) Init() tea.Cmd {
	return m.login.Init()
}

// Global struct to hold terminal dimensions
var WinSize = struct {
	Width  int
	Height int
}{}

type SwitchToSettingsMsg struct{}
type SwitchToFriendMsg struct{}
type SwitchToDashMsg struct{}

func SwitchToSettings() tea.Cmd {
	return func() tea.Msg {
		return SwitchToSettingsMsg{}
	}
}

func SwitchtoDash() tea.Cmd {
	return func() tea.Msg {
		return SwitchToDashMsg{}
	}
}

func SwitchtoFriend() tea.Cmd {
	return func() tea.Msg {
		return SwitchToFriendMsg{}
	}
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var newModel tea.Model

	// 1. GLOBAL WINDOW RESIZE (Safe)
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		WinSize.Width = msg.Width
		WinSize.Height = msg.Height
	}

	// 2. TOKEN CHECK (Restored)
	if !m.islogedin && m.loginAttampt < 3 {
		if tokenchekcing(baseURL) {
			m.islogedin = true
			m.state = DirectMsgState
			return m, m.directmsg.Init()
		}
	}

	// 3. SAFE STATE SWITCHING
	// We use "ctrl+o", "ctrl+s", etc. so typing normally doesn't break the app
	if msg, ok := msg.(tea.KeyMsg); ok {
		if m.islogedin {

			switch msg.String() {
			case "ctrl+w", "ctrl+W":
				m.state = FriendlistState
				return m, m.friendlist.Init()
			case "ctrl+s", "ctrl+S":
				m.state = SettingState
				return m, m.settings.Init()
			case "ctrl+d", "ctrl+D":
				m.state = DashState
				return m, m.dash.Init()
			}
		}
	}

	switch msg.(type) {
	case SwitchToSettingsMsg:
		m.state = SettingState
		return m, m.settings.Init()

	case SwitchToDashMsg:
		m.state = DashState
		return m, m.dash.Init()

	case SwitchToFriendMsg:
		m.state = FriendlistState
		return m, m.friendlist.Init()

	}
	// 4. ROUTING (The rest stays exactly the same)
	switch m.state {
	case LoginState:
		newModel, cmd = m.login.Update(msg)
		m.login = newModel.(LoginPageView)
		if m.login.loggedin {
			m.state = DashState
			return m, m.dash.Init()
		}

	case DashState:
		newModel, cmd = m.dash.Update(msg)
		m.dash = newModel.(DashboardView)

	case FriendlistState:
		newModel, cmd = m.friendlist.Update(msg)
		m.friendlist = newModel.(FriendlistView)

	case SettingState:
		newModel, cmd = m.settings.Update(msg)
		m.settings = newModel.(SettingsView)

	case DirectMsgState:
		newModel, cmd = m.directmsg.Update(msg)
		m.directmsg = newModel.(DirectMsgView)
	}

	return m, cmd
}

func (m rootModel) View() string {
	switch m.state {
	case LoginState:
		return m.login.View()

	case DashState:
		return m.dash.View()

	case FriendlistState:
		return m.friendlist.View()

	case SettingState:
		return m.settings.View()

	case DirectMsgState:
		return m.directmsg.View()

	default:
		return "Unknow state"
	}
}

type LoginPageView struct {
	loggedin        bool
	Width           int
	Hieght          int
	Homepageoptions []string
	choise          int
	Homeselected    bool
	Isselected      bool
	IsfullScreen    bool
	Quiting         bool
	back            bool

	// this are form home page opitons

	loginpage          bool
	registerpage       bool
	forgetpasswordpage bool

	// to take text input
	userinput string
	textinput textinput.Model
	err       error

	Riscarentinput      int
	iscarentinput       int
	username            string
	password            string
	email               string
	needotp             bool
	otp                 string
	recoveryEail        string
	recoveryOTP         bool
	NewPass             string
	forgetpasswordSteps int
	loginAttampt        int
	// client side warnings

	warning string

	// server warnings

	ServersideWarning string

	// for the color chaning

	colorstep        int
	animetedlog      bool
	currentlogoColor []string
	currentcolor     int
	glitchmode       bool
}

func (m LoginPageView) Init() tea.Cmd {
	return tick()
}

func InishialMOD() LoginPageView {
	var yes bool
	var currentcolore int
	if Currentcolor != "" || Animetedcolore != "" {
		number, _ := strconv.Atoi(Currentcolor)
		if Animetedcolore == "true" {
			yes = true
		}
		currentcolore = number
	}

	ti := textinput.New()
	ti.Placeholder = "Enter Your Username "
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 20
	return LoginPageView{
		textinput:       ti,
		err:             nil,
		Quiting:         false,
		IsfullScreen:    true,
		back:            false,
		Isselected:      false,
		Homeselected:    false,
		Homepageoptions: []string{"Login", "Register", "Forget Password", "Exit"},
		currentlogoColor: []string{
			"Red",
			"Orange",
			"Yellow",
			"Green",
			"Cyan",
			"Blue",
			"Purple",
			"Pink",
		},
		currentcolor: currentcolore,
		animetedlog:  yes,
	}
}

func iscicked(msg tea.Msg, key string) bool {
	k, ok := msg.(tea.KeyMsg)
	return ok && k.String() == key
}

type forgetResultMsg struct {
	success bool
}

func (m LoginPageView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.loggedin {
		m.Homeselected = true
		m.loginpage = false
		m.registerpage = false
		m.forgetpasswordpage = false
	}
	switch msg := msg.(type) {

	case tickMsg:
		m.colorstep = (m.colorstep + 5) % 1530
		return m, tick()

	case tea.KeyMsg:
		// Logic to check if we are currently typing in a form
		isTyping := m.loginpage || m.registerpage || m.forgetpasswordpage || m.needotp

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			// 1. Reset all page states immediately
			m.loginpage = false
			m.registerpage = false
			m.forgetpasswordpage = false
			m.needotp = false
			m.Homeselected = false // This lets the menu take control again

			// 2. FORCE the input to stop listening
			m.textinput.Blur()
			m.textinput.SetValue("")
			m.textinput.Placeholder = ""

			// 3. Reset internal counters
			m.iscarentinput = 0
			m.Riscarentinput = 0
			m.warning = ""

			// 4. Return nil for the command to stop any pending input updates
			return m, nil

		case "q", "Q":

			if !isTyping {
				m.Quiting = true
				return m, tea.Quit
			}

		case "up", "k", "K":
			if !isTyping && !m.Homeselected {
				if m.choise > 0 {
					m.choise--
				}
			}

		case "down", "j", "J":
			if !isTyping && !m.Homeselected {
				if m.choise < len(m.Homepageoptions)-1 {
					m.choise++
				}
			}

		case "i", "I", "tab":
			if !isTyping {
				if m.currentcolor >= 0 && m.currentcolor < len(m.currentlogoColor)-1 {
					m.currentcolor++
				} else {
					m.currentcolor = 0
				}
				savesettings(m.currentcolor, m.animetedlog)
			}

		case "Y", "y":
			if !isTyping {
				m.glitchmode = !m.glitchmode
			}

		case "g", "G":
			if !isTyping {
				m.animetedlog = !m.animetedlog
				savesettings(m.currentcolor, m.animetedlog)
			}

		case "enter":

			if !m.Homeselected {
				m.Homeselected = true
			}

			// --- LOGIN LOGIC ---
			if m.loginpage {
				if m.iscarentinput == 0 {
					m.warning = ""
					m.username = m.textinput.Value()
					m.textinput.SetValue("")
					m.textinput.Placeholder = "Enter Your Password "
					m.textinput.EchoMode = textinput.EchoPassword
					m.iscarentinput = 1
				} else if m.iscarentinput == 1 {
					m.password = m.textinput.Value()
					m.textinput.SetValue("")
					m.iscarentinput++

					if m.username != "" && m.password != "" {
						k, warnings := login(baseURL, m.username, m.password)
						if warnings != "" {
							m.ServersideWarning = warnings
						}
						if k {
							m.loginpage = false
							m.loggedin = true
							return m, nil
						} else {
							m.iscarentinput = 0
							m.warning = m.ServersideWarning
							if m.warning == "" {
								m.warning = "Invalid Credentials. Try Again."
							}
							m.username = ""
							m.password = ""
							m.textinput.SetValue("")
							m.textinput.Placeholder = "Enter Your Username"
							m.textinput.EchoMode = textinput.EchoNormal
						}
					}
					return m, nil
				}
			}

			// --- REGISTER LOGIC ---
			if m.registerpage {
				if m.needotp {
					otpVal := m.textinput.Value()
					if otpVal != "" {
						if register(baseURL, m.email, m.username, m.password) {
							m.loginpage = false // Or your post-registration state
						} else {
							m.needotp = false
							m.Riscarentinput = 0
							m.username = ""
							m.password = ""
							m.textinput.SetValue("")
							m.textinput.Placeholder = "Enter Your Username"
							m.textinput.EchoMode = textinput.EchoNormal
							return m, nil
						}
					}
				} else {
					if m.Riscarentinput == 0 {
						m.warning = ""
						m.username = m.textinput.Value()
						m.textinput.SetValue("")
						m.textinput.Placeholder = "Enter Your Email"
						m.Riscarentinput = 1
					} else if m.Riscarentinput == 1 {
						m.warning = ""
						m.email = m.textinput.Value()
						m.textinput.SetValue("")
						m.textinput.Placeholder = "Enter Your Password"
						m.textinput.EchoMode = textinput.EchoPassword
						m.Riscarentinput = 2
					} else if m.Riscarentinput == 2 {
						m.password = m.textinput.Value()
						m.textinput.SetValue("")

						ok, msg := emailcheck(baseURL, m.email, m.username, m.password)
						if msg != "" {
							m.warning = msg
						}
						if ok {
							m.warning = ""
							m.needotp = true
							m.textinput.Placeholder = "Enter Your OTP"
							m.textinput.EchoMode = textinput.EchoNormal
							m.Riscarentinput++
						} else {
							m.textinput.EchoMode = textinput.EchoNormal
							m.warning = "Failed to send OTP. Try again."
							m.Riscarentinput = 0
							m.textinput.Placeholder = "Enter Your Username"
						}
					}
					return m, nil
				}
			}

			// forget pass

			if m.forgetpasswordpage {

				if m.forgetpasswordSteps == 0 {

					m.recoveryEail = m.textinput.Value()
					m.textinput.SetValue("")
					m.textinput.Placeholder = "OTP"
					m.forgetpasswordSteps = 1

					if sentforgetpassreq(m.recoveryEail) {
						m.warning = fmt.Sprintf("SuccessFully Otp Send To : %s", m.recoveryEail)
					}

				} else if m.forgetpasswordSteps == 1 {
					m.otp = m.textinput.Value()
					m.textinput.SetValue("")
					m.textinput.Placeholder = " New Password"
					m.forgetpasswordSteps = 2

				} else if m.forgetpasswordSteps == 2 {
					m.NewPass = m.textinput.Value()
					m.textinput.SetValue("")

					// Create a copy of the values so they don't change while the request runs
					newPass := m.NewPass
					email := m.recoveryEail
					otp := m.otp

					// Return a command! This runs in the background.
					return m, func() tea.Msg {
						// Build the URL manually here to ensure it's clean
						resetURL := fmt.Sprintf("%s/forgetpass?otp=%s&user=%s&new=%s", baseURL, otp, email, newPass)
						resp, err := http.Post(resetURL, "application/json", nil)

						if err != nil {
							return forgetResultMsg{success: false}
						}
						defer resp.Body.Close()

						body, _ := io.ReadAll(resp.Body)
						// Check if server returned "success"
						if strings.HasPrefix(string(body), "success") {
							return forgetResultMsg{success: true}
						}
						return forgetResultMsg{success: false}
					}
				}

				return m, nil

			}
		}

	case tea.WindowSizeMsg:
		m.Hieght = msg.Height
		m.Width = msg.Width
		twidth = msg.Width

	case forgetResultMsg:
		if msg.success {
			m.forgetpasswordpage = false
			m.Homeselected = false
			m.warning = "Password Reset Successfully ! "
			m.forgetpasswordSteps = 0
			m.loginpage = true

		} else {
			m.forgetpasswordSteps = 0
			m.warning = "Faild to recover. Try again. "
		}
		return m, nil
	}

	// Menu Selection Logic
	if m.Homeselected && !m.loginpage && !m.registerpage && !m.forgetpasswordpage {
		switch m.choise {
		case 0:
			m.loginpage = true
		case 1:
			m.registerpage = true
		case 2:
			m.forgetpasswordpage = true
			m.textinput.Placeholder = "Enter Your Email"
		case 3:
			return m, tea.Quit
		}
	}

	// Ensure the text input actually processes characters
	if m.loginpage || m.registerpage || m.forgetpasswordpage || m.needotp {
		m.textinput.Focus()
		m.textinput, cmd = m.textinput.Update(msg)
	} else {
		// If no forms are open, the cursor should be off
		m.textinput.Blur()
	}

	return m, cmd

}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*40, func(t time.Time) tea.Msg {
		return tickMsg{} // Add the curly braces here!
	})
}

func animetedmakeGradientText(text string, step int, N bool) string {

	r, g, b := getRainbowColor(step)

	startColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255} //
	var endColor = color.RGBA{}

	if N {
		r2, g2, b2 := getRainbowColor(step + 500)
		endColor = color.RGBA{uint8(g2), uint8(r2), uint8(b2), 255}

	} else if !N {
		endColor = color.RGBA{255, 255, 255, 255}
	}

	runes := []rune(text)
	var out strings.Builder

	for i, r := range runes {

		f := float64(i) / float64(len(runes))

		currColor := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
			uint8(float64(startColor.R)*(1-f)+float64(endColor.R)*f),
			uint8(float64(startColor.G)*(1-f)+float64(endColor.G)*f),
			uint8(float64(startColor.B)*(1-f)+float64(endColor.B)*f)))

		out.WriteString(lipgloss.NewStyle().Foreground(currColor).Render(string(r)))
	}
	return out.String()
}

func makeGradientText(text string, colors []string, N int) string {

	// Default fallback colors (White to Grey)
	startColor := color.RGBA{255, 255, 255, 255}
	endColor := color.RGBA{100, 100, 100, 255}

	// Safety check for index out of bounds
	if N >= 0 && N < len(colors) {
		switch colors[N] {
		case "Red":
			// Red is deep, so we fade it to a bright "Glow"
			startColor = color.RGBA{255, 0, 0, 255}
			endColor = color.RGBA{255, 200, 200, 255} // Fades to a light pinkish-white
		case "Orange":
			// Orange looks best fading into a deep burnt shadow
			startColor = color.RGBA{255, 136, 0, 255}
			endColor = color.RGBA{50, 20, 0, 255} // Deep "Burnt" shadow
		case "Yellow":
			// Yellow is the brightest, so we fade to dark to make it readable
			startColor = color.RGBA{255, 255, 0, 255}
			endColor = color.RGBA{40, 40, 0, 255} // Dark Olive shadow
		case "Green":
			// Neon Green to deep forest shadow
			startColor = color.RGBA{0, 255, 0, 255}
			endColor = color.RGBA{0, 30, 0, 255} // Very dark green
		case "Cyan":
			// Cyan is bright, so fade it to a deep ocean blue/black
			startColor = color.RGBA{0, 255, 255, 255}
			endColor = color.RGBA{0, 20, 40, 255} // Deep midnight blue
		case "Blue":
			// Blue is dark, so we make it GLOW to pure white
			startColor = color.RGBA{0, 0, 255, 255}
			endColor = color.RGBA{255, 255, 255, 255} // Pure White Glow
		case "Purple":
			// Purple to White looks incredible for a "Pirate King" vibe
			startColor = color.RGBA{157, 0, 255, 255}
			endColor = color.RGBA{255, 255, 255, 255} // Pure White Glow
		case "Pink":
			// Pink to a deep velvet shadow
			startColor = color.RGBA{255, 0, 255, 255}
			endColor = color.RGBA{40, 0, 40, 255} // Deep Magenta shadow
		}
	}

	runes := []rune(text)
	var out strings.Builder

	for i, r := range runes {
		// Calculate the interpolation factor (0.0 to 1.0)
		f := float64(i) / float64(len(runes))

		// Standard RGB linear interpolation math
		currColor := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
			uint8(float64(startColor.R)*(1-f)+float64(endColor.R)*f),
			uint8(float64(startColor.G)*(1-f)+float64(endColor.G)*f),
			uint8(float64(startColor.B)*(1-f)+float64(endColor.B)*f)))

		out.WriteString(lipgloss.NewStyle().Foreground(currColor).Render(string(r)))
	}
	return out.String()
}

func (m LoginPageView) View() string {
	width := WinSize.Width
	var themeColor string

	if m.animetedlog {
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
	var selectedboxe = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(themeColor)).
		BorderForeground(lipgloss.Color(themeColor)).
		Border(lipgloss.RoundedBorder()).
		Width(50)

	// Update the version text color too!
	Versions := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Right).
		Foreground(lipgloss.Color(themeColor))

	// var selectedboxe = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0037")).
	// 		BorderForeground(lipgloss.Color("#ff0059")).
	// 		Border(lipgloss.RoundedBorder()).Width(30).Align(lipgloss.Center)
	// // inishializing rainbow color

	if !m.Homeselected {
		selectedboxe = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(themeColor)).       // Text matches theme
			BorderForeground(lipgloss.Color(themeColor)). // Border matches theme
			Border(lipgloss.RoundedBorder()).
			Width(30).
			Align(lipgloss.Center)
	}

	var boxrender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(themeColor)).
		Width(width-4).Padding(0, 0).Align(lipgloss.Center)
	v := "\n your welcome to chat init \n"

	Shortcut := lipgloss.NewStyle().Width((width - 11) / 2).Align(lipgloss.Left).
		Foreground(lipgloss.Color("#ffffff9b"))

	// Versions := lipgloss.NewStyle().Width((m.Width - 11) / 2).Align(lipgloss.Right).
	// 	Foreground(lipgloss.Color(themeColor))

	subtitle := cynetext.Render("BY ui_mik3y | YT && INSTA <3 ")
	Footther := lipgloss.NewStyle().Width(width - 10).Bold(true).
		Foreground(lipgloss.Color("rgb(0, 0, 0)"))

	var l string

	if m.animetedlog {

		l = animetedmakeGradientText(`
		
 ██████╗██╗  ██╗ █████╗ ████████╗    ██╗███╗   ██╗██╗████████╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██║████╗  ██║██║╚══██╔══╝
██║     ███████║███████║   ██║       ██║██╔██╗ ██║██║   ██║   
██║     ██╔══██║██╔══██║   ██║       ██║██║╚██╗██║██║   ██║   
╚██████╗██║  ██║██║  ██║   ██║       ██║██║ ╚████║██║   ██║   
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝       ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝   
                                    `, m.colorstep*2, m.glitchmode)

	} else {

		l = makeGradientText(`
		
 ██████╗██╗  ██╗ █████╗ ████████╗    ██╗███╗   ██╗██╗████████╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██║████╗  ██║██║╚══██╔══╝
██║     ███████║███████║   ██║       ██║██╔██╗ ██║██║   ██║   
██║     ██╔══██║██╔══██║   ██║       ██║██║╚██╗██║██║   ██║   
╚██████╗██║  ██║██║  ██║   ██║       ██║██║ ╚████║██║   ██║   
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝       ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝   
                                    `, m.currentlogoColor, m.currentcolor)

	}
	// 	logo := fmt.Sprintf(`

	//  ██████╗██╗  ██╗ █████╗ ████████╗    ██╗███╗   ██╗██╗████████╗
	// ██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██║████╗  ██║██║╚══██╔══╝
	// ██║     ███████║███████║   ██║       ██║██╔██╗ ██║██║   ██║
	// ██║     ██╔══██║██╔══██║   ██║       ██║██║╚██╗██║██║   ██║
	// ╚██████╗██║  ██║██║  ██║   ██║       ██║██║ ╚████║██║   ██║
	//  ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝       ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝
	//                                  %v   `, subtitle)

	var render string

	if !m.Homeselected {
		render += "\n"
		for i := range m.Homepageoptions {

			if m.choise == i {

				render += "\n" + selectedboxe.Render(m.Homepageoptions[i])

			} else {

				render += "\n" + boxe.Render(m.Homepageoptions[i])

			}
		}

	}

	var warningRender = ""

	if m.loginpage {

		render += "\n"

		smallbox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(50)
		if m.iscarentinput == 0 {
			render += selectedboxe.Render(wboldtext.Render(" USERNAME ", m.textinput.View()))

			render += "\n"

			render += smallbox.Render(wboldtext.Render(" PASSWORD : "))

		} else if m.iscarentinput == 1 {
			render += smallbox.Render(wboldtext.Render(" USERNAME : ", m.username))

			render += "\n"

			render += selectedboxe.Render(wboldtext.Render(" PASSWORD ", m.textinput.View()))

		} else if m.iscarentinput > 1 {

			render += smallbox.Render(wboldtext.Render(" USERNAME : ", m.username))

			render += "\n"

			render += smallbox.Render(wboldtext.Render(" PASSWORD : ", strings.Repeat("*", len(m.password))))

			render += "\n"
		}

		if m.warning != "" {
			warningRender = lipgloss.NewStyle().
				Width(50).
				Align(lipgloss.Center).
				MarginTop(1). // Adds space without breaking layout
				Render(Redtext.Render(m.warning))
		}

	}

	if m.registerpage {
		render += "\n"

		smallbox := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Width(50)

		if m.needotp {

			render += wboldtext.Render("Enter The OTP We Sented On : ", m.email)
			render += "\n"
			render += selectedboxe.Render(wboldtext.Render(" OTP ", m.textinput.View()))

		}

		if !m.needotp {

			if m.Riscarentinput == 0 {
				render += selectedboxe.Render(wboldtext.Render(" USERNAME ", m.textinput.View()))

				render += "\n"

				render += smallbox.Render(wboldtext.Render(" EMAIL : "))

				render += "\n"

				render += smallbox.Render(wboldtext.Render(" PASSWORD : "))

			} else if m.Riscarentinput == 1 {
				render += smallbox.Render(wboldtext.Render(" USERNAME : ", m.username))
				render += "\n"

				render += selectedboxe.Render(wboldtext.Render(" EMAIL ", m.textinput.View()))
				render += "\n"

				render += smallbox.Render(wboldtext.Render(" PASSWORD "))

			} else if m.Riscarentinput == 2 {
				render += smallbox.Render(wboldtext.Render(" USERNAME : ", m.username))
				render += "\n"

				render += smallbox.Render(wboldtext.Render(" EMAIL : ", m.email))
				render += "\n"

				render += selectedboxe.Render(wboldtext.Render(" PASSWORD ", m.textinput.View()))

			} else if m.Riscarentinput > 2 {

				render += smallbox.Render(wboldtext.Render(" USERNAME : ", m.username))
				render += "\n"

				render += smallbox.Render(wboldtext.Render(" EMAIL : ", m.email))
				render += "\n"

				render += smallbox.Render(wboldtext.Render(" PASSWORD : ", strings.Repeat("*", len(m.password))))
				render += "\n"

			}

			if m.warning != "" {

				warningRender = lipgloss.NewStyle().
					Width(50).
					Align(lipgloss.Center).
					MarginTop(1). // Adds space without breaking layout
					Render(Redtext.Render(m.warning))
			}
		}
	}

	if m.forgetpasswordpage {
		// render += "\n"
		// render += "\n"
		// render += wboldtext.Render(" THIS OPTION IS UNDER WORK :D ")
		render += "\n"

		switch m.forgetpasswordSteps {
		case 0:
			render += "\n"
			render += wboldtext.Render(" Enter Your Recovery Email ")
			render += "\n"
			render += selectedboxe.Render(" Email ", m.textinput.View())
			render += "\n"

		case 1:
			render += "\n"
			render += wboldtext.Render(" Enter Your OTP on : ", m.recoveryEail)
			render += "\n"
			render += selectedboxe.Render(" OTP ", m.textinput.View())
			render += "\n"

		case 2:
			render += "\n"
			render += wboldtext.Render(" Enter Your New Paword ")
			render += "\n"
			render += selectedboxe.Render(" PASSWORD ", m.textinput.View())
			render += "\n"
		}

		if m.warning != "" {

			warningRender = lipgloss.NewStyle().
				Width(50).
				Align(lipgloss.Center).
				MarginTop(1). // Adds space without breaking layout
				Render(Redtext.Render(m.warning))
		}

	}

	// if m.warning != "" {

	// 	warningRender = lipgloss.NewStyle().
	// 		Width(50).
	// 		Align(lipgloss.Center).
	// 		MarginTop(1). // Adds space without breaking layout
	// 		Render(Redtext.Render(m.warning))
	// }

	centerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		l+subtitle, render,
		warningRender, "\n",
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("'ESC' = Back 'Q' = Quit < 'I' & 'G' "), Versions.Render("v.1.02"))

	v = boxrender.Render(centerContent)

	return v

}

// this is where things starts
func main() {

	baseURL = "http://localhost:4040"

	// Try to load. If it fails (file deleted), mytoken/myuser will stay empty ""
	godotenv.Load(".env")
	mytoken = os.Getenv("token")
	myuser = os.Getenv("user")

	godotenv.Load(".setings")
	Currentcolor = os.Getenv("currentcolor")
	Animetedcolore = os.Getenv("animetedcolor")

	// 1. Setup the initial state of your models
	m := rootModel{
		state:      LoginState,
		login:      InishialMOD(), // Use a helper function to set up text inputs
		friendlist: FriendlistView{},
		dash:       NewDashboard(),
		settings:   SettingsView{},
		directmsg:  NewDirectMsg(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen()) // Use Altscreen for a "clean" look

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, a crash! Error: %v\n", err)
		os.Exit(1)
	}

}
