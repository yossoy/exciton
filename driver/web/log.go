package web

import (
	"fmt"
	"log"
)

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
