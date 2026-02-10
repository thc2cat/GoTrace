# GoTrace — a Go Network Analyzer, a diagnostic tool for latency and packet loss 🌐

Meta-description: "An open-source CLI tool developed in Go to analyze network performance. Measure latency, packet loss, and trace data paths with this network analyzer."

This project is a command-line tool written in Go for analyzing network connectivity and performance. It combines traceroute and ping features to provide a complete view of latency, packet loss, and the path that data takes to a destination.

![docs/GoTrace.png](docs/GoTrace.png)

## Features 🛠️

The program is organized into three distinct phases for detailed analysis:

### Phase 1: Router Discovery (Traceroute)

This phase identifies intermediate routers (hops) between your machine and the final destination. It uses ICMP packets with an incremental Time-to-Live (TTL) to map the full route.

### Phase 2: Performance Measurement (Ping)

After the route is traced, the program sends a specified number of ICMP packets to each router in the list. It collects latency data for each hop, allowing you to pinpoint weak spots or bottlenecks along the path.

### Phase 3: Statistics Display

Results are presented in a compact, easy-to-read table. For each router, the tool shows:

- The router's IP address.
- The median latency in microseconds (µs) (P50).
- The 90th percentile latency in microseconds (µs) (P90), useful to assess worst-case latency.
- The packet loss percentage, a key indicator of connection reliability.

## Requirements 📋

- Go (version 1.18 or newer)
- Administrator privileges (sudo on Linux/macOS) or equivalent on other systems, since the program requires access to raw ICMP sockets.

## Installation and Usage 🚀

Clone the repository:

```bash
git clone https://github.com/votre_utilisateur/go-network-analyzer.git
```

Run the program:
Execute the application by specifying the target host (domain name or IP) and the number of packets to send.

```bash
sudo go run main.go <hostname_or_ip> <number_of_packets> [delay_in_ms]
```

Example:
To analyze the path to google.com by sending 10 packets to each hop, run:

```bash
sudo go run main.go google.com 10
```

## Example Output 📊

Here is an example of the tool's final output:

```bash
   ----- Tracing routers to www.google.com (172.217.20.36) ----- 
Hop   | IP Address     | P50 (µs) | P90 (µs) | Loss (%)  
---------------------------------------------------------------------
...  
5     | 192.178.70.144 | 68       | 1996     | 0.00      
6     | 72.14.236.91   | 2825     | 3075     | 0.00      
7     | 142.251.253.35 | 2092     | 2215     | 0.00      
8     | 172.217.20.36  | 2030     | 2173     | 0.00      
```

## Notes and Caveats

- The program may require elevated privileges to send ICMP packets.
- Raw ICMP sockets are restricted on some platforms (e.g., native Windows). Use WSL on Windows or adapt the code to use platform-specific ICMP APIs for native Windows support.
- This tool is provided for educational and diagnostic purposes. Ensure you have permission to probe target hosts.
- Results may vary depending on network conditions, firewalls, or other security settings.

## Acknowledgments

Thanks to the Go community and contributors to the x/net package.
