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
