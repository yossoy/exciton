#include "driver.h"
#include "_cgo_export.h"
#include "json.h"
#include "sandbox.h"
//#include "menu.h"
#include "log.h"
#include "accelerator.h"

@implementation Driver
+ (instancetype)current {
  static Driver *driver = nil;

  @synchronized(self) {
    if (driver == nil) {
      driver = [[Driver alloc] init];
      NSApplication *app = [NSApplication sharedApplication];
      app.delegate = driver;
    }
  }
  return driver;
}

- (instancetype)init {
  self = [super init];
  if (self) {
    self.elements = [NSMutableDictionary dictionaryWithCapacity:256];
    self.respItems = [NSMutableArray array];
    self.lastUseRespItem = 0;
    self.dock = [[NSMenu alloc] initWithTitle:@""];

    self.eventHandlers = [@{} mutableCopy];

    [Accelerator current];

    [self addEventHandler:@"/app" name:@"about" handler:^(id argument, NSDictionary<NSString*, NSString*>* parameter, int responseNo) {
      defer([[NSApplication sharedApplication] orderFrontStandardAboutPanel:nil];);
    }];
    

  }

  return self;
}

- (BOOL)emitEvent:(NSString*) target name:(NSString *)name {
  return [self emitEvent:target name:name argument:[NSNull null]];
}

- (BOOL)emitEvent:(NSString *)target name:(NSString*)name argument:(id)argument {
  NSDictionary *driverEvent = @{
    @"target" : target,
    @"name" : name,
    @"argument" : argument,
    @"respCallbackNo" : @-1,
  };
  NSData *json = [JSONEncoder encodeFromObject:driverEvent];
  requestEventEmit((void *)[json bytes], [json length]);
  return TRUE;
}

- (BOOL)emitEvent:(NSString *)target name:(NSString*)name jsonEncodedArgument:(NSString *)argument {
  NSData *json =
      [[NSString stringWithFormat:
                     @"{\"target\":\"%@\",\"name\":\"%@\",\"argument\":%@,\"respCallbackNo\":-1}",
                     target, name, argument] dataUsingEncoding:NSUTF8StringEncoding];
  requestEventEmit((void *)[json bytes], [json length]);
  return TRUE;
}

- (void)emitEvent:(NSString* )target name:(NSString *)name argument:(id)argument respCallback:(NativeResponceCallback)respCallback {
  NSInteger respCallbackNo = -1;
  @synchronized(self.respItems) {
    NSInteger idx;
    for (idx = 0; idx < self.respItems.count; idx++) {
      if (!self.respItems[idx]) {
        respCallbackNo = idx;
        break;
      }
    }
    if (0 <= respCallbackNo) {
      self.respItems[respCallbackNo] = respCallback;
    } else {
      respCallbackNo = self.respItems.count;
      [self.respItems addObject:respCallback];
    }
  }
  NSDictionary *driverEvent = @{
    @"target" : target,
    @"name" : name,
    @"argument" : argument,
    @"respCallbackNo" : [NSNumber numberWithInteger:respCallbackNo],
  };
  NSData *json = [JSONEncoder encodeFromObject:driverEvent];
  requestEventEmit((void *)[json bytes], [json length]);
}

- (BOOL)emitEvent:(void *)bytes length:(NSUInteger)length {
  NSDictionary *driverEvent = [JSONDecoder decodeFromBytes:bytes length:length];
  NSString *targetPath = driverEvent[@"target"];
  NSString *name = driverEvent[@"name"];
  NativeEventHandler handler = nil;
  @synchronized(self.eventHandlers) {
    NSMutableDictionary* d = [self.eventHandlers objectForKey:targetPath];
    if (d != nil) {
      handler = [d objectForKey:name];
    }
  }
  if (handler == nil) {
    LOG_ERROR(@"driver.emitEvent: handler not found: %@/%@", targetPath, name);
    return FALSE;
  }
  LOG_DEBUG(@"driver.emitEvent: %p", handler);
  LOG_DEBUG(@"driver.emitEvent: dispatch event: %@", driverEvent[@"name"]);
  handler(driverEvent[@"argument"], driverEvent[@"parameter"],
          [driverEvent[@"respCallbackNo"] intValue]);

  return TRUE;
}


- (BOOL)responseEvent:(void* )bytes length:(NSUInteger)length respCallbackNo:(NSInteger)respCallbackNo {
  NativeResponceCallback callback = NULL;
  @synchronized(self.respItems) {
    if (respCallbackNo < self.respItems.count) {
      if (self.respItems[respCallbackNo] != [NSNull null]) {
        callback = self.respItems[respCallbackNo];
      }
      self.respItems[respCallbackNo] = [NSNull null];
    }
  }
  if (!callback) {
    return NO;
  }
  NSDictionary* result = [JSONDecoder decodeFromBytes:bytes length:length];
  id value = result[@"result"];
  NSString* err = result[@"error"];
  callback(value, err);
  return YES;
}

- (void)addEventHandler:(NSString *)path name:(NSString*)name handler:(NativeEventHandler)handler {
  // LOG_INFO(@"addEventHandler called: %@", name);
  @synchronized(self.eventHandlers) {
    NSMutableDictionary* d = self.eventHandlers[path];
    if (d == nil ){
      d = [@{} mutableCopy];
      self.eventHandlers[path] = d;
    }
    d[name] = handler;
  }
}

- (void)responceEventResult:(int)responceNo result:(id)result {
  NSData *resultData = [JSONEncoder encodeFromObject:result];
  responceEventResult(responceNo, (void *)[resultData bytes],
                      [resultData length]);
}
- (void)responceEventResult:(int)responceNo boolean:(BOOL)boolean {
  NSData *resultData = [JSONEncoder encodeBool:boolean];
  responceEventResult(responceNo, (void *)[resultData bytes],
                      [resultData length]);
}

- (void)responceEventResult:(int)responceNo jsonEncodedArgument:(NSString*)jsonResult {
  NSData *resultData = [jsonResult dataUsingEncoding:NSUTF8StringEncoding];
  responceEventResult(responceNo, (void *)[resultData bytes],
                      [resultData length]);
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
  [self emitEvent:@"/app" name:@"init"];
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  [self emitEvent:@"/app" name:@"focus"];
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  [self emitEvent:@"/app" name:@"blur"];
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  [self emitEvent:@"/app" name:@"reopen" argument:[NSNull null]];
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
  [self emitEvent:@"/app" name:@"filesopen" argument:filenames];
}

- (void)applicationWillFinishLaunching:(NSNotification *)aNotification {
  NSAppleEventManager *appleEventManager =
      [NSAppleEventManager sharedAppleEventManager];
  [appleEventManager
      setEventHandler:self
          andSelector:@selector(handleGetURLEvent:withReplyEvent:)
        forEventClass:kInternetEventClass
           andEventID:kAEGetURL];
}

- (void)handleGetURLEvent:(NSAppleEventDescriptor *)event
           withReplyEvent:(NSAppleEventDescriptor *)replyEvent {
  NSURL *url =
      [NSURL URLWithString:[[event paramDescriptorForKeyword:keyDirectObject]
                               stringValue]];
  [self emitEvent:@"/app" name:@"urlopen" argument:url.absoluteString];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  [self emitEvent:@"/app" name:@"terminate"];
  return NSTerminateNow; // TODO:
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
  [self emitEvent:@"/app" name:@"finalize"];
}

- (NSMenu *)applicationDockMenu:(NSApplication *)sender {
  return self.dock;
}
@end

void Driver_Run() {
  NSApplication *app = [NSApplication sharedApplication];
  app.delegate = [Driver current];
  [NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];
  [NSApp run];
}

void Driver_Terminate() { defer([NSApp terminate:NSApp];); }

BOOL Driver_EmitEvent(void *bytes, NSUInteger length) {
  return [[Driver current] emitEvent:bytes length:length];
}

BOOL Driver_ResponseEvent(NSInteger respNo, void* bytes, NSUInteger length) {
  return [[Driver current] responseEvent:bytes length:length respCallbackNo:respNo];
}

char *Driver_GetBundleResourcesPath() {
  NSBundle *mainBundle = [NSBundle mainBundle];
  return strdup(mainBundle.resourcePath.UTF8String);
}

const char* Driver_GetPreferrdLanguage()
{
  @autoreleasepool {
    NSArray<NSString*>* langs = [NSLocale preferredLanguages];
    NSString* langStrs = [langs componentsJoinedByString:@";"];
    return strdup(langStrs.UTF8String);
  }
}
