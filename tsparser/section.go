// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package tsparser

type ProgramAssociationSection struct {
	network    PID
	programMap map[uint16]PID
}

func ParseProgramAssociationSection(table Table) *ProgramAssociationSection {
	sec := new(ProgramAssociationSection)
	sec.programMap = make(map[uint16]PID)

	payload := table.Data()
	for i := 0; i < len(payload); i += 4 {
		programNumber := uint16(payload[i])<<8 | uint16(payload[i+1])
		if programNumber == 0 {
			sec.network = PID(payload[i+2]&0x1f)<<8 | PID(payload[i+3])
		} else {
			sec.programMap[programNumber] = PID(payload[i+2]&0x1f)<<8 | PID(payload[i+3])
		}
	}

	return sec
}

func (s *ProgramAssociationSection) NetworkPID() PID {
	return s.network
}

func (s *ProgramAssociationSection) ProgramMap() map[uint16]PID {
	return s.programMap
}
