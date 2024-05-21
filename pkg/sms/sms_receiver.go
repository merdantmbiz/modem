package sms

import (
	"flag"
	"io"
	"log"
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/pkg/config"
	rpc "github.com/warthog618/modem/pkg/grpc"
	"github.com/warthog618/modem/pkg/jwt"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
)

func StartSMSReciever(cfg *config.Config, grpcServer *rpc.Server) error {

	dev := flag.String("d", cfg.MODEM.PORT, "path to modem device")
	baud := flag.Int("b", 115200, "baud rate")
	period := flag.Duration("p", 10*time.Minute, "period to wait")
	timeout := flag.Duration("t", 400*time.Millisecond, "command timeout period")
	verbose := flag.Bool("v", false, "log modem interactions")
	hex := flag.Bool("x", false, "hex dump modem responses")

	flag.Parse()

	m, err := serial.New(serial.WithPort(*dev), serial.WithBaud(*baud))

	if err != nil {
		return err
	}

	defer m.Close()

	var mio io.ReadWriter = m

	if *hex {
		mio = trace.New(m, trace.WithReadFormat("r: %v"))
	} else if *verbose {
		mio = trace.New(m)
	}

	g := gsm.New(at.New(mio, at.WithTimeout(*timeout)))

	err = g.Init()

	if err != nil {
		return err
	}

	go pollSignalQuality(g, timeout)

	err = g.StartMessageRx(
		//message arrived
		func(msg gsm.Message) {

			//generate jwt token for phone number
			token, err := jwt.GenerateJWT(msg.Number)

			if err != nil {
				log.Println(err)
				return
			}
			go grpcServer.SendTokenToClient(msg.Number, token)
			log.Printf("%s: %s\n", msg.Number, token)
		},
		func(err error) {
			log.Printf("err: %v\n", err)
		})
	if err != nil {
		return err
	}
	defer g.StopMessageRx()

	for {
		select {
		case <-time.After(*period):
			log.Println("exiting...")
			return nil
		case <-g.Closed():
			log.Fatal("modem closed, exiting...")
		}
	}
}

// pollSignalQuality polls the modem to read signal quality every minute.
//
// This is run in parallel to SMS reception to demonstrate separate goroutines
// interacting with the modem.
func pollSignalQuality(g *gsm.GSM, timeout *time.Duration) {
	for {
		select {
		case <-time.After(time.Minute):
			i, err := g.Command("+CSQ")
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Signal quality: %v\n", i)
			}
		case <-g.Closed():
			return
		}
	}
}
