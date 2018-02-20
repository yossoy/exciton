#include <windows.h>
#include <shlobj.h>
#include <intshcut.h>
#include <mshtmlc.h>
#include "global.h"
#include "browsercontainer.h"

WCHAR g_szUserAgentW[512];

BOOL RegisterWindowClass(WNDPROC lpfnWndProc, LPTSTR lpszClassName)
{
	WNDCLASSEX wc;

	wc.cbSize        = sizeof(WNDCLASSEX);
	wc.style         = 0;
	wc.lpfnWndProc   = lpfnWndProc;
	wc.cbClsExtra    = 0;
	wc.cbWndExtra    = 0;
	wc.hInstance     = GetModuleHandle(NULL);
	wc.hIcon         = (HICON)LoadImage(NULL, IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_SHARED);
	wc.hCursor       = (HCURSOR)LoadImage(NULL, IDC_ARROW, IMAGE_CURSOR, 0, 0, LR_SHARED);
	wc.hbrBackground = (HBRUSH)GetStockObject(WHITE_BRUSH);
	wc.lpszMenuName  = NULL;
	wc.lpszClassName = lpszClassName;
	wc.hIconSm       = (HICON)LoadImage(NULL, IDI_APPLICATION, IMAGE_ICON, 0, 0, LR_SHARED);
	
	if (RegisterClassEx(&wc) == 0)
		return FALSE;

	return TRUE;
}

// void NavigateFromShortcut(IShellItem *psi, BOOL bNewTab)
// {
// 	LPWSTR                   lpszDisplayName, lpszUrl;
// 	IUniformResourceLocatorW *pUniformResourceLocator;
// 	IPersistFile             *pPersistFile;
// 	HRESULT                  hr;
	
// 	hr = CoCreateInstance(CLSID_InternetShortcut, NULL, CLSCTX_INPROC_SERVER, IID_PPV_ARGS(&pUniformResourceLocator));
// 	if (FAILED(hr))
// 		return;
	
// 	pUniformResourceLocator->QueryInterface(IID_PPV_ARGS(&pPersistFile));

// 	psi->GetDisplayName(SIGDN_FILESYSPATH, &lpszDisplayName);
// 	pPersistFile->Load(lpszDisplayName, STGM_READ);
	
// 	pUniformResourceLocator->GetURL(&lpszUrl);
// 	g_pWebBrowserContainer->Navigate(lpszUrl, bNewTab);
	
// 	CoTaskMemFree(lpszUrl);
// 	CoTaskMemFree(lpszDisplayName);
// 	pPersistFile->Release();
// 	pUniformResourceLocator->Release();
// }

void SetMenuItem(HMENU hmenu, int nId, BOOL bEnable, BOOL bCheck)
{
	MENUITEMINFO mii;
	
	mii.cbSize  = sizeof(MENUITEMINFO);
	mii.fMask   = MIIM_ID | MIIM_STATE;
	mii.fState  = bEnable ? MFS_ENABLED : MFS_DISABLED;
	mii.fState |= bCheck ? MFS_CHECKED : MFS_UNCHECKED;
	mii.wID     = nId;
	
	SetMenuItemInfo(hmenu, nId, FALSE, &mii);
}

void InitializeMenuItem(HMENU hmenu, LPCWSTR lpszItemName, int nId, HMENU hmenuSub)
{
	MENUITEMINFO mii;
	
	mii.cbSize = sizeof(MENUITEMINFO);
	mii.fMask  = MIIM_TYPE;

	if (lpszItemName != NULL) {
		mii.fMask     |= MIIM_ID;
		mii.wID        = nId;
		mii.fType      = MFT_STRING;
		mii.dwTypeData = const_cast<LPWSTR>(lpszItemName);
	}
	else
		mii.fType = MFT_SEPARATOR;

	if (hmenuSub != NULL) {
		mii.fMask   |= MIIM_SUBMENU;
		mii.hSubMenu = hmenuSub;
	}

	InsertMenuItem(hmenu, nId, FALSE, &mii);
}

// void ShowPrivacyDlg(HWND hwndParent, IWebBrowser2 *pWebBrowser2)
// {
// 	typedef HRESULT (WINAPI *LPFNDOPRIVACYDLG)(HWND, LPCWSTR, IEnumPrivacyRecords *, BOOL);

// 	HMODULE             hmod;
// 	LPFNDOPRIVACYDLG    lpfnDoPrivacyDlg;
// 	BSTR                URL;
// 	IDispatch           *pDispatch;
// 	IServiceProvider    *pServiceProvider;
// 	IEnumPrivacyRecords *pEnumPrivacyRecords;
// 	BOOL                bPrivacyImpacted;
	
// 	hmod = LoadLibrary(TEXT("shdocvw.dll"));
// 	if (hmod == NULL)
// 		return;

// 	lpfnDoPrivacyDlg = (LPFNDOPRIVACYDLG)GetProcAddress(hmod, "DoPrivacyDlg");
// 	if (lpfnDoPrivacyDlg == NULL) {
// 		FreeLibrary(hmod);
// 		return;
// 	}

// 	pWebBrowser2->get_Document(&pDispatch);
// 	pDispatch->QueryInterface(IID_PPV_ARGS(&pServiceProvider));
// 	pServiceProvider->QueryService(IID_IEnumPrivacyRecords, IID_PPV_ARGS(&pEnumPrivacyRecords));
// 	pServiceProvider->Release();
// 	pDispatch->Release();

// 	pWebBrowser2->get_LocationURL(&URL);
// 	pEnumPrivacyRecords->GetPrivacyImpacted(&bPrivacyImpacted);
// 	lpfnDoPrivacyDlg(hwndParent, URL, pEnumPrivacyRecords, !bPrivacyImpacted);

// 	SysFreeString(URL);
// 	pEnumPrivacyRecords->Release();
// 	FreeLibrary(hmod);
// }