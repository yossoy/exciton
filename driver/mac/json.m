#include "json.h"

@implementation JSONDecoder
+ (NSDictionary *)decodeFromBytes:(void *)bytes length:(NSUInteger)length {
  NSData *data =
      [NSData dataWithBytesNoCopy:bytes length:length freeWhenDone:TRUE];
  return [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
}
@end

@implementation JSONEncoder
+ (NSData *)encodeFromObject:(id)object {
  NSError *err = nil;
  if ([object isKindOfClass:[NSString class]]) {
    return [self encodeString:object];
  } else if ([object isKindOfClass:[NSNumber class]]) {
    return [self encodeNumber:object];
  }
  NSData *jsonData =
      [NSJSONSerialization dataWithJSONObject:object
                                      options:0 /*NSJSONWritingPrettyPrinted*/
                                        error:&err];
  if (err != nil) {
    @throw [NSException exceptionWithName:@"encoding to JSON failed"
                                   reason:err.localizedDescription
                                 userInfo:nil];
  }
  return jsonData;
}

+ (NSData *)encodeString:(NSString *)s {
  NSMutableString* str = [NSMutableString stringWithString:s];
  [str replaceOccurrencesOfString:@"\\" withString:@"\\\\" options:0 range:NSMakeRange(0, [str length])];
  [str replaceOccurrencesOfString:@"\"" withString:@"\\\"" options:0 range:NSMakeRange(0, [str length])];
  [str replaceOccurrencesOfString:@"/" withString:@"\\/" options:0 range:NSMakeRange(0, [str length])];
  [str replaceOccurrencesOfString:@"\n" withString:@"\\n" options:0 range:NSMakeRange(0, [str length])];
  [str replaceOccurrencesOfString:@"\r" withString:@"\\r" options:0 range:NSMakeRange(0, [str length])];
  [str replaceOccurrencesOfString:@"\t" withString:@"\\t" options:0 range:NSMakeRange(0, [str length])];

  return [[NSString stringWithFormat:@"\"%@\"", str] dataUsingEncoding:NSUTF8StringEncoding];
}

+ (NSData *)encodeNumber:(NSNumber *)n {
    return [[n stringValue] dataUsingEncoding:NSUTF8StringEncoding];
}

+ (NSData *)encodeBool:(BOOL)b {
  NSString* strb = b ? @"true" : @"false";
  return [strb dataUsingEncoding:NSUTF8StringEncoding];
}
@end
