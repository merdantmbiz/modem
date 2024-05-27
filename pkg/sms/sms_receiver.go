package sms

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/warthog618/modem/at"
	pb "github.com/warthog618/modem/gen"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	//	rpc "github.com/warthog618/modem/pkg/grpc"
	"github.com/warthog618/modem/pkg/jwt"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
)

func StartSMSReciever(cfg *config.Config) error {

	dev := flag.String("d", cfg.MODEM.PORT, "path to modem device")
	baud := flag.Int("b", 115200, "baud rate")
	//	period := flag.Duration("p", 10*time.Minute, "period to wait")
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

	g := gsm.New(at.New(mio))

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
			go notifyAuth(msg.Number)
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
		// case <-time.After(*period):
		// 	log.Println("exiting...")
		// 	return nil
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

func notifyAuth(clientId string) {
	// Set up a connection to the server.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOTPServiceClient(conn)

	r, err := c.PassOTP(ctx, &pb.OTPRequest{ClientId: clientId})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// Print the response from the server
	log.Printf("Greeting: %s", r.ClientId)
}
