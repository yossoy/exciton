#include "log.h"
#include "_cgo_export.h"

#if defined(DEBUG)
void Log_Init() {}
void GoLogOut(enum EXCITON_LOG_LEVEL lvl, NSString *str) {
  goDebugLog(lvl, (char *)[str UTF8String]);
}
#else
#import <os/log.h>
static os_log_t os_log_;
void Log_Init() {
  NSString *bundleIdentifier = [[NSBundle mainBundle] bundleIdentifier];
  if (bundleIdentifier) {
    os_log_ = os_log_create([bundleIdentifier UTF8String], "exciton_logging");
  } else {
    os_log_ = OS_LOG_DEFAULT;
  }
}
void GoLogOut(enum EXCITON_LOG_LEVEL lvl, NSString *str) {
  const char* s = [str UTF8String];
  GoLogOutC(lvl, s);
}
void GoLogOutC(enum EXCITON_LOG_LEVEL lvl, const char* s) {
  if (!s) s = "NULL";
  switch (lvl) {
  case EXCITON_LOG_ERROR:
    os_log_error(os_log_, "%{public}s", s);
    break;
  case EXCITON_LOG_INFO:
    os_log_info(os_log_, "%{public}s", s);
    break;
  case EXCITON_LOG_DEBUG:
    os_log_debug(os_log_, "%{public}s", s);
    break;
  case EXCITON_LOG_WARNING:
    os_log_error(os_log_, "%{public}s", s);
    break;
  }
}
#endif
