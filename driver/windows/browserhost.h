#ifndef __INC_WINDOWS_BROWSERHOST_H__
#define __INC_WINDOWS_BROWSERHOST_H__
#include <windows.h>

#include <exdisp.h>
#include <mshtmhst.h>
#include <memory>
#include <string>

// class CBandSite;
class CEventSink;
class CWebBrowserContainer;

class CWebBrowserHost : public IOleClientSite,
                        public IOleInPlaceSite,
                        public IDispatch,
                        public IOleCommandTarget,
                        public IDocHostUIHandler {
public:
  STDMETHODIMP QueryInterface(REFIID riid, void **ppvObject);
  STDMETHODIMP_(ULONG) AddRef();
  STDMETHODIMP_(ULONG) Release();

  STDMETHODIMP SaveObject();
  STDMETHODIMP GetMoniker(DWORD dwAssign, DWORD dwWhichMoniker,
                          IMoniker **ppmk);
  STDMETHODIMP GetContainer(IOleContainer **ppContainer);
  STDMETHODIMP ShowObject();
  STDMETHODIMP OnShowWindow(BOOL fShow);
  STDMETHODIMP RequestNewObjectLayout();

  STDMETHODIMP GetWindow(HWND *phwnd);
  STDMETHODIMP ContextSensitiveHelp(BOOL fEnterMode);
  STDMETHODIMP CanInPlaceActivate();
  STDMETHODIMP OnInPlaceActivate();
  STDMETHODIMP OnUIActivate();
  STDMETHODIMP GetWindowContext(IOleInPlaceFrame **ppFrame,
                                IOleInPlaceUIWindow **ppDoc, LPRECT lprcPosRect,
                                LPRECT lprcClipRect,
                                LPOLEINPLACEFRAMEINFO lpFrameInfo);
  STDMETHODIMP Scroll(SIZE scrollExtant);
  STDMETHODIMP OnUIDeactivate(BOOL fUndoable);
  STDMETHODIMP OnInPlaceDeactivate();
  STDMETHODIMP DiscardUndoState();
  STDMETHODIMP DeactivateAndUndo();
  STDMETHODIMP OnPosRectChange(LPCRECT lprcPosRect);

  STDMETHODIMP GetTypeInfoCount(UINT *pctinfo);
  STDMETHODIMP GetTypeInfo(UINT iTInfo, LCID lcid, ITypeInfo **ppTInfo);
  STDMETHODIMP GetIDsOfNames(REFIID riid, LPOLESTR *rgszNames, UINT cNames,
                             LCID lcid, DISPID *rgDispId);
  STDMETHODIMP Invoke(DISPID dispIdMember, REFIID riid, LCID lcid, WORD wFlags,
                      DISPPARAMS *pDispParams, VARIANT *pVarResult,
                      EXCEPINFO *pExcepInfo, UINT *puArgErr);

  STDMETHODIMP QueryStatus(const GUID *pguidCmdGroup, ULONG cCmds,
                           OLECMD prgCmds[], OLECMDTEXT *pCmdText);
  STDMETHODIMP Exec(const GUID *pguidCmdGroup, DWORD nCmdID, DWORD nCmdexecopt,
                    VARIANT *pvaIn, VARIANT *pvaOut);

  STDMETHODIMP ShowContextMenu(DWORD dwID, POINT *ppt, IUnknown *pcmdtReserved,
                               IDispatch *pdispReserved);
  STDMETHODIMP GetHostInfo(DOCHOSTUIINFO *pInfo);
  STDMETHODIMP ShowUI(DWORD dwID, IOleInPlaceActiveObject *pActiveObject,
                      IOleCommandTarget *pCommandTarget,
                      IOleInPlaceFrame *pFrame, IOleInPlaceUIWindow *pDoc);
  STDMETHODIMP HideUI(VOID);
  STDMETHODIMP UpdateUI(VOID);
  STDMETHODIMP EnableModeless(BOOL fEnable);
  STDMETHODIMP OnDocWindowActivate(BOOL fActivate);
  STDMETHODIMP OnFrameWindowActivate(BOOL fActivate);
  STDMETHODIMP ResizeBorder(LPCRECT prcBorder, IOleInPlaceUIWindow *pUIWindow,
                            BOOL fFrameWindow);
  STDMETHODIMP TranslateAccelerator(LPMSG lpMsg, const GUID *pguidCmdGroup,
                                    DWORD nCmdID);
  STDMETHODIMP GetOptionKeyPath(LPOLESTR *pchKey, DWORD dw);
  STDMETHODIMP GetDropTarget(IDropTarget *pDropTarget,
                             IDropTarget **ppDropTarget);
  STDMETHODIMP GetExternal(IDispatch **ppDispatch);
  STDMETHODIMP TranslateUrl(DWORD dwTranslate, OLECHAR *pchURLIn,
                            OLECHAR **ppchURLOut);
  STDMETHODIMP FilterDataObject(IDataObject *pDO, IDataObject **ppDORet);

public:
  CWebBrowserHost(std::shared_ptr<CWebBrowserContainer> pContainer);
  CWebBrowserHost(const CWebBrowserHost &) = delete;
  CWebBrowserHost(CWebBrowserHost &&) = delete;
  ~CWebBrowserHost();

public:
  BOOL Create(HWND hwndParent /*, HWND hwndRebar*/, LPCWSTR lpszUrl, int nId);
  BOOL Destroy();
  LRESULT WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam);
  void SetWindowSize(int x, int y, int nWidth, int nHeight);
  void Show(BOOL bShow);
  BOOL Navigate(LPCWSTR lpszUrl);
  HRESULT QueryBrowserInterface(REFIID riid, void **ppvObject);
  int GetId();
  void AmbientPropertyChange(int nId);
  BOOL IsContainAmbientProperty(int nId);
  void ForwardOrBack(BOOL bForward);
  void SetTravelState(BOOL bEnableForward, BOOL bEnableBack);
  void GetTravelState(LPBOOL lpbEnableForward, LPBOOL lpbEnableBack);
  // BOOL ShowBandMenu();
  HRESULT TranslateAccelerator(LPMSG lpMsg);
  void Exec(OLECMDID nCmdID, OLECMDEXECOPT nCmdexecopt, VARIANT *pvaIn,
            VARIANT *pvaOut);
  void ExecDocument(const GUID *pguid, DWORD nCmdID, DWORD nCmdexecopt,
                    VARIANT *pvaIn, VARIANT *pvaOut);
  void OnDocumentComplate(IDispatch* lpDisp);
  void EvaluateJavasScript(const std::wstring& funcName, const std::wstring& jsonArg, VARIANT* pRetValue);
  std::shared_ptr<CWebBrowserContainer> GetHostContainer() const { return m_pContainer; }
  void PutFullscreen(bool bEnter);
private:
  std::shared_ptr<CWebBrowserContainer> m_pContainer;
  std::string m_strInitialHtml;
  LONG m_cRef;
  HWND m_hwnd;
  IWebBrowser2 *m_pWebBrowser2;
  int m_nId;
  BOOL m_bEnableForward;
  BOOL m_bEnableBack;
  DWORD m_dwAmbientDLControl;
  CEventSink *m_pEventSink;
  // CBandSite *m_pBandSite;
};
#endif
