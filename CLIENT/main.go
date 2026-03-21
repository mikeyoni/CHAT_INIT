package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func register(url string, email string, username string, password string) {

	data := map[string]string{

		"email":    email,
		"username": username,
		"password": password,
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

	bodybyte, _ := io.ReadAll(resp.Body)

	masssage := string(bodybyte)

	if masssage == "otpsented" {
		fmt.Printf("\n otpsented \n")
		var userinput string
		fmt.Scan(&userinput)

		otpdata := map[string]string{

			"otp": userinput,
		}

		jsonopt, _ := json.Marshal(otpdata)

		otpreplay, err := http.Post("http://localhost:4040/otp", "/appicaltion/data", bytes.NewBuffer(jsonopt))

		if err != nil {
			fmt.Printf(" faild to sent otp ", err)
		}

		defer otpreplay.Body.Close()

		replaybytes, _ := io.ReadAll(otpreplay.Body)
		replayofotp := string(replaybytes)

		if replayofotp == "success" {
			fmt.Printf("successfully register ")
		} else {
			fmt.Printf(" top not match")
			return
		}

	}

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

	login("http://localhost:4040/login", "mikey", "mikey")
	register("http://localhost:4040/signup", "mda234343@gmail.com", "dra34ken", "paf32453")
	// forgetpass("http://localhost:4040/forgetpass", "mda35345345@gmail.com")

}
