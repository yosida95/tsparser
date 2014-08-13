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
	ok     bool
	data   []byte
	length int
}

func (b *tableScannerBuffer) extendData(right []byte) {
	left := b.data
	b.data = make([]byte, len(left)+len(right))
	copy(b.data, left)
	copy(b.data[len(left):], right)
}

func (b *tableScannerBuffer) Extend(packet Packet) {
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

func (b *tableScannerBuffer) OK() bool {
	return b.ok
}

func (b *tableScannerBuffer) Bytes() []byte {
	if b.OK() {
		return b.data[b.length:]
	}

	return nil
}

type TableScanner struct {
	s       PacketStream
	pid     PID
	buffers map[PID]*tableScannerBuffer
}

func NewTableScanner(s PacketStream) *TableScanner {
	buffers := make(map[PID]*tableScannerBuffer)
	return &TableScanner{
		s:       s,
		buffers: buffers,
	}
}

func (f *TableScanner) Scan() bool {
	for f.s.Scan() {
		packet := f.s.Packet()

		buffer, ok := f.buffers[packet.PID()]
		if !ok {
			f.buffers[packet.PID()] = new(tableScannerBuffer)
			buffer = f.buffers[packet.PID()]
		}

		buffer.Extend(packet)
		if buffer.OK() {
			f.pid = packet.PID()
			return true
		}
	}

	return false
}

func (f *TableScanner) Bytes() []byte {
	buffer := f.buffers[f.pid]
	if buffer.OK() {
		return buffer.Bytes()
	}

	return nil
}

func (f *TableScanner) Table() Table {
	return Table(f.Bytes())
}
