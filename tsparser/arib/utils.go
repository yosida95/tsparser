// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package arib

import (
	"math"
	"time"
)

func bcd2int(b byte) int {
	return int(b&0xf0)>>4*10 + int(b&0x0f)
}

func parseJSTTime(payload []byte) time.Time {
	mjd := float64(uint16(payload[0])<<8 | uint16(payload[1]))

	y := math.Floor((mjd - 15078.2) / 365.25)
	m := math.Floor((mjd - 14956.1 - math.Floor(y*365.25)) / 30.6001)
	d := mjd - 14956 - math.Floor(y*365.25) - math.Floor(m*30.6001)

	k := 0
	if m == 14 || m == 15 {
		k = 1
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")
	return time.Date(
		int(y)+k+1900, time.Month(int(m)-1-k*12), int(d),
		bcd2int(payload[2]), bcd2int(payload[3]), bcd2int(payload[4]), 0,
		loc)
}
