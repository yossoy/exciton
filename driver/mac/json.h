#ifndef __INC_DRIVER_MAC_JSON_H__
#define __INC_DRIVER_MAC_JSON_H__

#import <Cocoa/Cocoa.h>

@interface JSONDecoder : NSObject
+ (NSDictionary *)decodeFromBytes:(void *)bytes length:(NSUInteger)length;
@end

@interface JSONEncoder : NSObject
+ (NSData *)encodeFromObject:(id)object;
+ (NSData *)encodeString:(NSString *)s;
+ (NSData *)encodeNumber:(NSNumber *)n;
+ (NSData *)encodeBool:(BOOL)b;
@end

#endif
