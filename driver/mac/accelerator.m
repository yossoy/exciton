#include "accelerator.h"
#include <Carbon/Carbon.h>
#include "log.h"

@implementation Accelerator

+ (instancetype)current {
  static Accelerator *accelerator = nil;

  @synchronized(self) {
    if (accelerator == nil) {
      accelerator = [[Accelerator alloc] init];
    }
  }
  return accelerator;
}

- (instancetype)init {
  if ((self = [super init])) {
  }
  return self;
}

- (NSUInteger)parseModifier:(NSString *)strModifier {
  if ([strModifier isEqualToString:@"shift"]) {
    return NSEventModifierFlagShift;
  } else if ([strModifier isEqualToString:@"command"]) {
    return NSEventModifierFlagCommand;
  } else if ([strModifier isEqualToString:@"ctrl"] ||
             [strModifier isEqualToString:@"control"]) {
    return NSEventModifierFlagControl;
  } else if ([strModifier isEqualToString:@"command"] ||
             [strModifier isEqualToString:@"commandorcontrol"] ||
             [strModifier isEqualToString:@"cmdorctrl"] ||
             [strModifier isEqualToString:@"super"]) {
    return NSEventModifierFlagCommand;
  } else if ([strModifier isEqualToString:@"alt"] ||
             [strModifier isEqualToString:@"option"]) {
    return NSEventModifierFlagOption;
  }
  return 0;
}

- (BOOL)parseString:(NSString *)strAccelerator
        accelerator:(NSString **)accelerator
           modifier:(NSUInteger *)modifier {
  NSArray<NSString *> *components =
      [[strAccelerator lowercaseString] componentsSeparatedByString:@"+"];
  NSUInteger mod = 0;
  if (0 < components.count) {
    for (NSUInteger i = 0; i < (components.count - 1); i++) {
      NSString *strModifiter = [components[i]
          stringByTrimmingCharactersInSet:[NSCharacterSet
                                              whitespaceCharacterSet]];
      mod |= [self parseModifier:strModifiter];
    }
  }
  unichar ch;
  unichar ignorech;
  if (![self keyCodeFromString:[components.lastObject
                                   stringByTrimmingCharactersInSet:
                                       [NSCharacterSet whitespaceCharacterSet]]
                               flags:mod
                           character:&ch
          characterIgnoringModifiers:&ignorech]) {
    LOG_ERROR(@"parseString: parse error:%@(%@)", strAccelerator,
          components.lastObject);
    return FALSE;
  }
  *accelerator = [NSString stringWithCharacters:&ch length:1];
  *modifier = mod;

  return TRUE;
}

- (BOOL)keyCodeFromString:(NSString *)strKey
                         flags:(NSUInteger)flags
                     character:(unichar *)character
    characterIgnoringModifiers:(unichar *)characterIgnoringModifiers {
  static NSDictionary<NSString *, NSArray<NSNumber *> *> *codeMap = NULL;
  const char kShiftCharsForNumberKeys[] = ")!@#$%^&*(";
#define KM(k, m) @[ @(k), @(m) ]
#define K(k) @[ @(k), @0 ]

  if (!codeMap) {
    codeMap = @{
      @"backspace" : KM(kVK_Delete, kBackspaceCharCode),      //
      @"tab" : KM(kVK_Tab, kTabCharCode),                     //
      @"enter" : KM(kVK_Return, kReturnCharCode),             //
      @"return" : KM(kVK_Return, kReturnCharCode),            //
      @"shift" : K(kVK_Shift),                                //
      @"ctrl" : K(kVK_Control),                               //
      @"control" : K(kVK_Control),                            //
      @"alt" : K(kVK_Option),                                 //
      @"option" : K(kVK_Option),                              //
      @"esc" : KM(kVK_Escape, kEscapeCharCode),               //
      @"escape" : KM(kVK_Escape, kEscapeCharCode),            //
      @"space" : KM(kVK_Space, kSpaceCharCode),               //
      @"pageup" : KM(kVK_PageUp, NSPageUpFunctionKey),        //
      @"pagedown" : KM(kVK_PageDown, NSPageDownFunctionKey),  //
      @"end" : KM(kVK_End, NSEndFunctionKey),                 //
      @"home" : KM(kVK_Home, NSHomeFunctionKey),              //
      @"left" : KM(kVK_LeftArrow, NSLeftArrowFunctionKey),    //
      @"up" : KM(kVK_UpArrow, NSUpArrowFunctionKey),          //
      @"right" : KM(kVK_RightArrow, NSRightArrowFunctionKey), //
      @"down" : KM(kVK_DownArrow, NSDownArrowFunctionKey),    //
      @"printscreen" : KM(-1, NSPrintFunctionKey),            //
      @"delete" : KM(kVK_ForwardDelete, kDeleteCharCode),     //
      @"f1" : KM(kVK_F1, NSF1FunctionKey),                    //
      @"f2" : KM(kVK_F2, NSF2FunctionKey),                    //
      @"f3" : KM(kVK_F3, NSF3FunctionKey),                    //
      @"f4" : KM(kVK_F4, NSF4FunctionKey),                    //
      @"f5" : KM(kVK_F5, NSF5FunctionKey),                    //
      @"f6" : KM(kVK_F6, NSF6FunctionKey),                    //
      @"f7" : KM(kVK_F7, NSF7FunctionKey),                    //
      @"f8" : KM(kVK_F8, NSF8FunctionKey),                    //
      @"f9" : KM(kVK_F9, NSF9FunctionKey),                    //
      @"f10" : KM(kVK_F10, NSF10FunctionKey),                 //
      @"f11" : KM(kVK_F11, NSF11FunctionKey),                 //
      @"f12" : KM(kVK_F12, NSF12FunctionKey),                 //
      @"f13" : KM(kVK_F13, NSF13FunctionKey),                 //
      @"f14" : KM(kVK_F14, NSF14FunctionKey),                 //
      @"f15" : KM(kVK_F15, NSF15FunctionKey),                 //
      @"f16" : KM(kVK_F16, NSF16FunctionKey),                 //
      @"f17" : KM(kVK_F17, NSF17FunctionKey),                 //
      @"f18" : KM(kVK_F18, NSF18FunctionKey),                 //
      @"f19" : KM(kVK_F19, NSF19FunctionKey),                 //
      @"f20" : KM(kVK_F20, NSF20FunctionKey),                 //
      @";" : KM(kVK_ANSI_Semicolon, ';'),                     //
      @"=" : KM(kVK_ANSI_Equal, '='),                         //
      @"," : KM(kVK_ANSI_Comma, ','),                         //
      @"-" : KM(kVK_ANSI_Minus, '-'),                         //
      @"." : KM(kVK_ANSI_Period, '.'),                        //
      @"/" : KM(kVK_ANSI_Slash, '/'),                         //
      @"`" : KM(kVK_ANSI_Grave, '`'),                         //
      @"[" : KM(kVK_ANSI_LeftBracket, '['),                   //
      @"\\" : KM(kVK_ANSI_Backslash, '\\'),                   //
      @"]" : KM(kVK_ANSI_RightBracket, ']'),                  //
      @"\''" : KM(kVK_ANSI_Quote, '\''),                      //
    };
  }
  NSArray<NSNumber *> *km = codeMap[strKey];
  unichar ch;
  if (km) {
    ch = km[1].unsignedIntValue;
  } else {
    if (strKey.length == 1) {
      unichar c;
      [strKey getCharacters:&c];
      if (('0' <= c) && (c <= '9')) {
        ch = '0' + (c - '0');
      } else if (('a' <= c) && (c <= 'z')) {
        ch = 'a' + (c - 'a');
      } else {
        return FALSE;
      }
    } else {
      return FALSE;
    }
  }
  *character = *characterIgnoringModifiers = ch;

  if (flags & NSEventModifierFlagShift) {
    if (('0' <= ch) && (ch <= '9')) {
      *character = kShiftCharsForNumberKeys[ch - '0'];
    } else if (('a' <= ch) && (ch <= 'z')) {
      *character = 'A' + (ch - 'a');
    } else {
      switch (ch) {
      case '`':
        *character = '~';
        break;
      case '-':
        *character = '_';
        break;
      case '=':
        *character = '+';
        break;
      case '[':
        *character = '{';
        break;
      case ']':
        *character = '}';
        break;
      case '\\':
        *character = '|';
        break;
      case ';':
        *character = ':';
        break;
      case '\'':
        *character = '\"';
        break;
      case ',':
        *character = '<';
        break;
      case '.':
        *character = '>';
        break;
      case '/':
        *character = '?';
        break;
      }
    }
  }

  if (flags & NSEventModifierFlagControl) {
    if (('a' <= ch) && (ch <= 'z')) {
      *character = 1 + (ch - 'a');
    } else if (ch == '[') {
      *character = 27;
    } else if (ch == '\\') {
      *character = 28;
    } else if (ch == ']') {
      *character = 29;
    }
  }

  return TRUE;
}

@end
