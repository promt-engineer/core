package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("::1")

	ip = ip.To4()

	fmt.Println(ip)
}
