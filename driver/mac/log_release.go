// +build release

package mac

import (
	"fmt"
	"unsafe"
)

//TODO: ログ出力時に、goのログ出力を行う

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
	C.GoLogOutC(C.EXCITON_LOG_WARNING, str)
	C.free(unsafe.Pointer(str))
}
func driverLogError(msg string, args ...interface{}) {
	str := C.CString(fmt.Sprintf(msg, args...))
	C.GoLogOutC(C.EXCITON_LOG_ERROR, str)
	C.free(unsafe.Pointer(str))
}
