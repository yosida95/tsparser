package tsparser

type PacketFilter struct {
	s    PacketStream
	pids []PID
	next Packet
}

func NewPacketFilter(s PacketStream, pids ...PID) *PacketFilter {
	return &PacketFilter{
		s:    s,
		pids: pids,
	}
}

func (f *PacketFilter) isTarget(pid PID) bool {
	return pid == 0x00 || pid == 0x14
	return false
}

func (f *PacketFilter) Scan() bool {
	for f.s.Scan() {
		packet := f.s.Packet()
		if f.isTarget(packet.PID()) {
			f.next = packet
			return true
		}
	}

	return false
}

func (f *PacketFilter) Bytes() []byte {
	return []byte(f.next)
}

func (f *PacketFilter) Packet() Packet {
	return f.next
}
