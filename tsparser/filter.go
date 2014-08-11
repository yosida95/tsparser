package tsparser

type filterBuffer struct {
	ok     bool
	data   []byte
	length int
}

func (b *filterBuffer) extendData(right []byte) {
	left := b.data
	b.data = make([]byte, len(left)+len(right))
	copy(b.data, left)
	copy(b.data[len(left):], right)
}

func (b *filterBuffer) Extend(packet Packet) {
	isFirst := len(b.data) == 0

	payload := packet.Payload()
	if !packet.payloadUnitStartIndicator() {
		if !isFirst {
			b.extendData(payload)
		}
		b.ok = false
		return
	}

	if !isFirst {
		b.data = b.data[b.length:]
		b.length = len(b.data)
		b.ok = true
	}
	packetStartCodePrefix := payload[0]<<16 | payload[1]<<8 | payload[2]
	if packetStartCodePrefix == 0x000001 {
		// Packetized Elementary Stream
		b.extendData(payload)
	} else {
		// Program Specific Information
		pointerField := int(payload[0])
		if isFirst {
			b.extendData(payload[1+pointerField:])
		} else {
			b.length += pointerField
			b.extendData(payload[1:])
		}
	}

	return
}

func (b *filterBuffer) OK() bool {
	return b.ok
}

func (b *filterBuffer) Bytes() []byte {
	if b.OK() {
		return b.data[b.length:]
	}

	return nil
}

type Filter struct {
	s       *Scanner
	pid     PID
	buffers map[PID]*filterBuffer
}

func NewFilter(scanner *Scanner, pids ...PID) *Filter {
	buffers := make(map[PID]*filterBuffer)
	for _, pid := range pids {
		buffers[pid] = new(filterBuffer)
	}

	return &Filter{
		s:       scanner,
		buffers: buffers,
	}
}

func (f *Filter) Scan() bool {
	for f.s.Scan() {
		packet := f.s.Packet()

		if buffer, ok := f.buffers[packet.PID()]; ok {
			buffer.Extend(packet)
			if buffer.OK() {
				f.pid = packet.PID()
				return true
			}
		}
	}

	return false
}

func (f *Filter) Bytes() []byte {
	buffer := f.buffers[f.pid]
	if buffer.OK() {
		return buffer.Bytes()
	}

	return nil
}
