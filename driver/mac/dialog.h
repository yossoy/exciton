#pragma once

#import <Cocoa/Cocoa.h>

typedef enum _MESSAGE_BOX_TYPE {
	MESSAGE_BOX_TYPE_NONE = 0,
	MESSAGE_BOX_TYPE_INFO,
	MESSAGE_BOX_TYPE_WARNING,
	MESSAGE_BOX_TYPE_ERROR,
	MESSAGE_BOX_TYPE_QUESTION,
} MESSAGE_BOX_TYPE;

typedef enum _OPEN_DIALOG_TYPE {
	OPEN_DIALOG_FOR_OPEN_FILE = 1 << 0,
	OPEN_DIALOG_FOR_OPEN_DIRECTORY = 1 << 1,
	OPEN_DIALOG_WITH_MULTIPLE_SELECTIONS = 1 << 2,
	OPEN_DIALOG_WITH_CREATE_DIRECTORY = 1 << 3,
	OPEN_DIALOG_WITH_SHOW_HIDDEN_FILES = 1 << 4
} OPEN_DIALOG_TYPE;

@interface Dialog : NSObject
+(void) initEventHandlers;

@end

extern void Dialog_Init();
