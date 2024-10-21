package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")


	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	reqBuffer := make([]byte, 1024)
	if _, err = conn.Read(reqBuffer); err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	resBuffer := make([]byte, 8)
	binary.BigEndian.PutUint32(resBuffer[4:], 7)

	if _, err = conn.Write(resBuffer); err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
}
