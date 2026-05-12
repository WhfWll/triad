package network

import (
	"net"
)

func TelnetIsOpen(ip, port string) bool {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
