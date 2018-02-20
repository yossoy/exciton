#include <windows.h>

#include "log.h"
#include "util.h"

#include "_cgo_export.h"

void Log_Init() {}
#if defined(DEBUG)
void GoLogOut(enum EXCITON_LOG_LEVEL lvl, const char *str) {
  goDebugLog(lvl, (char *)str);
}
void GoLogOutW(enum EXCITON_LOG_LEVEL lvl, const WCHAR *str) {
  GoLogOut(lvl, exciton::util::ToUTF8String(str).c_str());
}
#else
void GoLogOut(enum EXCITON_LOG_LEVEL lvl, const char *str) {
  GoLogOutW(lvl, exciton::util::ToUTF16String(std::string(str)).c_str());
}
void GoLogOutW(enum EXCITON_LOG_LEVEL lvl, const WCHAR *str) {
  std::wstring wstr;
  switch (lvl) {
  case EXCITON_LOG_DEBUG:
    wstr = L"[DEBUG]\t";
    break;
  case EXCITON_LOG_INFO:
    wstr = L"[INFO]\t";
    break;
  case EXCITON_LOG_WARNING:
    wstr = L"[WARNING]\t";
    break;
  case EXCITON_LOG_ERROR:
    wstr = L"[ERROR]\t";
    break;
  }
  wstr += str;
  ::OutputDebugStringW(wstr.c_str());
}
#endif
