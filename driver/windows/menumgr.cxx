#include <windows.h>

#include <exdispid.h>
#include <mshtmcid.h>
#include <mshtmhst.h>
#include <mshtml.h>
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
const WCHAR MenuMgrWndClassName[] = L"exciton.MenuBar";
const int MENUBAR_ITEM_LABEL_MAXSIZE = 32;
const int MENUBAR_SEPARATOR_WIDTH = 10;
const UINT MENUBAR_DTFLAGS = (DT_CENTER | DT_VCENTER | DT_SINGLELINE);
} // namespace

#define BAND_MENUBAR 1

#define MENUBAR_WM_SET_MENU (WM_APP + 1)

HRESULT CALLBACK TaskDialogCallbackProc(HWND hwnd, UINT uNotification,
                                        WPARAM wParam, LPARAM lParam,
                                        LONG_PTR dwRefData) {
#if 0
	//TODO:
	if (uNotification == TDN_HYPERLINK_CLICKED)
		g_pWebBrowserContainer->Navigate((LPWSTR)lParam, TRUE);
#endif

  return S_OK;
}

// CMenuMgr

CMenuMgr::CMenuMgr(CWebBrowserContainer &container)
    : m_container(container), m_cRef(1), m_hwndParent(NULL), m_hwndRebar(NULL),
      m_hWnd(NULL), m_hNotifyWnd(NULL), m_hWndOldFocus(NULL),
      m_iPressedItem(-1), m_iHotItem(-1), m_bContinueHotTrace(false),
      m_bSelectFromKeyboard(false), m_bRTL(false), m_lDropDownCloseTime(0) {}

CMenuMgr::~CMenuMgr() {}

STDMETHODIMP CMenuMgr::QueryInterface(REFIID riid, void **ppvObject) {
  *ppvObject = NULL;

  if (IsEqualIID(riid, IID_IUnknown) ||
      IsEqualIID(riid, IID_IShellMenuCallback))
    *ppvObject = static_cast<IShellMenuCallback *>(this);
  else
    return E_NOINTERFACE;

  AddRef();

  return S_OK;
}

STDMETHODIMP_(ULONG) CMenuMgr::AddRef() {
  return InterlockedIncrement(&m_cRef);
}

STDMETHODIMP_(ULONG) CMenuMgr::Release() {
  return InterlockedDecrement(&m_cRef);
}

STDMETHODIMP CMenuMgr::CallbackSM(LPSMDATA psmd, UINT uMsg, WPARAM wParam,
                                  LPARAM lParam) {
#if 0
	if (uMsg == SMC_SFEXEC) {
		IShellItem *psi;

		SHCreateItemWithParent(psmd->pidlFolder, psmd->psf, psmd->pidlItem, IID_PPV_ARGS(&psi));
		NavigateFromShortcut(psi, FALSE);
		psi->Release();
		return S_OK;
	}
#endif

  return S_FALSE;
}

BOOL CMenuMgr::Create(HWND hwndParent) {

  /* Create ReBar window */
  m_hwndRebar = ::CreateWindowEx(WS_EX_TOOLWINDOW, L"ReBarWindow32", L"",
                                 WS_CHILD | WS_VISIBLE | WS_CLIPCHILDREN |
                                     WS_CLIPSIBLINGS | WS_BORDER |
                                     CCS_NODIVIDER | CCS_TOP | RBS_VARHEIGHT |
                                     RBS_BANDBORDERS | RBS_AUTOSIZE,
                                 0, 0, 0, 0, hwndParent, (HMENU)1000,
                                 Driver::Current().InstanceHandle(), NULL);

  /* Create the menubar control */
  HWND hwndMenubar =
      ::CreateWindowW(MenuMgrWndClassName, L"",
                      WS_CHILD | WS_VISIBLE | WS_CLIPCHILDREN |
                          WS_CLIPSIBLINGS | CCS_NORESIZE | CCS_NOPARENTALIGN,
                      0, 0, 0, 0, m_hwndRebar, (HMENU)1001,
                      Driver::Current().InstanceHandle(), (LPVOID)this);
  ::SendMessageW(hwndMenubar, TB_SETEXTENDEDSTYLE, 0,
                 TBSTYLE_EX_HIDECLIPPEDBUTTONS);
  SIZE szIdeal;
  ::SendMessageW(hwndMenubar, TB_GETIDEALSIZE, FALSE, (LPARAM)&szIdeal);

  /* Embed the menubar in the ReBar */
  REBARBANDINFO band = {0};
  band.cbSize = REBARBANDINFO_V6_SIZE;
  band.fMask = RBBIM_STYLE | RBBIM_CHILD | RBBIM_CHILDSIZE | RBBIM_SIZE |
               RBBIM_IDEALSIZE | RBBIM_ID;
  band.fStyle = RBBS_GRIPPERALWAYS | RBBS_USECHEVRON | RBBS_VARIABLEHEIGHT;
  band.hwndChild = hwndMenubar;
  DWORD dwBtnSize =
      (DWORD)::SendMessage(band.hwndChild, TB_GETBUTTONSIZE, 0, 0);
  WORD dwBtnHeight = HIWORD(dwBtnSize);
  band.cyChild = dwBtnHeight;
  band.cxMinChild = 0;
  band.cyMinChild = dwBtnHeight;
  band.cyMaxChild = dwBtnHeight;
  band.cyIntegral = dwBtnHeight;
  band.cx = 240;
  band.cxIdeal = szIdeal.cx;
  band.wID = BAND_MENUBAR;
  int r = SendMessage(m_hwndRebar, RB_INSERTBAND, -1, (LPARAM)&band);

  auto appMenu = CMenuModel::Instance().GetApplicationMenu();
  if (appMenu) {
    SetMenu(appMenu, false);
  }

  ::SendMessageW(m_hWnd, TB_GETIDEALSIZE, FALSE, (LPARAM)&szIdeal);
  dwBtnSize = (DWORD)::SendMessage(m_hWnd, TB_GETBUTTONSIZE, 0, 0);
  dwBtnHeight = HIWORD(dwBtnSize);
  band.fMask = RBBIM_CHILDSIZE | RBBIM_SIZE | RBBIM_IDEALSIZE;
  band.cyChild = dwBtnHeight;
  band.cxMinChild = 0;
  band.cyMinChild = dwBtnHeight;
  band.cyMaxChild = dwBtnHeight;
  band.cyIntegral = dwBtnHeight;
  band.cx = 240;
  band.cxIdeal = szIdeal.cx;
  r = SendMessage(m_hwndRebar, RB_SETBANDINFO, 0, (LPARAM)&band);

  //::PostMessageW(m_hWnd, )
  m_hwndParent = hwndParent;

  return TRUE;
}

void CMenuMgr::Destroy() { ::DestroyWindow(m_hWnd); }

void CMenuMgr::handleMenuEvent(
    std::shared_ptr<exciton::menu::Menu> pMenu,
    std::shared_ptr<exciton::menu::MenuItem> menuItem) {
  LOG_INFO("[%d] handleMenuEvent: %s", __LINE__, menuItem->EventName().c_str());

  picojson::value val(menuItem->ID());
  auto json = val.serialize();
  auto name = exciton::util::FormatString("/menu/%s", pMenu->ID().c_str());
  Driver::Current().emitEvent(name, "emit", json);
}

void CMenuMgr::OnMenuCommand(
    int nId, CWebBrowserHost *pWebBrowserHost,
    std::shared_ptr<exciton::menu::Menu> pMenu,
    std::shared_ptr<exciton::menu::MenuItem> menuItem) {
  using namespace exciton::menu;

  switch (static_cast<RoledCommandId>(nId)) {
  case RoledCommandId::About: {
    TASKDIALOGCONFIG config;
    int nResult;
    HWND hwndParent;
    pWebBrowserHost->GetWindow(&hwndParent);
    // HWND hwndParent = m_container.GetWindow();
    // TODO: fix implement

    ZeroMemory(&config, sizeof(TASKDIALOGCONFIG));
    auto appNameStr = Driver::Current().GetProductName();
    auto appVerStr = L"Version: " + Driver::Current().GetProductVersion();
    config.cbSize = sizeof(TASKDIALOGCONFIG);
    config.hwndParent = hwndParent;
    config.hInstance = NULL;
    config.dwFlags = TDF_ENABLE_HYPERLINKS;
    config.pszWindowTitle = L"Version information";
    config.pszMainIcon = TD_INFORMATION_ICON;
    config.pszMainInstruction = appNameStr.c_str();
    config.pszContent = appVerStr.c_str();
    config.pfCallback = TaskDialogCallbackProc;

    TaskDialogIndirect(&config, &nResult, NULL, NULL);
    break;
  }
  case RoledCommandId::Front:
    LOG_ERROR("[%d] OnMenuCommand(Fornt) not implement yet", __LINE__);
    break;
  case RoledCommandId::Cut:
    pWebBrowserHost->ExecDocument(&CMDSETID_Forms3, IDM_CUT,
                                  OLECMDEXECOPT_DODEFAULT, NULL, NULL);
    break;
  case RoledCommandId::Copy:
    pWebBrowserHost->ExecDocument(&CMDSETID_Forms3, IDM_COPY,
                                  OLECMDEXECOPT_DODEFAULT, NULL, NULL);
    break;
  case RoledCommandId::Paste:
    pWebBrowserHost->ExecDocument(&CMDSETID_Forms3, IDM_PASTE,
                                  OLECMDEXECOPT_DODEFAULT, NULL, NULL);
    break;
  case RoledCommandId::Delete:
    pWebBrowserHost->ExecDocument(&CMDSETID_Forms3, IDM_DELETE,
                                  OLECMDEXECOPT_DODEFAULT, NULL, NULL);
    break;
  case RoledCommandId::SelectAll:
    pWebBrowserHost->Exec(OLECMDID_SELECTALL, OLECMDEXECOPT_DODEFAULT, NULL,
                          NULL);
    break;
  case RoledCommandId::Minimize:
    ::SendMessageW(pWebBrowserHost->GetHostContainer()->GetWindow(),
                   WM_SYSCOMMAND, SC_MINIMIZE, 0);
    break;
  case RoledCommandId::Close:
    // TODO: close BrowserHost (tab?)
    // TODO: query close event?
    ::SendMessageW(pWebBrowserHost->GetHostContainer()->GetWindow(), WM_CLOSE,
                   0, 0);
    break;
  case RoledCommandId::Zoom:
    ::SendMessageW(pWebBrowserHost->GetHostContainer()->GetWindow(),
                   WM_SYSCOMMAND, SC_MAXIMIZE, 0);
    break;
  case RoledCommandId::Quit:
    // TODO: query quit message?
    ::PostQuitMessage(0);
    break;
  case RoledCommandId::ToggleFullscreen:
    pWebBrowserHost->PutFullscreen(true);
    break;
  case RoledCommandId::ViewSource: {
    const GUID CGID_IWebBrowser = {
        0xED016940L,
        0xBD5B,
        0x11cf,
        {0xBA, 0x4E, 0x00, 0xC0, 0x4F, 0xD7, 0x08, 0x16}};
    const int HTMLID_VIEWSOURCE = 2;
    pWebBrowserHost->ExecDocument(&CGID_IWebBrowser, HTMLID_VIEWSOURCE, 0, NULL,
                                  NULL);
    break;
  }
  case RoledCommandId::HistoryGoBack:
    pWebBrowserHost->ExecDocument(&CMDSETID_Forms3, IDM_GOBACKWARD,
                                  OLECMDEXECOPT_DODEFAULT, NULL, NULL);
    break;
  case RoledCommandId::HistoryGoForward:
    pWebBrowserHost->ExecDocument(&CMDSETID_Forms3, IDM_GOFORWARD,
                                  OLECMDEXECOPT_DODEFAULT, NULL, NULL);
    break;
  default: {
    // command
    handleMenuEvent(pMenu, menuItem);
    break;
  }
  }
}

void CMenuMgr::SetMenuState(HMENU hMenu, int nPos) {
  using namespace exciton::menu;
  MENUINFO mi;
  mi.cbSize = sizeof(mi);
  mi.fMask = MIM_MENUDATA;
  ::GetMenuInfo(hMenu, &mi);
  if (mi.dwMenuData == 0) {
    return;
  }
  auto pMenu = reinterpret_cast<Menu *>(mi.dwMenuData)->shared_from_this();

  int n = GetMenuItemCount(hMenu);
  for (int i = 0; i < n; i++) {
    int nId = GetMenuItemID(hMenu, i);
    if (nId == 0)
      continue;
    if (nId >= static_cast<int>(RoledCommandId::UserCommand)) {
      auto item = pMenu->GetMenuItem(i);
      if (item && item->cmdId_ == nId) {
        ::SetMenuItem(hMenu, nId, item->enabled_ ? TRUE : FALSE, FALSE);
      }
    } else {
      CWebBrowserHost *pWebBrowserHost = m_container.GetActiveBrowser();
      BOOL bEnable = pWebBrowserHost != NULL;
      BOOL bChecked = FALSE;
      if (pWebBrowserHost) {
        switch (static_cast<RoledCommandId>(nId)) {
        case RoledCommandId::Cut:
          pWebBrowserHost->QueryExecDocument(&CMDSETID_Forms3, IDM_CUT,
                                             &bEnable, &bChecked);
          break;
        case RoledCommandId::Copy:
          pWebBrowserHost->QueryExecDocument(&CMDSETID_Forms3, IDM_COPY,
                                             &bEnable, &bChecked);
          break;
        case RoledCommandId::Paste:
          pWebBrowserHost->QueryExecDocument(&CMDSETID_Forms3, IDM_PASTE,
                                             &bEnable, &bChecked);
          break;
        case RoledCommandId::Delete:
          pWebBrowserHost->QueryExecDocument(&CMDSETID_Forms3, IDM_DELETE,
                                             &bEnable, &bChecked);
          break;
        case RoledCommandId::HistoryGoBack:
          pWebBrowserHost->QueryExecDocument(&CMDSETID_Forms3, IDM_GOBACKWARD,
                                             &bEnable, &bChecked);
          break;
        case RoledCommandId::HistoryGoForward:
          pWebBrowserHost->QueryExecDocument(&CMDSETID_Forms3, IDM_GOFORWARD,
                                             &bEnable, &bChecked);
          break;
        default:
          break;
        }
      }
      ::SetMenuItem(hMenu, nId, bEnable, bChecked);
    }
  }
}

///-----------------------------------------------------------
WNDPROC CMenuMgr::s_origToolBarProc = NULL;
int CMenuMgr::s_wndClsExtraOffset;
CRITICAL_SECTION CMenuMgr::s_htCS;
HHOOK CMenuMgr::s_htHook = nullptr;
CMenuMgr *CMenuMgr::s_htMenuMgr = nullptr;
std::shared_ptr<exciton::menu::Menu> CMenuMgr::s_htSelMenu;
int CMenuMgr::s_htSelItem = -1;
UINT CMenuMgr::s_htSelFlags = 0;
POINT CMenuMgr::s_htLastPos = {0};

CMenuMgr *CMenuMgr::s_pActiveMenuBar = nullptr;

bool CMenuMgr::InitClass() {
  WNDCLASSW wc = {0};
  INITCOMMONCONTROLSEX icce = {0};

  icce.dwSize = sizeof(INITCOMMONCONTROLSEX);
  icce.dwICC = ICC_BAR_CLASSES | ICC_COOL_CLASSES;

  if (!InitCommonControlsEx(&icce)) {
    LOG_ERROR("[%d] CMenuMgr::InitClass: InitCommonControlsEx() failed",
              __LINE__);
    return false;
  }

  if (!::GetClassInfoW(NULL, L"ToolbarWindow32", &wc)) {
    LOG_ERROR("[%d] CMenuMgr::InitClass: GetClassInfo() failed", __LINE__);
    return false;
  }

  /* Remember needed values of standard toolbar window class */
  s_origToolBarProc = wc.lpfnWndProc;
  s_wndClsExtraOffset = wc.cbWndExtra;

  /* Create our subclass. */
  wc.lpfnWndProc = MenuBarWndProc;
  wc.cbWndExtra += sizeof(CMenuMgr *);
  wc.style |= CS_GLOBALCLASS;
  wc.hInstance = NULL;
  wc.lpszClassName = MenuMgrWndClassName;
  if (!::RegisterClassW(&wc)) {
    LOG_ERROR("[%d] CMenuMgr::InitClass: RegisterClass() failed", __LINE__);
    return false;
  }

  ::InitializeCriticalSection(&s_htCS);

  return true;
}

void CMenuMgr::FinalizeClass() {
  ::DeleteCriticalSection(&s_htCS);
  ::UnregisterClassW(MenuMgrWndClassName, NULL);
}

bool CMenuMgr::IsAppThemed() {
  // TODO:
  return true;
}

void CMenuMgr::UpdateUIState(bool keyboard_activity) {
  WORD action;

  if (keyboard_activity) {
    action = UIS_CLEAR;
  } else {
    BOOL show_accel_always;
    if (!SystemParametersInfoW(SPI_GETMENUUNDERLINES, 0, &show_accel_always,
                               0)) {
      show_accel_always = TRUE;
    }
    action = (show_accel_always ? UIS_CLEAR : UIS_SET);
  }

  ::PostMessage(m_hWnd, WM_CHANGEUISTATE, MAKELONG(action, UISF_HIDEACCEL), 0);
}

// int CMenuMgr::SetMenu(HMENU menu, BOOL is_refresh) {
int CMenuMgr::SetMenu(std::shared_ptr<exciton::menu::Menu> pMenu,
                      bool is_refresh) {
  int i, n;

  if (pMenu == m_pMenu && !is_refresh) {
    return 0;
  }

  /* If dropped down, cancel it */
  if (m_iPressedItem >= 0) {
    DisableHotTrace();
    MenuBarSendMsg(m_hWnd, WM_CANCELMODE, 0, 0);
  }

  /* Uninstall the old menu */
  if (m_pMenu) {
    n = MenuBarSendMsg(m_hWnd, TB_BUTTONCOUNT, 0, 0);
    for (i = 0; i < n; i++) {
      MenuBarSendMsg(m_hWnd, TB_DELETEBUTTON, 0, 0);
    }

    m_pMenu.reset();
  }

  /* Install the new menu */
  n = (pMenu ? pMenu->GetMenuItemCount() : 0);
  if (n < 0) {
    LOG_ERROR("[%d] CMenuMgr::SetMenu: GetMenuItemCount() failed.", __LINE__);
    return -1;
  }

  if (n == 0) {
    m_pMenu = pMenu;
    return 0;
  }

  BYTE *buffer = (BYTE *)_malloca(
      n * sizeof(TBBUTTON) + n * sizeof(TCHAR) * MENUBAR_ITEM_LABEL_MAXSIZE);
  TBBUTTON *buttons = (TBBUTTON *)buffer;
  WCHAR *labels = (TCHAR *)(buffer + n * sizeof(TBBUTTON));

  memset(buttons, 0, n * sizeof(TBBUTTON));

  for (i = 0; i < n; i++) {
    auto item = pMenu->GetMenuItem(i);

    buttons[i].iBitmap = I_IMAGENONE;
    buttons[i].fsState = 0;
    if (item->enabled_) {
      buttons[i].fsState |= TBSTATE_ENABLED;
    }
    if (item->separator_ && (i > 0)) {
      buttons[i - 1].fsState |= TBSTATE_WRAP;
    }

    if (item->subMenu_) {
      std::wstring wlabel = exciton::util::ToUTF16String(item->title_);
      WCHAR *label = labels + i * MENUBAR_ITEM_LABEL_MAXSIZE;
      wcscpy(label, wlabel.c_str());
      buttons[i].fsStyle = TBSTYLE_AUTOSIZE | TBSTYLE_DROPDOWN | BTNS_SHOWTEXT;
      buttons[i].dwData = i;
      buttons[i].iString = (INT_PTR)label;
      buttons[i].idCommand = i;
    } else {
      buttons[i].dwData = 0xffff;
      buttons[i].idCommand = 0xffff;
      if (item->separator_) {
        buttons[i].fsStyle = BTNS_SEP;
        buttons[i].iBitmap = MENUBAR_SEPARATOR_WIDTH;
      }
    }
  }

  if (!MenuBarSendMsg(m_hWnd, TB_ADDBUTTONS, n,
                      reinterpret_cast<LONG_PTR>(buttons))) {
    LOG_ERROR("[%d] CMenuMgr::SetMenu: TBL_ADDBUTTONS failed: 0x%x", __LINE__,
              ::GetLastError());
  }
  _freea(buffer);
  m_pMenu = pMenu;
  return 0;
}

void CMenuMgr::UpdateMenu(const std::string &menuId) {
  // TODO: not implement yet
}

void CMenuMgr::OnResize() { ::SendMessageW(m_hwndRebar, WM_SIZE, 0, 0); }

void CMenuMgr::ResetHotItem() {
  int item;
  POINT pt;

  ::GetCursorPos(&pt);
  ::MapWindowPoints(NULL, m_hWnd, &pt, 1);
  item = MenuBarSendMsg(m_hWnd, TB_HITTEST, 0, reinterpret_cast<LONG_PTR>(&pt));
  MenuBarSendMsg(m_hWnd, TB_SETHOTITEM, item, 0);
}

void CMenuMgr::PerformDropDown() {
  TPMPARAMS pmparams = {0};
  pmparams.cbSize = sizeof(TPMPARAMS);

  EnableHotTrace();
  ::SetFocus(m_hWnd);

  m_bContinueHotTrace = true;
  int iRetCmd = 0;
  while (m_bContinueHotTrace) {
    int item = m_iPressedItem;

    if (m_bSelectFromKeyboard) {
      ::keybd_event(VK_DOWN, 0, 0, 0);
      ::keybd_event(VK_DOWN, 0, KEYEVENTF_KEYUP, 0);
    }

    m_bSelectFromKeyboard = false;
    m_bContinueHotTrace = false;

    MenuBarSendMsg(m_hWnd, TB_SETHOTITEM, item, 0);
    DWORD btn_state = MenuBarSendMsg(m_hWnd, TB_GETSTATE, item, 0);
    MenuBarSendMsg(m_hWnd, TB_SETSTATE, item,
                   MAKELONG(btn_state | TBSTATE_PRESSED, 0));

    MenuBarSendMsg(m_hWnd, TB_GETITEMRECT, item,
                   reinterpret_cast<LONG_PTR>(&pmparams.rcExclude));

    if (true /*mc_win_version >= MC_WIN_VISTA && mcIsAppThemed()*/) {
      /* Fix for consistency with a native menu on newer Windows
       * when styles are enabled. */
      pmparams.rcExclude.bottom--;
    }

    ::MapWindowPoints(m_hWnd, HWND_DESKTOP, (POINT *)&pmparams.rcExclude, 2);

    UINT pmflags = TPM_LEFTBUTTON | TPM_VERTICAL | TPM_RETURNCMD;
    if (m_bRTL) {
      pmflags |= TPM_LAYOUTRTL;
    }

    HMENU hMenu = m_pMenu->GetHMenuAtIndex(item);
    iRetCmd = ::TrackPopupMenuEx(
        hMenu, pmflags,
        (m_bRTL ? pmparams.rcExclude.right : pmparams.rcExclude.left),
        pmparams.rcExclude.bottom, m_hWnd, &pmparams);
    DestroyMenu(hMenu);

    MenuBarSendMsg(m_hWnd, TB_SETSTATE, item, MAKELONG(btn_state, 0));
  }

  m_lDropDownCloseTime = static_cast<LONG>(::GetTickCount());

  std::shared_ptr<exciton::menu::MenuItem> pSelMenuItem;

  if (iRetCmd > 0) {
    auto pMenuItem = m_pMenu->GetMenuItem(m_iPressedItem);
    auto pMenu = pMenuItem->subMenu_;
    auto item = pMenu->FindMenuItemFromId(iRetCmd);
    if (item) {
      OnMenuCommand(iRetCmd, m_container.GetActiveBrowser(), m_pMenu, item);
    }
  }

  ResetHotItem();
  DisableHotTrace();
  ::SetFocus(m_hWndOldFocus);
}

void CMenuMgr::DropDown(int item, bool from_keyboard) {
  LONG period_since_last_dropdownclose;

  period_since_last_dropdownclose = ::GetMessageTime() - m_lDropDownCloseTime;
  if (-200 <= period_since_last_dropdownclose &&
      period_since_last_dropdownclose <= 0) {
    LOG_ERROR("[%d] CMenuMgr::DropDown: Ignoring a click which was responsible "
              "for the end of the most recent dropdown menu.",
              __LINE__);
    return;
  }

  m_iPressedItem = item;
  m_bSelectFromKeyboard = from_keyboard;

  ::PostMessageW(m_hWnd, TB_CUSTOMIZE, 0, 0);
}

LRESULT CMenuMgr::OnNotify(NMHDR *hdr) {
  switch (hdr->code) {
  case TBN_DROPDOWN: {
    NMTOOLBAR *info = (NMTOOLBAR *)hdr;
    DropDown(info->iItem, false);
    return TBDDRET_DEFAULT;
  }
  case TBN_HOTITEMCHANGE: {
    NMTBHOTITEM *info = (NMTBHOTITEM *)hdr;
    m_iHotItem = (info->dwFlags & HICF_LEAVING) ? -1 : info->idNew;
    return 0;
  }
  case NM_CUSTOMDRAW: {
    NMTBCUSTOMDRAW *info = (NMTBCUSTOMDRAW *)hdr;
    switch (info->nmcd.dwDrawStage) {
    case CDDS_PREPAINT:
      return (IsAppThemed()) ? CDRF_DODEFAULT : CDRF_NOTIFYITEMDRAW;
    case CDDS_ITEMPREPAINT: {
      int text_color_id;
      int bk_color_id;
      WCHAR buffer[MENUBAR_ITEM_LABEL_MAXSIZE];
      TBBUTTONINFOW btn;
      UINT flags = MENUBAR_DTFLAGS;
      HDC dc = info->nmcd.hdc;
      HFONT old_font;

      if (info->nmcd.uItemState & (CDIS_HOT | CDIS_SELECTED)) {
        text_color_id = COLOR_HIGHLIGHTTEXT;
        bk_color_id = COLOR_HIGHLIGHT;
      } else {
        text_color_id = COLOR_MENUTEXT;
        bk_color_id = -1;
      }

      btn.cbSize = sizeof(TBBUTTONINFO);
      btn.dwMask = TBIF_TEXT;
      btn.pszText = buffer;
      btn.cchText = sizeof(buffer) / sizeof(buffer[0]);
      MenuBarSendMsg(m_hWnd, TB_GETBUTTONINFO, info->nmcd.dwItemSpec,
                     reinterpret_cast<LONG_PTR>(&btn));

      if (MenuBarSendMsg(m_hWnd, WM_QUERYUISTATE, 0, 0) & UISF_HIDEACCEL) {
        flags |= DT_HIDEPREFIX;
      }

      if (bk_color_id >= 0) {
        ::FillRect(dc, &info->nmcd.rc, (HBRUSH)(INT_PTR)(bk_color_id + 1));
      }

      ::SetTextColor(dc, GetSysColor(text_color_id));
      old_font = reinterpret_cast<HFONT>(::SelectObject(
          dc,
          reinterpret_cast<HFONT>(MenuBarSendMsg(m_hWnd, WM_GETFONT, 0, 0))));
      ::SetBkMode(dc, TRANSPARENT);
      ::DrawText(dc, buffer, -1, &info->nmcd.rc, flags);
      ::SelectObject(dc, old_font);
      return CDRF_SKIPDEFAULT;
    }

    default:
      return CDRF_DODEFAULT;
    }
  }
  }

  return 0;
}

BOOL CMenuMgr::OnKeyDown(int vk, DWORD key_data) {
  /* Swap meaning of VK_LEFT and VK_RIGHT if having right-to-left layout. */
  if (m_bRTL) {
    if (vk == VK_LEFT) {
      vk = VK_RIGHT;
    } else if (vk == VK_RIGHT) {
      vk = VK_LEFT;
    }
  }

  switch (vk) {
  case VK_ESCAPE:
  case VK_F10:
  case VK_MENU:
    ::SetFocus(m_hWndOldFocus);
    ResetHotItem();
    s_pActiveMenuBar = NULL;
    UpdateUIState(false);
    return TRUE;

  case VK_LEFT:
    if (m_iHotItem >= 0) {
      int item = m_iHotItem - 1;
      if (item < 0) {
        item = MenuBarSendMsg(m_hWnd, TB_BUTTONCOUNT, 0, 0) - 1;
      }
      MenuBarSendMsg(m_hWnd, TB_SETHOTITEM, item, 0);
      UpdateUIState(true);
      return TRUE;
    }
    break;

  case VK_RIGHT:
    if (m_iHotItem >= 0) {
      int item = m_iHotItem + 1;
      if (item >= MenuBarSendMsg(m_hWnd, TB_BUTTONCOUNT, 0, 0)) {
        item = 0;
      }
      MenuBarSendMsg(m_hWnd, TB_SETHOTITEM, item, 0);
      UpdateUIState(true);
      return TRUE;
    }
    break;

  case VK_DOWN:
  case VK_UP:
  case VK_RETURN:
    if (m_iHotItem >= 0) {
      DropDown(m_iHotItem, true);
      UpdateUIState(true);
      return TRUE;
    }
    break;
  }

  /* If we have not consume the key, report it to the caller. */
  return FALSE;
}

CMenuMgr *CMenuMgr::OnNcCreate(HWND hWnd, UINT msg, WPARAM wParam,
                               LPARAM lParam) {
  if (FALSE == MenuBarSendMsg(hWnd, msg, wParam, lParam)) {
    LOG_ERROR("CMenuMgr::OnNcCreate: MenuBarSendMsg(WM_NCCREATE) failed");
    return nullptr;
  }

  const CREATESTRUCTW *cs = reinterpret_cast<const CREATESTRUCTW *>(lParam);
  CMenuMgr *mgr = static_cast<CMenuMgr *>(cs->lpCreateParams);
  mgr->m_hWnd = hWnd;
  WCHAR parent_class[16];

  ::GetClassNameW(cs->hwndParent, parent_class, 16);
  mgr->m_hNotifyWnd = (::wcscmp(parent_class, L"ReBarWindow32") == 0)
                          ? ::GetAncestor(cs->hwndParent, GA_PARENT)
                          : cs->hwndParent;
  mgr->m_iHotItem = -1;
  mgr->m_iPressedItem = -1;
  mgr->m_bRTL = (cs->dwExStyle & (WS_EX_LAYOUTRTL | WS_EX_RTLREADING)) != 0;

  ::SetWindowLongPtr(hWnd, s_wndClsExtraOffset,
                     reinterpret_cast<LONG_PTR>(mgr));
  mgr->AddRef();

  return mgr;
}

int CMenuMgr::OnCreate(const CREATESTRUCTW *cs) {
  if (MenuBarSendMsg(
          m_hWnd, WM_CREATE, 0,
          reinterpret_cast<LONG_PTR>(const_cast<CREATESTRUCTW *>(cs)))) {
    return -1;
  }

  MenuBarSendMsg(m_hWnd, TB_SETPARENT, reinterpret_cast<UINT_PTR>(m_hWnd), 0);
  MenuBarSendMsg(m_hWnd, TB_BUTTONSTRUCTSIZE, sizeof(TBBUTTON), 0);
  MenuBarSendMsg(m_hWnd, TB_SETBITMAPSIZE, 0, MAKELONG(0, -2));
  MenuBarSendMsg(m_hWnd, TB_SETPADDING, 0, MAKELONG(10, 6));
  MenuBarSendMsg(m_hWnd, TB_SETDRAWTEXTFLAGS, MENUBAR_DTFLAGS, MENUBAR_DTFLAGS);

  ::SetWindowLongPtr(m_hWnd, GWL_STYLE,
                     cs->style | TBSTYLE_FLAT | TBSTYLE_TRANSPARENT |
                         CCS_NODIVIDER);

  // TODO?
  UpdateUIState(false);

  return 0;
}

void CMenuMgr::OnDestroy() { LOG_DEBUG("CMenuMgr::OnDestroy"); }

void CMenuMgr::OnNcDestroy(CMenuMgr *mgr) {
  LOG_DEBUG("CMenuMgr::OnNcDestroy");
  if (mgr) {
    mgr->Release();
  }
}

LRESULT CALLBACK CMenuMgr::MenuBarWndProc(HWND win, UINT msg, WPARAM wp,
                                          LPARAM lp) {
  CMenuMgr *pMgr = reinterpret_cast<CMenuMgr *>(
      ::GetWindowLongPtr(win, s_wndClsExtraOffset));

  switch (msg) {
  case TB_SETPARENT:
  // fall through
  case CCM_SETNOTIFYWINDOW: {
    HWND old = pMgr->m_hNotifyWnd;
    pMgr->m_hNotifyWnd = (wp ? (HWND)wp : GetAncestor(win, GA_PARENT));
    return (LRESULT)old;
  }

  case WM_COMMAND:
    LOG_DEBUG("[%d] CMenuMgr::MenuBarWndProc(WM_COMMAND): code=%d; wid=%d; "
              "control=%p",
              __LINE__, (int)HIWORD(wp), (int)LOWORD(wp), (HWND)lp);
    return 0;

  case WM_NOTIFY: {
    NMHDR *hdr = (NMHDR *)lp;
    if (hdr->hwndFrom == win)
      return pMgr->OnNotify(hdr);
    break;
  }

  case WM_ENTERMENULOOP:
  case WM_EXITMENULOOP:
  case WM_CONTEXTMENU:
  case WM_INITMENU:
  case WM_INITMENUPOPUP:
  case WM_UNINITMENUPOPUP:
  case WM_MENUSELECT:
  case WM_MENUCHAR:
  case WM_MENURBUTTONUP:
  case WM_MENUCOMMAND:
  case WM_MENUDRAG:
  case WM_MENUGETOBJECT:
  case WM_MEASUREITEM:
  case WM_DRAWITEM:
    return ::SendMessageW(pMgr->m_hNotifyWnd, msg, wp, lp);

  case WM_KEYDOWN:
  case WM_SYSKEYDOWN:
    if (pMgr->OnKeyDown(wp, lp)) {
      return 0;
    }
    break;

  case WM_GETDLGCODE:
    return (MenuBarSendMsg(win, msg, wp, lp) | DLGC_WANTALLKEYS |
            DLGC_WANTARROWS);

  case WM_SETFOCUS:
    if (win != reinterpret_cast<HWND>(wp)) {
      pMgr->m_hWndOldFocus = reinterpret_cast<HWND>(wp);
    }
    s_pActiveMenuBar = pMgr;
    break;

  case WM_KILLFOCUS:
    pMgr->m_hWndOldFocus = NULL;
    MenuBarSendMsg(pMgr->m_hWnd, TB_SETHOTITEM, -1, 0);
    pMgr->UpdateUIState(false);
    s_pActiveMenuBar = nullptr;
    break;

  case WM_STYLECHANGED:
    if (wp == GWL_EXSTYLE) {
      STYLESTRUCT *ss = reinterpret_cast<STYLESTRUCT *>(lp);
      pMgr->m_bRTL = (ss->styleNew & (WS_EX_LAYOUTRTL | WS_EX_RTLREADING)) != 0;
      ::InvalidateRect(pMgr->m_hWnd, nullptr, TRUE);
    }
    break;

  case WM_NCCREATE:
    pMgr = OnNcCreate(win, msg, wp, lp);
    if (!pMgr) {
      return FALSE;
    }
    return TRUE;

  case WM_CREATE:
    return pMgr->OnCreate(reinterpret_cast<const CREATESTRUCTW *>(lp));

  case WM_DESTROY:
    pMgr->OnDestroy();
    break;

  case WM_NCDESTROY:
    OnNcDestroy(pMgr);
    break;

  /* Disable those standard toolbar messages, which modify contents of
   * the toolbar, as it is our internal responsibility to set it
   * according to the menu. */
  case TB_ADDBITMAP:
  case TB_ADDSTRING:
    LOG_DEBUG("[%d] CMenuMgr::MenuBarWndProc: Suppressing message TB_xxxx (%d)",
              __LINE__, msg);
    ::SetLastError(ERROR_CALL_NOT_IMPLEMENTED);
    return -1;
  case TB_ADDBUTTONS:
  case TB_BUTTONSTRUCTSIZE:
  case TB_CHANGEBITMAP:
  case TB_DELETEBUTTON:
  case TB_ENABLEBUTTON:
  case TB_HIDEBUTTON:
  case TB_INDETERMINATE:
  case TB_INSERTBUTTON:
  case TB_LOADIMAGES:
  case TB_MARKBUTTON:
  case TB_MOVEBUTTON:
  case TB_PRESSBUTTON:
  case TB_REPLACEBITMAP:
  case TB_SAVERESTORE:
  case TB_SETANCHORHIGHLIGHT:
  case TB_SETBITMAPSIZE:
  case TB_SETBOUNDINGSIZE:
  case TB_SETCMDID:
  case TB_SETDISABLEDIMAGELIST:
  case TB_SETHOTIMAGELIST:
  case TB_SETIMAGELIST:
  case TB_SETINSERTMARK:
  case TB_SETPRESSEDIMAGELIST:
  case TB_SETSTATE:
    LOG_DEBUG("[%d] CMenuMgr::MenuBarWndProc: Suppressing message TB_xxxx (%d)",
              __LINE__, msg);
    ::SetLastError(ERROR_CALL_NOT_IMPLEMENTED);
    return 0; /* FALSE == NULL == 0 */

  case TB_CUSTOMIZE:
    /* Show the popup menu */
    pMgr->PerformDropDown();
    return 0;
  }

  return MenuBarSendMsg(win, msg, wp, lp);
}

void CMenuMgr::HotTraceChangeDropDown(int item, bool from_keyboard) {
  m_iPressedItem = item;
  m_bSelectFromKeyboard = from_keyboard;
  m_bContinueHotTrace = true;
  MenuBarSendMsg(m_hWnd, WM_CANCELMODE, 0, 0);
}

LRESULT CALLBACK CMenuMgr::MenuBarHotTraceProc(int code, WPARAM wp, LPARAM lp) {
  if (code >= 0) {
    MSG *msg = (MSG *)lp;
    CMenuMgr *pMgr = s_htMenuMgr;

    switch (msg->message) {
    case WM_MENUSELECT: {
      HMENU hMenu = reinterpret_cast<HMENU>(msg->lParam);
      if (!hMenu) {
        s_htSelMenu.reset();
      } else {
        MENUINFO mi;
        mi.cbSize = sizeof(mi);
        mi.fMask = MIM_MENUDATA;
        ::GetMenuInfo(hMenu, &mi);
        auto pMenu = reinterpret_cast<exciton::menu::Menu *>(mi.dwMenuData);
        s_htSelMenu = pMenu->shared_from_this();
      }
      s_htSelItem = LOWORD(msg->wParam);
      s_htSelFlags = HIWORD(msg->wParam);
      break;
    }

    case WM_MOUSEMOVE: {
      POINT pt = msg->pt;

      ::MapWindowPoints(NULL, pMgr->m_hWnd, &pt, 1);
      int item = MenuBarSendMsg(pMgr->m_hWnd, TB_HITTEST, 0,
                                reinterpret_cast<LONG_PTR>(&pt));
      if (s_htLastPos.x != pt.x || s_htLastPos.y != pt.y) {
        s_htLastPos = pt;
        if ((item != pMgr->m_iPressedItem) && (0 <= item) &&
            item < MenuBarSendMsg(pMgr->m_hWnd, TB_BUTTONCOUNT, 0, 0)) {
          pMgr->HotTraceChangeDropDown(item, false);
        }
      }
      break;
    }

    case WM_KEYDOWN:
    case WM_SYSKEYDOWN: {
      int vk = msg->wParam;

      /* Swap meaning of VK_LEFT and VK_RIGHT if having right-to-left layout. */
      if (pMgr->m_bRTL) {
        if (vk == VK_LEFT)
          vk = VK_RIGHT;
        else if (vk == VK_RIGHT)
          vk = VK_LEFT;
      }

      switch (vk) {
      case VK_MENU:
      case VK_F10:
        pMgr->HotTraceChangeDropDown(-1, true);
        return 1; /* Consume the message. */

      case VK_LEFT:
        if (!s_htSelMenu ||
            s_htSelMenu == pMgr->m_pMenu->GetSubMenu(pMgr->m_iPressedItem)) {
          int item = pMgr->m_iPressedItem - 1;
          if (item < 0) {
            item = MenuBarSendMsg(pMgr->m_hWnd, TB_BUTTONCOUNT, 0, 0) - 1;
          }
          if (item != pMgr->m_iPressedItem) {
            pMgr->HotTraceChangeDropDown(item, true);
          }
          return 1; /* Consume the message. */
        }
        break;

      case VK_RIGHT:
        if (!s_htSelMenu || !(s_htSelFlags & MF_POPUP) ||
            (s_htSelFlags & (MF_GRAYED | MF_DISABLED))) {
          int item = pMgr->m_iPressedItem + 1;
          if (item >= MenuBarSendMsg(pMgr->m_hWnd, TB_BUTTONCOUNT, 0, 0)) {
            item = 0;
          }
          if (item != pMgr->m_iPressedItem) {
            pMgr->HotTraceChangeDropDown(item, true);
          }
          return 1; /* Consume the message. */
        }
        break;
      }
      break;
    }
    }
  }

  return CallNextHookEx(s_htHook, code, wp, lp);
}

void CMenuMgr::EnableHotTrace() {
  ::EnterCriticalSection(&s_htCS);

  if (s_htMenuMgr) {
    LOG_DEBUG("[%d] CMenuMgr::EnableHotTrace: Another menubar hot tracking???",
              __LINE__);
    PerformDisableHotTrace();
  }

  s_htHook = ::SetWindowsHookExW(WH_MSGFILTER, MenuBarHotTraceProc,
                                 Driver::Current().InstanceHandle(),
                                 GetCurrentThreadId());
  if (!s_htHook) {
    LOG_ERROR("[%d] CMenuMgr::EnableHotTrace: SetWindowsHookEx() failed",
              __LINE__);
    goto err_hook;
  }

  s_htMenuMgr = this;
  ::GetCursorPos(&s_htLastPos);
  ::MapWindowPoints(NULL, m_hWnd, &s_htLastPos, 1);

err_hook:
  ::LeaveCriticalSection(&s_htCS);
}

void CMenuMgr::DisableHotTrace() {
  ::EnterCriticalSection(&s_htCS);

  if (s_htMenuMgr != this) {
    LOG_DEBUG("[%d] CMenuMgr::DisableHotTrace: Another menubar hot tracking???",
              __LINE__);
  } else {
    PerformDisableHotTrace();
  }

  ::LeaveCriticalSection(&s_htCS);
}

void CMenuMgr::PerformDisableHotTrace(void) {
  if (s_htHook) {
    ::UnhookWindowsHookEx(s_htHook);
    s_htHook = NULL;
    s_htMenuMgr = nullptr;
    s_htSelMenu = NULL;
    s_htSelItem = -1;
    s_htSelFlags = 0;
  }
}
