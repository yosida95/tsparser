package tsparser

type DescriptorTag uint8

type Descriptor []byte

func ParseDescriptor(data []byte) (d Descriptor, length int) {
	length = int(data[1]) + 2
	d = Descriptor(data[0:length])
	return
}

func (d Descriptor) Tag() DescriptorTag {
	return DescriptorTag(d[0])
}

func (d Descriptor) Length() uint8 {
	return uint8(d[1])
}

func (d Descriptor) Payload() []byte {
	return d[2 : d.Length()+2]
}
