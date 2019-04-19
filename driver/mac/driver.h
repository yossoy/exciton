//-*- objc -*-
#pragma once
#import <Cocoa/Cocoa.h>

#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

typedef void (^NativeEventHandler)(
    id argument, NSDictionary<NSString *, NSString *> *parameter,
    int responceNo);
typedef void (^NativeResponceCallback)(id value, NSString* error);

@interface Driver : NSObject <NSApplicationDelegate>
@property NSMenu *dock;
@property NSMutableDictionary<NSString *, id> *elements;
@property NSMutableDictionary<NSString *, NSMutableDictionary<NSString*, NativeEventHandler>*> *eventHandlers;
@property NSMutableArray *respItems;
@property int lastUseRespItem;

+ (instancetype)current;
- (instancetype)init;
- (BOOL)emitEvent:(NSString* )target name:(NSString *)name;
- (BOOL)emitEvent:(NSString* )target name:(NSString *)name argument:(id)argument;
- (BOOL)emitEvent:(NSString* )target name:(NSString *)name jsonEncodedArgument:(NSString *)argument;
- (void)emitEvent:(NSString* )target name:(NSString *)name argument:(id)argument respCallback:(NativeResponceCallback)respCallback;
- (void)addEventHandler:(NSString *)path name:(NSString *)name handler:(NativeEventHandler)handler;
- (void)responceEventResult:(int)responceNo result:(id)result;
- (void)responceEventResult:(int)responceNo boolean:(BOOL)boolean;
- (void)responceEventResult:(int)responceNo jsonEncodedArgument:(NSString* )jsonResult;
@end

void Driver_Run();
void Driver_Terminate();
BOOL Driver_EmitEvent(void *bytes, NSUInteger length);
BOOL Driver_ResponseEvent(NSInteger respNo, void* bytes, NSUInteger length);
char *Driver_GetBundleResourcesPath();
const char* Driver_GetPreferrdLanguage();
