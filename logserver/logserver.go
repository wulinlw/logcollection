package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:5000")
	checkError(err)
	defer l.Close()
	for {
		conn, err := l.Accept()
		checkError(err)

		// Handle connections in a new goroutine.
		go handleRequest(conn)

	}
}
func handleRequest(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(buf[:reqLen]))
	}
}
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
