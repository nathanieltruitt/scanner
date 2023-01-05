package scanner

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func PortScan(addr string, port int, proto string) {
	d := net.Dialer{Timeout: time.Second * 5}
	defer wg.Done()
	conn, err := d.Dial(proto, addr + ":" + strconv.Itoa(port))
	// run if the connection opens successfully
	if err != nil {
		return
	}

	if conn != nil {
		defer conn.Close()
		fmt.Printf("%v  ", port)
	}
}

func Scan(addr string, port string, proto string) {
	fmt.Printf("port scanning %s \n\n------------------------------------\n------------------------------------\n\n", addr)
	// start and end parameters here are for starting port and ending port
	if strings.Contains(port, "-") {
		start, err := strconv.Atoi(strings.Split(port, "-")[0])
		if err != nil {
			log.Println(err)
		}
		end, err := strconv.Atoi(strings.Split(port, "-")[1])
		if err != nil {
			log.Println(err)
		} 

		fmt.Printf("scanning port %d through port %d for %s\n", start, end, addr)
		fmt.Println("Open ports: ")
		for i := start; i <= end; i++ {
			wg.Add(1)
			go PortScan(addr, i, proto)
		}
		wg.Wait()
	} else if strings.Contains(port, ",") {
		ports := strings.Split(port, ",")
		fmt.Println("Open ports: ")
		for _, pt := range ports {
			portNum, err := strconv.Atoi(pt)
			if err != nil {
				log.Println(err)
			}
			wg.Add(1)
			go PortScan(addr, portNum, proto)
		}
		wg.Wait()
	} else {
		// assume that their is only one IP address
		fmt.Println("Open ports: ")
		portNum, err := strconv.Atoi(port)
		if err != nil {
			log.Println(err)
		}
		wg.Add(1)
		go PortScan(addr, portNum, proto)
		wg.Wait()
	}
}