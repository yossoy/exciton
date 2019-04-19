#include "menu.h"
#include "accelerator.h"
#include "driver.h"
#include "log.h"
#include "menu.h"
#include "window.h"

enum {
  ditNone = 0,
  ditCreateNode,
  ditCreateNodeWithNS,
  ditCreateTextNode,
  ditSelectCurNode,
  ditSelectArg1Node,
  ditSelectArg2Node,
  ditPropertyValue,
  ditDelProperty,
  ditAttributeValue,
  ditDelAttributeValue,
  ditAddClassList,
  ditDelClassList,
  ditAddDataSet,
  ditDelDataSet,
  ditAddStyle,
  ditDelStyle,
  ditNodeValue,
  ditInnerHTML,
  ditAppendChild,
  ditInsertBefore,
  ditRemoveChild,
  ditReplaceChild,
  ditAddEventListener,
  ditRemoveEventListener,
  ditSetRootItem,
  ditNodeUUID,
  ditAddClientEvent,
  ditMountComponent,
  ditUnmountComponent
};

@implementation Menu

#define VS(s) ([NSValue valueWithPointer:s])

+ (NSString *)getRoleDefaultLabel:(NSString *)role {
  static NSDictionary<NSString *, NSString *> *labelMap = NULL;
  if (!labelMap) {
    labelMap = @{
      @"about" : @"About %@",                           //
      @"hide" : @"Hide %@",                             //
      @"hideothers" : @"Hide Others",                   //
      @"unhide" : @"Show All",                          //
      @"front" : @"Bring All to Front",                 //
      @"undo" : @"Undo",                                //
      @"redo" : @"Redo",                                //
      @"cut" : @"Cut",                                  //
      @"copy" : @"Copy",                                //
      @"paste" : @"Paste",                              //
      @"delete" : @"Delete",                            //
      @"pasteandmatchstyle" : @"Paste and Match Style", //
      @"selectall" : @"Select All",                     //
      @"startspeaking" : @"Start Speaking",             //
      @"stopspeaking" : @"Stop Speaking",               //
      @"minimize" : @"Minimize",                        //
      @"close" : @"Close Window",                       //
      @"zoom" : @"Zoom",                                //
      @"quit" : @"Quit %@",                             //
      @"togglefullscreen" : @"Toggle Full Screen",      //
      @"back" : @"Back",
      @"forward" : @"Forward",
    };
  }
  NSString *label = labelMap[role];
  if (!label) {
    return NULL;
  }
  return [NSString
      stringWithFormat:label, [[NSRunningApplication currentApplication]
                                  localizedName]];
}

+ (NSString *)getRoleDefaultAccelerator:(NSString *)role {
  static NSDictionary<NSString *, NSString *> *accelMap = NULL;

  if (!accelMap) {
    accelMap = @{
      @"hide" : @"Command+H",                            //
      @"hideothers" : @"Option+Command+H",               //
      @"undo" : @"Command+Z",                            //
      @"redo" : @"Shift+Command+Z",                      //
      @"cut" : @"Command+X",                             //
      @"copy" : @"Command+C",                            //
      @"paste" : @"Command+V",                           //
      @"delete" : @"Delete",                             //
      @"pasteandmatchstyle" : @"Option+Shift+Command+V", //
      @"selectall" : @"Command+A",                       //
      @"minimize" : @"Command+M",                        //
      @"close" : @"Command+W",                           //
      @"quit" : @"Command+Q",                            //
      @"togglefullscreen" : @"Control+Command+F",        //
      @"back" : @"Command+[",
      @"forward" : @"Command+]",
    };
  }
  return accelMap[role];
}

+ (SEL)getRole:(NSString *)role {
  static NSDictionary<NSString *, NSValue *> *roleMap = NULL;

  if (!roleMap) {
    roleMap = [NSDictionary
        dictionaryWithObjectsAndKeys:                                  //
            VS(@selector(orderFrontStandardAboutPanel:)), @"about",    //
            VS(@selector(hide:)), @"hide",                             //
            VS(@selector(hideOtherApplications:)), @"hideothers",      //
            VS(@selector(unhideAllApplications:)), @"unhide",          //
            VS(@selector(arrangeInFront:)), @"front",                  //
            VS(@selector(undo:)), @"undo",                             //
            VS(@selector(redo:)), @"redo",                             //
            VS(@selector(cut:)), @"cut",                               //
            VS(@selector(copy:)), @"copy",                             //
            VS(@selector(paste:)), @"paste",                           //
            VS(@selector(delete:)), @"delete",                         //
            VS(@selector(pasteAndMatchStyle:)), @"pasteandmatchstyle", //
            VS(@selector(selectAll:)), @"selectall",                   //
            VS(@selector(startSpeaking:)), @"startspeaking",           //
            VS(@selector(stopSpeaking:)), @"stopspeaking",             //
            VS(@selector(performMiniaturize:)), @"minimize",           //
            VS(@selector(performClose:)), @"close",                    //
            VS(@selector(performZoom:)), @"zoom",                      //
            VS(@selector(terminate:)), @"quit",                        //
            VS(@selector(toggleFullScreen:)), @"togglefullscreen",     //
            VS(@selector(goBack:)), @"back", VS(@selector(goForward:)),
            @"forward", nil];
  }
  NSValue *val = [roleMap objectForKey:role];
  if (val) {
    return val.pointerValue;
  }
  return NULL;
}

+ (void)initEventHandlers {
  Driver *d = [Driver current];

  [d addEventHandler:@"/menu/:menu"
                name:@"new"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"menu"];
               Menu *menu = [[Menu alloc] initWithID:idstr];
               Driver *d = [Driver current];
               d.elements[idstr] = menu;
               [d responceEventResult:responceNo boolean:TRUE];
               LOG_DEBUG(@"menu::new callback end");
             }];
  [d addEventHandler:@"/menu/:menu"
                name:@"newApplicationMenu"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"menu"];
               NSArray<NSDictionary *> *json =
                   (NSArray<NSDictionary *> *)argument;
               Menu *m = [[Menu alloc] initWithID:idstr];
               MenuContainer *mc =
                   [[MenuContainer alloc] initWithAppMenuTemplate:json
                                                            title:@""
                                                         delegate:m];
               Driver *d = [Driver current];
               if (mc == nil) {
                 [d responceEventResult:responceNo boolean:NO];
               } else {
                 mc.ID = idstr;
                 m.pMenu = mc;
                 d.elements[idstr] = m;
                 // [NSApp setMainMenu:m.pMenu];
                 [d responceEventResult:responceNo boolean:TRUE];
               }
               LOG_DEBUG(@"menu::new callback end");
             }];
  [d addEventHandler:@"/menu/:menu"
                name:@"newPopupMenu"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               NSString *idstr = parameter[@"menu"];
               NSArray<NSDictionary *> *json =
                   (NSArray<NSDictionary *> *)argument;
               Menu *m = [[Menu alloc] initWithID:idstr];
               MenuContainer *mc =
                   [[MenuContainer alloc] initWithMenuTemplate:json
                                                         title:@""
                                                      delegate:m];
               Driver *d = [Driver current];
               if (mc == nil) {
                 [d responceEventResult:responceNo boolean:NO];
               } else {
                 mc.ID = idstr;
                 m.pMenu = mc;
                 d.elements[idstr] = m;
                 [d responceEventResult:responceNo boolean:TRUE];
               }
               LOG_DEBUG(@"menu::new callback end");
             }];
  [d addEventHandler:@"/menu/:menu"
                name:@"updateDiffSetHandler"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               LOG_DEBUG(@"menu::updateDiffSetHandler is called, call defer");
               defer(NSString *idstr = parameter[@"menu"];
                     NSDictionary *diff = argument;
                     Driver *driver = [Driver current];
                     Menu *m = driver.elements[idstr];
                     LOG_INFO(@"updateDiffSetHandler: %@", idstr);
                     [m populateWithDiffset:diff];
                     [driver responceEventResult:responceNo boolean:TRUE];);
             }];
  [d addEventHandler:@"/menu/:menu"
                name:@"setApplicationMenu"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               defer(NSString *idstr = parameter[@"menu"];
                     Driver *d = [Driver current]; Menu *m = d.elements[idstr];
                     LOG_INFO(@"setApplicationMenu: %@", idstr);
                     [NSApp setMainMenu:m.pMenu];);
             }];
  [d addEventHandler:@"/menu/:menu"
                name:@"popupContextMenu"
             handler:^(id argument,
                       NSDictionary<NSString *, NSString *> *parameter,
                       int responceNo) {
               defer(
                   NSString *idstr = parameter[@"menu"];
                   Driver *d = [Driver current]; Menu *m = d.elements[idstr];
                   float posX = [argument[@"position"][@"x"] floatValue];
                   float posY = [argument[@"position"][@"y"] floatValue];
                   LOG_INFO(@"popupContextMenu: %@, %f, %f", idstr, posX, posY);
                   NSString *winidstr = argument[@"windowId"];
                   Window *parentWindow = d.elements[winidstr];
                   NSWindow *parent = parentWindow.window;
                   NSView *contentView = parent.contentView;
                   NSRect scrRect = NSMakeRect(
                       posX, parent.screen.frame.size.height - posY, 0.0, 0.0);
                   NSRect winRect = [parent convertRectFromScreen:scrRect];
                   NSPoint pos = [contentView convertPoint:winRect.origin
                                                  fromView:nil];
                   [m.pMenu popUpMenuPositioningItem:m.pMenu.itemArray[0]
                                          atLocation:pos
                                              inView:parent.contentView];);
             }];
}

- (instancetype)initWithID:(NSString *)menuId {
  if ((self = [super init])) {
    self.ID = menuId;
  }
  return self;
}

- (instancetype)initWithID:(NSString *)menuId
         popupMenuTemplate:(NSArray<NSDictionary *> *)json {
  if ((self = [super init])) {
    self.ID = menuId;
  }
  return self;
}

- (void)dealloc {
  Driver *d = [Driver current];
  [d emitEvent:[NSString stringWithFormat:@"/menu/%@", self.ID]
          name:@"finalize"];
  if (self.pMenu) {
    [self.pMenu setDelegate:nil];
  }
}

- (NSMenuItem *)resolveMenuNode:(NSArray *)items {
  const NSUInteger cnt = [items count];
  NSMenuItem *ret = NULL;
  for (NSUInteger i = 0; i < cnt; i++) {
    NSUInteger idx = [items[i] unsignedIntegerValue];
    if (ret) {
      ret = [[ret submenu] itemAtIndex:idx];
    } else {
      ret = [self.pMenu itemAtIndex:idx];
    }
  }
  return ret;
}

- (BOOL)populateWithDiffset:(NSDictionary *)diffset {
  NSMutableArray<id> *creNodes = [NSMutableArray arrayWithCapacity:16];
  id curNode;
  id arg1Node;
  id arg2Node;
  for (id item in diffset[@"items"]) {
    // LOG_DEBUG(@"diff: %@ <== %@", item[@"type"], item[@"v"]);
    const int key = [item[@"t"] intValue];
    NSString *k = item[@"k"];
    id v = item[@"v"];
    NSString *str1;
    switch (key) {
    case ditCreateNode: {
      str1 = (NSString *)v;
      if ([str1 isEqualToString:@"menu"]) {
        MenuContainer *menu = [[MenuContainer alloc] initWithTitle:@""];
        [menu setDelegate:self];
        if (creNodes.count != 0 || self.pMenu) {
          // child menu
          MenuItem *mi = [[MenuItem alloc] init];
          mi.submenu = menu;
          menu.hostItem = mi;
          curNode = mi;
          [creNodes addObject:mi];
        } else {
          curNode = menu;
          [creNodes addObject:menu];
          // TODO: more cleanup
          self.pMenu = menu;
        }
      } else if ([str1 isEqualToString:@"menuitem"]) {
        MenuItem *mi = [[MenuItem alloc] init];
        [mi setTarget:self];
        [mi setTag:1];
        [mi setEnabled:YES];
        [creNodes addObject:mi];
        curNode = mi;
      } else if ([str1 isEqualToString:@"hr"]) {
        NSMenuItem *si = [NSMenuItem separatorItem];
        [creNodes addObject:si];
        curNode = si;
      } else {
        LOG_ERROR(@"ditCreateNode: unsupported tag: %@", str1);
        return FALSE;
      }
      break;
    }
    case ditSelectCurNode:
      if (v == nil || v == [NSNull null]) {
        curNode = self.pMenu;
      } else if ([v isKindOfClass:[NSNumber class]]) {
        NSUInteger idx = [v unsignedIntegerValue];
        curNode = [creNodes objectAtIndex:[v unsignedIntegerValue]];
      } else {
        curNode = [self resolveMenuNode:v];
      }
      break;
    case ditSelectArg1Node:
      if (v == nil || v == [NSNull null]) {
        arg1Node = self.pMenu;
      } else if ([v isKindOfClass:[NSNumber class]]) {
        NSUInteger idx = [v unsignedIntegerValue];
        arg1Node = [creNodes objectAtIndex:[v unsignedIntegerValue]];
      } else {
        arg2Node = [self resolveMenuNode:v];
      }
      break;
    case ditSelectArg2Node:
      if (v == nil || v == [NSNull null]) {
        arg2Node = self.pMenu;
      } else if ([v isKindOfClass:[NSNumber class]]) {
        NSUInteger idx = [v unsignedIntegerValue];
        arg2Node = [creNodes objectAtIndex:[v unsignedIntegerValue]];
      } else {
        arg2Node = [self resolveMenuNode:v];
      }
      break;
    case ditAttributeValue: {
      if (![curNode isKindOfClass:[MenuItem class]]) {
        if (![k isEqualToString:@"type"]) {
          LOG_ERROR(@"ditAttributeValue: invalid curNode: %@", curNode);
          return FALSE;
        }
      } else {
        MenuItem *mi = curNode;
        if ([k isEqualToString:@"label"]) {
          mi.title = v;
          if (mi.submenu) {
            mi.submenu.title = v;
          }
        }
      }
      break;
    }
    case ditDelAttributeValue:
      LOG_WARNING(@"Not implement yet: ditDelAttributeValue");
      break;
    case ditAddDataSet: {
      if (!curNode || ![curNode isKindOfClass:[MenuItem class]]) {
        LOG_ERROR(@"ditAddDataSet: invalid target: %@", curNode);
        return FALSE;
      }
      NSString *name = k;
      NSString *val = v;
      if ([name isEqualToString:@"menuRole"]) {
        SEL sel = [Menu getRole:val];
        if (!sel) {
          if ([curNode submenu]) {
            if ([val isEqualToString:@"window"]) {
              [NSApp setWindowsMenu:[curNode submenu]];
              break;
            } else if ([val isEqualToString:@"help"]) {
              [NSApp setHelpMenu:[curNode submenu]];
              break;
            } else if ([val isEqualToString:@"services"]) {
              [NSApp setServicesMenu:[curNode submenu]];
              break;
            }
          }
          LOG_ERROR(@"ditAddDataSet: unsupported role name: %@", val);
          return FALSE;
        }
        [curNode setTarget:nil];
        [curNode setAction:sel];
        NSString *s = [Menu getRoleDefaultAccelerator:val];
        if (s) {
          NSString *accel = NULL;
          NSUInteger modifier;
          if ([[Accelerator current] parseString:s
                                     accelerator:&accel
                                        modifier:&modifier]) {
            [curNode setKeyEquivalent:accel];
            [curNode setKeyEquivalentModifierMask:modifier];
          }
        }
        s = [Menu getRoleDefaultLabel:val];
        if (s) {
          [curNode setTitle:s];
          if ([curNode submenu]) {
            [[curNode submenu] setTitle:s];
          }
        }
      } else if ([name isEqualToString:@"menuAcclerator"]) {
        NSString *accel = NULL;
        NSUInteger modifier;
        if ([[Accelerator current] parseString:val
                                   accelerator:&accel
                                      modifier:&modifier]) {
          [curNode setKeyEquivalent:accel];
          [curNode setKeyEquivalentModifierMask:modifier];
        }
      } else {
        LOG_ERROR(@"ditAddDataSet: unknwon dataSet Name:%@", v);
        return FALSE;
      }
      break;
    }
    case ditDelDataSet:
      LOG_WARNING(@"Not implement yet: ditDelDataSet");
      break;
    case ditAppendChild: {
      NSMenu *target;
      if (!curNode) {
        target = self.pMenu;
      } else if ([curNode isKindOfClass:[MenuContainer class]]) {
        target = curNode;
      } else if ([curNode isKindOfClass:[MenuItem class]]) {
        MenuItem *item = curNode;
        target = item.submenu;
      }
      NSMenuItem *mi = arg1Node;
      if (!target || !mi) {
        LOG_ERROR(@"ditAppendChild: invalid arg1");
        return FALSE;
      }
      // LOG_DEBUG(@"addItem %@ <<- %@", target, mi);
      if (target != arg1Node) { // TODO: more cleanup
        [target addItem:mi];
      }
      break;
    }
    case ditInsertBefore:
      LOG_WARNING(@"Not implement yet: ditInsertBefore");
      break;
    case ditRemoveChild:
      LOG_WARNING(@"Not implement yet: ditRemoveChild");
      break;
    case ditAddEventListener: {
      if (!curNode || ![curNode isKindOfClass:[MenuItem class]]) {
        LOG_ERROR(@"ditAddEventListener: invalid target: %@", curNode);
        return FALSE;
      }
      if (![k isEqualToString:@"click"]) {
        LOG_ERROR(@"ditAddEventListener: unsupported event: %@", v);
        return FALSE;
      }
      MenuItem *mi = curNode;
      [mi setAction:@selector(itemSelected:)];
      mi.onClick = v[@"id"];
      break;
    }
    case ditRemoveEventListener:
      LOG_WARNING(@"Not implement yet: ditRemoveEventListener");
      break;
    case ditSetRootItem: {
      if (!curNode) {
        LOG_ERROR(@"ditSetRootItem: current node is null!");
        return FALSE;
      }
      if (![curNode isKindOfClass:[MenuContainer class]]) {
        LOG_ERROR(@"ditSetRootItem: invalid current node!");
        return FALSE;
      }
      if (self.pMenu) {
        LOG_ERROR(@"ditSetRootItem: rootItem is already exist");
        return FALSE;
      }
      LOG_INFO(@"NewRootItem: %@", curNode);
      self.pMenu = curNode;
      break;
    }
    case ditNodeUUID: {
      if ([curNode isKindOfClass:[MenuContainer class]]) {
        MenuContainer *menu = curNode;
        menu.ID = v;
        if (menu.hostItem) {
          menu.hostItem.ID = v;
        }
      } else if ([curNode isKindOfClass:[MenuItem class]]) {
        MenuItem *mi = curNode;
        mi.ID = v;
      } else if (![curNode isKindOfClass:[NSMenuItem class]]) {
        LOG_ERROR(@"ditNodeUUID: current node is invalid: %@", curNode);
        return FALSE;
      }
      break;
    }
    case ditAddClassList:
    case ditDelClassList:
    case ditMountComponent:
    case ditUnmountComponent:
      break;
    case ditCreateNodeWithNS:
    case ditCreateTextNode:
    case ditPropertyValue:
    case ditDelProperty:
    case ditAddStyle:
    case ditDelStyle:
    case ditNodeValue:
    case ditInnerHTML:
    case ditReplaceChild:
    case ditAddClientEvent:
    default:
      LOG_ERROR(@"Unsupported item type: %d", key);
      return FALSE;
    }
  }

  LOG_INFO(@"menu root = %@", self.pMenu);
  return TRUE;
}

- (BOOL)validateUserInterfaceItem:(id<NSValidatedUserInterfaceItem>)item {
  LOG_INFO(@"validateUserInterfaceItem: %@", (id)item);
  if ([(id)item isKindOfClass:[MenuItem class]]) {
    MenuItem *mi = (MenuItem *)item;
    dispatch_semaphore_t sem = dispatch_semaphore_create(0);
    __block BOOL result = YES;
    LOG_INFO(@"==> %@", mi.ID);
    [[Driver current] emitEvent:[NSString stringWithFormat:@"/menu/%@", self.ID] name:@"is-item-enabled" argument:mi.ID respCallback:^(id value, NSString* error) {
      LOG_INFO(@"emitEvent: id:%@value:%@ error:%@", mi.ID, value, error);
      result = [value boolValue];
      dispatch_semaphore_signal(sem);
    }];
    dispatch_semaphore_wait(sem, DISPATCH_TIME_FOREVER);
#if !OS_OBJECT_USE_OBJC
    dispatch_release(sem);
#endif
    return result;
  }
  // toolbar is not supported yet..
  LOG_ERROR(@"validateUserInterfaceItem not match: %@", (id)item);

  return NO;
}

- (void)itemSelected:(id)sender {
  LOG_DEBUG(@"itemSelected called: %@", sender);
  if (![sender isKindOfClass:[MenuItem class]]) {
    return;
  }
  LOG_DEBUG(@"itemSelected called2: %@", sender);
  MenuItem *mi = sender;
  NSPoint mousePt = [NSEvent mouseLocation];
  NSEventModifierFlags modifiers = [NSEvent modifierFlags];

  LOG_DEBUG(@"itemSelected called3: %@", mi);

  // NSDictionary *target = @{
  //   @"menuId" : self.ID,
  //   @"elementId" : mi.ID,
  // };

  // LOG_DEBUG(@"itemSelected called4: %@", sender);

  // NSDictionary *fakeEvent = @{
  //   // Event
  //   @"bubbles" : @NO,
  //   @"cancelBubble" : @NO,
  //   @"cancelable" : @NO,
  //   @"composed" : @NO,
  //   @"currentTarget" : target,
  //   @"defaultPrevented" : @NO,
  //   @"eventPhase" : @2,
  //   @"target" : target,
  //   @"timeStamp" : @0,
  //   @"type" : @"click",
  //   @"isTrusted" : @NO,

  //   // UIEvent
  //   @"detail" : @0,
  //   // TODO: アクティブウィンドウを取得出来るようにするべきか
  //   //@"view":???

  //   // MouseEvent
  //   @"altKey" : (modifiers & NSEventModifierFlagOption) ? @YES : @NO,
  //   @"button" : @0,
  //   @"buttons" : @1,
  //   @"clientX" : @0,
  //   @"clientY" : @0,
  //   @"ctrlKey" : (modifiers & NSEventModifierFlagControl) ? @YES : @NO,
  //   @"metaKey" : (modifiers & NSEventModifierFlagCommand) ? @YES : @NO,
  //   @"movementX" : @0.0,
  //   @"movementY" : @0.0,
  //   @"region" : [NSNull null],
  //   //"relatedTarget":???
  //   @"screenX" : [NSNumber numberWithFloat:mousePt.x],
  //   @"screenY" : [NSNumber numberWithFloat:mousePt.y],
  //   @"shiftKey" : (modifiers & NSEventModifierFlagShift) ? @YES : @NO
  // };
  LOG_DEBUG(@"itemSelected called5: %@", sender);
  LOG_INFO(@"itemSelected: %@: %@,%@", mi, mi.ID, mi.onClick);
  // [[Driver current] emitEvent:[NSString stringWithFormat:@"/menu/%@/html/%@",
  //                                                        self.ID, mi.onClick]
  //                        name:@"click"
  //                    argument:fakeEvent];
  [[Driver current] emitEvent:[NSString stringWithFormat:@"/menu/%@", self.ID] name:@"emit" argument:mi.ID];
}

- (void)menuDidClose:(NSMenu *)menu {
  LOG_INFO(@"***MenuDidClose:%@", menu);
}

@end

@implementation MenuContainer
- (instancetype)initWithMenuTemplate:(NSArray<NSDictionary *> *)items
                               title:(NSString *)title
                            delegate:(Menu *)delegate {
  if ((self = [super initWithTitle:title])) {
    [self setDelegate:delegate];
    __block BOOL isError = NO;
    [items enumerateObjectsUsingBlock:^(NSDictionary *item, NSUInteger idx,
                                        BOOL *isStop) {
      BOOL isSeparator = [[item objectForKey:@"separator"] boolValue];
      if (isSeparator) {
        NSMenuItem *mi = [NSMenuItem separatorItem];
        [self addItem:mi];
      } else {
        MenuItem *mi = [[MenuItem alloc] initWithItemTemplate:item
                                                     delegate:delegate];
        if (mi == nil) {
          isError = YES;
          *isStop = YES;
          return;
        }
        [self addItem:mi];
      }
    }];
    if (isError) {
      return nil;
    }
  }
  return self;
}

- (instancetype)initWithAppMenuTemplate:(NSArray<NSDictionary *> *)items
                                  title:(NSString *)title
                               delegate:(Menu *)delegate {
  if ((self = [super initWithTitle:title])) {
    [self setDelegate:delegate];
    __block BOOL isError = NO;
    [items enumerateObjectsUsingBlock:^(NSDictionary *item, NSUInteger idx,
                                        BOOL *isStop) {
      MenuItem *mi = [[MenuItem alloc] initWithItemTemplate:item
                                                   delegate:delegate];
      if (mi == nil) {
        isError = YES;
        *isStop = YES;
        return;
      }
      [self addItem:mi];
    }];
    if (isError) {
      return nil;
    }
  }
  return self;
}
@end

@implementation MenuItem
- (instancetype)initWithItemTemplate:(NSDictionary *)item
                            delegate:(Menu *)delegate {
  if ((self = [super init])) {
    [self setTarget:delegate];
    [self setTag:1];
    [self setEnabled:YES];
    NSString *label = [item objectForKey:@"label"];
    NSArray<NSDictionary *> *subMenu = [item objectForKey:@"subMenu"];
    if (subMenu != nil) {
      MenuContainer *mc = [[MenuContainer alloc] initWithMenuTemplate:subMenu
                                                                title:@""
                                                             delegate:delegate];
      if (mc == nil) {
        return nil;
      }
      self.submenu = mc;
      mc.hostItem = self;
    }
    NSString *accel = [item objectForKey:@"acclerator"];
    NSString *role = [item objectForKey:@"role"];
    self.ID = [item objectForKey:@"id"];
    if (0 < [role length]) {
      SEL sel = [Menu getRole:role];
      if (!sel) {
        if ([role isEqualToString:@"window"]) {
          if (self.submenu == nil) {
            MenuContainer *mc = [[MenuContainer alloc] initWithTitle:@""];
            self.submenu = mc;
            mc.hostItem = self;
          }
          [NSApp setWindowsMenu:self.submenu];
        } else if ([role isEqualToString:@"help"]) {
          if (self.submenu == nil) {
            MenuContainer *mc = [[MenuContainer alloc] initWithTitle:@""];
            self.submenu = mc;
            mc.hostItem = self;
          }
          [NSApp setHelpMenu:self.submenu];
        } else if ([role isEqualToString:@"services"]) {
          if (self.submenu == nil) {
            MenuContainer *mc = [[MenuContainer alloc] initWithTitle:@""];
            self.submenu = mc;
            mc.hostItem = self;
          }
          [NSApp setServicesMenu:self.submenu];
        } else {
          LOG_ERROR(@"ditAddDataSet: unsupported role name: %@", role);
          return nil;
        }
      } else {
        [self setTarget:nil];
        [self setAction:sel];
        NSString *acc = [Menu getRoleDefaultAccelerator:role];
        if (acc != nil && (0 == [accel length])) {
          accel = acc;
        }

        NSString *rl = [Menu getRoleDefaultLabel:role];
        if (rl != nil && (0 == [label length])) {
          label = rl;
        }
      }
    } else if (!self.submenu) {
      LOG_DEBUG(@"----------------> itemSelecteed set: %@\n", label);
      [self setAction:@selector(itemSelected:)];
    }
    if (0 < [label length]) {
      [self setTitle:label];
      if (self.submenu != nil) {
        [self.submenu setTitle:label];
      }
    }
    if (0 < [accel length]) {
      NSUInteger modifier;
      NSString *a = nil;
      if ([[Accelerator current] parseString:accel
                                 accelerator:&a
                                    modifier:&modifier]) {
        [self setKeyEquivalent:a];
        [self setKeyEquivalentModifierMask:modifier];
      }
    }
  }
  return self;
}
@end

void Menu_Init() { [Menu initEventHandlers]; }
