package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func hexToIPv4(hexAddr string) string {
	addr, err := strconv.ParseUint(hexAddr, 16, 32)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d.%d", addr&0xff, (addr>>8)&0xff, (addr>>16)&0xff, (addr>>24)&0xff)
}

func hexToIPv6(hexAddr string) string {
	if len(hexAddr) != 32 {
		return ""
	}
	var ip net.IP
	for i := 0; i < 32; i += 2 {
		byteVal, _ := strconv.ParseUint(hexAddr[i:i+2], 16, 8)
		ip = append(ip, byte(byteVal))
	}
	return ip.String()
}

func hexToPort(hexPort string) int {
	port, err := strconv.ParseUint(hexPort, 16, 16)
	if err != nil {
		return 0
	}
	return int(port)
}

func parseProcNetFile(filePath string, isIPv6 bool) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening %s: %s\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过文件头部
		if strings.HasPrefix(line, "  sl") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 10 {
			localAddress := fields[1]
			state := fields[3]

			// TCP/UDP 监听状态码为 "0A"，对于 UDP6 "07" 也代表活跃端口
			if state == "0A" || (filePath == "/proc/net/udp6" && state == "07") {
				addrPort := strings.Split(localAddress, ":")
				if len(addrPort) == 2 {
					ip := hexToIPv4(addrPort[0])
					if isIPv6 {
						ip = hexToIPv6(addrPort[0])
					}
					port := hexToPort(addrPort[1])
					fmt.Printf("Listening on %s:%d\n", ip, port)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading %s: %s\n", filePath, err)
	}
}

func parseUnixSockets() {
	file, err := os.Open("/proc/net/unix")
	if err != nil {
		fmt.Printf("Error opening /proc/net/unix: %s\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 跳过文件头部
		if strings.HasPrefix(line, "Num") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 8 {
			refCount := fields[2]
			protocol := fields[3]
			flags := fields[4]
			typeSock := fields[5]
			state := fields[6]
			path := fields[7]

			fmt.Printf("Unix Socket Path: %s, RefCount: %s, Protocol: %s, Flags: %s, Type: %s, State: %s\n",
				path, refCount, protocol, flags, typeSock, state)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading /proc/net/unix: %s\n", err)
	}
}

func main() {
	fmt.Println("========= tcp =========")
	parseProcNetFile("/proc/net/tcp", false) // IPv4 TCP
	parseProcNetFile("/proc/net/tcp6", true) // IPv6 TCP
	fmt.Println("========= tcp =========")
	fmt.Println()
	fmt.Println("========= udp =========")
	parseProcNetFile("/proc/net/udp", false) // IPv4 UDP
	parseProcNetFile("/proc/net/udp6", true) // IPv6 UDP
	fmt.Println("========= udp =========")
	fmt.Println()
	fmt.Println("========= unix socket =========")
	parseUnixSockets() // Unix Sockets
	fmt.Println("========= unix socket =========")
}
