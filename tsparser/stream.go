package tsparser

type PacketStream interface {
	Scan() bool
	Bytes() []byte
	Packet() Packet
}

type TableStream interface {
	Scan() bool
	Table() Table
	Bytes() []byte
}
