package arib

import (
	"time"

	"github.com/yosida95/tsparser/tsparser"
)

func ParseTimeDateSection(table tsparser.Table) time.Time {
	JSTTime := table.Data()
	return parseJSTTime(JSTTime)
}
