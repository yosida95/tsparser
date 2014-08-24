package tsparser

var (
	crc32Table [256]uint32
)

func AsUint32(data []byte) (n uint32) {
	bytes := len(data)
	if bytes > 4 {
		panic("Overflow")
	}

	for i := 0; i < bytes; i++ {
		n |= uint32(data[bytes-i-1]) << uint(i*8)
	}
	return
}

func CheckCRC32(payload []byte) bool {
	return updateCRC32(0xffffffff, &crc32Table, payload) == 0
}

func updateCRC32(crc uint32, table *[256]uint32, payload []byte) uint32 {
	for _, b := range payload {
		crc = crc32Table[byte(crc>>24)^b] ^ (crc << 8)
	}
	return crc
}

func fillCRC32Table(polynomial uint32, table *[256]uint32) {
	for i := 0; i < 256; i++ {
		crc := uint32(i) << 24
		for j := 0; j < 8; j++ {
			if crc&0x80000000 == 0x80000000 {
				crc = (crc << 1) ^ polynomial
			} else {
				crc <<= 1
			}
		}
		table[i] = crc
	}

	return
}

func init() {
	fillCRC32Table(0x04c11db7, &crc32Table)
}
