//-*- objc -*-
#pragma once
#import <Cocoa/Cocoa.h>

#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

typedef void (^NativeEventHandler)(
    id argument, NSDictionary<NSString *, NSString *> *parameter,
    int responceNo);

@interface Driver : NSObject <NSApplicationDelegate>
@property NSMenu *dock;
@property NSMutableDictionary<NSString *, id> *elements;
@property NSMutableDictionary<NSString *, NativeEventHandler> *eventHandlers;
@property NSArray *respItems;
@property int lastUseRespItem;

+ (instancetype)current;
- (instancetype)init;
- (BOOL)emitEvent:(NSString *)name;
- (BOOL)emitEvent:(NSString *)name argument:(id)argument;
- (BOOL)emitEvent:(NSString *)name jsonEncodedArgument:(NSString *)argument;
- (void)addEventHandler:(NSString *)name handler:(NativeEventHandler)handler;
- (void)responceEventResult:(int)responceNo result:(id)result;
- (void)responceEventResult:(int)responceNo boolean:(BOOL)boolean;
- (void)responceEventResult:(int)responceNo jsonEncodedArgument:(NSString* )jsonResult;
@end

void Driver_Run();
void Driver_Terminate();
BOOL Driver_EmitEvent(void *bytes, NSUInteger length);
char *Driver_GetBundleResourcesPath();
const char* Driver_GetPreferrdLanguage();
