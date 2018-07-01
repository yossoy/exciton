// +build !release

package log

import (
	"github.com/yossoy/exciton/driver"
)

func PrintDebug(fmt string, args ...interface{}) {
	driver.Log(driver.LogLevelDebug, fmt, args...)
}

func PrintInfo(fmt string, args ...interface{}) {
	driver.Log(driver.LogLevelInfo, fmt, args...)
}
