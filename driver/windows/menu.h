#ifndef __INC_WINDOWS_WINDOW_H__
#define __INC_WINDOWS_WINDOW_H__
#ifdef __cplusplus
#include <map>
#include <memory>
#include <set>
#include <string>
#include <vector>

#include "myjson.h"

// menu model
namespace exciton {
namespace menu {

enum class RoledCommandId : int {
  About = 1,
  Front,
  //  Undo,
  //  Redo,
  Cut,
  Copy,
  Paste,
  Delete,
  //  PasteAndMatchStyle,
  SelectAll,
  Minimize,
  Close,
  Zoom,
  Quit,
  ToggleFullscreen,
  ViewSource,
  UserCommand,
};
class MenuItem;
struct Menu : public std::enable_shared_from_this<Menu> {
  std::string id_;
  std::vector<std::shared_ptr<MenuItem>> items_;
  std::weak_ptr<MenuItem> hostItem;
  std::string title_;
  Menu();
  Menu(const std::string &id);
  Menu(const Menu &) = delete;
  Menu(Menu &&) = delete;
  std::size_t GetMenuItemCount() const { return items_.size(); }
  std::shared_ptr<MenuItem> GetMenuItem(size_t idx) const { return items_[idx]; }
  void AddMenuItem(std::shared_ptr<MenuItem> item);
  std::shared_ptr<Menu> GetSubMenu(size_t idx) const;
  std::shared_ptr<MenuItem> itemAtIndex(size_t idx) const;
  HMENU GetHMenuAtIndex(size_t idx) const;
  HMENU GetHMenu() const;
  std::shared_ptr<MenuItem> FindMenuItemFromId(int cmdId) const;
  const std::string& ID() const { return id_; }
};
struct MenuItem : public std::enable_shared_from_this<MenuItem> {
  std::string id_;
  std::string onClick_;
  std::shared_ptr<Menu> subMenu_;
  std::string title_;
  int cmdId_;
  bool enabled_;
  bool separator_;
  MenuItem();

public:
  void setSubMenu(std::shared_ptr<Menu> menu);
  const std::string ID() const { return id_; }
  const std::string EventName() const { return onClick_; }
};

struct MenuData {
  std::shared_ptr<Menu> pMenu;
  std::string id;
  bool populateWithDiffset(const picojson::value &diffSet);
};

} // namespace menu
} // namespace exciton

class CMenuModel {
  CMenuModel();
  std::map<std::string, std::shared_ptr<exciton::menu::MenuData>> menus_;
  std::string applicationMenuId_;
  std::map<std::string, int> menuCommands_;
  std::set<int> removedMenuCommands_;
  CRITICAL_SECTION cs_;

public:
  static CMenuModel &Instance();

public:
  ~CMenuModel();

public:
  void NewMenu(const picojson::value &argument,
               const std::map<std::string, std::string> &parameter,
               int responceNo);
  void UpdateDiffSetHandler(const picojson::value &argument,
                            const std::map<std::string, std::string> &parameter,
                            int responceNo);
  void SetApplicationMenu(const picojson::value &argument,
                          const std::map<std::string, std::string> &parameter,
                          int responceNo);
  void PopupContextMenu(const picojson::value &argument,
                          const std::map<std::string, std::string> &parameter,
                          int responceNo);
  int GetNewCommandId(const std::string& id);

public:
  std::shared_ptr<exciton::menu::Menu> GetApplicationMenu() const;

private:
  std::shared_ptr<exciton::menu::Menu>
  getMenuFromId(const std::string &id) const;
};

extern "C" {
#endif

void Menu_Init();

#ifdef __cplusplus
};
#endif

#endif