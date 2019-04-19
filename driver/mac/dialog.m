#include "dialog.h"
#include "driver.h"
#include "window.h"

static NSAlert *createNSAlert(NSDictionary *params) {
  Driver *d = [Driver current];
  NSString *parentId = params[@"windowId"];
  Window *parent = NULL;
  if (parentId) {
    parent = d.elements[parentId];
  }
  MESSAGE_BOX_TYPE type = [params[@"type"] integerValue];
  NSArray<NSString *> *buttons = params[@"buttons"];
  int defaultId = [params[@"defaultId"] integerValue];
  int i;

  NSAlert *alert = [[NSAlert alloc] init];
  [alert setMessageText:params[@"message"]];
  [alert setInformativeText:params[@"detail"]];

  switch (type) {
  case MESSAGE_BOX_TYPE_INFO:
    [alert setAlertStyle:NSAlertStyleInformational];
    break;
  case MESSAGE_BOX_TYPE_WARNING:
    [alert setAlertStyle:NSAlertStyleWarning];
    break;
  default:
    break;
  }
  for (i = 0; i < buttons.count; i++) {
    NSString *title = buttons[i];
    if (!title || ((id)title == [NSNull null]) || (title.length == 0)) {
      title = @"(empty)";
    }
    NSButton *button = [alert addButtonWithTitle:title];
    button.tag = i;
  }
  NSArray<NSButton *> *ns_buttons = [alert buttons];
  if ((defaultId >= 0) && defaultId < (int)(ns_buttons.count)) {
    ns_buttons[0].keyEquivalent = @"";
    ns_buttons[defaultId].keyEquivalent = @"\r";
  }

  // TODO: icon

  return alert;
}

void setupAllowedFileTypes(NSSavePanel *dialog,
                           NSArray<NSDictionary *> *filters) {
  NSMutableSet<NSString *> *filetypes = [NSMutableSet set];
  for (NSDictionary *filter in filters) {
    NSString *name = filter[@"name"];
    NSArray<NSString *> *extensions = filter[@"extensions"];
    for (NSString *ext in extensions) {
      if ([ext isEqualToString:@"*"]) {
        [dialog setAllowsOtherFileTypes:YES];
        return;
      }
      [filetypes addObject:ext];
    }
  }
  NSArray<NSString *> *file_types = nil;
  if (filetypes.count) {
    file_types = [filetypes allObjects];
  }
  [dialog setAllowedFileTypes:file_types];
}

void setupDialog(NSSavePanel *dialog, NSDictionary *params) {
  NSString *str;
  if ((str = params[@"title"])) {
    [dialog setTitle:str];
  }
  if ((str = params[@"buttonLabel"])) {
    [dialog setPrompt:str];
  }
  NSString *defaultDir = nil;
  NSString *defaultFileName = nil;
  str = params[@"defaultPath"];
  if (str) {
    BOOL isDir = FALSE;
    if ([[NSFileManager defaultManager] fileExistsAtPath:str
                                             isDirectory:&isDir]) {
      if (isDir) {
        defaultDir = str;
      } else {
        defaultDir = [str stringByDeletingLastPathComponent];
        defaultFileName = [str lastPathComponent];
      }
    }
  }
  if (defaultDir) {
    [dialog setDirectoryURL:[NSURL fileURLWithPath:defaultDir isDirectory:YES]];
  }
  if (defaultFileName) {
    [dialog setNameFieldStringValue:defaultFileName];
  }
  NSArray<NSDictionary *> *filters = params[@"filters"];
  if (!filters || !filters.count) {
    [dialog setAllowsOtherFileTypes:YES];
  } else {
    setupAllowedFileTypes(dialog, filters);
  }
}

@implementation Dialog
+ (void)initEventHandlers {
  Driver *d = [Driver current];

  [d addEventHandler:@"/app" name:@"showMessageBox"
              handler:^(id argument,
                        NSDictionary<NSString *, NSString *> *parameter,
                        int responceNo) {
                defer([Dialog showMessageBox:argument
                                   parameter:parameter
                                  responceNo:responceNo];);
              }];
  [d addEventHandler:@"/app" name:@"showOpenDialog"
              handler:^(id argument,
                        NSDictionary<NSString *, NSString *> *parameter,
                        int responceNo) {
                defer([Dialog showOpenDialog:argument
                                   parameter:parameter
                                  responceNo:responceNo];);

              }];
  [d addEventHandler:@"/app" name:@"showSaveDialog"
              handler:^(id argument,
                        NSDictionary<NSString *, NSString *> *parameter,
                        int responceNo) {
                defer([Dialog showSaveDialog:argument
                                   parameter:parameter
                                  responceNo:responceNo];);

              }];
}

+ (NSWindow *)resolveParntWindow:(NSString *)parentId {
  NSWindow *parent = nil;
  Driver *d = [Driver current];
  if (parentId) {
    Window *parentWindow = d.elements[parentId];
    if (parentWindow) {
      parent = parentWindow.window;
    }
  }
  if (!parent) {
    parent = [NSApplication sharedApplication].mainWindow;
  }
  return parent;
}

+ (void)showMessageBox:(id)argument
             parameter:(NSDictionary<NSString *, NSString *> *)parameter
            responceNo:(int)responceNo {
  NSString *parentId = [argument objectForKey:@"windowId"];
  NSWindow *parent = [Dialog resolveParntWindow:parentId];
  NSAlert *alert = createNSAlert(argument);

  if (!parent) {
    NSModalResponse resp = [alert runModal];
    [[Driver current] responceEventResult:responceNo
                                   result:[NSNumber numberWithInteger:resp]];
    return;
  }
  [alert
      beginSheetModalForWindow:parent
             completionHandler:^(NSModalResponse returnCode) {
               [[Driver current]
                   responceEventResult:responceNo
                                result:[NSNumber numberWithInteger:returnCode]];
             }];
}

+ (void)showOpenDialog:(id)argument
             parameter:(NSDictionary<NSString *, NSString *> *)parameter
            responceNo:(int)responceNo {
  NSOpenPanel *dialog = [NSOpenPanel openPanel];
  setupDialog(dialog, argument);
  int properties = [argument[@"properties"] intValue];
  dialog.canChooseFiles = (properties & OPEN_DIALOG_FOR_OPEN_FILE) ? YES : NO;
  if (properties & OPEN_DIALOG_FOR_OPEN_DIRECTORY) {
    dialog.canChooseDirectories = YES;
  }
  if (properties & OPEN_DIALOG_WITH_CREATE_DIRECTORY) {
    dialog.canCreateDirectories = YES;
  }
  if (properties & OPEN_DIALOG_WITH_MULTIPLE_SELECTIONS) {
    dialog.allowsMultipleSelection = YES;
  }
  if (properties & OPEN_DIALOG_WITH_SHOW_HIDDEN_FILES) {
    dialog.showsHiddenFiles = YES;
  }
  NSWindow *parent = [Dialog resolveParntWindow:argument[@"windowId"]];
  if (!parent) {
    NSInteger chosen = [dialog runModal];
    NSMutableArray<NSString *> *files = [[NSMutableArray alloc] init];
    if (chosen != NSModalResponseCancel) {
      for (NSURL *url in dialog.URLs) {
        [files addObject:[url path]];
      }
    }
    [[Driver current] responceEventResult:responceNo result:files];
    return;
  }
  [dialog
      beginSheetModalForWindow:parent
             completionHandler:^(NSInteger chosen) {
               NSMutableArray<NSString *> *files =
                   [[NSMutableArray alloc] init];
               if (chosen != NSModalResponseCancel) {
                 for (NSURL *url in dialog.URLs) {
                   [files addObject:[url path]];
                 }
               }
               [[Driver current] responceEventResult:responceNo result:files];
             }];
}

+ (void)showSaveDialog:(id)argument
             parameter:(NSDictionary<NSString *, NSString *> *)parameter
            responceNo:(int)responceNo {
  NSSavePanel *dialog = [NSSavePanel savePanel];
  setupDialog(dialog, argument);
  dialog.showsHiddenFiles = YES;

  NSWindow *parent = [Dialog resolveParntWindow:argument[@"windowId"]];
  if (!parent) {
    NSInteger chosen = [dialog runModal];
    NSString *strReturn;
    if (chosen != NSModalResponseCancel) {
      strReturn = [[dialog URL] path];
    } else {
      strReturn = @"";
    }
    [[Driver current] responceEventResult:responceNo result:strReturn];
    return;
  }
  [dialog beginSheetModalForWindow:parent
                 completionHandler:^(NSInteger chosen) {
                   NSString *strReturn;
                   if (chosen != NSModalResponseCancel) {
                     strReturn = [[dialog URL] path];
                   } else {
                     strReturn = @"";
                   }
                   [[Driver current] responceEventResult:responceNo
                                                  result:strReturn];
                 }];
}
@end

void Dialog_Init() { [Dialog initEventHandlers]; }