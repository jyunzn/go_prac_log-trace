package utils

import (
	"net"
	"fmt"
	"strings"
)

func GetOutboundIP()(ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:8080")
	if err != nil {
		fmt.Println("dial err", err)
		return
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(addr.IP.String(), ":")[0]
	return
}
