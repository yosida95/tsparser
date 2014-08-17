// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

import (
	"errors"
)

var (
	ErrInvalidTableLength   = errors.New("Invalid table length")
	ErrInvalidSectionLength = errors.New("Invalid value of section_length")
)

type TableId uint8

const (
	ProgramAssociationTable TableId = 0x00
)

type Table []byte

func (t Table) TableId() TableId {
	return TableId(t[0])
}

func (t Table) SectionSyntaxIndicator() bool {
	return t[1]&0x80 > 0
}

func (t Table) PrivateIndicator() bool {
	return t[1]&0x40 > 0
}

func (t Table) dataStartsAt() int {
	start := 3
	if t.SectionSyntaxIndicator() {
		start += 5
	}

	return start
}

func (t Table) dataEndsAt() int {
	end := t.SectionLength() + 3
	if t.SectionSyntaxIndicator() {
		end -= 4
	}

	return end
}

func (t Table) SectionLength() int {
	length := int(t[1]&0x0f)<<8 | int(t[2])
	return length
}

func (t Table) TableIdExtension() uint16 {
	if t.SectionSyntaxIndicator() {
		return uint16(t[3])<<8 | uint16(t[4])
	}

	return 0
}

func (t Table) VersionNumber() uint8 {
	if t.SectionSyntaxIndicator() {
		return uint8(t[5]&0x3e) >> 1
	}

	return 0
}

func (t Table) CurrentNextIndicator() bool {
	if t.SectionSyntaxIndicator() {
		return t[5]&0x01 > 1
	}

	return false
}

func (t Table) SectionNumber() uint8 {
	if t.SectionSyntaxIndicator() {
		return 0
	}

	return uint8(t[6])
}

func (t Table) LastSectionNumber() uint8 {
	if t.SectionSyntaxIndicator() {
		return 0
	}

	return uint8(t[7])
}

func (t Table) CRC32() []byte {
	if t.SectionSyntaxIndicator() {
		end := t.dataEndsAt()
		return t[end : end+4]
	}

	return nil
}

func (t Table) Data() []byte {
	return t[t.dataStartsAt():t.dataEndsAt()]
}

func (t Table) validate() error {
	if len(t) < 3 {
		return ErrInvalidTableLength
	}

	start := t.dataStartsAt()
	end := t.dataEndsAt()
	if start > end || len(t) < end {
		return ErrInvalidSectionLength
	}

	return nil
}
