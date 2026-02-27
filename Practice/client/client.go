package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main() {
	PORT := ":5000"
	conn, err := net.Dial("tcp", "localhost"+PORT)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	msg := "Hello From Client!"

	data := []byte(msg)

	binary.Write(conn, binary.LittleEndian, int32(len(data)))
	binary.Write(conn, binary.LittleEndian, data)

	fmt.Println("Message sent:", msg)

	var respLen int32
	binary.Read(conn, binary.LittleEndian, &respLen)

	respBuf := make([]byte, respLen)
	io.ReadFull(conn, respBuf)

	fmt.Println("Received from server!")

}
