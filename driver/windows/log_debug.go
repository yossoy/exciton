// +build !release

package windows

//TODO: ログ出力時に、goのログ出力を行う
import (
	"fmt"
	"log"
)

/*
#include "log.h"
*/
import "C"

//export goDebugLog
func goDebugLog(lvl C.enum_EXCITON_LOG_LEVEL, cstr *C.char) {
	str := C.GoString(cstr)
	switch lvl {
	case C.EXCITON_LOG_DEBUG:
		driverLogDebug(str)
	case C.EXCITON_LOG_INFO:
		driverLogInfo(str)
	case C.EXCITON_LOG_WARNING:
		driverLogWarning(str)
	case C.EXCITON_LOG_ERROR:
		driverLogError(str)
	}
}

func driverLogDebug(str string, args ...interface{}) {
	log.Print("[DEBUG]", fmt.Sprintf(str, args...))
}
func driverLogInfo(str string, args ...interface{}) {
	log.Print("[INFO]", fmt.Sprintf(str, args...))
}
func driverLogWarning(str string, args ...interface{}) {
	log.Print("[WARNING]", fmt.Sprintf(str, args...))
}
func driverLogError(str string, args ...interface{}) {
	log.Print("[ERROR]", fmt.Sprintf(str, args...))
}
