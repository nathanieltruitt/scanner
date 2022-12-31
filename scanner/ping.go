package scanner

import (
	"runtime"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Response struct {
	Len int
	Address net.Addr
}

type scanData struct {
	sync.RWMutex
	OnlineHosts []string
	wg sync.WaitGroup
}


func PingScan(subnet string, start int, end int) {
	// Get my IP that is on this subnet
	data := scanData{}
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
		for i := start; i < end; i++ {
			data.wg.Add(1)
			dst := ip + strconv.Itoa(i)
			go icmpPing(myIp.String(), dst, &data)
		}
		data.wg.Wait()
		fmt.Println(data.OnlineHosts)
	}
}


func getResponse(ctx context.Context, responseChan chan bool, packet []byte, listener *icmp.PacketConn) {
	n, _, _ := listener.ReadFrom(packet)


	msg, _ := icmp.ParseMessage(1, packet[:n])

	// return response if the length is greater than 0
	if msg != nil {
		body, err := msg.Body.Marshal(1)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(string(body))
		if strings.Contains(string(body), "HELLO-R-U-THERE") {
			responseChan <- true
		} else {
			responseChan <- false
		}
	} else {
		responseChan <- false
	}
}


func icmpPing(src string, dst string, data *scanData) {
	// sends an icmp echo request and listens for 1 second before timing out.
  // create context
	data.Lock()
	defer data.Unlock()
	defer data.wg.Done()
	connected := false
	protocol := "udp4"

	// retreive local IP
	ip, err := localIP()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond * 100)
	defer cancel()

	// response chan
	responseChan := make(chan bool)

	if dst == "" {
		log.Fatal("Please specify an IP Address!")
	}

	// check if the OS is windows
	if runtime.GOOS == "windows" {
		protocol = "ip4:icmp"
	}

	// TODO: needs to dynamically retreive the correct interface IP.
	listener, err := icmp.ListenPacket(protocol, ip.String())
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
	listener.WriteTo(encoded, &net.UDPAddr{IP: net.ParseIP(dst)}); 
	rb := make([]byte, 1500)

	go getResponse(ctx, responseChan, rb, listener)

	select {
	case connected = <- responseChan:
	case <- ctx.Done():
		connected = false;
	}

	if connected {
		data.OnlineHosts = append(data.OnlineHosts, dst)
	}

}
