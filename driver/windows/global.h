extern BOOL RegisterWindowClass(WNDPROC lpfnWndProc, LPTSTR lpszClassName);
//extern void NavigateFromShortcut(IShellItem *psi, BOOL bNewTab);
extern void SetMenuItem(HMENU hmenu, int nId, BOOL bEnable, BOOL bCheck);
extern void InitializeMenuItem(HMENU hmenu, LPCWSTR lpszItemName, int nId, HMENU hmenuSub);
//extern void ShowPrivacyDlg(HWND hwndParent, IWebBrowser2 *pWebBrowser2);

const TCHAR g_szWinodowName[] = TEXT("sample");
const TCHAR g_szContainerClassName[] = TEXT("container-class");
const TCHAR g_szHostClassName[] = TEXT("host-class");
const TCHAR g_szSidebarClassName[] = TEXT("sidebar-class");
const TCHAR g_szSuggestClassName[] = TEXT("suggest-class");

extern WCHAR g_szUserAgentW[512];

#define ID_REBAR 600
#define ID_ADDRESS_TOOLBAR 601
#define ID_ADDRESS_COMBOBOX 602
#define ID_TAB_TOOLBAR 603

#define ID_STATUS 700

#define ID_SEARCH_TOOLBAR 800
#define ID_SEARCH_EDIT 801
#define ID_SUGGEST 802

#define ID_SIDEBAR 900
#define ID_SIDEBAR_TAB  901

#define ID_FEED 1000

#define ID_HISTORY 1100