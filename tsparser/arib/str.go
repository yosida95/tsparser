// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.
//
// This source code is a porting of aribstr.py distributed at https://github.com/murakamiy/epgdump_py.
// License of original software can be found in README.rst.

package arib

import (
	"bytes"
	"sort"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
)

type code uint16

const (
	kanji code = iota
	alphanumeric
	hiragana
	katakana
	mosaicA
	mosaicB
	mosaicC
	mosaicD
	propAlphanumeric
	propHiragana
	propKatakana
	jisX0201Katakana
	jisKanjiPlane1
	jisKanjiPlane2
	additionalSymboles
	unsupported
)

type byteSlice []byte

func (bytes byteSlice) Len() int {
	return len(bytes)
}

func (bytes byteSlice) Less(i, j int) bool {
	return bytes[i] < bytes[j]
}

func (bytes byteSlice) Swap(i, j int) {
	bytes[i], bytes[j] = bytes[j], bytes[i]
}

func (bytes byteSlice) Search(code byte) int {
	return sort.Search(len(bytes), func(x int) bool {
		return bytes[x] >= code
	})
}

type codeSet struct {
	code   code
	length int
}

var (
	codeSetG = map[byte]codeSet{
		0x42: codeSet{
			code:   kanji,
			length: 2,
		},
		0x4a: codeSet{
			code:   alphanumeric,
			length: 1,
		},
		0x30: codeSet{
			code:   hiragana,
			length: 1,
		},
		0x31: codeSet{
			code:   katakana,
			length: 1,
		},
		0x32: codeSet{
			code:   mosaicA,
			length: 1,
		},
		0x33: codeSet{
			code:   mosaicB,
			length: 1,
		},
		0x34: codeSet{
			code:   mosaicC,
			length: 1,
		},
		0x35: codeSet{
			code:   mosaicD,
			length: 1,
		},
		0x36: codeSet{
			code:   propAlphanumeric,
			length: 1,
		},
		0x37: codeSet{
			code:   propHiragana,
			length: 1,
		},
		0x38: codeSet{
			code:   propKatakana,
			length: 1,
		},
		0x49: codeSet{
			code:   jisX0201Katakana,
			length: 1,
		},
		0x39: codeSet{
			code:   jisKanjiPlane1,
			length: 1,
		},
		0x3a: codeSet{
			code:   jisKanjiPlane2,
			length: 1,
		},
		0x3b: codeSet{
			code:   additionalSymboles,
			length: 1,
		},
	}

	codeSetDrcs = map[byte]codeSet{
		0x40: codeSet{
			code:   unsupported,
			length: 2,
		},
		0x41: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x42: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x43: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x44: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x45: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x46: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x47: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x48: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x49: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x4a: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x4b: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x4c: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x4d: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x4e: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x4f: codeSet{
			code:   unsupported,
			length: 1,
		},
		0x70: codeSet{
			code:   unsupported,
			length: 1,
		},
	}

	codeSetKeys = make(byteSlice, len(codeSetG)+len(codeSetDrcs))

	aribBASE = map[byte]byte{
		0x79: 0x3c,
		0x7a: 0x23,
		0x7b: 0x56,
		0x7c: 0x57,
		0x7d: 0x22,
		0x7e: 0x26,
	}
	aribHiraganaMap = map[byte]byte{
		0x77: 0x35,
		0x78: 0x36,
	}
	aribKatakanaMap = map[byte]byte{
		0x77: 0x33,
		0x78: 0x34,
	}
)

func init() {
	i := 0
	for key := range codeSetG {
		codeSetKeys[i] = key
		i += 1
	}
	for key := range codeSetDrcs {
		codeSetKeys[i] = key
		i += 1
	}
	sort.Sort(codeSetKeys)

	for key, value := range aribBASE {
		aribHiraganaMap[key] = value
		aribKatakanaMap[key] = value
	}
}

type codeArea uint8

const (
	codeAreaLeft codeArea = iota
	codeAreaRight
)

type escapeSequence []byte

var (
	escSeqAscii   = escapeSequence{0x1b, 0x28, 0x42}
	escSeqZenkaku = escapeSequence{0x1b, 0x24, 0x42}
	seqSeqHankaku = escapeSequence{0x1b, 0x28, 0x49}
)

type bufferIndex int8

const (
	bufferUnset bufferIndex = iota - 1
	bufferG0
	bufferG1
	bufferG2
	bufferG3
)

type strBuffer struct {
	data   []byte
	escSeq []byte
}

func (b *strBuffer) extend(data []byte) {
	n := len(b.data)
	m := n + len(data)

	if m > cap(b.data) {
		newSlice := make([]byte, m*2)
		copy(newSlice, b.data)
		b.data = newSlice
	}

	b.data = b.data[0:m]
	copy(b.data[n:m], data)
}

func (b *strBuffer) appendStr(escSeq escapeSequence, data ...byte) {
	if !bytes.Equal(b.escSeq, escSeq) { // XXX: ???
		b.extend(escSeq)
	}

	b.escSeq = escSeq
	b.extend(data)
}

func (b *strBuffer) String() string {
	d := japanese.ISO2022JP.NewDecoder()

	result, _, err := transform.Bytes(d, b.data)
	if err != nil {
		panic(err)
	}

	return string(result)
}

type strController struct {
	escSeqCount    int
	escDracs       bool
	escBufferIndex bufferIndex

	shingleShift bufferIndex
	graphicLeft  bufferIndex
	graphicRight bufferIndex

	vBuffer map[bufferIndex]codeSet
}

func newController() *strController {
	return &strController{
		escSeqCount:    0,
		escDracs:       false,
		escBufferIndex: bufferG0,

		shingleShift: bufferUnset,
		graphicLeft:  bufferG0,
		graphicRight: bufferG2,

		vBuffer: map[bufferIndex]codeSet{
			bufferG0: codeSetG[0x42],
			bufferG1: codeSetG[0x4a],
			bufferG2: codeSetG[0x30],
			bufferG3: codeSetG[0x31],
		},
	}
}

func (c *strController) getCurrentCode(b byte) (codeSet, bool) {
	var indicator bufferIndex

	if 0x21 <= b && b <= 0x7e {
		if c.shingleShift == bufferUnset {
			indicator = c.graphicLeft
		} else {
			indicator = c.shingleShift
			c.shingleShift = bufferUnset
		}
	} else if 0xa1 <= b && b <= 0xfe {
		indicator = c.graphicRight
	} else {
		return codeSet{}, false
	}

	codeSet, ok := c.vBuffer[indicator]
	return codeSet, ok
}

func (c *strController) invoke(bufferIndex bufferIndex, area codeArea, lockingShift bool) {
	if area == codeAreaLeft {
		if lockingShift {
			c.graphicLeft = bufferIndex
		} else {
			c.shingleShift = bufferIndex
		}
	} else if area == codeAreaRight {
		c.graphicRight = bufferIndex
	}

	c.escSeqCount = 0
	return
}

func (c *strController) degignate(b byte) {
	i := codeSetKeys.Search(b)
	if codeSetKeys.Len() == i || codeSetKeys[i] != b {
		panic("not found")
	}

	if c.escDracs {
		c.vBuffer[c.escBufferIndex] = codeSetDrcs[b]
	} else {
		c.vBuffer[c.escBufferIndex] = codeSetG[b]
	}

	c.escSeqCount = 0
}

func (c *strController) setEscape(bufferIndex bufferIndex, drcs bool) {
	if bufferIndex != bufferUnset {
		c.escBufferIndex = bufferIndex
	}

	c.escDracs = drcs
	c.escSeqCount += 1
}

type strConverter struct {
	controller *strController
	jisArray   *strBuffer
}

func NewStrConverter() *strConverter {
	return &strConverter{
		controller: newController(),
		jisArray:   new(strBuffer),
	}
}

func (c *strConverter) Convert(data []byte) {
	bytes := len(data)

	for i := 0; i < bytes; i++ {
		b := data[i]
		if c.controller.escSeqCount > 0 {
			c.doEscape(b)
			continue
		}

		if 0x21 <= b && b <= 0x7e || 0xa1 <= b && b <= 0xfe {
			code, ok := c.controller.getCurrentCode(b)
			if !ok {
				panic("!ok")
			}

			char := b
			char2 := byte(0x00)
			if code.length == 2 {
				i += 1
				if i == bytes {
					return
				}
				char2 = data[i]
			}
			if 0xa1 <= char && char <= 0xfe {
				char &= 0x7f
				char2 &= 0x7f
			}

			c.doConvert(code.code, char, char2)
			continue
		}

		switch b {
		case 0x20, 0xa0, 0x09:
			c.jisArray.appendStr(escSeqAscii, 0x20)
		case 0x0d, 0x0a:
			c.jisArray.appendStr(escSeqAscii, 0x0a)
		default:
			c.doControl(b)
		}
	}
}

func (c *strConverter) doEscape(b byte) {
	switch c.controller.escSeqCount {
	case 1:
		switch b {
		case 0x6e:
			c.controller.invoke(bufferG2, codeAreaLeft, true)
		case 0x6f:
			c.controller.invoke(bufferG3, codeAreaLeft, true)
		case 0x7e:
			c.controller.invoke(bufferG1, codeAreaRight, true)
		case 0x7d:
			c.controller.invoke(bufferG2, codeAreaRight, true)
		case 0x7c:
			c.controller.invoke(bufferG3, codeAreaRight, true)
		case 0x24, 0x28:
			c.controller.setEscape(bufferG0, false)
		case 0x29:
			c.controller.setEscape(bufferG1, false)
		case 0x2a:
			c.controller.setEscape(bufferG2, false)
		case 0x2b:
			c.controller.setEscape(bufferG3, false)
		default:
			panic("unknown")
		}
	case 2:
		switch b {
		case 0x20:
			c.controller.setEscape(bufferUnset, true)
		case 0x28:
			c.controller.setEscape(bufferG0, false)
		case 0x29:
			c.controller.setEscape(bufferG1, false)
		case 0x2a:
			c.controller.setEscape(bufferG2, false)
		case 0x2b:
			c.controller.setEscape(bufferG3, false)
		default:
			c.controller.degignate(b)
		}
	case 3:
		if b == 0x20 {
			c.controller.setEscape(bufferUnset, true)
		} else {
			c.controller.degignate(b)
		}
	case 4:
		c.controller.degignate(b)
	}
}

func (c *strConverter) doConvert(code code, char, char2 byte) {
	switch code {
	case kanji, jisKanjiPlane1, jisKanjiPlane2:
		c.jisArray.appendStr(escSeqZenkaku, char, char2)
	case alphanumeric, propAlphanumeric:
		c.jisArray.appendStr(escSeqAscii, char)
	case hiragana, propHiragana:
		if char >= 0x77 {
			c.jisArray.appendStr(escSeqZenkaku, 0x21, aribHiraganaMap[char])
		} else {
			c.jisArray.appendStr(escSeqZenkaku, 0x24, char)
		}
	case katakana, propKatakana:
		if char >= 0x77 {
			c.jisArray.appendStr(escSeqZenkaku, 0x21, aribKatakanaMap[char])
		} else {
			c.jisArray.appendStr(escSeqZenkaku, 0x25, char)
		}
	case jisX0201Katakana:
		c.jisArray.appendStr(seqSeqHankaku, char)
	case additionalSymboles:
		// TODO
	}
}

func (c *strConverter) doControl(b byte) {
	switch b {
	case 0x0f:
		c.controller.invoke(bufferG0, codeAreaLeft, true)
	case 0x0e:
		c.controller.invoke(bufferG1, codeAreaLeft, true)
	case 0x19:
		c.controller.invoke(bufferG2, codeAreaLeft, false)
	case 0x1d:
		c.controller.invoke(bufferG3, codeAreaLeft, false)
	case 0x1b:
		c.controller.escSeqCount = 1
	}
}
