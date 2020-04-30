#include <windows.h>

#include <oleauto.h>

#include "util.h"
#include <cstdio>
#include <memory>
#include <stdarg.h>

namespace exciton {
namespace util {
std::string ToUTF8StringBSTR(BSTR bstr) {
  int iLen;
  iLen =
      ::WideCharToMultiByte(CP_UTF8, 0, static_cast<LPCWSTR>(bstr),
                            ::SysStringLen(bstr), nullptr, 0, nullptr, nullptr);
  std::string data;
  data.resize(iLen, '\0');
  iLen = ::WideCharToMultiByte(CP_UTF8, 0, static_cast<LPCWSTR>(bstr),
                               ::SysStringLen(bstr), data.data(), data.length(),
                               nullptr, nullptr);
  return data;
}

std::string ToUTF8StringWCHAR(const WCHAR *lpStr) {
  int iStrLen = ::wcslen(lpStr);
  int iLen;
  iLen = ::WideCharToMultiByte(CP_UTF8, 0, lpStr, iStrLen, nullptr, 0, nullptr,
                               nullptr);
  std::string data;
  data.resize(iLen, '\0');
  iLen = ::WideCharToMultiByte(CP_UTF8, 0, lpStr, iStrLen, data.data(),
                               data.length(), nullptr, nullptr);
  return data;
}

std::string ToUTF8StringWstr(const std::wstring& str)
{
  int iLen;
  iLen = ::WideCharToMultiByte(CP_UTF8, 0, str.c_str(), str.length(), nullptr, 0, nullptr, nullptr);
  std::string data;
  data.resize(iLen, '\0');
  iLen = ::WideCharToMultiByte(CP_UTF8, 0, str.c_str(), str.length(), data.data(), data.length(), nullptr, nullptr);
  return data;
}

std::wstring ToUTF16String(BSTR bstr) {
  std::wstring data(static_cast<LPCWSTR>(bstr), ::SysStringLen(bstr));
  return data;
}

std::wstring ToUTF16String(const std::string &str) {
  int iLen;
  iLen =
      ::MultiByteToWideChar(CP_UTF8, 0, str.data(), str.length(), nullptr, 0);
  std::wstring data(iLen, '\0');
  iLen = ::MultiByteToWideChar(CP_UTF8, 0, str.data(), str.length(),
                               data.data(), data.length());
  return data;
}

namespace {
std::string FormatStringV(const char *fmt, va_list ap) {
  int n = static_cast<int>(strlen(fmt) * 2 + 1);
  std::unique_ptr<char[]> result;
  for (;;) {
    result.reset(new char[n]);
    strcpy(&result[0], fmt);
    int nn = vsnprintf(&result[0], n, fmt, ap);
    if (nn < 0 || nn >= n) {
      n += abs(nn - n + 1);
    } else {
      break;
    }
  }
  return std::string(&result[0]);
}
std::wstring FormatStringV(const WCHAR *fmt, va_list ap) {
  int n = static_cast<int>(wcslen(fmt) * 2 + 1);
  std::unique_ptr<WCHAR[]> result;
  for (;;) {
    result.reset(new WCHAR[n]);
    wcscpy(&result[0], fmt);
    int nn = vswprintf(&result[0], n, fmt, ap);
    if (nn < 0 || nn >= n) {
      n += abs(nn - n + 1);
    } else {
      break;
    }
  }
  return std::wstring(&result[0]);
}
} // namespace

std::string FormatString(const char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  std::string ret = FormatStringV(fmt, ap);
  va_end(ap);
  return ret;
}

std::string FormatString(const WCHAR *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  std::wstring ret = FormatStringV(fmt, ap);
  va_end(ap);
  return ToUTF8StringWstr(ret);
}

std::wstring FormatStringW(const char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  std::string ret = FormatStringV(fmt, ap);
  va_end(ap);
  return ToUTF16String(ret);
}

std::wstring FormatStringW(const WCHAR *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  std::wstring ret = FormatStringV(fmt, ap);
  va_end(ap);
  return ret;
}

} // namespace util
} // namespace exciton