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

type Scanner struct {
	r      io.Reader
	buffer [PacketSize * BufferedPacketCount]byte
	seek   int
	eof    bool
	err    error
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:    r,
		seek: 0,
		eof:  false,
		err:  nil,
	}
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) leftShift(n int) {
	copy(s.buffer[:BufferSize-n], s.buffer[n:])
	s.seek -= n
}

func (s *Scanner) rightShift(n int) {
	copy(s.buffer[n:], s.buffer[:BufferSize-n])
	s.seek += n
}

func (s *Scanner) fillBuffer(n int) bool {
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

func (s *Scanner) isSynced() bool {
	synced := true
	for i := s.seek; i < BufferSize; i += PacketSize {
		if s.buffer[i] != SyncByte {
			synced = false
			break
		}
	}

	return synced
}

func (s *Scanner) sync() bool {
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

func (s *Scanner) Scan() bool {
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

func (s *Scanner) Bytes() []byte {
	if s.seek < BufferSize {
		return s.buffer[s.seek : s.seek+PacketSize]
	}

	return nil
}
