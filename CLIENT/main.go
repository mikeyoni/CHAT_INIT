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
			Render("\n   ❯ ")

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

	selectedboxe = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff0037")).
			BorderForeground(lipgloss.Color("#ff0059")).
			Border(lipgloss.RoundedBorder()).Width(30).Align(lipgloss.Center)
)

var yellotext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd900"))
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
	mytoken string
	myuser  string
)

// this is login function it do post login info in a url in json form package

func login(url string, username string, password string) (bool, string) {
	var treturn string
	url = fmt.Sprintf("%s/login", string(url))

	data := map[string]string{

		"username": username,
		"password": password,
	}

	jsondata, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		return false, "Network Problem Faild To Canect To Server !"
	} else {

		treturn = fmt.Sprintf(" \n sucessfully sented : %v ", resp.Status)
	}

	defer resp.Body.Close()

	bodybytes, _ := io.ReadAll(resp.Body)
	massage := string(bodybytes)

	parts := strings.Split(massage, ":")

	if len(parts) > 0 && parts[0] == "success" {

		savecradenshial(username, parts[2])
		mytoken = parts[2]
		myuser = username

		return true, treturn
	}

	if parts[0] == "not" {
		fmt.Printf(" \n %v ", parts[1])
	}

	return false, ""
}

// this is register function this do post register info to the srever in json form

func emailcheck(url string, email string, username string, password string) bool {

	urle := fmt.Sprintf("%s/signup", url)

	data := map[string]string{

		"email":    email,
		"username": username,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {

		msg := fmt.Sprintf("\n Faild to Marshel json data : %v ", err)
		fmt.Print(purpultext.Render(msg))
		return false
	}
	resp, err := http.Post(urle, "application/json-data", bytes.NewBuffer(jsondata))

	if err != nil {
		msg := fmt.Sprintf("\n Faild to post register data : %v ", err)
		fmt.Print(Redtext.Render(msg))

		return false
	} else {

		msg := fmt.Sprintf("\n Sucessfully Posted Register data : %v ", resp.Status)
		fmt.Print(greentext.Render(msg))

	}

	defer resp.Body.Close()

	newbytes, _ := io.ReadAll(resp.Body)

	massage := string(newbytes)

	if massage == "success" {
		if register(url, email, username, password) {

			return true
		}
	}

	return false
}

func register(url string, email string, username string, password string) bool {

	url = fmt.Sprintf("%s/confarmregister", url)

	var userinput string
	fmt.Print(cynetext.Render(" \n Enter the otp : "))
	fmt.Scan(&userinput)

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

func forgetpass(baseURL string, email string) bool {
	// 1. REQUEST THE OTP
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
		fmt.Printf("\n %s", parts[1])

		var otpInput string
		var newPassInput string
		fmt.Printf("\n Enter OTP: ")
		fmt.Scan(&otpInput)
		fmt.Printf("\n Enter New Password: ")
		fmt.Scan(&newPassInput)

		resetURL := fmt.Sprintf("%s/forgetpass?otp=%s&user=%s&new=%s", baseURL, otpInput, email, newPassInput)

		resp2, err := http.Post(resetURL, "application/json", nil)
		if err != nil {
			return false
		}
		defer resp2.Body.Close()

		finalBody, _ := io.ReadAll(resp2.Body)
		fmt.Printf("\n Final Result: %s", string(finalBody))
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

	url = fmt.Sprintf("%v/checking?token=%s", url, mytoken)

	resp, _ := http.Post(url, "application/checking", nil)

	replaybyte, _ := io.ReadAll(resp.Body)

	replay := string(replaybyte)

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

type model struct {
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

	iscarentinput int
	username      string
	password      string

	// client side warnings

	warning string

	// server warnings

	ServersideWarning string

	// for the color chaning

	colorstep int
}

func (m model) Init() tea.Cmd {
	return tick()
}

func InishialMOD() model {
	ti := textinput.New()
	ti.Placeholder = "Enter Your Username "
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 20
	return model{
		textinput:       ti,
		err:             nil,
		Quiting:         false,
		IsfullScreen:    true,
		back:            false,
		Isselected:      false,
		Homeselected:    false,
		Homepageoptions: []string{"Login", "Register", "Forget Password", "Exit"},
	}
}

func iscicked(msg tea.Msg, key string) bool {
	k, ok := msg.(tea.KeyMsg)
	return ok && k.String() == key
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// colure chanign

	switch msg := msg.(type) {

	case tickMsg:
		m.colorstep = (m.colorstep + 5) % 1530
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":

			if myuser == "" || mytoken == "" {
				m.Homeselected = false
				m.loginpage = false
				m.registerpage = false
				m.forgetpasswordpage = false
				m.warning = ""
				m.ServersideWarning = ""
			}

		case "q", "Q":
			m.Quiting = true
			return m, tea.Quit

		case "up", "k", "K":
			// this controles are for home page
			if !m.Homeselected {
				if m.choise > 0 {
					m.choise--
				}

			}
		case "down", "j", "J":
			// this controles are home home page
			if !m.Homeselected {
				if m.choise < len(m.Homepageoptions)-1 {
					m.choise++
				}

			}

		case "enter":

			if !m.Homeselected {
				m.Homeselected = true
			}

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
							// successfuly login
							return m, nil
						} else {
							m.iscarentinput = 0
							if m.ServersideWarning != "" {

								m.warning = m.ServersideWarning
							} else {

								m.warning = "Invalid Cradenshial Try Againg . "

							}

							m.username = ""
							m.password = ""
							m.textinput.SetValue("")
							m.textinput.Placeholder = "Enter Your Username "
							m.textinput.EchoMode = textinput.EchoNormal
						}
					}

					return m, nil
				} else {

					m.iscarentinput = 0
					m.textinput.EchoMode = textinput.EchoNormal
				}

			}

		}

	case tea.WindowSizeMsg:
		m.Hieght = msg.Height
		m.Width = msg.Width
		twidth = msg.Width
	}

	// logic of home page

	if m.Homeselected {
		if m.choise == 0 {
			m.loginpage = true
		}
		if m.choise == 1 {
			m.registerpage = true
		}
		if m.choise == 2 {
			m.forgetpasswordpage = true
		}
		if m.choise == 3 {

			return m, tea.Quit
		}
	}

	if m.loginpage && (m.iscarentinput == 1 || m.iscarentinput == 0) {
		m.textinput, cmd = m.textinput.Update(msg)
	}

	return m, cmd
}

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*20, func(t time.Time) tea.Msg {
		return tickMsg{} // Add the curly braces here!
	})
}

func makeGradientText(text string , step int) string {

	r, g, b:= getRainbowColor(step)

	startColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}  //
	endColor := color.RGBA{255, 255, 255, 0} // Purple

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

func (m model) View() string {

	// inishializing rainbow color


	var boxrender = lipgloss.NewStyle().Border(lipgloss.ThickBorder()).Width(m.Width-4).Padding(0, 0).Align(lipgloss.Center)
	v := "\n your welcome to chat init \n"

	Shortcut := lipgloss.NewStyle().Width((m.Width - 11) / 2).Align(lipgloss.Left).
		Foreground(lipgloss.Color("#ffffff9b"))

	Versions := lipgloss.NewStyle().Width((m.Width - 11) / 2).Align(lipgloss.Right).
		Foreground(lipgloss.Color("#ff0000"))

	subtitle := cynetext.Render("BY ui_mik3y | YT && INSTA <3 ")
	Footther := lipgloss.NewStyle().Width(m.Width - 10).Bold(true).
		Foreground(lipgloss.Color("#ffffff00"))

	l := makeGradientText(`
		
 ██████╗██╗  ██╗ █████╗ ████████╗    ██╗███╗   ██╗██╗████████╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝    ██║████╗  ██║██║╚══██╔══╝
██║     ███████║███████║   ██║       ██║██╔██╗ ██║██║   ██║   
██║     ██╔══██║██╔══██║   ██║       ██║██║╚██╗██║██║   ██║   
╚██████╗██║  ██║██║  ██║   ██║       ██║██║ ╚████║██║   ██║   
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝       ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝   
                                    ` , m.colorstep * 2 )

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
			render += smallbox.Render(wboldtext.Render(" USERNAME ", m.textinput.View()))

			render += "\n"

			render += smallbox.Render(wboldtext.Render(" PASSWORD : "))

		} else if m.iscarentinput == 1 {
			render += smallbox.Render(wboldtext.Render(" USERNAME : ", m.username))

			render += "\n"

			render += smallbox.Render(wboldtext.Render(" PASSWORD ", m.textinput.View()))

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

	}

	if m.forgetpasswordpage {

	}

	centerContent := lipgloss.JoinVertical(
		lipgloss.Center,
		l+subtitle, render,
		warningRender,
	)

	centerContent += "\n" + Footther.Render(Shortcut.Render("'ESC' = Back 'Q' = Quit  "), Versions.Render("v.1.02"))

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

	App := tea.NewProgram(InishialMOD(), tea.WithAltScreen())

	if _, err := App.Run(); err != nil {
		fmt.Printf("%v", err)
	}

}
