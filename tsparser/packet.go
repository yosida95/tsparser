// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

type PID uint16

type Packet []byte

func (p Packet) transportErrorIndicator() bool {
	return uint8(p[1]&0x80)>>7 == 1
}

func (p Packet) payloadUnitStartIndicator() bool {
	return uint8(p[1]&0x40)>>6 == 1
}

func (p Packet) transportPriority() uint8 {
	return uint8(p[1]&0x20) >> 5
}

func (p Packet) PID() PID {
	return PID(p[1]&0x1f)<<8 | PID(p[2])
}

func (p Packet) transportScramblingControl() uint8 {
	return uint8(p[3]&0xc0) >> 6
}

func (p Packet) adaptationFieldControl() uint8 {
	return uint8(p[3]&0x30) >> 4
}

func (p Packet) HasAdaptationField() bool {
	return p.adaptationFieldControl() > 1
}

func (p Packet) HasPayload() bool {
	return p.adaptationFieldControl()&1 > 0
}

func (p Packet) continuityCounter() uint8 {
	return uint8(p[3] & 0x0f)
}

func (p Packet) AdaptationField() []byte {
	if !p.HasAdaptationField() {
		return nil
	}

	length := p[4]
	if p.HasPayload() {
		if length > 182 {
			return nil
		}
	} else {
		if length != 183 {
			return nil
		}
	}
	return p[4 : 5+length]
}

func (p Packet) Payload() []byte {
	if !p.HasPayload() {
		return nil
	}

	offset := 4
	if p.HasAdaptationField() {
		af := p.AdaptationField()
		if af == nil {
			return nil
		}
		offset += len(p.AdaptationField())
	}

	return p[offset:]
}
