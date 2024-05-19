package main

import (
	"log"

	"github.com/warthog618/modem/pkg/config"
	"github.com/warthog618/modem/pkg/sms"
)

var version = "undefined"

func main() {

	err := config.InitTomlConf("config", "./pkg/config")

	if err != nil {
		log.Println(err)
	}

	err = sms.StartSMSReciever(&config.TomlConf)

	if err != nil {
		log.Println(err)
	}
}
