package tsparser

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

func (t Table) SectionLength() uint16 {
	return uint16(t[1]&0x0f)<<8 | uint16(t[2])
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
		length := t.SectionLength()
		return t[length-1 : length+3]
	}

	return nil
}

func (t Table) Data() []byte {
	start := 3
	if t.SectionSyntaxIndicator() {
		start += 5
	}

	end := int(t.SectionLength()) + 3
	if t.SectionSyntaxIndicator() {
		end -= 4
	}

	return t[start:end]
}
