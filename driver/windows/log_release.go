// +build release

package windows

import (
	"fmt"
	"unsafe"
)

/*
#include "log.h"
#include <stdlib.h>
*/
import "C"

func driverLogDebug(msg string, args ...interface{}) {
}

func driverLogInfo(msg string, args ...interface{}) {
}

func driverLogWarning(msg string, args ...interface{}) {
	str := C.CString(fmt.Sprintf(msg, args...))
	C.GoLogOut(C.EXCITON_LOG_WARNING, str)
	C.free(unsafe.Pointer(str))
}

func driverLogError(msg string, args ...interface{}) {
	str := C.CString(fmt.Sprintf(msg, args...))
	C.GoLogOut(C.EXCITON_LOG_ERROR, str)
	C.free(unsafe.Pointer(str))
}
