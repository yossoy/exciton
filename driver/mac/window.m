#include "window.h"
#include "color.h"
#include "driver.h"
#include "json.h"
#include "log.h"
#import <Foundation/Foundation.h>

@implementation Window
+ (void)initEventHandlers {
  Driver *d = [Driver current];
  [d addEventHandler:@"/window/:id/new"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"id"];
               NSDictionary *cfg = (NSDictionary *)argument;
               if (![Window newWindow:idstr config:cfg]) {
                 LOG_ERROR(@"[Window newWindow failed\n");
               }
               [[Driver current] responceEventResult:responceNo boolean:TRUE];
             }];
  [d addEventHandler:@"/window/:id/requestAnimationFrame"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"id"];
               Driver *driver = [Driver current];
               Window *w = driver.elements[idstr];
               [w requestAnimationFrame];
             }];
  [d addEventHandler:@"/window/:id/updateDiffSetHandler"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"id"];
               Driver *driver = [Driver current];
               Window *w = driver.elements[idstr];
               [w updateDiffSetHandler:argument];
             }];
  [d addEventHandler:@"/window/:id/browserSync"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"id"];
               Driver *driver = [Driver current];
               Window *w = driver.elements[idstr];
               [w browserSyncRequest:argument responceNo:responceNo];
             }];
  [d addEventHandler:@"/window/:id/redirectTo"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"id"];
               Driver *driver = [Driver current];
               Window *w = driver.elements[idstr];
               [w redirectTo:argument];
             }];
}

+ (BOOL)newWindow:(NSString *)id config:(NSDictionary *)cfg {
  NSString *title = cfg[@"title"];
  NSDictionary *pos = [cfg objectForKey:@"position"];
  float left = [pos[@"x"] floatValue];
  float top = [pos[@"y"] floatValue];
  NSDictionary *size = cfg[@"size"];
  float width = [size[@"width"] floatValue];
  float height = [size[@"height"] floatValue];
  NSDictionary *minSize = cfg[@"minSize"];
  NSDictionary *maxSize = cfg[@"maxSize"];
  NSString *backgroundColor = cfg[@"backgroundColor"];
  BOOL noResizable = [cfg[@"noResizable"] boolValue];
  BOOL noClosable = [cfg[@"noClosable"] boolValue];
  BOOL noMinimizable = [cfg[@"noMinimizable"] boolValue];
  BOOL titlebarHidden = [cfg[@"titlebarHidden"] boolValue];
  NSString *defaultURL = cfg[@"default-url"];
  NSNumber *backgroundVibrancy = cfg[@"mac"][@"background-vibrancy"];

  dispatch_block_t block = ^{
    // Configuring raw window.
    NSRect rect = NSMakeRect(left, top, width, height);
    NSUInteger styleMask =
        NSWindowStyleMaskTitled | NSWindowStyleMaskFullSizeContentView;
    if (!noResizable) {
      styleMask |= NSWindowStyleMaskResizable;
    }
    if (!noClosable) {
      styleMask |= NSWindowStyleMaskClosable;
    }
    if (!noMinimizable) {
      styleMask |= NSWindowStyleMaskMiniaturizable;
    }

    NSWindow *rawWindow =
        [[NSWindow alloc] initWithContentRect:rect
                                    styleMask:styleMask
                                      backing:NSBackingStoreBuffered
                                        defer:NO];

    Window *win = [[Window alloc] initWithWindow:rawWindow];
    win.ID = id;
    if (title) {
      win.windowFrameAutosaveName = title;
    }
    win.window.delegate = win;

    if (minSize) {
      win.window.minSize = NSMakeSize([minSize[@"width"] doubleValue],
                                      [minSize[@"height"] doubleValue]);
    }
    if (maxSize) {
      win.window.maxSize = NSMakeSize([maxSize[@"width"] doubleValue],
                                      [maxSize[@"height"] doubleValue]);
    }

    [win configBackgroundColor:backgroundColor
                      vibrancy:backgroundVibrancy.integerValue];
    [win configWebview];
    [win configTitlebar:title hidden:titlebarHidden];

    // Registering window.
    Driver *driver = [Driver current];
    driver.elements[id] = win;

    [win showWindow:nil];

    [win.webview
        loadRequest:[NSURLRequest
                        requestWithURL:[NSURL URLWithString:cfg[@"url"]]]];
  };

  if ([NSThread isMainThread]) {
    block();
  } else {
    dispatch_async(dispatch_get_main_queue(), block);
  }

  return TRUE;
}

- (void)dealloc {
  LOG_DEBUG(@"window::dealloc\n");
  Driver *d = [Driver current];
  [d emitEvent:[NSString stringWithFormat:@"/window/%@/finalize", self.ID]];
}

- (void)configBackgroundColor:(NSString *)color
                     vibrancy:(NSVisualEffectMaterial)vibrancy {
  if (vibrancy != NSVisualEffectMaterialAppearanceBased) {
    NSVisualEffectView *visualEffectView =
        [[NSVisualEffectView alloc] initWithFrame:self.window.frame];
    visualEffectView.material = vibrancy;
    visualEffectView.blendingMode = NSVisualEffectBlendingModeBehindWindow;
    visualEffectView.state = NSVisualEffectStateActive;

    self.window.contentView = visualEffectView;
    return;
  }

  if (color.length == 0) {
    return;
  }
  self.window.backgroundColor =
      [NSColor colorWithCIColor:[CIColor colorWithHexString:color]];
}

- (void)configWebview {
  WKUserContentController *userContentController =
      [[WKUserContentController alloc] init];
  [userContentController addScriptMessageHandler:self name:@"golangRequest"];

  WKWebViewConfiguration *conf = [[WKWebViewConfiguration alloc] init];
  conf.userContentController = userContentController;

#if defined(DEBUG)
  [conf.preferences setValue:@YES forKey:@"developerExtrasEnabled"];
#endif

  WKWebView *webview = [[WKWebView alloc] initWithFrame:NSMakeRect(0, 0, 0, 0)
                                          configuration:conf];
  webview.translatesAutoresizingMaskIntoConstraints = NO;
  webview.navigationDelegate = self;
  webview.UIDelegate = self;

  // Make background transparent.
  [webview setValue:@(NO) forKey:@"drawsBackground"];

  [self.window.contentView addSubview:webview];
  webview.translatesAutoresizingMaskIntoConstraints = NO;
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[webview]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webview)]];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[webview]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              webview)]];
  self.webview = webview;
}

- (void)userContentController:(WKUserContentController *)userContentController
      didReceiveScriptMessage:(WKScriptMessage *)message {
  if (![message.name isEqual:@"golangRequest"]) {
    return;
  }
  NSDictionary *arg = message.body;

  //  LOG_DEBUG(@"userContentController:%@, data:%@", self.ID, message.body);
  [[Driver current] emitEvent:arg[@"path"] jsonEncodedArgument:arg[@"arg"]];
}

- (void)webView:(WKWebView *)webView
    decidePolicyForNavigationAction:(WKNavigationAction *)navigationAction
                    decisionHandler:
                        (void (^)(WKNavigationActionPolicy))decisionHandler {
  if (navigationAction.navigationType == WKNavigationTypeReload ||
      navigationAction.navigationType == WKNavigationTypeOther) {
    if (navigationAction.targetFrame.request != nil) {
      decisionHandler(WKNavigationActionPolicyCancel);
      return;
    }

    decisionHandler(WKNavigationActionPolicyAllow);
    return;
  }

  NSURL *url = navigationAction.request.URL;
  LOG_DEBUG(@"decidePolicyForNavigationAction:%@", url);
  // TO DO:
  // Call go request to navigate to anoter component.
  decisionHandler(WKNavigationActionPolicyCancel);
  // decisionHandler(WKNavigationActionPolicyAllow);
}

- (void)configTitlebar:(NSString *)title hidden:(BOOL)isHidden {
  if (title) {
    self.window.title = title;
  }

  if (!isHidden) {
    return;
  }

  self.window.titleVisibility = NSWindowTitleHidden;
  self.window.titlebarAppearsTransparent = isHidden;

  WindowTitleBar *titlebar = [[WindowTitleBar alloc] init];
  titlebar.translatesAutoresizingMaskIntoConstraints = NO;

  [self.window.contentView addSubview:titlebar];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"|[titlebar]|"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              titlebar)]];
  [self.window.contentView
      addConstraints:
          [NSLayoutConstraint
              constraintsWithVisualFormat:@"V:|[titlebar(==22)]"
                                  options:0
                                  metrics:nil
                                    views:NSDictionaryOfVariableBindings(
                                              titlebar)]];
}

- (void)requestAnimationFrame {
  // LOG_INFO(@"requestAnimationFrame");
  defer([self.webview evaluateJavaScript:@"window.exciton.requestBrowserEvent('"
                                         @"requestAnimationFrame', null);"
                       completionHandler:nil];);
}

- (void)updateDiffSetHandler:(id)diff {
  // LOG_INFO(@"updateDiffSetHandler: %@", diff);
  NSData *jsonData = [JSONEncoder encodeFromObject:diff];
  NSString *jsonStr =
      [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
  NSString *jsonStr2 =
      [[jsonStr stringByReplacingOccurrencesOfString:@"\\" withString:@"\\\\"]
          stringByReplacingOccurrencesOfString:@"\'"
                                    withString:@"\\\'"];
  NSString *cmdstr =
      [NSString stringWithFormat:@"window.exciton.requestBrowserEvent('"
                                 @"updateDiffSetHandler', '%@');",
                                 jsonStr2];
  defer([self.webview evaluateJavaScript:cmdstr completionHandler:nil];);
}

- (void)browserSyncRequest:(id)argument responceNo:(int)responceNo {
  NSData *jsonData = [JSONEncoder encodeFromObject:argument];
  NSString *jsonStr =
      [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
  NSString *jsonStr2 =
      [[jsonStr stringByReplacingOccurrencesOfString:@"\\" withString:@"\\\\"]
          stringByReplacingOccurrencesOfString:@"\'"
                                    withString:@"\\\'"];
  NSString *cmdstr =
      [NSString stringWithFormat:@"window.exciton.requestBrowerEventSync('"
                                 @"browserSync', '%@');",
                                 jsonStr2];
  defer([self.webview
            evaluateJavaScript:cmdstr
             completionHandler:^(id object, NSError *error) {
               LOG_INFO(@"requestBrowerEventSync: (%d) responce ==> %@",
                        responceNo, object);
               Driver *d = [Driver current];
               [d responceEventResult:responceNo jsonEncodedArgument:object];
             }];);
}

- (void)redirectTo:(id)args {
  // LOG_INFO(@"updateDiffSetHandler: %@", diff);
  NSData *jsonData = [JSONEncoder encodeFromObject:args];
  NSString *jsonStr =
      [[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
  NSString *jsonStr2 =
      [[jsonStr stringByReplacingOccurrencesOfString:@"\\" withString:@"\\\\"]
          stringByReplacingOccurrencesOfString:@"\'"
                                    withString:@"\\\'"];
  NSString *cmdstr =
      [NSString stringWithFormat:@"window.exciton.requestBrowserEvent('"
                                 @"redirectTo', '%@');",
                                 jsonStr2];
  LOG_DEBUG(@"redirectTo: ==> %@", cmdstr);
  defer([self.webview evaluateJavaScript:cmdstr completionHandler:nil];);
}

- (void)windowDidResize:(NSNotification *)notification {
  Driver *driver = [Driver current];

  NSDictionary<NSString *, id> *size = @{
    @"width" : [NSNumber numberWithDouble:self.window.frame.size.width],
    @"height" : [NSNumber numberWithDouble:self.window.frame.size.height]
  };
  [driver emitEvent:[NSString stringWithFormat:@"/window/%@/resize", self.ID]
           argument:size];
}

- (void)windowDidBecomeKey:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver emitEvent:[NSString stringWithFormat:@"/window/%@/focus", self.ID]];
}

- (void)windowDidResignKey:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver emitEvent:[NSString stringWithFormat:@"/window/%@/blur", self.ID]];
}

- (void)windowDidEnterFullScreen:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver
      emitEvent:[NSString stringWithFormat:@"/window/%@/fullscreen", self.ID]
       argument:@YES];
}

- (void)windowDidExitFullScreen:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver
      emitEvent:[NSString stringWithFormat:@"/window/%@/fullscreen", self.ID]
       argument:@NO];
}

- (void)windowDidMiniaturize:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver emitEvent:[NSString stringWithFormat:@"/window/%@/minimize", self.ID]
           argument:@YES];
}

- (void)windowDidDeminiaturize:(NSNotification *)notification {
  Driver *driver = [Driver current];

  [driver emitEvent:[NSString stringWithFormat:@"/window/%@/minimize", self.ID]
           argument:@NO];
}

- (BOOL)windowShouldClose:(NSWindow *)sender {
  Driver *driver = [Driver current];

  LOG_INFO(@"windowShouldClose: %@", self.ID);

  // TOOD: return value
  [[Driver current]
      emitEvent:[NSString stringWithFormat:@"/window/%@/close", self.ID]
       argument:[NSNull null]];

  return TRUE;
}

- (void)windowWillClose:(NSNotification *)notification {
  @autoreleasepool {
    self.window = nil;
    [self.webview.configuration.userContentController
        removeScriptMessageHandlerForName:@"golangRequest"];
    self.webview = nil;
  }

  LOG_INFO(@"windowWillClose: %@", self.ID);

  Driver *driver = [Driver current];
  [driver emitEvent:[NSString stringWithFormat:@"/window/%@/closed", self.ID]
           argument:[NSNull null]];
  [driver.elements removeObjectForKey:self.ID];
}
@end

@implementation WindowTitleBar
- (void)mouseDragged:(nonnull NSEvent *)theEvent {
  [self.window performWindowDragWithEvent:theEvent];
}

- (void)mouseUp:(NSEvent *)event {
  Window *win = (Window *)self.window.windowController;
  [win.webview mouseUp:event];

  if (event.clickCount == 2) {
    [win.window zoom:nil];
  }
}

- (WKNavigation *)goBack {
  Window *win = (Window *)self.window.windowController;
  LOG_DEBUG(@"goBack");
  return [win.webview goBack];
}

- (WKNavigation *)goForward {
  Window *win = (Window *)self.window.windowController;
  LOG_DEBUG(@"goForward");
  return [win.webview goForward];
}

- (BOOL)validateMenuItem:(NSMenuItem *)anItem {
  Window *win = (Window *)self.window.windowController;

  if (anItem.action == @selector(goBack:)) {
    return win.webview.canGoBack;
  }
  if (anItem.action == @selector(goFront:)) {
    return win.webview.canGoForward;
  }

  return [super validateMenuItem:anItem];
}
@end

void Window_Init(void) { [Window initEventHandlers]; }
