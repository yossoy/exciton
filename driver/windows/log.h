#pragma once

#ifdef __cplusplus

extern "C" {
#endif

enum EXCITON_LOG_LEVEL{
    EXCITON_LOG_DEBUG = 0,
    EXCITON_LOG_INFO,
    EXCITON_LOG_WARNING,
    EXCITON_LOG_ERROR,
};

extern void Log_Init(void);
extern void GoLogOut(enum EXCITON_LOG_LEVEL lvl, const char* str);

#ifdef __cplusplus
extern void GoLogOutW(enum EXCITON_LOG_LEVEL lvl, const WCHAR* str);
}
#include "util.h"

#if defined(DEBUG)
#define LOG_DEBUG(...)      GoLogOutW(EXCITON_LOG_DEBUG,      exciton::util::FormatStringW(__VA_ARGS__).c_str())
#define LOG_INFO(...)       GoLogOutW(EXCITON_LOG_INFO,       exciton::util::FormatStringW(__VA_ARGS__).c_str())
#define LOG_WARNING(...)    GoLogOutW(EXCITON_LOG_WARNING,    exciton::util::FormatStringW(__VA_ARGS__).c_str())
#define LOG_ERROR(...)      GoLogOutW(EXCITON_LOG_ERROR,      exciton::util::FormatStringW(__VA_ARGS__).c_str())
#elif defined(NDEBUG)
#define LOG_DEBUG(...)
#define LOG_INFO(...)
#define LOG_WARNING(...)    GoLogOutW(EXCITON_LOG_WARNING,    exciton::util::FormatStringW(__VA_ARGS__).c_str())
#define LOG_ERROR(...)      GoLogOutW(EXCITON_LOG_ERROR,      exciton::util::FormatStringW(__VA_ARGS__).c_str())
#else
#error need define DEBUG or NDEBUG
#endif

#endif

