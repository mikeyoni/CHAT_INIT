package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
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

func main() {

	fmt.Printf("hllo \n")

	otp := generateOTP()

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
