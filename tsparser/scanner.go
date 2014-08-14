package tsparser

import (
	"io"
)

const (
	PacketSize          = 188
	BufferedPacketCount = 5
	BufferSize          = PacketSize * BufferedPacketCount

	SyncByte byte = 0x47
)

type PacketScanner struct {
	r      io.Reader
	buffer [PacketSize * BufferedPacketCount]byte
	seek   int
	eof    bool
	err    error
}

func NewPacketScanner(r io.Reader) *PacketScanner {
	return &PacketScanner{
		r:    r,
		seek: 0,
		eof:  false,
		err:  nil,
	}
}

func (s *PacketScanner) Err() error {
	return s.err
}

func (s *PacketScanner) leftShift(n int) {
	copy(s.buffer[:BufferSize-n], s.buffer[n:])
	s.seek -= n
}

func (s *PacketScanner) rightShift(n int) {
	copy(s.buffer[n:], s.buffer[:BufferSize-n])
	s.seek += n
}

func (s *PacketScanner) fillBuffer(n int) bool {
	switch m, err := io.ReadFull(s.r, s.buffer[BufferSize-n:]); err {
	case io.ErrUnexpectedEOF, io.EOF:
		s.eof = true

		restBytes := (BufferSize - n + m) / PacketSize * PacketSize
		shiftBytes := BufferSize - restBytes
		s.rightShift(shiftBytes)
		return s.seek < BufferSize
	case nil:
		return true
	default:
		panic(err)
	}
}

func (s *PacketScanner) isSynced() bool {
	synced := true
	for i := s.seek; i < BufferSize; i += PacketSize {
		if s.buffer[i] != SyncByte {
			synced = false
			break
		}
	}

	return synced
}

func (s *PacketScanner) sync() bool {
	for i := 0; i < PacketSize; i++ {
		if !s.isSynced() {
			s.seek++
			continue
		}

		if i == 0 {
			return true
		}

		s.leftShift(i)
		return s.fillBuffer(i)
	}

	return false
}

func (s *PacketScanner) Scan() bool {
	if s.seek >= BufferSize {
		return false
	}

	s.seek += PacketSize
	if s.seek < BufferSize && (s.eof || s.isSynced()) {
		return true
	} else if s.eof {
		return false
	}

	s.seek = 0
	if !s.fillBuffer(BufferSize) {
		return false
	}

	for !s.sync() {
		s.leftShift(PacketSize)
		if !s.fillBuffer(PacketSize) {
			return false
		}
	}

	return true
}

func (s *PacketScanner) Bytes() []byte {
	if s.seek < BufferSize {
		return s.buffer[s.seek : s.seek+PacketSize]
	}

	return nil
}

func (s *PacketScanner) Packet() Packet {
	return Packet(s.Bytes())
}

type tableScannerBuffer struct {
	data    []byte
	current []byte
	isPES   bool
}

func (b *tableScannerBuffer) extend(right []byte) {
	n := len(b.data)
	m := n + len(right)

	if m > cap(b.data) {
		newSlice := make([]byte, m*2)
		copy(newSlice, b.data)
		b.data = newSlice
	}

	b.data = b.data[0:m]
	copy(b.data[n:m], right)
}

func (b *tableScannerBuffer) clear() {
	b.data = make([]byte, 0, cap(b.data))
}

func (b *tableScannerBuffer) freeze() {
	b.current = make([]byte, len(b.data))
	copy(b.current, b.data)
}

func (b *tableScannerBuffer) Begin(payload []byte) {
	packetStartCodePrefix := payload[0]<<16 | payload[1]<<8 | payload[2]
	b.isPES = packetStartCodePrefix == 0x000001

	if b.isPES {
		b.freeze()
		b.clear()
		b.extend(payload)
	} else {
		pointerField := int(payload[0])
		if len(b.data) > 0 {
			b.extend(payload[1 : 1+pointerField])
			b.freeze()
		}

		b.clear()
		b.extend(payload[1+pointerField:])
	}
}

func (b *tableScannerBuffer) Extend(payload []byte) {
	if len(b.data) == 0 {
		return
	}

	b.extend(payload)
}

func (b *tableScannerBuffer) isFull() bool {
	return len(b.current) > 0
}

func (b *tableScannerBuffer) Bytes() []byte {
	return b.current
}

type TableScanner struct {
	s       PacketStream
	pid     PID
	buffers map[PID]*tableScannerBuffer
}

func NewTableScanner(s PacketStream) *TableScanner {
	return &TableScanner{
		s:       s,
		buffers: make(map[PID]*tableScannerBuffer),
	}
}

func (s *TableScanner) Scan() bool {
	for s.s.Scan() {
		packet := s.s.Packet()

		buffer, ok := s.buffers[packet.PID()]
		if !ok {
			s.buffers[packet.PID()] = new(tableScannerBuffer)
			buffer = s.buffers[packet.PID()]
		}

		if packet.payloadUnitStartIndicator() {
			buffer.Begin(packet.Payload())
			if buffer.isFull() {
				s.pid = packet.PID()
				return true
			}
		} else if ok {
			buffer.Extend(packet.Payload())
		}
	}

	return false
}

func (s *TableScanner) Bytes() []byte {
	buffer, ok := s.buffers[s.pid]
	if ok && buffer.isFull() {
		return buffer.Bytes()
	}

	return nil
}

func (s *TableScanner) Table() Table {
	return Table(s.Bytes())
}
