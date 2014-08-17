// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package arib

import (
	"time"

	"github.com/yosida95/tsparser/tsparser"
)

func ParseTimeDateSection(table tsparser.Table) time.Time {
	JSTTime := table.Data()
	return parseJSTTime(JSTTime)
}
