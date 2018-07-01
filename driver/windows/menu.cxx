#include <windows.h>

#include <mshtmcid.h>
#include <shdeprecated.h>
#include <shlobj.h>
#include <wininet.h>

#include "browsercontainer.h"
#include "browserhost.h"
#include "driver.h"
#include "global.h"
#include "log.h"
#include "menu.h"
#include "menumgr.h"
#include "myjson.h"
#include "util.h"

namespace {
std::string
getIdFromParam(const std::map<std::string, std::string> &parameter) {
  auto fiter = parameter.find("id");
  if (fiter == parameter.end()) {
    return "";
  }
  return (*fiter).second;
}
CWebBrowserHost *resolveParentWindow(const std::string &id) {
  Driver &d = Driver::Current();
  if (!id.empty()) {
    auto fiter = d.Hosts().find(id);
    if (fiter != d.Hosts().end()) {
      CWebBrowserHost *p = (*fiter).second.ptr_;
      return p;
    }
  }
  auto pActiveContainer = CWebBrowserContainer::GetActiveContainer();
  return pActiveContainer->GetActiveBrowser();
}
} // namespace

namespace exciton {
namespace menu {
// TODO: commonにしたい。。。
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
  ditNodeUUID
};
struct RoleInfo {
  const char *label;
  const char *accel;
  RoledCommandId command;
  RoleInfo(RoledCommandId cmd, const char *l, const char *a = nullptr)
      : command(cmd), label(l), accel(a) {}
};
namespace {
static std::map<std::string, RoleInfo> labelMap = {
    {"about", {RoledCommandId::About, "&About %s"}},
    {"front", {RoledCommandId::Front, "Bring All to Front"}}, //
    //{"undo", { "Undo"}},                                //
    //{"redo", {"Redo"}},                                //
    {"cut", {RoledCommandId::Cut, "Cut"}},          //
    {"copy", {RoledCommandId::Copy, "Copy"}},       //
    {"paste", {RoledCommandId::Paste, "Paste"}},    //
    {"delete", {RoledCommandId::Delete, "Delete"}}, //
    //      {"pasteandmatchstyle", {"Paste and Match Style"}}, //
    {"selectall", {RoledCommandId::SelectAll, "Select All"}}, //
    //      {"startspeaking", {"Start Speaking"}},             //
    //      {"stopspeaking", {"Stop Speaking"}},               //
    {"minimize", {RoledCommandId::Minimize, "Minimize"}}, //
    {"close", {RoledCommandId::Close, "Close Window"}},   //
    {"zoom", {RoledCommandId::Zoom, "Zoom"}},             //
    {"quit", {RoledCommandId::Quit, "Quit"}},             //
    {"togglefullscreen",
     {RoledCommandId::ToggleFullscreen, "Toggle Full Screen"}}, //
    {"viewsource", {RoledCommandId::ViewSource, "View Source..."}},
};
}
const RoleInfo *getDefaultRoleInfo(const std::string &role) {
  auto fiter = labelMap.find(role);
  if (fiter == labelMap.end()) {
    return nullptr;
  }
  return &(*fiter).second;
};

Menu::Menu() {}

Menu::Menu(const std::string &id) : id_(id) {}

void Menu::AddMenuItem(std::shared_ptr<MenuItem> item) {
  items_.push_back(item);
}

std::shared_ptr<Menu> Menu::GetSubMenu(size_t idx) const {
  std::shared_ptr<Menu> ret;
  if (idx < items_.size()) {
    auto item = items_[idx];
    if (item->subMenu_) {
      ret = item->subMenu_;
    }
  }
  return ret;
}

std::shared_ptr<MenuItem> Menu::itemAtIndex(size_t idx) const {
  std::shared_ptr<MenuItem> ret;
  if (idx < items_.size()) {
    ret = items_[idx];
  }
  return ret;
}

HMENU Menu::GetHMenuAtIndex(size_t idx) const {
  auto menuItem = itemAtIndex(idx);
  if (!menuItem) {
    return NULL;
  }
  if (!menuItem->subMenu_) {
    return NULL;
  }
  return menuItem->subMenu_->GetHMenu();
}

HMENU Menu::GetHMenu() const {
  HMENU hMenu = ::CreatePopupMenu();
  MENUINFO mi;
  mi.cbSize = sizeof(mi);
  mi.fMask = MIM_MENUDATA;
  mi.dwMenuData = reinterpret_cast<ULONG_PTR>(this);
  ::SetMenuInfo(hMenu, &mi);
  for (auto item : items_) {
    MENUITEMINFOW mii;
    mii.cbSize = sizeof(MENUITEMINFOW);
    mii.fMask = MIIM_TYPE;
    std::wstring wtitle;
    if (item->separator_) {
      mii.fType = MFT_SEPARATOR;
    } else {
      mii.fMask |= MIIM_ID;
      mii.wID = item->cmdId_;
      mii.fType = MFT_STRING;
      wtitle = exciton::util::ToUTF16String(item->title_);
      mii.dwTypeData = const_cast<LPWSTR>(wtitle.c_str());
    }
    if (item->subMenu_) {
      mii.fMask |= MIIM_SUBMENU;
      mii.hSubMenu = item->subMenu_->GetHMenu();
    }
    if (!::InsertMenuItemW(hMenu, -1, TRUE, &mii)) {
      LOG_ERROR("[%d] Menu::GetHMenu: InsertMenuItem failed(0x%08x)\n",
                __LINE__, ::GetLastError());
      ::DestroyMenu(hMenu);
      return NULL;
    }
  }
  return hMenu;
}

std::shared_ptr<MenuItem> Menu::FindMenuItemFromId(int cmdId) const {
  for (auto item : items_) {
    if (item->cmdId_ == cmdId) {
      return item;
    }
    if (item->subMenu_) {
      auto ret = item->subMenu_->FindMenuItemFromId(cmdId);
      if (ret) {
        return ret;
      }
    }
  }
  return std::shared_ptr<MenuItem>();
}

MenuItem::MenuItem() : cmdId_(-1), enabled_(false), separator_(false) {}

void MenuItem::setSubMenu(std::shared_ptr<Menu> menu) {
  subMenu_ = menu;
  auto pThis = shared_from_this();
  menu->hostItem = pThis;
  enabled_ = true;
}

struct MenuHolder {
  enum HOLD_TYPE {
    NONE,
    MENU,
    ITEM,
  };
  MenuHolder() : type(NONE) {}
  MenuHolder(std::shared_ptr<Menu> m) : type(MENU), menu(m) {}
  MenuHolder(std::shared_ptr<MenuItem> i) : type(ITEM), item(i) {}
  explicit operator bool() const { return type != NONE; }
  bool isMenu() const { return type == MENU; }
  bool isMenuItem() const { return type == ITEM; }
  HOLD_TYPE type;
  std::shared_ptr<Menu> menu;
  std::shared_ptr<MenuItem> item;
};

std::shared_ptr<MenuItem> resolveMenuNode(std::shared_ptr<Menu> pMenu,
                                          const picojson::array &items) {
  std::shared_ptr<MenuItem> ret;
  for (auto &item : items) {
    auto idx = item.get<int64_t>();
    if (ret) {
      ret = ret->subMenu_->itemAtIndex(idx);
    } else {
      ret = pMenu->itemAtIndex(idx);
    }
  }
  return ret;
}

bool MenuData::populateWithDiffset(const picojson::value &diffSet) {
  using namespace exciton::menu;
  auto &items = diffSet.get("items").get<picojson::array>();
  std::vector<MenuHolder> creNodes;
  MenuHolder curNode;
  MenuHolder arg1Node;
  MenuHolder arg2Node;
  for (auto &item : items) {
    // LOG_DEBUG("[%d] MenuData::populateWithDiffset: item: %s\n", __LINE__,
    // item.to_str().c_str());
    auto key = item.get("t").get<int64_t>();
    auto &k0 = item.get("k");
    std::string k;
    if (!k0.is<picojson::null>()) {
      k = k0.get<std::string>();
    }
    auto &v = item.get("v");
    switch (key) {
    case ditCreateNode: {
      auto &s = v.get<std::string>();
      if (s == "menu") {
        auto menu = std::make_shared<Menu>();
        if (!creNodes.empty() || pMenu) {
          auto mi = std::make_shared<MenuItem>();
          mi->setSubMenu(menu);
          menu->hostItem = mi;
          curNode = mi;
          creNodes.emplace_back(mi);
        } else {
          curNode = menu;
          creNodes.emplace_back(menu);
          pMenu = menu;
          pMenu->id_ = id;
        }
      } else if (s == "menuitem") {
        auto mi = std::make_shared<MenuItem>();
        creNodes.emplace_back(mi);
        curNode = mi;
      } else if (s == "hr") {
        auto mi = std::make_shared<MenuItem>();
        mi->separator_ = true;
        creNodes.emplace_back(mi);
        curNode = mi;
      } else {
        LOG_ERROR("[%d] MenuData::populateWithDiffset: ditCreateNode: "
                  "unsupported tag: %s\n",
                  __LINE__, s.c_str());
        return false;
      }
      break;
    }
    case ditSelectCurNode:
      if (v.is<picojson::null>()) {
        curNode = pMenu;
      } else if (v.is<std::int64_t>()) {
        curNode = creNodes[v.get<int64_t>()];
      } else {
        curNode = resolveMenuNode(pMenu, v.get<picojson::array>());
      }
      break;
    case ditSelectArg1Node:
      if (v.is<picojson::null>()) {
        arg1Node = pMenu;
      } else if (v.is<std::int64_t>()) {
        arg1Node = creNodes[v.get<int64_t>()];
      } else {
        arg1Node = resolveMenuNode(pMenu, v.get<picojson::array>());
      }
      break;
    case ditSelectArg2Node:
      if (v.is<picojson::null>()) {
        arg2Node = pMenu;
      } else if (v.is<std::int64_t>()) {
        arg2Node = creNodes[v.get<int64_t>()];
      } else {
        arg2Node = resolveMenuNode(pMenu, v.get<picojson::array>());
      }
      break;
    case ditAttributeValue:
      if (!curNode.isMenuItem()) {
        if (k != "type") {
          LOG_ERROR("[%d] MenuData::populateWithDiffset: ditAttributeValue: "
                    "invalid attribute: %s",
                    __LINE__, k.c_str());
          return false;
        }
        // do nothing yet...
      } else {
        auto mi = curNode.item;
        if (k == "label") {
          mi->title_ = v.get<std::string>();
          if (mi->subMenu_) {
            mi->subMenu_->title_ = v.get<std::string>();
          }
        }
      }
      break;
    case ditDelAttributeValue:
      LOG_WARNING("[%d] MenuData::populateWithDiffset: Not implement yet: "
                  "ditDelAttributeValue",
                  __LINE__);
      break;
    case ditAddDataSet: {
      if (!curNode || curNode.type != MenuHolder::ITEM) {
        LOG_ERROR("[%d] MenuData::populateWithDiffset: ditAddEventListener: "
                  "invalid target:",
                  __LINE__);
        return false;
      }
      auto mi = curNode.item;
      auto &val = v.get<std::string>();
      if (k == "menuRole") {
        auto role = getDefaultRoleInfo(val);
        if (!role) {
          // TODO: ルートメニューのroleはどうする?
          LOG_WARNING("[%d] MenuData::populateWithDiffset: ditAddDataSet: "
                      "unsupported role name: %s",
                      __LINE__, val.c_str());
          // TODO: 未知のロールのルートメニューは単に無視するだけにする?
          break;
          // return false;
        }
        // TODO: roleの場合にItemに設定する値は??
        if (role->label) {
          auto appNameW = Driver::Current().GetProductName();
          auto appName = exciton::util::ToUTF8String(appNameW.c_str());
          auto labelStr =
              exciton::util::FormatString(role->label, appName.c_str());
          mi->title_ = labelStr;
          if (mi->subMenu_) {
            mi->subMenu_->title_ = labelStr;
          }
        }
        mi->cmdId_ = static_cast<int>(role->command);
        mi->enabled_ = true;
      } else if (k == "menuAcclerator") {
        // TODO : あとで書く
        LOG_WARNING("[%d] MenuData::populateWithDiffset: ditAddDataSet: "
                    "menuRole[menuAccelerator] not implement yet.",
                    __LINE__);
      } else {
        LOG_WARNING("[%d] MenuData::populateWithDiffset: ditAddDataSet: "
                    "unknwon dataSet Name:%s",
                    __LINE__, k.c_str());
      }
      break;
    }
    case ditDelDataSet:
      LOG_WARNING("[%d] MenuData::populateWithDiffset: ditDelDataSet: Not "
                  "implement yet.",
                  __LINE__);
      break;
    case ditAppendChild: {
      std::shared_ptr<Menu> target;
      if (!curNode) {
        target = pMenu;
      } else if (curNode.isMenu()) {
        target = curNode.menu;
      } else if (curNode.isMenuItem()) {
        target = curNode.item->subMenu_;
      }
      if (!target) {
        LOG_ERROR(
            "[%d] MenuData::populateWithDiffset: ditAppendChild: invalid arg1.",
            __LINE__);
        return false;
      }
      if (!arg1Node.isMenuItem()) {
        if (arg1Node.menu != target) {
          LOG_ERROR("[%d] MenuData::populateWithDiffset: ditAppendChild: "
                    "invalid arg1.",
                    __LINE__);
          return false;
        }
      } else {
        target->AddMenuItem(arg1Node.item);
      }
      break;
    }
    case ditInsertBefore:
      LOG_WARNING("[%d] MenuData::populateWithDiffset: ditInsertBefore: Not "
                  "implement yet.",
                  __LINE__);
      break;
    case ditRemoveChild:
      LOG_WARNING("[%d] MenuData::populateWithDiffset: ditRemoveChild: Not "
                  "implement yet.",
                  __LINE__);
      break;
    case ditAddEventListener: {
      if (!curNode.isMenuItem()) {
        LOG_ERROR("[%d] MenuData::populateWithDiffset: ditAddEventListener: "
                  "invalid target.",
                  __LINE__);
        return false;
      }
      if (k != "click") {
        LOG_ERROR("[%d] MenuData::populateWithDiffset: ditAddEventListener: "
                  "unsupported event:[%s]",
                  __LINE__, k.c_str());
        return FALSE;
      }
      auto mi = curNode.item;
      auto id = v.get("id").get<std::string>();
      mi->cmdId_ = CMenuModel::Instance().GetNewCommandId(id);
      mi->onClick_ = id;
      mi->enabled_ = true;
      break;
    }
    case ditRemoveEventListener:
      LOG_WARNING("[%d] MenuData::populateWithDiffset: ditRemoveEventListener: "
                  "Not implement yet.",
                  __LINE__);
      break;
    case ditNodeUUID: {
      auto id = v.get<std::string>();
      if (curNode.isMenu()) {
        auto menu = curNode.menu;
        if (menu != pMenu) {
          menu->id_ = id;
          auto hi = menu->hostItem.lock();
          if (hi) {
            hi->id_ = id;
          }
        }
      } else if (curNode.isMenuItem()) {
        auto item = curNode.item;
        item->id_ = id;
      } else {
        LOG_ERROR("[%d] MenuData::populateWithDiffset: ditNodeUUID: current "
                  "node is invalid.",
                  __LINE__);
        return false;
      }
      break;
    }
    case ditAddClassList:
    case ditDelClassList:
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
    default:
      LOG_ERROR(
          "[%d] MenuData::populateWithDiffset:: Unsupported item type[%d]",
          __LINE__, key);
      return FALSE;
    }
  }
  return true;
}

} // namespace menu
} // namespace exciton

CMenuModel::CMenuModel() { ::InitializeCriticalSection(&cs_); }
CMenuModel::~CMenuModel() { ::DeleteCriticalSection(&cs_); }

CMenuModel &CMenuModel::Instance() {
  static CMenuModel s_instance;

  return s_instance;
}

std::shared_ptr<exciton::menu::Menu>
CMenuModel::getMenuFromId(const std::string &id) const {
  auto fiter = menus_.find(id);
  if (fiter == menus_.end()) {
    return nullptr;
  }
  return (*fiter).second->pMenu;
}

void CMenuModel::NewMenu(const picojson::value &argument,
                         const std::map<std::string, std::string> &parameter,
                         int responceNo) {
  Driver &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR("[%d] CMenuModel::NewMenu: parameter['id'] not found", __LINE__);
    Driver::Current().responceEventBoolResult(responceNo, false);
    return;
  }
  auto menuData = std::make_shared<exciton::menu::MenuData>();
  menus_[id] = menuData;
  menuData->id = id;
  d.responceEventBoolResult(responceNo, true);
}

void CMenuModel::UpdateDiffSetHandler(
    const picojson::value &argument,
    const std::map<std::string, std::string> &parameter, int responceNo) {
  Driver &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR(
        "[%d] CMenuModel::UpdateDiffSetHandler: parameter['id'] not found.",
        __LINE__);
    d.responceEventBoolResult(responceNo, false);
    return;
  }
  auto fiter = menus_.find(id);
  if (fiter == menus_.end()) {
    LOG_ERROR("[%d] CMenuModel::UpdateDiffSetHandler: menu not found['%s']",
              __LINE__, id.c_str());
    d.responceEventBoolResult(responceNo, false);
    return;
  }

  auto menuData = (*fiter).second;

  if (!menuData->populateWithDiffset(argument)) {
    LOG_ERROR("[%d] CMenuModel::UpdateDiffSetHandler: populateWithDiffset "
              "failed['%s']",
              __LINE__, id.c_str());
    d.responceEventBoolResult(responceNo, false);
    return;
  }

  // notify menu update
  d.notifyUpdateMenu(id);

  Driver::Current().responceEventBoolResult(responceNo, true);
}

void CMenuModel::SetApplicationMenu(
    const picojson::value &argument,
    const std::map<std::string, std::string> &parameter, int responceNo) {
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR(
        "[%d] CMenuModel::UpdateDiffSetHandler: parameter['id'] not found.",
        __LINE__);
    return;
  }
  applicationMenuId_ = id;
}

void CMenuModel::PopupContextMenu(
    const picojson::value &argument,
    const std::map<std::string, std::string> &parameter, int responceNo) {
  auto &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR(
        "[%d] CMenuModel::UpdateDiffSetHandler: parameter['id'] not found.",
        __LINE__);
    return;
  }
  auto menu = getMenuFromId(id);
  auto posX = argument.get("position").get("x").get<double>();
  auto posY = argument.get("position").get("y").get<double>();
  auto winidstr = argument.get("windowId").get<std::string>();
  CWebBrowserHost *pWindow = resolveParentWindow(winidstr);
  HWND hWnd = NULL;
  pWindow->GetWindow(&hWnd);
  HMENU hMenu = menu->GetHMenu();
  UINT pmflags = TPM_RETURNCMD | TPM_NONOTIFY;
  // if (m_bRTL) {
  //   pmflags |= TPM_LAYOUTRTL;
  // }
  int iRetCmd = ::TrackPopupMenu(hMenu, pmflags, static_cast<int>(posX),
                                 static_cast<int>(posY), 0, hWnd, nullptr);

  ::DestroyMenu(hMenu);
  if (iRetCmd != 0) {
    auto item = menu->FindMenuItemFromId(iRetCmd);
    if (item) {
      CMenuMgr::OnMenuCommand(iRetCmd, pWindow, menu, item);
    }
  }
}

std::shared_ptr<exciton::menu::Menu> CMenuModel::GetApplicationMenu() const {
  return getMenuFromId(applicationMenuId_);
}

int CMenuModel::GetNewCommandId(const std::string &id) {
  int iRet = -1;
  ::EnterCriticalSection(&cs_);
  auto fiter = menuCommands_.find(id);
  if (fiter != menuCommands_.end()) {
    iRet = (*fiter).second;
  } else {
    auto f = removedMenuCommands_.begin();
    if (f != removedMenuCommands_.end()) {
      iRet = (*f);
      removedMenuCommands_.erase(f);
    } else {
      iRet = static_cast<int>(exciton::menu::RoledCommandId::UserCommand) +
             menuCommands_.size();
      menuCommands_[id] = iRet;
    }
  }
  ::LeaveCriticalSection(&cs_);
  return iRet;
}

namespace {
void newMenu(const picojson::value &argument,
             const std::map<std::string, std::string> &parameter,
             int responceNo) {
  CMenuModel::Instance().NewMenu(argument, parameter, responceNo);
}
void updateDiffSetHandler(const picojson::value &argument,
                          const std::map<std::string, std::string> &parameter,
                          int responceNo) {
  CMenuModel::Instance().UpdateDiffSetHandler(argument, parameter, responceNo);
}

void setApplicationMenu(const picojson::value &argument,
                        const std::map<std::string, std::string> &parameter,
                        int responceNo) {
  CMenuModel::Instance().SetApplicationMenu(argument, parameter, responceNo);
}

void popupContextMenu(const picojson::value &argument,
                      const std::map<std::string, std::string> &parameter,
                      int responceNo) {
  CMenuModel::Instance().PopupContextMenu(argument, parameter, responceNo);
}

} // namespace

void Menu_Init() {
  auto &d = Driver::Current();
  d.addEventHandler("/menu/:id/new", newMenu);
  d.addEventHandler("/menu/:id/updateDiffSetHandler", updateDiffSetHandler);
  d.addEventHandler("/menu/:id/setApplicationMenu", setApplicationMenu);
  d.addDeferEventHandler("/menu/:id/popupContextMenu", popupContextMenu);
}