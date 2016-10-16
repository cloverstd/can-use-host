package main

import (
	"flag"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	defer func() {
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
		}
	}()

	start := flag.String("start", "127.0.0.1", "start ip")
	end := flag.String("end", "127.0.0.1", "end ip")
	flag.Parse()
	startIP := net.ParseIP(*start)
	endIP := net.ParseIP(*end)
	if startIP == nil || endIP == nil {
		fmt.Println("start or end is a invalid IP address.")
		return
	}
	fmt.Println(startIP, endIP)
	scan_ip(startIP, endIP)
}

func ping(ip string, jobs chan bool, results map[string]interface{}) {
	defer func() {
		// Println executes normally even if there is a panic
		if err := recover(); err != nil {
			// results[ip] = "error"
			fmt.Println("runtime error")
		}
	}()
	// ticker := time.NewTicker(10 * time.Second)

	// go func() {
	// for _ = range ticker.C {
	// }
	// jobs <- true
	// results[ip] = "timeout"
	// }()

	fmt.Printf("PING %s\n", ip)
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		jobs <- false
		results[ip] = err.Error()
		return
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {

		defer func() {
			// Println executes normally even if there is a panic
			if err := recover(); err != nil {
				results[ip] = "error"
			}
		}()

		// fmt.Printf("ping %s: receive, RTT: %v\n", addr.String(), rtt)
		results[ip] = rtt
	}
	p.OnIdle = func() {
		defer func() {
			// Println executes normally even if there is a panic
			if err := recover(); err != nil {
				results[ip] = "error"
			}
		}()
		jobs <- true
		// fmt.Printf("ping %s finish\n", ip)
	}
	err = p.Run()
	if err != nil {
		jobs <- false
		results[ip] = err.Error()
	}
}

func scan_ip(startIP, endIP net.IP) map[string]interface{} {
	results := make(map[string]interface{})
	start := inet_aton(startIP)
	end := inet_aton(endIP)
	if start > end {
		temp := end
		end = start
		start = temp
	}
	jobs := make(chan bool, end)
	for i := start; i <= end; i++ {
		ipAddr := inet_ntoa(i)
		ip := ipAddr.String()
		results[ip] = "down"
		go ping(ip, jobs, results)
	}

	for j := start; j <= end; j++ {
		<-jobs
	}

	fmt.Println("\n\n\nping done.\n\n\n")
	for j := start; j <= end; j++ {
		ipAddr := inet_ntoa(j)
		ip := ipAddr.String()
		result := results[ip]
		fmt.Printf("%s is %s\n", ip, result)
	}
	return results
}

func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}
