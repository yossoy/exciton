#pragma once

#include <string>
#include <memory>

class CWebBrowserHost;
class CMenuMgr;

class CWebBrowserContainer : public std::enable_shared_from_this<CWebBrowserContainer>
{
public:
  CWebBrowserContainer();
  ~CWebBrowserContainer();
  BOOL NewWindow(HINSTANCE hinst /*, int nCmdShow*/);
  CWebBrowserHost* NewBrowser(const std::string& initialHtml);
  LRESULT WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);
  HWND GetWindow();
  CWebBrowserHost *GetActiveBrowser();
  void SetActiveBrowser(CWebBrowserHost *pWebBrowserHost);
  static CWebBrowserContainer *GetActiveContainer();
  // BOOL IsShowSidebar();
  // void ShowSidebar(BOOL bShow);
  void Navigate(LPWSTR lpszUrl, BOOL bNewTab);
  void OnCommandStateChange(long Command, VARIANT_BOOL Enable);
  void OnTitleChange(BSTR Text);
  void OnNavigateComplete2(BSTR URL);
  // void OnNewWindow3(BSTR URL);
  // void OnStatusTextChange(BSTR Text);
  BOOL TranslateAccelerator(LPMSG lpMsg);
  public:
  void SetInitialHTML(const std::string& str) {
	  m_strInitialHTML = str;
  }
  std::string GetInitialHTML() {
	  std::string ret;
	  if (!m_strInitialHTML.empty()) {
		  std::swap(ret, m_strInitialHTML);
	  }
	  return ret;
  }
  void UpdateMenu(const std::string& menuId);

private:
  void SetUserAgent();
  BOOL Create(HWND hwnd);
  void Destroy();
  void ResizeWindow();

private:
  static BOOL s_bClassInitialized;
  static CWebBrowserContainer *s_pActiveContainer;
  HWND m_hwnd;
  CWebBrowserHost *m_pActiveWebBrowserHost;
  CMenuMgr *m_pMenuMgr;
  std::string m_strInitialHTML;
};
