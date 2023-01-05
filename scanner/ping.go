package scanner

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

type ScanData struct {
	sync.RWMutex
	OnlineHosts []string
	wg sync.WaitGroup
}

func Ping(addr string, data *ScanData) {
	data.RWMutex.Lock()
	defer data.RWMutex.Unlock()
	defer data.wg.Done()
	pinger, err := ping.NewPinger(addr); 
	if err != nil {
		log.Fatal(err)
	}
	pinger.Timeout = time.Millisecond * 100
	pinger.Count = 1
	

	pinger.Run()

	if pinger.PacketsRecv > 0 {
		data.OnlineHosts = append(data.OnlineHosts, addr)
	}
}

func PingSingle(addr string, data *ScanData) {
	// used for when pinging a single IP
	pinger, err := ping.NewPinger(addr); 
	if err != nil {
		log.Fatal(err)
	}
	pinger.Timeout = time.Millisecond * 100
	pinger.Count = 1
	

	pinger.Run()

	if pinger.PacketsRecv > 0 {
		data.OnlineHosts = append(data.OnlineHosts, addr)
	}
}

func PingRange(scope string, start int, end int, data *ScanData) {
	// can pass a network address with CIDR or comma delimited IPs
	// if string of comma delimited IPs
	if strings.Contains(scope, ",") {
		// runs scans for multiple IPs
		addrs := strings.Split(scope, ",")
		for _, addr := range addrs {
			data.wg.Add(1)
			go Ping(addr, data)
		}
		data.wg.Wait()
		fmt.Println("Completed scanning alive hosts.")
		
	} else if !strings.Contains(scope, "/") {
		// execute a ping on a single host
		data.wg.Add(1)
		go Ping(scope, data)
		data.wg.Wait()
		fmt.Printf("Looks like %s is alive\n", scope)
	} else {
		// run this code if the scope does not contain commas
		ip := strings.TrimSuffix(strings.Split(scope, "/")[0], "0")
		mask, err := strconv.Atoi(strings.Split(scope, "/")[1])
		if err != nil {
			log.Fatal(err)
		}
	
		// TODO: need to change to be able to handle more than one network class
		if mask == 24 {
			for i := start; i <= end; i++ {
				data.wg.Add(1)
				dst := ip + strconv.Itoa(i)
				go Ping(dst, data)
			}
			data.wg.Wait()
			fmt.Println("Completed scanning alive hosts.")
		}
	}
}


// func PingScan(subnet string, start int, end int) {
// 	// Get my IP that is on this subnet
// 	data := scanData{}
// 	myIp, err := localIP()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ip := strings.TrimSuffix(strings.Split(subnet, "/")[0], "0")
// 	mask, err := strconv.Atoi(strings.Split(subnet, "/")[1]) 
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if mask == 24 {
// 		for i := start; i < end; i++ {
// 			data.wg.Add(1)
// 			dst := ip + strconv.Itoa(i)
// 			go icmpPing(myIp.String(), dst, &data)
// 		}
// 		data.wg.Wait()
// 		fmt.Println(data.OnlineHosts)
// 	}
// }


// func getResponse(ctx context.Context, responseChan chan bool, packet []byte, listener *icmp.PacketConn) {
// 	n, _, _ := listener.ReadFrom(packet)


// 	msg, _ := icmp.ParseMessage(1, packet[:n])

// 	// return response if the length is greater than 0
// 	if msg != nil {
// 		body, err := msg.Body.Marshal(1)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		fmt.Println(string(body))
// 		if strings.Contains(string(body), "HELLO-R-U-THERE") {
// 			responseChan <- true
// 		} else {
// 			responseChan <- false
// 		}
// 	} else {
// 		responseChan <- false
// 	}
// }


// func icmpPing(src string, dst string, data *scanData) {
// 	// sends an icmp echo request and listens for 1 second before timing out.
//   // create context
// 	data.Lock()
// 	defer data.Unlock()
// 	defer data.wg.Done()
// 	connected := false
// 	protocol := "udp4"

// 	// retreive local IP
// 	ip, err := localIP()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx := context.Background()
// 	ctx, cancel := context.WithTimeout(ctx, time.Millisecond * 100)
// 	defer cancel()

// 	// response chan
// 	responseChan := make(chan bool)

// 	if dst == "" {
// 		log.Fatal("Please specify an IP Address!")
// 	}

// 	// check if the OS is windows
// 	if runtime.GOOS == "windows" {
// 		protocol = "ip4:icmp"
// 	}

// 	// TODO: needs to dynamically retreive the correct interface IP.
// 	listener, err := icmp.ListenPacket(protocol, ip.String())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer listener.Close()

// 	message := icmp.Message{
// 		Type: ipv4.ICMPTypeEcho,
// 		Code: 0,
// 		Body: &icmp.Echo{
// 			ID: os.Getpid() & 0xffff, Seq: 1,
// 			Data: []byte("HELLO-R-U-THERE"),
// 		},
// 	}
// 	encoded, err := message.Marshal(nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	listener.WriteTo(encoded, &net.UDPAddr{IP: net.ParseIP(dst)}); 
// 	rb := make([]byte, 1500)

// 	go getResponse(ctx, responseChan, rb, listener)

// 	select {
// 	case connected = <- responseChan:
// 	case <- ctx.Done():
// 		connected = false;
// 	}

// 	if connected {
// 		data.OnlineHosts = append(data.OnlineHosts, dst)
// 	}

// }
