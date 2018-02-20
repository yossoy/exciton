#pragma once
#import <Cocoa/Cocoa.h>

enum EXCITON_LOG_LEVEL{
    EXCITON_LOG_DEBUG = 0,
    EXCITON_LOG_INFO,
    EXCITON_LOG_WARNING,
    EXCITON_LOG_ERROR,
};

extern void Log_Init(void);
extern void GoLogOut(enum EXCITON_LOG_LEVEL lvl, NSString* str);
#if defined(DEBUG)
#define LOG_DEBUG(...)      GoLogOut(EXCITON_LOG_DEBUG,      [NSString stringWithFormat:__VA_ARGS__])
#define LOG_INFO(...)       GoLogOut(EXCITON_LOG_INFO,       [NSString stringWithFormat:__VA_ARGS__])
#define LOG_WARNING(...)    GoLogOut(EXCITON_LOG_WARNING,    [NSString stringWithFormat:__VA_ARGS__])
#define LOG_ERROR(...)      GoLogOut(EXCITON_LOG_ERROR,      [NSString stringWithFormat:__VA_ARGS__])
#elif defined(NDEBUG)
extern void GoLogOutC(enum EXCITON_LOG_LEVEL lvl, const char* str);
#define LOG_DEBUG(...)
#define LOG_INFO(...)
#define LOG_WARNING(...)    GoLogOut(EXCITON_LOG_WARNING,    [NSString stringWithFormat:__VA_ARGS__])
#define LOG_ERROR(...)      GoLogOut(EXCITON_LOG_ERROR,      [NSString stringWithFormat:__VA_ARGS__])
#else
#error need define DEBUG or NDEBUG
#endif

