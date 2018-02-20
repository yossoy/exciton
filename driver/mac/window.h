#ifndef __INC_DRIVER_MAC_WINDOW_H__
#define __INC_DRIVER_MAC_WINDOW_H__

#import <Cocoa/Cocoa.h>
#import <WebKit/WebKit.h>

typedef void bridge_result;

@interface Window : NSWindowController <NSWindowDelegate, WKNavigationDelegate,
                                        WKUIDelegate, WKScriptMessageHandler>
@property NSString *ID;
@property(weak) WKWebView *webview;

+ (void)initEventHandlers;
+ (BOOL)newWindow:(NSString *)id config:(NSDictionary *)cfg;
- (void)configBackgroundColor:(NSString *)color
                     vibrancy:(NSVisualEffectMaterial)vibrancy;
- (void)configWebview;
- (void)configTitlebar:(NSString *)title hidden:(BOOL)isHidden;

#if 0
- (bridge_result)position:(NSURLComponents *)url payload:(NSString *)payload;
- (bridge_result)move:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)center:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)size:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)resize:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)focus:(NSURLComponents *)url payload:(NSString *)payload;
+ (bridge_result)toggleFullScreen:(NSURLComponents *)url
                          payload:(NSString *)payload;
+ (bridge_result)toggleMinimize:(NSURLComponents *)url
                        payload:(NSString *)payload;
+ (bridge_result)close:(NSURLComponents *)url payload:(NSString *)payload;
#endif
@end

@interface WindowTitleBar : NSView
@end

void Window_Init();

#endif
