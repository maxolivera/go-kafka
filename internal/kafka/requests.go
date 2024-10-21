package kafka

import (
	"encoding/binary"
	"net"
)

const maxRawMessage = 1024

type KafkaRequest int

const (
	APIVersions KafkaRequest = iota // get the supported API versions by the broker
	Fetch                           // get data from the broker
	Produce                         // send data to the broker
)

// Requests and responses start with 4 bytes stating the
// length of the entire message (including header and body),
// followed by the message

// First 4 bytes are int32, big-endian

// NOTE(maolivera): NULLABLE_STRING Represents a sequence of characters
// or null. For non-null strings, first the length N is given as an
// INT16. Then N bytes follow which are the UTF-8 encoding of the
// character sequence. A null value is encoded with length of -1
// and there are no following bytes.

type Request struct {
	Length uint32
	Header *RequestHeader
}

type RequestHeader struct {
	ApiKey        uint16
	ApiVersion    uint16
	CorrelationID uint32
	ClientID      string
	TaggedFields  string
}

func ReadRequest(conn net.Conn) (*Request, error) {
	rawRequest := make([]byte, maxRawMessage)
	if _, err := conn.Read(rawRequest); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(rawRequest[:4])
	apiKey := binary.BigEndian.Uint16(rawRequest[4:6])
	apiVersion := binary.BigEndian.Uint16(rawRequest[6:8])
	correlationID := binary.BigEndian.Uint32(rawRequest[8:12])

	header := RequestHeader{
		ApiKey: apiKey,
		ApiVersion: apiVersion,
		CorrelationID: correlationID,
	}
	req := Request{
		Length: length,
		Header: &header,
	}
	return &req, nil
}
