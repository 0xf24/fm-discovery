package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	//create new listener on all addresses on the port 1030
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 1030})

	//panic if listener failed
	if err != nil {
		log.Fatal(err)
		return
	}

	//TODO: replace with actual feedback
	println("listening on 0.0.0.0:1030")

	for {
		conn, err := ln.AcceptTCP()

		if err != nil {
			//ignore errors
			break
		}

		//handle connection on new go routine
		go handle(conn)
	}
}

// handle incoming tcp connections
func handle(conn *net.TCPConn) {
	//get remote address
	addr := conn.RemoteAddr()

	//convert address into IP & port combo
	tcpaddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())

	if err != nil {
		//this shouldn't happen, but just in case we close the socket
		_ = conn.Close()
		return
	}

	//create writer so we can chuck bytes
	w := bufio.NewWriter(conn)

	//write version number
	_ = w.WriteByte(0x01)

	_, _ = w.Write(tcpaddr.IP.To4())

	//write port to buffer

	port := uint16(tcpaddr.Port)

	_ = w.WriteByte(byte(port))
	_ = w.WriteByte(byte(port >> 8))

	//flush buffer
	_ = w.Flush()

	//close connection
	_ = conn.Close()
}
