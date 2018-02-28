#include <windows.h>

#include <shlobj.h>

#include "browsercontainer.h"
#include "browserhost.h"

//#include "sidebar.h"
#include "menumgr.h"
//#include "rebarmgr.h"
//#include "statusbar.h"
//#include "opensearch.h"
#include "global.h"
#include "log.h"
#include "util.h"


CWebBrowserContainer *CWebBrowserContainer::s_pActiveContainer = nullptr;
BOOL CWebBrowserContainer::s_bClassInitialized = FALSE;

LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);

LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam,
                            LPARAM lParam) {
  if (uMsg == WM_NCCREATE) {
    CREATESTRUCT *cs = reinterpret_cast<CREATESTRUCT *>(lParam);
    ::SetWindowLongPtr(hwnd, GWLP_USERDATA, (LONG_PTR)(cs->lpCreateParams));
    return TRUE;
  }
  CWebBrowserContainer *c =
      (CWebBrowserContainer *)::GetWindowLongPtr(hwnd, GWLP_USERDATA);
  return c->WindowProc(hwnd, uMsg, wParam, lParam);
}

LRESULT CALLBACK HostProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
  CWebBrowserHost *p = (CWebBrowserHost *)GetWindowLongPtr(hwnd, GWLP_USERDATA);

  if (p != NULL)
    return p->WindowProc(hwnd, uMsg, wParam, lParam);
  else
    return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

// CWebBrowserContainer

CWebBrowserContainer::CWebBrowserContainer() {
  //  m_cRef = 1;
  m_hwnd = NULL;
  m_pActiveWebBrowserHost = NULL;
  m_pMenuMgr = NULL;
}

CWebBrowserContainer::~CWebBrowserContainer() {}

BOOL CWebBrowserContainer::NewWindow(HINSTANCE hinst, int width, int height) {
  HWND hwnd;
  //  MSG msg;
  if (!s_bClassInitialized) {
    WNDCLASSEX wc;

    wc.cbSize = sizeof(WNDCLASSEX);
    wc.style = 0;
    wc.lpfnWndProc = ::WindowProc;
    wc.cbClsExtra = 0;
    wc.cbWndExtra = 0;
    wc.hInstance = hinst;
    wc.hIcon =
        (HICON)LoadImage(NULL, IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_SHARED);
    wc.hCursor =
        (HCURSOR)LoadImage(NULL, IDC_ARROW, IMAGE_CURSOR, 0, 0, LR_SHARED);
    wc.hbrBackground = (HBRUSH)GetStockObject(GRAY_BRUSH);
    wc.lpszMenuName = NULL;
    wc.lpszClassName = g_szContainerClassName;
    wc.hIconSm =
        (HICON)LoadImage(NULL, IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_SHARED);

    if (RegisterClassEx(&wc) == 0) {
      return FALSE;
    }

    wc.hbrBackground = (HBRUSH)GetStockObject(WHITE_BRUSH);
    wc.lpfnWndProc = HostProc;
    wc.lpszClassName = g_szHostClassName;
    if (RegisterClassEx(&wc) == 0) {
      return 0;
    }
    s_bClassInitialized = TRUE;
  }

  RECT r = {0, 0, width, height};
  AdjustWindowRect(&r, WS_OVERLAPPEDWINDOW | WS_CLIPCHILDREN, TRUE);

  hwnd = CreateWindowEx(0, g_szContainerClassName, g_szWinodowName,
                        WS_OVERLAPPEDWINDOW | WS_CLIPCHILDREN, CW_USEDEFAULT,
                        CW_USEDEFAULT, r.right - r.left, r.bottom - r.top, NULL, NULL,
                        hinst, this);
  if (hwnd == NULL) {
    return FALSE;
  }

  ShowWindow(hwnd, SW_SHOW);
  UpdateWindow(hwnd);
  return TRUE;
}

BOOL CWebBrowserContainer::Create(HWND hwnd) {
  INITCOMMONCONTROLSEX ic;

  ic.dwSize = sizeof(INITCOMMONCONTROLSEX);
  ic.dwICC = ICC_BAR_CLASSES | ICC_COOL_CLASSES | ICC_TAB_CLASSES |
             ICC_TREEVIEW_CLASSES;
  ::InitCommonControlsEx(&ic);

  m_hwnd = hwnd;

  m_pMenuMgr = new CMenuMgr(*this);
  m_pMenuMgr->Create(hwnd);

  ::CoInternetSetFeatureEnabled(FEATURE_TABBED_BROWSING, SET_FEATURE_ON_PROCESS,
                                TRUE);
  ::CoInternetSetFeatureEnabled(FEATURE_FEEDS, SET_FEATURE_ON_PROCESS, TRUE);
  ::CoInternetSetFeatureEnabled(FEATURE_SECURITYBAND, SET_FEATURE_ON_PROCESS,
                                TRUE);
  ::CoInternetSetFeatureEnabled(FEATURE_WEBOC_POPUPMANAGEMENT,
                                SET_FEATURE_ON_PROCESS, TRUE);
  ::CoInternetSetFeatureEnabled(FEATURE_RESTRICT_FILEDOWNLOAD,
                                SET_FEATURE_ON_PROCESS, TRUE);

  SetUserAgent();
  ResizeWindow();

  return TRUE;
}

CWebBrowserHost *
CWebBrowserContainer::NewBrowser(const std::string &initialHtml) {
  m_pActiveWebBrowserHost = new CWebBrowserHost(shared_from_this(), initialHtml);
  if (!m_pActiveWebBrowserHost->Create(m_hwnd, L"about:blank", 30000)) {
    return FALSE;
  }

  if (m_pActiveWebBrowserHost != NULL) {
    BOOL bEnableForward, bEnableBack;

    m_pActiveWebBrowserHost->Show(TRUE);

    m_pActiveWebBrowserHost->GetTravelState(&bEnableForward, &bEnableBack);
    ResizeWindow();
  }
  return m_pActiveWebBrowserHost;
}

void CWebBrowserContainer::Destroy() {
  m_pActiveWebBrowserHost = NULL;

  if (m_pMenuMgr != NULL) {
    m_pMenuMgr->Destroy();
    m_pMenuMgr->Release();
  }

  // TODO: emit destroy message?
}

LRESULT CWebBrowserContainer::WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam,
                                         LPARAM lParam) {
  switch (uMsg) {

  case WM_CREATE:
    LOG_DEBUG("[%d] CWebBrowserContainer::WindowProc: WM_CREATE", __LINE__);
    Create(hwnd);
    return 0;

  case WM_NOTIFY:
    return 0;

  case WM_INITMENUPOPUP:
    m_pMenuMgr->SetMenuState((HMENU)wParam, LOWORD(lParam));
    return 0;

  case WM_SIZE:
    if (m_pMenuMgr == NULL)
      return 0;
    ResizeWindow();
    return 0;

  case WM_DESTROY:
    Destroy();
    PostQuitMessage(0); // TODO: fixme post to go /window/:id/closed
    return 0;
  case WM_ACTIVATE:
    if (LOWORD(wParam) != WA_INACTIVE) {
      s_pActiveContainer = this;
    }
    break;
  default:
    break;
  }

  return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

HWND CWebBrowserContainer::GetWindow() { return m_hwnd; }

void CWebBrowserContainer::ResizeWindow() {
  int nWidthClient, nHeightClient;
  RECT rc;

  if (!m_pMenuMgr)
    return;

  m_pMenuMgr->OnResize();

  ::GetClientRect(m_hwnd, &rc);
  nWidthClient = rc.right - rc.left;
  nHeightClient = rc.bottom - rc.top;
  ::GetWindowRect(m_pMenuMgr->GetWindow(), &rc);
  int nHeightMenu = rc.bottom - rc.top;

  if (m_pActiveWebBrowserHost != NULL)
    m_pActiveWebBrowserHost->SetWindowSize(0, nHeightMenu, nWidthClient, nHeightClient - nHeightMenu);
}

void CWebBrowserContainer::SetUserAgent() {
  CHAR szUserAgent[512];
  DWORD dwLength;

  UrlMkGetSessionOption(URLMON_OPTION_USERAGENT, szUserAgent,
                        sizeof(szUserAgent), &dwLength, 0);

  szUserAgent[dwLength - 2] = '\0';
  lstrcatA(szUserAgent, "; Exciton WebBrowser/1.0)");
  UrlMkSetSessionOption(URLMON_OPTION_USERAGENT, (char *)szUserAgent,
                        lstrlenA(szUserAgent), 0);

  MultiByteToWideChar(CP_ACP, 0, szUserAgent, -1, g_szUserAgentW, 512);
}

BOOL CWebBrowserContainer::TranslateAccelerator(LPMSG lpMsg) {
  if (m_pActiveWebBrowserHost == NULL)
    return FALSE;

  return m_pActiveWebBrowserHost->TranslateAccelerator(lpMsg) == S_OK;
}

CWebBrowserHost *CWebBrowserContainer::GetActiveBrowser() {
  return m_pActiveWebBrowserHost;
}

void CWebBrowserContainer::SetActiveBrowser(CWebBrowserHost *pWebBrowserHost) {
  m_pActiveWebBrowserHost = pWebBrowserHost;
}

CWebBrowserContainer *CWebBrowserContainer::GetActiveContainer() {
  if (s_pActiveContainer && ::IsWindow(s_pActiveContainer->m_hwnd)) {
    return s_pActiveContainer;
  }
  return nullptr;
}

void CWebBrowserContainer::Navigate(LPWSTR lpszUrl, BOOL bNewTab) {
  if (m_pActiveWebBrowserHost != NULL) {
    m_pActiveWebBrowserHost->Navigate(lpszUrl);
  }
}

void CWebBrowserContainer::OnCommandStateChange(long Command,
                                                VARIANT_BOOL Enable) {
  //TODO: need implement??
}

void CWebBrowserContainer::OnTitleChange(BSTR Text) {
  std::wstring str = exciton::util::ToUTF16String(Text);
  ::SendMessageW(GetWindow(), WM_SETTEXT, 0,
                 reinterpret_cast<LPARAM>(str.c_str()));
}

void CWebBrowserContainer::OnNavigateComplete2(BSTR URL) {
  //TODO: need implement?
}

void CWebBrowserContainer::UpdateMenu(const std::string &menuId) {
  if (m_pMenuMgr) {
    m_pMenuMgr->UpdateMenu(menuId);
  }
}
