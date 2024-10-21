package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/kafka-starter-go/internal/kafka"
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

	req, err := kafka.ReadRequest(conn)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	resBuffer := make([]byte, 12)
	binary.BigEndian.PutUint32(resBuffer[4:8], req.Header.CorrelationID)

	if req.Header.ApiVersion < 0 || req.Header.ApiVersion > 4 {
		binary.BigEndian.PutUint16(resBuffer[8:12], uint16(kafka.UnsupportedVersion))
	}

	if _, err = conn.Write(resBuffer); err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
}
