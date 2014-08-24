// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

type DescriptorTag uint8

type Descriptor []byte

func ParseDescriptor(data []byte) (d Descriptor, length int) {
	length = int(data[1]) + 2
	d = Descriptor(data[0:length])
	return
}

func ParseDescriptors(data []byte) []Descriptor {
	result := make([]Descriptor, 0)
	bytes := len(data)

	for i := 0; i < bytes; {
		d, spent := ParseDescriptor(data[i:bytes])
		i += spent

		result = append(result, d)
	}

	return result
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
