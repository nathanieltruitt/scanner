package scanner

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Response struct {
	Len int
	Address net.Addr
}

func getResponse(ctx context.Context, responseChan chan Response, packet []byte, listener *icmp.PacketConn) {
	n, peer, err := listener.ReadFrom(packet)
	if err != nil {
		listener.Close()
		log.Fatal(err)
	}
	responseChan <- Response{Len: n, Address: peer}
}

func IcmpPing(dst string) {
  // create context
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
	defer cancel()

	// response chan
	responseChan := make(chan Response)

	if dst == "" {
		log.Fatal("Please specify an IP Address!")
	}

	// TODO: needs to dynamically retreive the correct interface IP.
	listener, err := icmp.ListenPacket("udp4", "10.0.0.208")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	encoded, err := message.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = listener.WriteTo(encoded, &net.UDPAddr{IP: net.ParseIP(dst)}); 
	if err != nil {
		log.Fatal(err)
	}
	rb := make([]byte, 1500)

	go getResponse(ctx, responseChan, rb, listener)

	var response Response
	select {
	case response = <- responseChan:
		fmt.Println("Response acquired")
	case <- ctx.Done():
		log.Printf("ping to %v timed out", dst)
	}

	// run if there is a response length
	if response.Len != 0 {
		rm, err := icmp.ParseMessage(1, rb[:response.Len])
		if err != nil {
			log.Fatal(err)
		}
		switch rm.Type {
		case ipv4.ICMPTypeEchoReply:
			log.Printf("pong received from %v", response.Address)
		default:
			log.Printf("got %+v; want echo reply", rm)
		}
	}
}
