package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/kafka-starter-go/internal/kafka"
)

func main() {
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

	fmt.Printf("Received a %d request\n", req.Header.ApiKey)

	resBuffer := make([]byte, 8)
	binary.BigEndian.PutUint32(resBuffer[4:8], req.Header.CorrelationID)
	switch req.Header.ApiKey {
	case 0, 1, 2, 3, 4:
		fmt.Println("Unsupported request: ", req.Header.ApiVersion)
		os.Exit(1)
	case 18:
		// Error
		resBuffer = binary.BigEndian.AppendUint16(resBuffer, uint16(kafka.NoError))
		fmt.Println("Cumulative response buffer")
		printHex(resBuffer)
		// ApiKey
		resBuffer = binary.BigEndian.AppendUint16(resBuffer, uint16(kafka.APIVersions))
		fmt.Println("Cumulative response buffer")
		printHex(resBuffer)
		// MinVersion
		resBuffer = binary.BigEndian.AppendUint16(resBuffer, 0)
		fmt.Println("Cumulative response buffer")
		printHex(resBuffer)
		// MaxVersion
		resBuffer = binary.BigEndian.AppendUint16(resBuffer, 4)
		fmt.Println("Cumulative response buffer")
		printHex(resBuffer)
	default:
		resBuffer = binary.BigEndian.AppendUint16(resBuffer, uint16(kafka.UnsupportedVersion))
	}

	// Write length

	length := uint32(len(resBuffer)) - 4
	binary.BigEndian.PutUint32(resBuffer[0:4], length)
	fmt.Println("Cumulative response buffer")
	printHex(resBuffer)
	fmt.Printf("Response length %d\n", length)

	if _, err = conn.Write(resBuffer); err != nil {
		fmt.Println("Error sending response: ", err.Error())
		os.Exit(1)
	}
}

func printHex(bytes []byte) {
	for n, byte := range bytes {
		if n%10 == 0 {
			fmt.Println()
		}
		fmt.Printf(" %x", byte)
	}
}
