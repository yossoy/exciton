#pragma once

#include <string>
#include <wtypes.h>

namespace exciton {
namespace util {
std::string ToUTF8StringBSTR(BSTR bstr);
std::string ToUTF8StringWCHAR(const WCHAR* lpstr);
std::string ToUTF8StringWstr(const std::wstring& str);

std::wstring ToUTF16String(BSTR bstr);
std::wstring ToUTF16String(const std::string& str);

std::string FormatString(const char* fmt, ...);
std::string FormatString(const WCHAR* fmt, ...);
std::wstring FormatStringW(const char* fmt, ...);
std::wstring FormatStringW(const WCHAR* fmt, ...);
}
} // namespace exciton