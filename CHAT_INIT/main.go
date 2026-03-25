package main

import (
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
	UserName string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"passowrd"`
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func sentOPTEmail(targetEmail string, otp string) error {

	validemail := isEmailValid(targetEmail)

	if !validemail {
		fmt.Printf(" \n invalid email \n ")
		return nil
	}

	form := "uimikey1@gmail.com"
	password := os.Getenv("pass")

	// setup smtp server settings for gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// create massage

	subject := "Subject: Pirate King Verification Code\r\n"

	fromHeader := "From: uimikey1@gmail.com\r\n"
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
		Email:    email,
		UserName: username,
		Password: hashpass,
		Token:    token,
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
				usertoke := user.Token
				fmt.Fprintf(w, "success:token:%s", usertoke)
			}

			break

		}

	}

	if !userfound {
		fmt.Fprintf(w, " user not found ")
	} else if !passwordfund {
		fmt.Fprintf(w, " passsowrd not found ")
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

	sentOPTEmail(email, otp)

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

// this is forget password email capture

func forgetpass(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/forgetpass" {
		fmt.Printf(" faild show forget passpage : %+v", http.StatusNotFound)
		return
	} else {

		fmt.Fprintf(w, " \n Forgetpassrod gape is active :D \n")
	}

	var imcamingdata struct {
		Email string `json:"email"`
	}

	json.NewDecoder(r.Body).Decode(&imcamingdata)

	fmt.Printf(" \n we got the email that forget passwrod : %+v ", imcamingdata)

}

func dcstyletokengen(username string, email string, password string) string {

	hashpass := password[len(password)-10:]

	rowtoken := fmt.Sprintf("%s:%s:%s", username, email, hashpass)

	return base64.StdEncoding.EncodeToString([]byte(rowtoken))

}

func chating(w http.ResponseWriter, r *http.Request) {

	username := r.URL.Query().Get("user")
	token := r.URL.Query().Get("token")

	if username == "" || token == "" {

		fmt.Fprintf(w, "\n user info misisng \n ")
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	clientMu.Lock()
	client[username] = conn
	clientMu.Unlock()

	fmt.Printf(" \n %s has cannected the server : \n" , username )

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
			parts := strings.SplitN(row, ":" , 4)
			
			if len(parts) >= 0 && parts[3] != "" {
				
				message := fmt.Sprintf(" [ %s ] : %s ", username, parts[3])
				sendto(parts[1], message)
			}
		}

	}

}

func sendto(tusr string, msg string) {

	clientMu.Lock()
	defer clientMu.Unlock()

	target, ok := client[tusr]

	if !ok {
		fmt.Printf("\n the user is offline \n ")
		return
	}

	err := target.WriteMessage(websocket.TextMessage, []byte(msg))

	if err != nil {
		fmt.Printf("\n cannection problem : %s ", &err)

		target.Close()
		delete(client, tusr)
	}

}

func main() {

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
	http.HandleFunc("/chat" , chating)

	// http.HandleFunc("/chat-init" , long)

	// server startig indicatin
	fmt.Printf("%s", text.Render("\n\n The server is starting... \n"))

	// this the last stage and the biginig of the server the code will be stay here in a loop and the server will be start
	fmt.Printf("%s", oktext.Render("\n ( CTRL + C ) TO STOP THE SRVER <3 \n\n"))

	if err := http.ListenAndServe(":4040", nil); err != nil {

		fmt.Println(" local Fort 4040 is unusebal right this moment ! ", err)

	}

}
