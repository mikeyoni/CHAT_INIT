package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/smtp"
	"time"
)

func generateOTP() string {

	b := make([]byte, 2)

	rand.Read(b)

	number := binary.BigEndian.Uint16(b)

	return fmt.Sprintf("%04d", number%10000)

}

type record struct {
	Code      string
	Createdat time.Time
}

var otpstorage = make(map[string]record)

// EMAIL SEINGING SYSTREAWM

func sentOPTEmail(targetEmail string, otp string) error {

	form := "uimikey1@gmail.com"
	password := "kygcvtcqcyjwxebl"

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

func main() {

	fmt.Printf("hllo \n")

	otp := generateOTP()

	sentOPTEmail("mda891526@gmail.com", otp)

	newentry := record{
		Code:      otp,
		Createdat: time.Now(),
	}

	otpstorage["mikey"] = newentry

	fmt.Printf(" %s \n", otp)

	fmt.Printf(" %s is the otp and its created on %v ", otpstorage["mikey"].Code, otpstorage["mikey"].Createdat)

	go func() {

		time.Sleep(15 * time.Second)
		delete(otpstorage, "mikey")
		fmt.Printf(" \n the stored otp is delated ! ")
	}()

	for {
		// 1. The "comma, ok" check: 'exists' is a boolean (true/false)
		record, exists := otpstorage["mikey"]

		if !exists {
			fmt.Println("\n[LOOP] The record is gone! Breaking out.")
			break
		}

		// 2. Just to show it's working
		fmt.Printf("\rChecking... OTP %s is still in RAM", record.Code)

		// 3. ESSENTIAL: Give the CPU a 1-second break
		time.Sleep(1 * time.Second)
	}

}
