package tsparser

type PacketStream interface {
	Scan() bool
	Packet() Packet
}

type TableStream interface {
	Scan() bool
	Table() Table
}
