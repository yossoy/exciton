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
    /*  self.eventHandlers = [NSMutableDictionary dictionaryWithCapacity:16];*/
    self.respItems = [NSArray array];
    self.lastUseRespItem = 0;
    self.dock = [[NSMenu alloc] initWithTitle:@""];

    self.eventHandlers = [@{
      // @"run" : ^(id argument, NSDictionary<NSString *, NSString *>
      // *parameter,
      //            NSInteger responceNo){
      // }
    } mutableCopy];
    //TODO:
    [Accelerator current];

  }

  return self;
}

- (BOOL)emitEvent:(NSString *)name {
  return [self emitEvent:name argument:[NSNull null]];
}

- (BOOL)emitEvent:(NSString *)name argument:(id)argument {
  NSDictionary *driverEvent = @{
    @"name" : name,
    @"argument" : argument,
    @"respCallbackNo" : @-1,
  };
  NSData *json = [JSONEncoder encodeFromObject:driverEvent];
  requestEventEmit((void *)[json bytes], [json length]);
  return TRUE;
}

- (BOOL)emitEvent:(NSString *)name jsonEncodedArgument:(NSString *)argument {
  NSData *json =
      [[NSString stringWithFormat:
                     @"{\"name\":\"%@\",\"argument\":%@,\"respCallbackNo\":-1}",
                     name, argument] dataUsingEncoding:NSUTF8StringEncoding];
  requestEventEmit((void *)[json bytes], [json length]);
  return TRUE;
}

- (BOOL)emitEvent:(void *)bytes length:(NSUInteger)length {
  NSDictionary *driverEvent = [JSONDecoder decodeFromBytes:bytes length:length];
  NSString *name = driverEvent[@"name"];
  NativeEventHandler handler;
  @synchronized(self.eventHandlers) {
    handler = [self.eventHandlers objectForKey:name];
  }
  if (handler == nil) {
    LOG_ERROR(@"driver.emitEvent: handler not found: %@", name);
    return FALSE;
  }
  LOG_DEBUG(@"driver.emitEvent: %p", handler);
  LOG_DEBUG(@"driver.emitEvent: dispatch event: %@", driverEvent[@"name"]);
  handler(driverEvent[@"argument"], driverEvent[@"parameter"],
          [driverEvent[@"respCallbackNo"] intValue]);

  return TRUE;
}

- (void)addEventHandler:(NSString *)name handler:(NativeEventHandler)handler {
  // LOG_INFO(@"addEventHandler called: %@", name);
  @synchronized(self.eventHandlers) {
    self.eventHandlers[name] = handler;
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
  [self emitEvent:@"/app/init"];
}

- (void)applicationDidBecomeActive:(NSNotification *)aNotification {
  [self emitEvent:@"/app/focus"];
}

- (void)applicationDidResignActive:(NSNotification *)aNotification {
  [self emitEvent:@"/app/blur"];
}

- (BOOL)applicationShouldHandleReopen:(NSApplication *)sender
                    hasVisibleWindows:(BOOL)flag {
  [self emitEvent:@"/app/reopen" argument:[NSNull null]];
  return YES;
}

- (void)application:(NSApplication *)sender
          openFiles:(NSArray<NSString *> *)filenames {
  [self emitEvent:@"/app/filesopen" argument:filenames];
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
  [self emitEvent:@"/app/urlopen" argument:url.absoluteString];
}

- (NSApplicationTerminateReply)applicationShouldTerminate:
    (NSApplication *)sender {
  [self emitEvent:@"/app/terminate"];
  return NSTerminateNow; // TODO:
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
  [self emitEvent:@"/app/finalize"];
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
