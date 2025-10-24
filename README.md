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
- The average latency in microseconds (µs).
- The latency standard deviation in microseconds (µs), indicating variability.
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
   ----- Tracing routers to www.google.com (142.250.178.132) ----- 
Hop   | IP Address       | Avg (µs)     | Std Dev (µs)    | Loss (%)  
---------------------------------------------------------------------
1     | 192.168.1.1      | 528.25       | 104.78          | 0.00
2     | 10.12.0.1        | 1256.74      | 258.91          | 0.00
3     | 172.16.25.1      | 2530.12      | 450.32          | 0.00
...
10    | 142.250.75.14    | 25687.55     | 1205.80         | 0.00
```

## Notes and Caveats

- The program may require elevated privileges to send ICMP packets.
- Raw ICMP sockets are restricted on some platforms (e.g., native Windows). Use WSL on Windows or adapt the code to use platform-specific ICMP APIs for native Windows support.
- This tool is provided for educational and diagnostic purposes. Ensure you have permission to probe target hosts.
- Results may vary depending on network conditions, firewalls, or other security settings.

## Acknowledgments

Thanks to the Go community and contributors to the x/net package.
