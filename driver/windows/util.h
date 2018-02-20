#pragma once

#include <string>
#include <wtypes.h>

namespace exciton {
namespace util {
std::string ToUTF8String(BSTR bstr);
std::string ToUTF8String(const WCHAR* lpstr);

std::wstring ToUTF16String(BSTR bstr);
std::wstring ToUTF16String(const std::string& str);

std::string FormatString(const char* fmt, ...);
std::string FormatString(const WCHAR* fmt, ...);
std::wstring FormatStringW(const char* fmt, ...);
std::wstring FormatStringW(const WCHAR* fmt, ...);
}
} // namespace exciton