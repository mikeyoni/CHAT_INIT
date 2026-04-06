package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
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
)

var yellotext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd900"))
var Redtext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)
var greentext = lipgloss.NewStyle().Foreground(lipgloss.Color("#3cff00")).Bold(true)
var purpultext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0062")).Bold(true)
var cynetext = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
var wboldtext = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)

var (
	mytoken string
	myuser  string
)

// this is login function it do post login info in a url in json form package

func login(url string, username string, password string) bool {

	url = fmt.Sprintf("%s/login", string(url))

	data := map[string]string{

		"username": username,
		"password": password,
	}

	jsondata, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {

		fmt.Printf("\n Faild to sent post req : %v ", err)
		return false
	} else {

		fmt.Printf(" \n sucessfully sented : %v ", resp.Status)
	}

	defer resp.Body.Close()

	bodybytes, _ := io.ReadAll(resp.Body)
	massage := string(bodybytes)

	parts := strings.Split(massage, ":")

	if len(parts) > 0 && parts[0] == "success" {

		savecradenshial(username, parts[2])
		mytoken = parts[2]
		myuser = username

		return true
	}

	if parts[0] == "not" {
		fmt.Printf(" \n %v ", parts[1])
	}

	return false
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

	showModernLogo()

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

func showModernLogo() {
	// 1. Create a stylish box for the Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#5A189A")). // Deep Purple
		Padding(0, 2).
		MarginLeft(2)

	// 2. Create a "Tag" for the version
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#00F5D4")). // Neon Teal
		Padding(0, 1)

	// 3. Create a subtitle line
	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9D4EDD")).
		Italic(true).
		MarginLeft(2).
		Render("Ship: Thousand Sunny • Location: New World")

	// Combine them
	header := lipgloss.JoinHorizontal(lipgloss.Top,
		titleStyle.Render(" CHAT-INIT "),
		versionStyle.Render(" v2.0 "),
	)

	fmt.Println("\n" + header)
	fmt.Println(subtitle + "\n")
}

func showPirateLogo() {
	poster := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#E8833A")). // Bounty Gold
		Padding(1, 4).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#FAFAFA")).
		Render("  ⚓ CHAT-INIT ⚓  \n  THE GREAT PIRATE ERA  ")

	fmt.Println(poster)
}

func showBigModernLogo() {
	// 1. Create two separate figures for the two words
	// "block" is the boldest clean font available
	chatFig := figure.NewFigure(" CHAT-", "block", true).String()
	initFig := figure.NewFigure("INIT ", "block", true).String()

	// 2. Create Styles for each word
	chatStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fff700")). // Luffy Red
		Bold(true)

	initStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00c3ff")). // Bounty Gold
		Bold(true)

	// 3. Join them horizontally so they sit on the same line
	// This makes the logo two different colors side-by-side
	fullLogo := lipgloss.JoinHorizontal(lipgloss.Top,
		chatStyle.Render(chatFig),
		initStyle.Render(initFig),
	)

	fmt.Println(fullLogo)

	// 4. Subtitle with a nice border or highlight
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00F5D4")).
		Background(lipgloss.Color("#1A1A1A")). // Darker background for contrast
		Padding(0, 1).
		MarginLeft(4).
		Render(" Ship: Thousand Sunny • Location: New World v2.0 by @mikeyoni < GH ")

	fmt.Println(subtitleStyle)
}

func getUsername() string {

	fmt.Print(userpromptStyle + arrowStyle)

	// 3. Take the input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func getPassword() string {

	fmt.Print(passpromptStyle + REDarrowStyle)

	// 3. Take the input
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
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
			showBigModernLogo()
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
					freindsetting()
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

		manue(baseURL)
		return
	}

}

func freindsetting() {

	var next bool

	for {

		cls()
		showBigModernLogo()

		title := " \n            Friend Manage            \n"
		cmds := " \n 1. Add Friend  2. Delate Friend 3. Accept Friend \n 4. Reject Friend   5. Next  6. Friend list 0. Exit    \n"
		var listVer string

		Reqlist = viewReqlist(baseURL)
		titless := greentext.Render("    Friends     ")
		Rtitless := greentext.Render("    Requested     ")
		Friendlist := flist
		Reclist := Reqlist

		var reqLisevew string

		for i := range Friendlist {
			listVer += fmt.Sprintf("\n • %v  \n", Friendlist[i])
		}

		for i := range Reclist {
			reqLisevew += fmt.Sprintf("\n • %v  \n", Reclist[i])
		}

		switch next {
		case true:

			fmt.Println(texts.Render(title + texts.Render(Rtitless+wboldtext.Render(reqLisevew)) + purpultext.Render(cmds)))

		default:

			fmt.Println(texts.Render(title + texts.Render(titless+wboldtext.Render(listVer)) + purpultext.Render(cmds)))

		}

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
				DM(baseURL)
				return
			}

			if choice > 0 && choice <= 6 {

				if choice == 1 {
					user := getUsername()

					todo(baseURL, mytoken, myuser, "SRQ", user)
				}

				if choice == 2 {
					user := getUsername()

					todo(baseURL, mytoken, myuser, "DLF", user)
				}

				if choice == 3 {
					user := getUsername()

					todo(baseURL, mytoken, myuser, "AFQ", user)
				}

				if choice == 4 {
					user := getUsername()

					todo(baseURL, mytoken, myuser, "RFQ", user)
				}

				if choice == 5 {
					next = true
					continue
				}

				if choice == 6 {
					next = false
					continue
				}
			}

		} else {
			fmt.Println("❌ Invalid Option number!")
			time.Sleep(1 * time.Second)
			continue

		}

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

func manue(url string) {
	for {

		cls()
		showBigModernLogo()

		// If globals are empty, we go straight to Login/Register
		if mytoken == "" || myuser == "" {
			fmt.Println(headerStyle.Render(" WELCOME TO CHAT-INIT "))
			fmt.Println(texts.Render(" CHOOSE OPTION 0-3   \n 1. Login  \n 2. Register \n 3. Forget-password   \n 0. Exit "))

			userinput := bufio.NewScanner(os.Stdin)
			fmt.Printf("%v", cynetext.Render("  > "))

			if userinput.Scan() {
				text := strings.TrimSpace(userinput.Text())
				switch text {
				case "1":
					username := getUsername()
					password := getPassword()
					if login(url, username, password) {
						DM(url)
						return
					} else {
						fmt.Println(Redtext.Render("\n Failed to login. Try again!"))
						time.Sleep(1 * time.Second)
					}

				case "2":
					username := getUsername()
					password := getPassword()
					email := getEmail()
					if emailcheck(url, email, username, password) {

						DM(url)
						return
					} else {
						fmt.Println(Redtext.Render("\n Failed to Register"))
						time.Sleep(1 * time.Second)
					}

				case "3":
					email := getEmail()
					if forgetpass(url, email) {
						fmt.Println(greentext.Render("\n Password Reset! Login Now"))
						time.Sleep(2 * time.Second)
					}

				case "0":
					cls()
					fmt.Println("Bye Bye! Have a Good Day <3")
					os.Exit(0)

				default:
					fmt.Println("Invalid Option")
					time.Sleep(1 * time.Second)
				}
			}
		} else {
			// We have a token saved; check if it's still good on the server
			if tokenchekcing(url) {
				DM(url)
				return
			} else {
				// Server rejected the token (Session expired or bad data)
				fmt.Println(Redtext.Render("! Session expired. Clearing credentials..."))
				mytoken = ""
				myuser = ""
				os.Remove(".env")
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}
}

var Reqlist []string

func main() {

	baseURL = "http://localhost:4040"

	// Try to load. If it fails (file deleted), mytoken/myuser will stay empty ""
	godotenv.Load(".env")

	mytoken = os.Getenv("token")
	myuser = os.Getenv("user")

	// Start the menu
	manue(baseURL)
}
