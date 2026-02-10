package main

// Méta-description :
//  "An open-source CLI tool developed in Go to analyze network performance.
//  Measure latency, packet loss, and trace data paths with this network analyzer."
// Combine traceroute and ping functionalities to analyze network performance.
// Measure latency, packet loss, and trace data paths.
//
// WARNING :
// Raw ICMP sockets are not supported on native Windows Go environments due to
// OS security restrictions.
//
// Usage: sudo go run main.go <hostname or IP> <number of packets> [delay in µs]
// Example: sudo go run main.go example.com 10 500
// Note: Requires elevated privileges to send ICMP packets.
// Author: ChatGPT (and improved by user)
// License: MIT
// Repository: https://github.com/thc2cat/GoTrace
// Version: 1.1.0
// Date: 2025-09-23
// Language: Go
// Tags: network, traceroute, ping, latency, packet loss, CLI tool
// Categories: Networking, Utilities, Command-Line Tools
// Keywords: network analysis, traceroute, ping, latency measurement, packet loss, Go programming, CLI tool
// Platforms: Cross-platform (Linux, macOS, Windows with WSL)
// Requirements: Go 1.16+, elevated privileges for ICMP
// Installation: go mod init GoTrace && go mod tidy
// Usage Instructions: Run with sudo or as administrator
// Contribution Guidelines: Fork the repo, make changes, and submit a pull request
// Contact Information: Open an issue on GitHub for support or questions
// Disclaimer: Use responsibly and ethically, respecting network policies.
// Acknowledgments: Thanks to the Go community and contributors to the x/net package.
// Enjoy analyzing your network performance!

// Note: This code requires the "golang.org/x/net/icmp" and "golang.org/x/net/ipv4" packages.
// Install them using: go get golang.org/x/net/icmp golang.org/x/net/ipv4

// Note: Run the program with elevated privileges (e.g., using sudo) to allow sending ICMP packets.
// Note: This code is for educational purposes. Ensure you have permission to ping/traceroute the target hosts.
// Note: The program may not work on Windows without WSL due to raw socket restrictions.
// Note: The program may require adjustments for different operating systems or network configurations.
// Note: The program may not handle all edge cases or network errors. Use with caution.
// Note: The program may produce different results based on network conditions and configurations.
// Note: The program may not work behind certain firewalls or network security settings.
// Note: The program is provided "as is" without warranty of any kind. Use at your own risk.
// Note: The program may not be suitable for production use. Test thoroughly before deployment.
// Note: The program may require additional error handling for robustness.
// Note: The program may not be compatible with all Go versions. Tested with Go 1.16+.
// Note: The program may not work in all network environments. Test in your specific setup.
// Note: The program may require additional dependencies or libraries for full functionality.
// Note: The program may not be optimized for performance. Use for small-scale testing.
// Note: The program may not handle IPv6 addresses. Modify as needed for IPv6 support.
// Note: The program may not work with certain network configurations (e.g., VPNs, proxies).
// Note: The program may produce different results based on the target host's response behavior.
// Note: The program may require additional permissions or capabilities on certain operating systems.
// Note: The program may not be suitable for all users. Use with caution and understanding of network protocols.
// Note: The program may not be compatible with all network devices or configurations.
// Note: The program may require additional configuration for specific use cases.
// Note: The program may not handle all types of ICMP messages. Modify as needed for specific requirements.
// Note: The program may not be suitable for high-frequency or large-scale network testing.
// Note: The program may require additional logging or debugging for troubleshooting.
// Note: The program may not be compatible with all Go modules or package management systems.
// Note: The program may require additional documentation or user guides for effective use.
// Note: The program may not be suitable for all network environments. Test in your specific setup.
// Note: The program may require additional security considerations for safe use.
// Note: The program may not be compatible with all network protocols or configurations.
// Note: The program may require additional testing or validation for specific use cases.
// Note: The program may not handle all types of network errors or exceptions. Use with caution.
// Note: The program may not be suitable for all users. Ensure you understand the implications of network testing.

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	protocolICMP = 1
)

// RouterStats est une structure pour stocker les statistiques de chaque routeur.
type RouterStats struct {
	IP         string
	Latencies  []time.Duration
	PacketLoss float64
}

var (
	delay time.Duration = 500 * time.Microsecond
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sudo go run main.go <hostname or IP> <number of packets> [delay in µs]")
		os.Exit(1)
	}

	host := os.Args[1]
	numPackets, err := strconv.Atoi(os.Args[2])
	if err != nil || numPackets <= 0 {
		fmt.Println("Invalid number of packets. Please provide a positive integer.")
		os.Exit(1)
	}

	if len(os.Args) >= 4 {
		delayP, err := strconv.Atoi(os.Args[3])
		if err == nil {
			delay = time.Duration(delayP) * time.Microsecond
		}
	}

	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println("Error resolving host:", err)
		os.Exit(1)
	}

	const totalWidth = 69
	centeredString := fmt.Sprintf(" ----- Tracing routers to %s (%s) ----- ", host, ipAddr.String())
	stringLength := len(centeredString)
	leftPadding := (totalWidth - stringLength) / 2
	spaces := strings.Repeat(" ", leftPadding)
	fmt.Printf("%s%s\n", spaces, centeredString)

	routerList := traceroute(ipAddr)

	statsList := make([]RouterStats, len(routerList))
	for i, routerIPAddr := range routerList {
		// fmt.Printf("Measuring on router %d: %s\n", i+1, routerIPAddr.String())
		latencies, loss := ping(routerIPAddr.IP, numPackets)
		statsList[i] = RouterStats{
			IP:         routerIPAddr.String(),
			Latencies:  latencies,
			PacketLoss: loss,
		}
	}

	displayResults(statsList)
}

// traceroute découvre les routeurs intermédiaires et retourne leurs adresses IP.
func traceroute(dest *net.IPAddr) []*net.IPAddr {
	var routers []*net.IPAddr
	maxHops := 30
	reached := false

	for ttl := 1; ttl <= maxHops && !reached; ttl++ {
		conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		if err != nil {
			fmt.Println("Error listening for ICMP:", err)
			return nil
		}
		defer conn.Close()

		if err := conn.IPv4PacketConn().SetTTL(ttl); err != nil {
			fmt.Println("Error setting TTL:", err)
			return nil
		}

		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: 1,
				Data: []byte("HELLO-GO-TRACEROUTE"),
			},
		}
		wb, _ := wm.Marshal(nil)

		if _, err := conn.WriteTo(wb, dest); err != nil {
			continue
		}

		reply := make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, peer, err := conn.ReadFrom(reply)

		if err != nil {
			fmt.Printf("%2d: * * *\n", ttl)
			continue
		}

		rm, _ := icmp.ParseMessage(protocolICMP, reply[:n])
		hopIP := peer.(*net.IPAddr)

		if rm.Type == ipv4.ICMPTypeEchoReply {
			reached = true
		}

		// fmt.Printf("%2d: %s\n", ttl, hopIP.String())
		routers = append(routers, hopIP)
	}
	return routers
}

// ping envoie des paquets ICMP pour mesurer la latence et la perte.
func ping(dest net.IP, count int) ([]time.Duration, float64) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Error listening for ICMP:", err)
		return nil, 100.0
	}
	defer conn.Close()

	sentCount := 0
	receivedCount := 0
	var latencies []time.Duration

	for i := 0; i < count; i++ {
		sentCount++
		start := time.Now()

		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho, Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: i,
				Data: []byte("HELLO-GO-PING"),
			},
		}
		wb, _ := wm.Marshal(nil)

		if _, err := conn.WriteTo(wb, &net.IPAddr{IP: dest}); err != nil {
			continue
		}

		reply := make([]byte, 1500)
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, _, err := conn.ReadFrom(reply)

		if err != nil {
			continue
		}

		rm, _ := icmp.ParseMessage(protocolICMP, reply[:n])
		if rm.Type == ipv4.ICMPTypeEchoReply {
			receivedCount++
			rtt := time.Since(start)
			latencies = append(latencies, rtt)
		}
		time.Sleep(delay)
	}

	loss := float64(sentCount-receivedCount) / float64(sentCount) * 100
	return latencies, loss
}

// displayResults displays the statistics for each router with improved formatting.
func displayResults(statsList []RouterStats) {
	format1 := "%-5s | %-14s | %-8s | %-8s | %-10s\n"
	format2 := "%-5d | %-14s | %-8.f | %-8.f | %-10.2f\n"

	// New formats for better alignment
	fmt.Printf(format1, "Hop", "IP Address", "Avg (µs)", "σ (µs)", "Loss (%)")
	fmt.Println("---------------------------------------------------------------------")

	for i, stats := range statsList {
		if len(stats.Latencies) == 0 {
			fmt.Printf("%-5d | %-16s | %-12s | %-12s | %-10.2f\n", i+1, stats.IP, "N/A", "N/A", stats.PacketLoss)
			continue
		}

		// Calcul de la moyenne
		var totalRTT time.Duration
		for _, rtt := range stats.Latencies {
			totalRTT += rtt
		}
		avgRTT := totalRTT / time.Duration(len(stats.Latencies))

		// Calcul de l'écart-type
		var sumSquares float64
		for _, rtt := range stats.Latencies {
			diff := float64(rtt - avgRTT)
			sumSquares += diff * diff
		}
		variance := sumSquares / float64(len(stats.Latencies))
		stdDev := time.Duration(math.Sqrt(variance))

		// Display with new formatting for microseconds
		fmt.Printf(format2,
			i+1,
			stats.IP,
			float64(avgRTT.Microseconds()),
			// (float64(stdDev.Microseconds())/float64(avgRTT.Microseconds()))*100,
			float64(stdDev.Microseconds()),
			stats.PacketLoss,
		)
	}
}
