package tsparser

import (
	"errors"
	"fmt"
	"io"
	"log"
)

var (
	ErrInvalidPointer  = errors.New("Invalid value of pointer_field")
	ErrPacketScrambled = errors.New("Scrambled")
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
	logger *log.Logger
}

func NewPacketScanner(r io.Reader, logger *log.Logger) *PacketScanner {
	return &PacketScanner{
		r:      r,
		seek:   0,
		eof:    false,
		logger: logger,
	}
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
	current Table
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

func (b *tableScannerBuffer) freeze() (err error) {
	if len(b.data) > 0 {
		table := make(Table, len(b.data))
		copy(table, b.data)
		if err = table.validate(); err != nil {
			return
		}

		b.current = table
	}

	return
}

func (b *tableScannerBuffer) Begin(payload []byte) (err error) {
	packetStartCodePrefix := payload[0]<<16 | payload[1]<<8 | payload[2]
	b.isPES = packetStartCodePrefix == 0x000001

	if !b.isPES {
		pointerField := int(payload[0])
		if pointerField > 182 {
			b.clear()

			err = ErrInvalidPointer
			return
		}

		if len(b.data) > 0 {
			b.extend(payload[1 : 1+pointerField])
		}
		payload = payload[1+pointerField:]
	}

	err = b.freeze()
	b.clear()
	b.extend(payload)
	return
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

func (b *tableScannerBuffer) Table() Table {
	return b.current
}

type TableScanner struct {
	s       PacketStream
	pid     PID
	buffers map[PID]*tableScannerBuffer
	logger  *log.Logger
}

func NewTableScanner(s PacketStream, l *log.Logger) *TableScanner {
	return &TableScanner{
		s:       s,
		buffers: make(map[PID]*tableScannerBuffer),
		logger:  l,
	}
}

func (s *TableScanner) log(p Packet, v ...interface{}) {
	if s.logger == nil {
		return
	}

	values := make([]interface{}, len(v)+1)
	values[0] = fmt.Sprintf("pid=0x%04x: ", p.PID())
	copy(values[1:], v)
	s.logger.Print(values...)
}

func (s *TableScanner) Scan() bool {
	for s.s.Scan() {
		packet := s.s.Packet()
		if !packet.HasPayload() {
			continue
		} else if packet.transportScramblingControl() > 0 {
			s.log(packet, ErrPacketScrambled)
			continue
		}

		buffer, ok := s.buffers[packet.PID()]
		if !ok {
			s.buffers[packet.PID()] = new(tableScannerBuffer)
			buffer = s.buffers[packet.PID()]
		}

		if packet.payloadUnitStartIndicator() {
			if err := buffer.Begin(packet.Payload()); err != nil {
				s.log(packet, err)
				continue
			}

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

func (s *TableScanner) Table() Table {
	buffer, ok := s.buffers[s.pid]
	if ok && buffer.isFull() {
		return Table(buffer.Table())
	}

	return nil
}
