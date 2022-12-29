package scanner

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func IcmpPing(dst string) {
	if dst == "" {
		log.Fatal("Please specify an IP Address!")
	}

	listener, err := icmp.ListenPacket("udp4", "10.0.0.208")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeExtendedEchoRequest,
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
	if _, err := listener.WriteTo(encoded, &net.IPAddr{IP: net.ParseIP(dst), Zone: "en0"}); err != nil {
		log.Fatal(err)
	}

	rb := make([]byte, 1500)
	n, peer, err := listener.ReadFrom(rb)
	if err != nil {
		log.Fatal(err)
	}
	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		log.Fatal(err)
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("pong received from %v", peer)
	default:
		log.Printf("got %+v; want echo reply", rm)
	}
}
