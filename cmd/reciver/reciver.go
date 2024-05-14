package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "COM12", Baud: 9600, ReadTimeout: time.Second * 5}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	defer s.Close()

	buf := make([]byte, 128)
	for {
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		if n == 0 {
			continue
		}

		data := string(buf[:n])
		fmt.Printf("Received: %s\n", data)

		if data == "AT+CMGL=4\r" {
			fmt.Println("Reading SMS messages...")
			s.Write([]byte("AT+CMGL=\"ALL\"\r"))
			time.Sleep(time.Second * 2)
			sms, err := s.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("SMS messages: %s\n", string(buf[:sms]))
		}
	}
}
