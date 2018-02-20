#pragma once

#import <Cocoa/Cocoa.h>

@interface Accelerator : NSObject
+ (instancetype)current;
- (instancetype)init;

- (BOOL)parseString:(NSString *)strAccelerator
        accelerator:(NSString **)accelerator
           modifier:(NSUInteger *)modifier;

- (BOOL)keyCodeFromString:(NSString *)strKey
                         flags:(NSUInteger)flags
                     character:(unichar *)character
    characterIgnoringModifiers:(unichar *)characterIgnoringModifiers;

@end
