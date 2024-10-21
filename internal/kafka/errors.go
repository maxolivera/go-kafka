package kafka

type ErrorCode uint16

const (
	NoError            ErrorCode = 0
	UnsupportedVersion ErrorCode = 35
)
