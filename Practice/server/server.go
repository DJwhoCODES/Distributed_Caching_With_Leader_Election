package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main() {
	PORT := ":5000"
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening on port:", PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting request:", err)
			continue
		}

		fmt.Println("Client Connected!")

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	const MaxSize = 10 * 1024 * 1024

	var strLen int32

	err := binary.Read(conn, binary.LittleEndian, &strLen)
	if err != nil {
		fmt.Println("Error while reading length:", err)
		return
	}

	if strLen <= 0 {
		fmt.Println("Invalid length:", strLen)
		return
	}

	if strLen > MaxSize {
		fmt.Println("Length too large:", strLen)
		return
	}

	buf := make([]byte, strLen)
	_, err2 := io.ReadFull(conn, buf)

	if err2 != nil {
		fmt.Println("Error while reading data:", err2)
		return
	}

	fmt.Println("Received data:", string(buf))

	binary.Write(conn, binary.LittleEndian, int32(len(buf)))
	binary.Write(conn, binary.LittleEndian, buf)
}
