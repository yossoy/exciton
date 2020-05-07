#pragma once

#include <string>
#include <vector>
class CWebBrowserContainer;
class CWebBrowserHost;
namespace exciton {
  namespace menu {
    struct Menu;
    struct MenuItem;
  }
}

class CMenuMgr : public IShellMenuCallback {
public:
  STDMETHODIMP QueryInterface(REFIID riid, void **ppvObject);
  STDMETHODIMP_(ULONG) AddRef();
  STDMETHODIMP_(ULONG) Release();

  STDMETHODIMP CallbackSM(LPSMDATA psmd, UINT uMsg, WPARAM wParam,
                          LPARAM lParam);

  CMenuMgr(CWebBrowserContainer &container);
  CMenuMgr(const CMenuMgr &) = delete;
  CMenuMgr(CMenuMgr &&) = delete;
  ~CMenuMgr();
  BOOL Create(HWND hwndParent);
  void Destroy();
  static void OnMenuCommand(int nId, CWebBrowserHost* pWebBrowserHost, std::shared_ptr<exciton::menu::Menu> pMenu, std::shared_ptr<exciton::menu::MenuItem> menuItem);
  void SetMenuState(HMENU hmenu, int nPos);
  static void SetMenuState(CWebBrowserHost* pWebBrowserHost, HMENU hmenu);
  HWND GetWindow() const { return m_hwndRebar; }

public:
  static bool InitClass();
  static void FinalizeClass();
  static bool IsAppThemed();
  static LRESULT CALLBACK MenuBarWndProc(HWND win, UINT msg, WPARAM wp,
                                         LPARAM lp);
  static LRESULT CALLBACK MenuBarHotTraceProc(int code, WPARAM wp, LPARAM lp);

  static int MenuBarSendMsg(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam) {
    return ::CallWindowProcW(s_origToolBarProc, hWnd, msg, wParam, lParam);
  }
  void UpdateUIState(bool keyboard_activity);
//  int SetMenu(HMENU menu, BOOL is_refresh);
  int SetMenu(std::shared_ptr<exciton::menu::Menu> pMenu, bool is_refresh);
  void UpdateMenu(const std::string& menuId);

  void OnResize();

  // HotTrace
  void EnableHotTrace();
  void DisableHotTrace();
  static void PerformDisableHotTrace(void);

private:
  void ResetHotItem();
  void PerformDropDown();
  void DropDown(int item, bool from_keyboard);
  LRESULT OnNotify(NMHDR* hdr);
  BOOL OnKeyDown(int vk, DWORD key_data);
  static CMenuMgr* OnNcCreate(HWND hWnd, UINT msg, WPARAM wParam, LPARAM lParam);
  int OnCreate(const CREATESTRUCTW* hWnd);
  void OnDestroy();
  static void OnNcDestroy(CMenuMgr* mgr);
  void HotTraceChangeDropDown(int item, bool from_keyboard);
  static void handleMenuEvent(std::shared_ptr<exciton::menu::Menu> pMenu, std::shared_ptr<exciton::menu::MenuItem> menuItem);


private:
  CWebBrowserContainer &m_container;
  LONG m_cRef;
  HWND m_hwndParent;
  HWND m_hwndRebar;

  HWND m_hWnd;
  HWND m_hNotifyWnd;
  HWND m_hWndOldFocus;
  std::shared_ptr<exciton::menu::Menu> m_pMenu;
  int m_iPressedItem; /* = -1 */
  short m_iHotItem;
  bool m_bContinueHotTrace;
  bool m_bSelectFromKeyboard;
  bool m_bRTL;
  LONG m_lDropDownCloseTime;
  static CMenuMgr* s_pActiveMenuBar;
  // HotTrace
  static WNDPROC s_origToolBarProc;
  static int s_wndClsExtraOffset;
  static CRITICAL_SECTION s_htCS;
  static HHOOK s_htHook;
  static CMenuMgr *s_htMenuMgr;
  static std::shared_ptr<exciton::menu::Menu> s_htSelMenu;
  static int s_htSelItem;
  static UINT s_htSelFlags;
  static POINT s_htLastPos;
};
