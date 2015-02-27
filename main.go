package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/cloudfoundry/noaa"
	"github.com/cloudfoundry/noaa/events"
)

var (
	dopplerAddress = flag.String("dopplerAddress", "ws://192.168.6.19.xip.io:51200", "doppler agent address")
)

var authToken = os.Getenv("CF_ACCESS_TOKEN")

const firehoseSubscriptionId = "firehose-a"

func main() {
	flag.Parse()
	connection := noaa.NewConsumer(*dopplerAddress, &tls.Config{InsecureSkipVerify: true}, nil)
	connection.SetDebugPrinter(ConsoleDebugPrinter{})

	fmt.Println("===== Streaming Firehose (will only succeed if you have admin credentials)")

	msgChan := make(chan *events.Envelope)
	go func() {
		defer close(msgChan)
		errorChan := make(chan error)
		go connection.Firehose(firehoseSubscriptionId, authToken, msgChan, errorChan, nil)

		for err := range errorChan {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
	}()

	for msg := range msgChan {
		fmt.Printf("%v \n", msg)
	}
}

type ConsoleDebugPrinter struct{}

func (c ConsoleDebugPrinter) Print(title, dump string) {
	println(title)
	println(dump)
}
