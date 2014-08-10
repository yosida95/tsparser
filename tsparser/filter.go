package tsparser

type Filter struct {
	s      *Scanner
	pid    uint16
	buffer []byte
	length int
}

func NewFilter(scanner *Scanner, pid uint16) *Filter {
	return &Filter{
		s:      scanner,
		pid:    pid,
		length: 0,
	}
}

func (f *Filter) extendBuffer(data []byte) {
	buffer := f.buffer
	f.buffer = make([]byte, len(buffer)+len(data))
	copy(f.buffer, buffer)
	copy(f.buffer[len(buffer):], data)
}

func (f *Filter) Scan() bool {
	if f.length > 0 {
		f.buffer = f.buffer[f.length:]
	}

	for f.s.Scan() {
		packet := f.s.Packet()
		if packet.PID() != f.pid {
			continue
		}

		payload := packet.Payload()
		if packet.payloadUnitStartIndicator() == 1 {
			f.length = len(f.buffer)
			packetStartCodePrefix := payload[0]<<16 | payload[1]<<8 | payload[2]
			if packetStartCodePrefix == 0x000001 {
				// Packetized Elementary Stream
				f.extendBuffer(payload)
			} else {
				// Program Specific Information
				pointerField := payload[0]
				f.length += int(pointerField)
				f.extendBuffer(payload[1:])
			}

			if f.length == 0 {
				return f.Scan()
			}
			return true
		} else if len(f.buffer) > 0 {
			f.extendBuffer(payload)
		}
	}

	return false
}

func (f *Filter) Bytes() []byte {
	if len(f.buffer) > f.length {
		return f.buffer[:f.length]
	}

	return nil
}
