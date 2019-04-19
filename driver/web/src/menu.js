'use strict';
const ExcitonEventData = 'exciton-event-data';
const ExcitonComponentID = 'exciton-component-id';

const ditNone = 0;
const ditCreateNode = 1;
const ditCreateNodeWithNS = 2;
const ditCreateTextNode = 3;
const ditSelectCurNode = 4;
const ditSelectArg1Node = 5;
const ditSelectArg2Node = 6;
const ditPropertyValue = 7;
const ditDelProperty = 8;
const ditAttributeValue = 9;
const ditDelAttributeValue = 10;
const ditAddClassList = 11;
const ditDelClassList = 12;
const ditAddDataSet = 13;
const ditDelDataSet = 14;
const ditAddStyle = 15;
const ditDelStyle = 16;
const ditNodeValue = 17;
const ditInnerHTML = 18;
const ditAppendChild = 19;
const ditInsertBefore = 20;
const ditRemoveChild = 21;
const ditReplaceNode = 22;
const ditAddEventListener = 23;
const ditRemoveEventListener = 24;
const ditSetRootItem = 25
const ditNodeUUID = 26;
const ditAddClientEvent = 27;
const ditMountComponent = 28;
const ditUnmountComponent = 29;
const ditAttributeValueNS = 30;
const ditDelAttributeValueNS = 31;


function resolveMenuNode(menu, items) {
  let ret;
  for (let idx of items) {
    if (ret) {
      ret = ret.subMenu.itemAtIndex(idx);
    } else {
      ret = menu.itemAtIndex(idx)
    }
  }
  return ret;
}

class Menu {
  constructor(menuData) {
    this.id = null;
    this.items = [];
    this.hostItem = null;
    this.title = "";
    this.menuData = menuData;
  }

  initWithMenuTemplate(nsobj, templ) {
    for (var t of templ) {
      const mi = new MenuItem(this.menuData);
      if (t['separator']) {
        mi.separator = true;
      } else {
        mi.initWithMenuTemplate(nsobj, t);
      }
      this.addMenuItem(mi);
    }
  }

  addMenuItem(item) {
    this.items.push(item);
  }
  getSubMenu(idx) {
    if (idx < this.items.length) {
      const item = this.items[idx];
      if (item.subMenu) {
        return item.subMenu;
      }
    }
    return null;
  }
  itemAtIndex(idx) {
    if (idx < this.items.length) {
      return this.items[idx];
    }
    return null;
  }
  getNodeAtIndex(idx) {
    const item = this.itemAtIndex(idx);
    if (!item) {
      return null;
    }
    if (!item.subMenu) {
      return null;
    }
    return item.subMenu.getNode();
  }
  getNode() {
    if (!this.items || this.items.length == 0) {
      return null;
    }
    const menu = document.createElement('ul');
    for (let item of this.items) {
      let mi;
      if (item.separator) {
        if (menu.lastChild) {
          menu.lastChild.style.marginBottom = '10px';
        }
        //mi = document.createElement('hr');
      } else if (item.subMenu) {
        const sm = item.subMenu.getNode();
        if (sm) {
          mi = document.createElement('dl');
          mi.classList.add('column');
          const t = document.createTextNode(item.title);
          const dt = document.createElement('dt');
          dt.appendChild(t);
          mi.appendChild(dt);
          const dd = document.createElement('dd');
          dd.appendChild(sm);
          mi.appendChild(dd);
          dd.style.display = 'none';
        }
      } else {
        if (item.title !== '') {
          const ma = document.createElement('a');
          const mt = document.createTextNode(item.title);
          if (item.enabled && item.handler) {
            ma.onclick = item.handler;
            ma.setAttribute('href', '#');
          }
          ma.appendChild(mt);
          //mi.appendChild(ma);
          mi = ma;
        }
      }
      if (mi) {
        const li = document.createElement('li');
        li.classList.add('menuItem');
        li.appendChild(mi);
        menu.appendChild(li);
      }
    }
    return menu;
  }
};

function roleCmdAbout(nsobj, e) {
  nsobj.showAboutDialog();
}

function roleCmdCut(nsobj, e) {

}

function roleCmdCopy(nsobj, e) {

}

function roleCmdPaste(nsobj, e) {

}

function roleCmdDelete(nsobj, e) {

}

function roleCmdSelectAll(nsobj, e) {

}

function roleCmdZoom(nsobj, e) {

}

function roleCmdClose(nsobj, e) {}

function roleCmdQuit(nsobj, e) {}

function roleCmdToggleFullscreen(nsobj, e) {

}

function roleCmdHistoryBack(nsobj, e) {

}

function roleCmdHistoryForward(nsobj, e) {

}

const defaultRoleInfo = {
  'about': {
    command: roleCmdAbout,
    label: 'About...'
  },
  'front': {
    command: null,
    label: 'Bring All to Front'
  },
  'cut': {
    command: roleCmdCut,
    label: 'Cut'
  },
  'copy': {
    command: roleCmdCopy,
    label: 'Copy'
  },
  'paste': {
    command: roleCmdPaste,
    label: 'Paste'
  },
  'delete': {
    command: roleCmdDelete,
    label: 'Delete'
  },
  'selectall': {
    command: roleCmdSelectAll,
    label: 'Select All'
  },
  'minimize': {
    command: null,
    label: 'Minimize'
  },
  // 'close': {
  //   command: roleCmdClose,
  //   label: 'Close Window'
  // },
  'zoom': {
    command: roleCmdZoom,
    label: 'Zoom'
  },
  // 'quit': {
  //   command: roleCmdQuit,
  //   label: 'Quit'
  // },
  'togglefullscreen': {
    command: roleCmdToggleFullscreen,
    label: 'Toggle Full Screen'
  },
  'viewsource': {
    command: null,
    label: 'View Source'
  },
  'back': {
    command: roleCmdHistoryBack,
    label: 'Back'
  },
  'forward': {
    command: roleCmdHistoryForward,
    label: 'Forward'
  },
};

function translateEvent(nsobj, e, id, menuData) {
  const eventProps = [
    // Event
    'bubbles', 'cancelBubble', 'cancelable', 'composed',
    //'@currentTarget',
    'defaultPrevented', 'eventPhase',
    // '@target', 
    'timeStamp', 'type',
    'isTrusted',
    // UIEvent
    'detail', //'@view',
    // MouseEvent
    'altKey', 'button', 'buttons', 'clientX', 'clientY', 'ctrlKey', 'metaKey',
    'movementX', 'movementY', 'region', '@relatedTarget', 'screenX', 'screenY',
    'shiftKey'
  ];
  const ee = {};
  for (let pn of eventProps) {
    const p = e[pn];
    if (p) {
      ee[pn] = p;
    }
  }
  // target
  const target = {
    'menuId': menuData.id,
    'elementId': id,
    'appId': nsobj.ID,
  };
  ee['currentTarget'] = target;
  ee['target'] = target;

  return ee;
}

class MenuData {
  constructor() {
    this.menu = null;
    this.id = "";
  }

  getAppMenuNode(bar) {
    for (let item of this.menu.items) {
      const dd = document.createElement('div');
      dd.classList.add('dropdown');
      const btn = document.createElement('button');
      btn.classList.add('dropbtn');
      btn.appendChild(document.createTextNode(item.title));
      dd.appendChild(btn);
      if (item.subMenu) {
        const sm = item.subMenu.getNode();
        if (sm) {
          const mc = document.createElement('div');
          mc.classList.add('dropdown-content');
          mc.appendChild(sm);
          dd.appendChild(mc);
        }
      }
      bar.appendChild(dd);
    }
    return bar;
  }
  getPopupMenuNode() {
    const popup = document.createElement('div');
    popup.classList.add('popupMenu');
    popup.appendChild(this.menu.getNode());
    return popup;
  }

  populateWithAppMenuTemplate(nsobj, temp) {
    const menu = new Menu(this);
    this.menu = menu;
    menu.initWithMenuTemplate(nsobj, temp);
  }

  polulateWithDiffset(nsobj, diffSet) {
    const items = diffSet.items;
    const creNodes = [];
    let curNode;
    let arg1Node;
    let arg2Node;
    console.log('DiffSet:', JSON.stringify(items));
    for (let item of items) {
      const key = item.t;
      const k = item.k;
      const v = item.v;
      switch (key) {
        case ditCreateNode:
          if (v === 'menu') {
            const menu = new Menu()
            if (creNodes.length > 0 || this.menu) {
              const mi = new MenuItem();
              mi.setSubMenu(menu);
              curNode = mi;
              creNodes.push(mi);
            } else {
              curNode = menu;
              creNodes.push(menu);
              this.menu = menu;
              this.menu.id = this.id;
            }
          } else if (v === 'menuitem') {
            const mi = new MenuItem();
            creNodes.push(mi);
            curNode = mi;
          } else if (v === 'hr') {
            const mi = new MenuItem();
            mi.separator = true;
            creNodes.push(mi);
            curNode = mi;
          } else {
            throw 'unsupported tag: ' + v;
          }
          break;
        case ditSelectCurNode:
          if (!v) {
            curNode = this.menu;
          } else if (typeof (v) === 'number') {
            curNode = creNodes[v];
          } else {
            curNode = resolveMenuNode(this.menu, v)
          }
          break;
        case ditSelectArg1Node:
          if (!v) {
            arg1Node = this.menu;
          } else if (typeof (v) === 'number') {
            arg1Node = creNodes[v];
          } else {
            arg1Node = resolveMenuNode(this.menu, v);
          }
          break;
        case ditSelectArg2Node:
          if (!v) {
            arg2Node = this.menu;
          } else if (typeof (v) === 'number') {
            arg2Node = creNodes[v];
          } else {
            arg2Node = resolveMenuNode(this.menu, v);
          }
          break;
        case ditAttributeValue:
          if (!(curNode instanceof MenuItem)) {
            if (k !== 'type') {
              throw 'invalid attribute: ' + k + '/' + v;
            }
          } else {
            const mi = curNode;
            if (k === 'label') {
              mi.title = v;
              if (mi.subMenu) {
                mi.subMenu.title = v;
              }
            }
          }
          break;
        case ditDelAttributeValue:
          console.log('ditDelAttributeValue');
          break;
        case ditAddDataSet:
          if (!curNode || !(curNode instanceof MenuItem)) {
            throw 'invalid target: ' + curNode;
          } else {
            const mi = curNode;
            if (k === 'menuRole') {
              const role = defaultRoleInfo[v];
              if (!role) {
                console.warn('unsupported role name:', v);
                break;
              }
              if (role.label) {
                mi.title = role.label;
                if (mi.subMenu) {
                  mi.subMenu.title = role.label;
                }
              }
              if (role.command) {
                mi.handler = (e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  role.command(nsobj, e);
                };
                mi.enabled = true;
              }
            } else if (k == 'menuAcclerator') {
              console.log('menuAcclerator not implement yet.');
            } else {
              console.warn('unknown dataset name:', k);
            }
          }
          break;
        case ditDelDataSet:
          console.warn('ditDelDataSet: not implement yet.');
          break;
        case ditAppendChild:
          {
            let target;
            if (!curNode) {
              target = this.menu;
            } else if (curNode instanceof Menu) {
              target = curNode;
            } else if (curNode instanceof MenuItem) {
              target = curNode.subMenu;
            }
            if (!target) {
              throw 'ditAppendChild: invalid arg: ' + curNode;
            }
            if (!(arg1Node instanceof MenuItem)) {
              if (arg1Node != target) {
                throw 'ditAppendChild: invalid arg1: ' + arg1Node;
              }
            } else {
              target.addMenuItem(arg1Node);
            }
          }
          break;
        case ditInsertBefore:
          console.warn('ditInsertBefore: not implement yet.');
          break;
        case ditRemoveChild:
          console.warn('ditRemoveChild: not implement yet.');
          break;
        case ditAddEventListener:
          if (!(curNode instanceof MenuItem)) {
            throw 'ditAddEventListener: invalid target: ' + curNode;
          }
          if (k !== 'click') {
            throw 'ditAddEventListener: unsupported event';
          } else {
            const mi = curNode;
            const id = v['id'];
            mi.handler = (e) => {
              e.preventDefault();
              e.stopPropagation();
              const ee = translateEvent(nsobj, e, id, this.menu);
              console.log('fakeEvent ==> ', ee);
              nsobj.callNativeMethod('/menu/' + this.menu.id + '/html/' +id + '/click', ee);
            };
            mi.enabled = true;
          }
          break;
        case ditRemoveEventListener:
          console.warn('ditRemoveEventListener: not implement yet.');
          break;
        case ditNodeUUID:
          {
            const id = v;
            if (curNode instanceof Menu) {
              if (curNode != this.menu) {
                curNode.id = id;
                if (curNode.hostItem) {
                  curNode.hostItem.id = id;
                }
              }
            } else if (curNode instanceof MenuItem) {
              curNode.id = id;
            } else {
              throw 'node is invalid';
            }
          }
          break;
        case ditAddClassList:
          break;
        case ditDelClassList:
          break;
        case ditMountComponent:
          break;
        case ditUnmountComponent:
          break;
        default:
          throw 'Unsupported diff type:' + key;
      }
    }
  }
};

class MenuItem {
  constructor(menuData) {
    this.id = "";
    this.subMenu = null;
    this.title = "";
    this.cmdId = -1;
    this.enabled = false;
    this.separator = false;
    this.handler = null;
    this.menuData = menuData;
  }

  initWithMenuTemplate(nsobj, templ) {
    const label = templ['label'];
    const subMenu = templ['subMenu'];
    this.id = templ['id'];
    if (subMenu) {
      const m = new Menu(this.menuData);
      m.initWithMenuTemplate(nsobj, subMenu);
      this.setSubMenu(m);
    }
    const role0 = templ['role'];
    if (role0) {
      const role = defaultRoleInfo[role0];
      if (!role) {
        console.warn('unsupported role name:', role0);
        return;
      }
      if (role.label) {
        this.title = role.label;
        if (this.subMenu) {
          this.subMenu.title = role.label;
        }
      }
      if (role.command) {
        this.handler = (e) => {
          e.preventDefault();
          e.stopPropagation();
          role.command(nsobj, e);
        };
        this.enabled = true;
      }
    } else if (!this.subMenu) {
      this.handler = (e) => {
        e.preventDefault();
        e.stopPropagation();
        // const ee = translateEvent(nsobj, e, id, this.menu);
        // console.log('fakeEvent ==> ', ee);
        nsobj.callNativeMethod('/menu/' + this.menuData.id, 'emit', this.id);
      };
      this.enabled = true;
    }
    if (label) {
      this.title = label;
    }
  }

  setSubMenu(menu) {
    this.subMenu = menu;
    menu.hostItem = this;
    this.enabled = true;
  }
};
const menuDatas = {};

function newMenu(nsobj, dd) {
  const params = dd.parameter;
  const id = params['menu'];
  if (!id) {
    throw 'parameter[menu] not found';
  }
  const menuData = new MenuData();
  menuData.id = id;
  menuDatas[id] = menuData;
  nsobj.responceValue(true, dd.respCallbackNo);
}

function newAppMenu(nsobj, dd) {
  const params = dd.parameter;
  const id = params['menu'];
  const argument = dd.argument;
  if (!id) {
    throw 'parameter[menu] not found';
  }
  const menuData = new MenuData();
  menuData.id = id;
  menuData.populateWithAppMenuTemplate(nsobj,argument);
  menuDatas[id] = menuData;
  
  nsobj.responceValue(true, dd.respCallbackNo);
}

function updateDiffSetHandler(nsobj, dd) {
  const params = dd.parameter;
  const argument = dd.argument;
  const id = params['menu'];
  if (!id) {
    throw 'parameter[menu] not found';
  }
  const menuData = menuDatas[id];
  if (!menuData) {
    throw 'invalid menu';
  }
  menuData.polulateWithDiffset(nsobj, argument);
  nsobj.responceValue(true, dd.respCallbackNo);
}

function getMenuData(params) {
  const id = params['menu'];
  if (!id) {
    throw 'parameter[menu] not found';
  }
  const menuData = menuDatas[id];
  if (!menuData) {
    throw 'invalid menu';
  }
  return menuData;
}

function setApplicationMenu(nsobj, dd) {
  const params = dd.parameter;
  const menuData = getMenuData(params);
  const menuBar = document.getElementById('menubar');
  while (menuBar.firstChild) {
    menuBar.removeChild(menuBar.firstChild);
  }
  menuData.getAppMenuNode(menuBar);
}

function popupContextMenu(nsobj, dd) {
  const params = dd.parameter;
  const menuData = getMenuData(params);
  const pos = dd.argument.position;
  const winid = 'win' + dd.argument.windowId;
  const win = document.getElementById(winid);
  console.log('popupContextMenu is not supported ==> pos: ', pos, ',  win:', win);
}

export default {
  newMenu,
  newAppMenu,
  updateDiffSetHandler,
  setApplicationMenu,
  getMenuData,
  popupContextMenu,
};