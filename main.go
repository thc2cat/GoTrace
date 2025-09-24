package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"strconv"
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sudo go run main.go <hostname or IP> <number of packets>")
		os.Exit(1)
	}

	host := os.Args[1]
	numPackets, err := strconv.Atoi(os.Args[2])
	if err != nil || numPackets <= 0 {
		fmt.Println("Invalid number of packets. Please provide a positive integer.")
		os.Exit(1)
	}

	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		fmt.Println("Error resolving host:", err)
		os.Exit(1)
	}

	fmt.Printf("Phase 1: Découverte des routeurs vers %s (%s)\n", host, ipAddr.String())
	fmt.Println("---------------------------------------")

	routerList := traceroute(ipAddr)

	fmt.Println("\n---------------------------------------")
	fmt.Println("Phase 2: Mesure de la latence et de la perte de paquets pour chaque routeur")
	fmt.Println("---------------------------------------")

	statsList := make([]RouterStats, len(routerList))
	for i, routerIPAddr := range routerList {
		fmt.Printf("Mesure sur le routeur %d: %s\n", i+1, routerIPAddr.String())
		latencies, loss := ping(routerIPAddr.IP, numPackets)
		statsList[i] = RouterStats{
			IP:         routerIPAddr.String(),
			Latencies:  latencies,
			PacketLoss: loss,
		}
	}

	fmt.Println("\n---------------------------------------")
	fmt.Println("Phase 3: Affichage des résultats")
	fmt.Println("---------------------------------------")
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

		fmt.Printf("%2d: %s\n", ttl, hopIP.String())
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
	}

	loss := float64(sentCount-receivedCount) / float64(sentCount) * 100
	return latencies, loss
}

// displayResults affiche les statistiques de chaque routeur avec un formatage amélioré.
func displayResults(statsList []RouterStats) {
	// Nouveaux formats pour un meilleur alignement
	fmt.Printf("%-5s | %-16s | %-12s | %-12s | %-10s\n", "Hop", "Adresse IP", "Moyenne (µs)", "Écart-type (µs)", "Perte (%)")
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

		// Affichage avec le nouveau formatage pour les microsecondes
		fmt.Printf("%-5d | %-16s | %-12.2f | %-15.2f | %-10.2f\n",
			i+1,
			stats.IP,
			float64(avgRTT.Microseconds()),
			float64(stdDev.Microseconds()),
			stats.PacketLoss,
		)
	}
}
