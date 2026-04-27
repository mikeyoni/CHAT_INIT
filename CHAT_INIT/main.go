package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var (
	upgrader = websocket.Upgrader{ReadBufferSize: 2040,
		WriteBufferSize: 2040, CheckOrigin: func(r *http.Request) bool { return true }}

	client   = make(map[string]*websocket.Conn)
	clientMu sync.Mutex
)

func Hashpassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func Compareshass(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// structure of the json data base

type user struct {
	UserName   string   `json:"username"`
	Email      string   `json:"email"`
	Token      string   `json:"token"`
	Password   string   `json:"passowrd"`
	Friendlist []string `json:"friendlist"`
	Reqestlist []string `json:"reqestlist"`
}

// otp genaretor

func generateOTP() string {

	b := make([]byte, 2)

	rand.Read(b)

	newNumber := binary.BigEndian.Uint16(b)

	return fmt.Sprintf("%04d", newNumber%10000)

}

// valid email

func isEmailValid(e string) bool {
	// This is the standard Regex for email addresses
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	// Convert to lowercase first (Fedora/Linux style)
	e = strings.ToLower(strings.TrimSpace(e))

	return emailRegex.MatchString(e)
}

// otp save

type record struct {
	Code         string
	Codesavetime time.Time
}

var Otpstorage = make(map[string]record)

// funcitn that save and delate otp automaticly

func otpsaveanddelate(otp string, user string) record {

	newotp := record{
		Code:         otp,
		Codesavetime: time.Now(),
	}

	Otpstorage[user] = newotp

	go func() {

		time.Sleep(2 * time.Minute)
		delete(Otpstorage, user)
		fmt.Printf(" \n The otp is expired ! ")

	}()

	return newotp

}

// otp sender

func init() {
	// Try to load .env if it exists, but don't crash if it's missing.
	_ = godotenv.Load()
}

// setupCredentials checks if GMAIL_USER and GMAIL_PASS are configured.
// If not, it prompts interactively and saves them to .env for later runs.
func setupCredentials() {
	promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00F5D4")).Bold(true)
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3cff00")).Bold(true)
	warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9900")).Bold(true)

	gmailUser := os.Getenv("GMAIL_USER")
	gmailPass := os.Getenv("GMAIL_PASS")

	if gmailUser != "" && gmailPass != "" {
		return
	}

	fmt.Println(warnStyle.Render("\n SMTP credentials not found! The server needs a Gmail account to send OTP emails."))
	fmt.Println(promptStyle.Render(" You only need to do this once; credentials will be saved locally.\n"))

	reader := bufio.NewReader(os.Stdin)

	if gmailUser == "" {
		fmt.Print(promptStyle.Render(" Enter your Gmail address: "))
		input, _ := reader.ReadString('\n')
		gmailUser = strings.TrimSpace(input)
	}

	if gmailPass == "" {
		fmt.Print(promptStyle.Render(" Enter your Google App Password: "))
		input, _ := reader.ReadString('\n')
		gmailPass = strings.TrimSpace(input)
	}

	os.Setenv("GMAIL_USER", gmailUser)
	os.Setenv("GMAIL_PASS", gmailPass)

	envContent := fmt.Sprintf("GMAIL_USER=%s\nGMAIL_PASS=%s\n", gmailUser, gmailPass)
	if err := os.WriteFile(".env", []byte(envContent), 0600); err != nil {
		log.Printf("warning: could not save credentials to .env: %v", err)
		log.Println("credentials are loaded for this session but won't persist")
		return
	}

	fmt.Println(successStyle.Render("\n Credentials saved to .env.\n"))
}

func sentOPTEmail(targetEmail string, otp string) error {

	validemail := isEmailValid(targetEmail)

	if !validemail {
		fmt.Printf(" \n invalid email \n ")
		return nil
	}

	form := os.Getenv("GMAIL_USER")
	password := os.Getenv("GMAIL_PASS")

	// setup smtp server settings for gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// create massage

	subject := "Subject: Pirate King Verification Code\r\n"

	fromHeader := fmt.Sprintf("From: %s\r\n", form)
	toHeader := fmt.Sprintf("To: %s\r\n", targetEmail)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf("<html><body><h1>Your Code: %s</h1><p>Expires in 60 seconds.</p></body></html>", otp)
	message := []byte(subject + fromHeader + toHeader + mime + body)
	// authentication
	auth := smtp.PlainAuth("", form, password, smtpHost)

	// sending the actual emaiil

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, form, []string{targetEmail}, message)

	if err != nil {
		return err
	}
	return nil

}

// json file read

func jsonreade() ([]user, error) {

	var jsondata []user
	bytes, err := os.ReadFile("database.json")
	if err != nil {

		return nil, err

	}

	err = json.Unmarshal(bytes, &jsondata)
	if err != nil {

		return nil, err
	}

	return jsondata, nil

}

// data save to the data base

func jsondatasave(email string, username string, password string) ([]user, error) {

	var jsondata []user
	bytes, err := os.ReadFile("database.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &jsondata)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%+v",jsondata)
	hashpass, _ := Hashpassword(password)
	token := dcstyletokengen(username, email, hashpass)

	newuser := user{
		Email:      email,
		UserName:   username,
		Password:   hashpass,
		Token:      token,
		Friendlist: []string{},
		Reqestlist: []string{},
	}

	jsondata = append(jsondata, newuser)

	updateadata, err := json.MarshalIndent(jsondata, "", "\t")

	if err != nil {
		return nil, err
	}

	err = os.WriteFile("database.json", updateadata, 0644)

	if err != nil {
		return nil, err
	}

	return jsondata, nil

}

// this funcion add "addinfriend" in "theuser" friendlist
// and also remove "add friendin" form "theuser" requested list
func addfriend(theuser string, addinfriend string) {

	user, err := jsonreade()

	if err != nil {
		return
	}

	isdone := false
	me := false
	u := false
	for i := range user {

		if user[i].UserName == theuser {

			user[i].Friendlist = append(user[i].Friendlist, addinfriend)

			var cleanReqs []string
			for _, r := range user[i].Reqestlist {
				if r != addinfriend {
					cleanReqs = append(cleanReqs, r)
				}
			}

			user[i].Reqestlist = cleanReqs
			me = true

		}

		if user[i].UserName == addinfriend {

			user[i].Friendlist = append(user[i].Friendlist, theuser)
			u = true

		}

		if me && u {
			isdone = saveuserdata(user)
			break
		}
	}

	if isdone {
		fmt.Printf(" \n Successfullty added \n")
	}

	fmt.Printf("\n user not found is request list ")
}

// this return u the friendlist slice
func friendlistview(user string) []string {

	jsondata, err := jsonreade()

	if err != nil {
		return nil
	}

	for i := range jsondata {
		if jsondata[i].UserName == user {

			for _, freinds := range jsondata[i].Friendlist {
				fmt.Printf("%+v \n", freinds)
			}

			return jsondata[i].Friendlist
		}
	}

	fmt.Printf(" \n User not found \n ")
	return nil
}

func onlinefriendlist(username string) []string {
	friendlist := friendlistview(username)
	if len(friendlist) == 0 {
		return nil
	}

	clientMu.Lock()
	defer clientMu.Unlock()

	var online []string
	for _, friend := range friendlist {
		if _, ok := client[friend]; ok {
			online = append(online, friend)
		}
	}

	return online
}

// this return u the requestedlist slice
func requestedlistview(user string) []string {

	jsondata, err := jsonreade()

	if err != nil {
		return nil
	}

	for i := range jsondata {
		if jsondata[i].UserName == user {
			fmt.Println(" Done sending friend req ")
			return jsondata[i].Reqestlist
		}
	}

	fmt.Printf(" \n User not found \n ")
	return nil
}

// this save the hole userd data that u modify by jsonreader
func saveuserdata(savealluser []user) bool {

	updateadata, err := json.MarshalIndent(savealluser, "", "\t")

	if err != nil {
		return false
	}

	err = os.WriteFile("database.json", updateadata, 0644)

	if err != nil {
		return false
	}

	return true

}

// this add "theuser" to the "requesteduser" request list
func requestfirendadd(theuser string, requesteduser string) string {

	if theuser == requesteduser {
		return "invalid"
	}

	userdata, _ := jsonreade()
	isdone := false
	for i, userindo := range userdata {

		if userindo.UserName == requesteduser {

			if slices.Contains(userdata[i].Friendlist, theuser) {

				return "alrady freinds!"
			}

			if slices.Contains(userdata[i].Reqestlist, theuser) {
				return "requested alradysent"
			}

			userdata[i].Reqestlist = append(userdata[i].Reqestlist, theuser)
			isdone = saveuserdata(userdata)

			break

		}

	}

	if isdone {
		fmt.Printf(" \n success fully sented the request ")
		return "done"
	}

	return "user not found "

}

// this remove the "theuser" form the "toremoveform" user request list
func removeformreq(theusername string, toremoveform string) {

	if theusername == toremoveform {
		return
	}

	userdata, _ := jsonreade()

	for i, stringuser := range userdata {

		if stringuser.UserName == toremoveform {

			var updatereqlist []string

			for _, reqfreind := range userdata[i].Reqestlist {

				if reqfreind != theusername {

					updatereqlist = append(updatereqlist, reqfreind)
				}

			}

			userdata[i].Reqestlist = updatereqlist

			saveuserdata(userdata)
			fmt.Printf("\n removed form reqlist ")
			break
		}
	}
}

// me: The current logged-in user (e.g., "Luffy")
// target: The friend you want to delete (e.g., "Zoro")
func removeformfriendlist(me string, target string) {
	// is it ok i gess yes its ok

	if me == target {
		return
	}

	userdata, _ := jsonreade()

	for i, sclises := range userdata {

		if sclises.UserName == me {

			var updatefriendlist []string

			for _, friends := range userdata[i].Friendlist {

				if friends != target {

					updatefriendlist = append(updatefriendlist, friends)

				}

			}

			userdata[i].Friendlist = updatefriendlist
			saveuserdata(userdata)

			fmt.Printf("\n Removed form the friendlist \n")
			break

		}
	}
}

// the login funciton
func login(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/login" {

		fmt.Printf(" login page faild ! : %v ", http.StatusNotFound)
		return
	}

	var incamingdata struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(r.Body).Decode(&incamingdata)

	// loading the data base

	jsondata, _ := jsonreade()

	userfound := false
	passwordfund := false

	// checking the user name is in data base or not

	for _, user := range jsondata {
		if user.UserName == incamingdata.Username {
			userfound = true
			passwordfund = Compareshass(incamingdata.Password, user.Password)
			if passwordfund {
				fmt.Fprintf(w, "success:token:%s", user.Token)
				return // FIX: Stop execution here
			} else {
				fmt.Fprintf(w, "not:Invalid Password")
				return // FIX: Stop execution here
			}
		}
	}

	if !userfound {
		fmt.Fprintf(w, "not:User Not Found") // FIX: Use 'not:' prefix
		return
	}
	fmt.Printf(" %+v ", incamingdata)

}

// this is register function
func checkup(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/signup" {

		fmt.Printf(" signup page faild : %v ", http.StatusNotFound)
		return
	}

	var incamingdats struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	json.NewDecoder(r.Body).Decode(&incamingdats)

	user, _ := jsonreade()

	for _, axistinguser := range user {

		if axistinguser.UserName == incamingdats.Username {

			fmt.Fprintf(w, " Username alrady exist ")
			return
		}
		if axistinguser.Email == incamingdats.Email {
			fmt.Fprintf(w, " email is alrady exist ")
			return
		}
	}

	otp := generateOTP()

	email := incamingdats.Email

	err := sentOPTEmail(email, otp)

	if err != nil {

		fmt.Fprintf(w, "nosuccess")
		return
	}

	otpsaveanddelate(otp, incamingdats.Username)

	fmt.Fprintf(w, "success")

}

func register(w http.ResponseWriter, r *http.Request) {
	// 1. Always check the Method (Safety first!)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var incoming struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
		Otp      string `json:"otp"`
	}

	// 2. Check if the Decoding actually worked
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		fmt.Printf("JSON Decode Error: %v\n", err)
		return
	}

	// 3. Debug Print: See what the Ryzen 5600G is actually seeing
	fmt.Printf("Received: User=%s, Email=%s, OTP=%s\n", incoming.Username, incoming.Email, incoming.Otp)

	// ... (Your existing user-exists loop) ...

	savedata := Otpstorage[incoming.Username]

	// 4. THE FIX: Only save if the OTP matches!
	if savedata.Code != "" && savedata.Code == incoming.Otp {
		fmt.Printf("\n[SUCCESS] OTP Matched for %s\n", incoming.Username)

		_, err := jsondatasave(incoming.Email, incoming.Username, incoming.Password)
		if err != nil {
			fmt.Printf("Save Error: %v\n", err)
			fmt.Fprintf(w, "save_error")
			return
		}

		usersdata, _ := jsonreade()

		for _, userdatas := range usersdata {

			if userdatas.UserName == incoming.Username {

				fmt.Fprintf(w, "success:token:%s", userdatas.Token)

			}
		}

		delete(Otpstorage, incoming.Username) // Clean up the RAM
	} else {
		fmt.Printf("[FAIL] OTP Mismatch. Got: %s, Expected: %s\n", incoming.Otp, savedata.Code)
		fmt.Fprintf(w, "top not match")
	}
}

func forgetpass(w http.ResponseWriter, r *http.Request) {
	// 1. CHECK FOR OTP AND NEW PASS FIRST
	incamingotp := r.URL.Query().Get("otp")
	incamingemail := r.URL.Query().Get("user")
	incaminnewpass := r.URL.Query().Get("new")

	if incamingotp != "" && incamingemail != "" {
		savedotp := Otpstorage[incamingemail]

		if savedotp.Code == incamingotp {
			if incaminnewpass != "" {
				// SUCCESS: Update the password in your JSON here!
				fmt.Printf("\n SUCCESS: Changing password for %s to %s", incamingemail, incaminnewpass)
				userdatas, _ := jsonreade()

				for i := range userdatas {

					if userdatas[i].Email == incamingemail {
						newpass, _ := Hashpassword(incaminnewpass)
						userdatas[i].Password = newpass
						isdone := saveuserdata(userdatas)

						if isdone {
							fmt.Print(" successfully reset password ")
						}
						break

					}

				}
				fmt.Fprintf(w, "success:Password has been reset")
				return // STOP HERE! Don't send another OTP.
			}

		}
		fmt.Fprintf(w, "error:Invalid or expired OTP")
		return
	}

	// 2. IF NO OTP, HANDLE EMAIL CAPTURE (SENDING OTP)
	var incomingdata struct {
		Email string `json:"email"`
	}

	// Decode the JSON email
	err := json.NewDecoder(r.Body).Decode(&incomingdata)
	if err != nil || incomingdata.Email == "" {
		return // Silently fail if no email provided
	}

	// Check database for email
	jsondata, _ := jsonreade()
	for i := range jsondata {
		if jsondata[i].Email == incomingdata.Email {
			otp := generateOTP()
			sentOPTEmail(incomingdata.Email, otp)
			otpsaveanddelate(otp, incomingdata.Email)
			fmt.Printf("\n OTP SENTED TO THE EMAIL : %v ", incomingdata.Email)
			fmt.Fprintf(w, "done:OTP sent to your email")
			return // STOP HERE!
		}
	}
	fmt.Fprintf(w, "not:Email not found")
}

func dcstyletokengen(username string, email string, password string) string {

	hashpass := password[len(password)-10:]

	rowtoken := fmt.Sprintf("%s:%s:%s", username, email, hashpass)

	return base64.StdEncoding.EncodeToString([]byte(rowtoken))

}

func tokencheck(token string) bool {

	jsondata, _ := jsonreade()

	for i := range jsondata {

		if jsondata[i].Token == token {
			return true
		}
	}

	return false
}

func checktoken(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	isok := tokencheck(token)

	if isok {
		fmt.Printf(" \n Token is valid ! ")
		fmt.Fprintf(w, "ok")
	} else if !isok {

		fmt.Fprintf(w, "no")
	}
}

func chating(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("user")
	token := r.URL.Query().Get("token")

	if username == "" || token == "" {

		fmt.Fprintf(w, "\n user info misisng \n ")
	}

	istokenvalid := tokencheck(token)

	if istokenvalid {

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			return
		}

		clientMu.Lock()
		client[username] = conn
		clientMu.Unlock()

		fmt.Printf(" \n %s has cannected the server : \n", username)

		defer func() {

			clientMu.Lock()
			fmt.Printf(" %s has left form canection ", username)
			delete(client, username)
			clientMu.Unlock()

		}()

		for {

			_, p, err := conn.ReadMessage()

			if err != nil {
				break
			}

			row := string(p)
			if strings.HasPrefix(row, "tusr:") {
				parts := strings.SplitN(row, ":", 4)

				if len(parts) >= 0 && parts[3] != "" {

					message := fmt.Sprintf(" [ %s ] : %s ", username, parts[3])

					var isfriend = false

					friendlist := friendlistview(parts[1])

					isfriend = slices.Contains(friendlist, username)

					if isfriend {

						success := sendto(parts[1], message)

						if !success {

							offlineMsg := fmt.Sprintf("System: %s is currently offline. Your message was not delivered.", parts[1])
							conn.WriteMessage(websocket.TextMessage, []byte(offlineMsg))
						}
					} else {

						conn.WriteMessage(websocket.TextMessage, []byte("System: You are not in this user's friend list."))
					}

				}
			}

		}

	}

}

func sendto(tusr string, msg string) bool {
	clientMu.Lock()
	defer clientMu.Unlock()

	target, ok := client[tusr]

	if !ok {

		fmt.Printf("\n User %s is offline \n", tusr)
		return false
	}

	err := target.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		fmt.Printf("\n Connection problem: %v \n", err)
		target.Close()
		delete(client, tusr)
		return false
	}

	return true
}

func todo(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("user")
	token := r.URL.Query().Get("token")
	action := r.URL.Query().Get("act")
	targeteduser := r.URL.Query().Get("tar")

	if username == "" || token == "" {
		return
	}

	userdata, _ := jsonreade()

	authorized := false
	targetfound := false
	userfound := false
	for i := range userdata {

		if userdata[i].UserName == targeteduser {
			targetfound = true
		}
		if userdata[i].UserName == username && userdata[i].Token == token {

			switch action {

			case "sentfreq":

				fmt.Printf(" \n its sentdreq to %v ", targeteduser)
				requestfirendadd(username, targeteduser)

			case "rejectfreq":

				fmt.Printf(" \n its rejectreq %v form %v req list ", targeteduser, username)
				removeformreq(targeteduser, username)

			case "acceptfreq":

				fmt.Printf(" \n its acceptreq the freind req of %v ", targeteduser)
				addfriend(username, targeteduser)

			case "delatfre":

				fmt.Printf(" \n its delatereq %v form %v list", targeteduser, username)
				removeformfriendlist(username, targeteduser)
				removeformfriendlist(targeteduser, username)

			default:
				break

			}

			userfound = true

		}

		if targetfound && userfound {

			authorized = true
		}

	}

	if !authorized {

		a := ""
		b := ""

		if targetfound && !userfound {
			a = "target found but user missing"
		} else if !targetfound && userfound {
			a = "User found but targeted user missing"
		} else if !targetfound && !userfound {
			b = "target and user missing"
		}

		msg := fmt.Sprintf("Invalid Token or User . %v , %v ", a, b)

		fmt.Fprintf(w, " %v ", msg)

	}

}

func flistgiver(w http.ResponseWriter, r *http.Request) {

	user := r.URL.Query().Get("user")
	token := r.URL.Query().Get("token")

	isvalid := tokencheck(token)

	var freindlist []string

	if isvalid {
		freindlist = friendlistview(user)

		if freindlist == nil {
			fmt.Printf("\n Frliend list is empty ")
			return
		}

		jsondata, err := json.Marshal(freindlist)

		if err != nil {
			fmt.Printf("\n Data marsheling proble : %v ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsondata)
	}

}

func GIveReqlist(w http.ResponseWriter, r *http.Request) {

	user := r.URL.Query().Get("user")
	token := r.URL.Query().Get("token")

	isvalid := tokencheck(token)

	var Rreindlist []string

	if isvalid {
		Rreindlist = requestedlistview(user)

		if Rreindlist == nil {
			fmt.Printf("\n Frliend list is empty ")
			return
		}

		jsondata, err := json.Marshal(Rreindlist)

		if err != nil {
			fmt.Printf("\n Data marsheling proble : %v ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsondata)
	}

}

func GiveOnlineFriends(w http.ResponseWriter, r *http.Request) {

	user := r.URL.Query().Get("user")
	token := r.URL.Query().Get("token")

	if user == "" || token == "" {
		return
	}

	if !tokencheck(token) {
		return
	}

	onlinefriends := onlinefriendlist(user)
	if onlinefriends == nil {
		onlinefriends = []string{}
	}

	jsondata, err := json.Marshal(onlinefriends)
	if err != nil {
		fmt.Printf("\n Data marsheling proble : %v ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsondata)
}

func main() {

	// addfriend("dra34ken" , "mikey3")
	// jsondata, err := jsonreade()
	// if err != nil {
	// 	fmt.Printf(" Faild to Load json file : " , err)
	// 	return
	// }

	// fmt.Printf("\n %+v \n", jsondata )

	// this are only style of text colors in terminal

	var text = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Italic(true)
	var oktext = lipgloss.NewStyle().Foreground(lipgloss.Color("#3cff00")).Bold(true)

	// this are the main functin that canect logical funcion to the url quaris

	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", checkup)
	http.HandleFunc("/confarmregister", register)
	http.HandleFunc("/forgetpass", forgetpass)
	http.HandleFunc("/chat", chating)
	http.HandleFunc("/do", todo)
	http.HandleFunc("/checking", checktoken)
	http.HandleFunc("/viewflist", flistgiver)
	http.HandleFunc("/viewReqlist", GIveReqlist)
	http.HandleFunc("/viewonlinefriends", GiveOnlineFriends)

	// http.HandleFunc("/chat-init" , long)

	setupCredentials()

	// server startig indicatin
	fmt.Printf("%s", text.Render("\n\n The server is starting... \n"))

	// this the last stage and the biginig of the server the code will be stay here in a loop and the server will be start
	fmt.Printf("%s", oktext.Render("\n ( CTRL + C ) TO STOP THE SRVER <3 \n\n"))

	if err := http.ListenAndServe(":4040", nil); err != nil {

		fmt.Println(" local Fort 4040 is unusebal right this moment ! ", err)

	}

}
