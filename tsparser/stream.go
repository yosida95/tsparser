// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

type PacketStream interface {
	Scan() bool
	Packet() Packet
}

type TableStream interface {
	Scan() bool
	Table() Table
}
