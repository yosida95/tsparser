// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

import (
	"sort"
)

type PIDSlice []PID

func (pids PIDSlice) Len() int {
	return len(pids)
}

func (pids PIDSlice) Less(i, j int) bool {
	return pids[i] < pids[j]
}

func (pids PIDSlice) Swap(i, j int) {
	pids[i], pids[j] = pids[j], pids[i]
}

func (pids PIDSlice) Search(pid PID) int {
	return sort.Search(len(pids), func(x int) bool {
		return pids[x] >= pid
	})
}

type PacketFilter struct {
	s    PacketStream
	pids PIDSlice
	next Packet
}

func NewPacketFilter(s PacketStream, pids ...PID) *PacketFilter {
	pidSlice := PIDSlice(pids)
	sort.Sort(pidSlice)

	return &PacketFilter{
		s:    s,
		pids: pidSlice,
	}
}

func (f *PacketFilter) isTarget(pid PID) bool {
	i := f.pids.Search(pid)
	return f.pids.Len() > i && f.pids[i] == pid
}

func (f *PacketFilter) Scan() bool {
	for f.s.Scan() {
		packet := f.s.Packet()
		if f.isTarget(packet.PID()) {
			f.next = packet
			return true
		}
	}

	return false
}

func (f *PacketFilter) Bytes() []byte {
	return []byte(f.next)
}

func (f *PacketFilter) Packet() Packet {
	return f.next
}
