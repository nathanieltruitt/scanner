package scanner

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Response struct {
	Len int
	Address net.Addr
}


func PingScan(subnet string) {
	// Get my IP that is on this subnet
	// var onlineHosts []string
	myIp, err := localIP()
	if err != nil {
		log.Fatal(err)
	}

	ip := strings.TrimSuffix(strings.Split(subnet, "/")[0], "0")
	mask, err := strconv.Atoi(strings.Split(subnet, "/")[1]) 
	if err != nil {
		log.Fatal(err)
	}

	if mask == 24 {
		for i := 1; i < 254; i++ {
			dst := ip + strconv.Itoa(i)
			if icmpPing(myIp.String(), dst) {
				fmt.Println(dst)
				// onlineHosts = append(onlineHosts, dst)
			}
		}
	}
}


func getResponse(ctx context.Context, responseChan chan Response, packet []byte, listener *icmp.PacketConn) {
	n, peer, err := listener.ReadFrom(packet)
	if err != nil {
		listener.Close()
	}
	responseChan <- Response{Len: n, Address: peer}
}


func icmpPing(src string, dst string) bool {
	// sends an icmp echo request and listens for 1 second before timing out.
  // create context
	var connected bool

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond * 300)
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

	select {
	case <- responseChan:
		connected = true
	case <- ctx.Done():
		connected = false
	}

	// run if there is a response length
	// if response.Len != 0 {
	// 	rm, err := icmp.ParseMessage(1, rb[:response.Len])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	switch rm.Type {
	// 	case ipv4.ICMPTypeEchoReply:
	// 		log.Printf("pong received from %v", response.Address)
	// 	default:
	// 		log.Printf("got %+v; want echo reply", rm)
	// 	}
	// }

	return connected;
}
