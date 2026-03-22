package main

import (
	"bytes"
	// "encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "golang.org/x/tools/go/analysis/checker"
	// "net/url"
	// "net/http"
)

// this is login function it do post login info in a url in json form package

func login(url string, username string, password string) {

	data := map[string]string{

		"username": username,
		"password": password,
	}

	jsondata, _ := json.Marshal(data)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {

		fmt.Printf("\n Faild to sent post req : ", err)
		return
	} else {

		fmt.Printf(" \n sucessfully sented : ", resp.Status)
	}

	defer resp.Body.Close()

	bodybytes, _ := io.ReadAll(resp.Body)
	massage := string(bodybytes)

	fmt.Printf(massage)

}

// this is register function this do post register info to the srever in json form

func emailcheck(url string, email string, username string, password string) {

	data := map[string]string{

		"email":    email,
		"username": username,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {

		fmt.Printf(" \n faild to marshel json data : ", err)
		return
	}
	resp, err := http.Post(url, "application/json-data", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Printf(" \n failt to post register info : ", err)
		return
	} else {

		fmt.Printf(" \n successfuly posted register data : ", err)

	}

	defer resp.Body.Close()

	newbytes, _ := io.ReadAll(resp.Body)

	massage := string(newbytes)

	if massage == "success" {
		register("http://localhost:4040/confarmregister", email, username, password)
	}

}

func register(url string, email string, username string, password string) {

	var userinput string
	fmt.Printf(" \n Enter the otp : ")
	fmt.Scan(&userinput)

	data := map[string]string{

		"email":    email,
		"username": username,
		"password": password,
		"otp":      userinput,
	}

	jsondata, err := json.Marshal(data)
	if err != nil {

		fmt.Printf(" \n faild to marshel json data : ", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsondata))

	if err != nil {
		fmt.Printf(" \n failt to post register info : ", err)
		return
	} else {

		fmt.Printf(" \n successfuly posted register data : ", err)

	}

	defer resp.Body.Close()

	newbystes, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n [SERVER RESPONSE]: %s\n", string(newbystes))

}

// forget password

func forgetpass(url string, email string) {

	data := map[string]string{

		"email": email,
	}

	jsondata, err := json.Marshal(data)

	if err != nil {

		fmt.Printf(" \n faild to marshal json data : %+s ", err)
		return
	}

	resp, err := http.Post(url, "/application/paasforget", bytes.NewBuffer(jsondata))

	if err != nil {

		fmt.Printf(" \n failt to sent the email to server : %+s ", err)
		return
	} else {
		fmt.Printf(" \n success fully sented forget password email : %+s", resp.Status)

	}

	defer resp.Body.Close()

}

func main() {

	// login("http://localhost:4040/login", "mikey", "mikey")
	emailcheck("http://localhost:4040/signup", "uimikey78@gmail.com", "dra34ken", "paf32453")
	// forgetpass("http://localhost:4040/forgetpass", "mda35345345@gmail.com")

}
