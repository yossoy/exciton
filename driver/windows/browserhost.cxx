#include <windows.h>

#include <cstdio>

#include <mshtmcid.h>
#include <mshtmdid.h>
#include <mshtmhst.h>
#include <mshtml.h>
#include <shlobj.h>
#include <shlwapi.h>

#include "browsercontainer.h"
#include "browserhost.h"
#include "driver.h"
#include "eventsink.h"
#include "global.h"
#include "htmlmoniker.h"
#include "log.h"
#include "util.h"


#define MY_DEFINE_GUID(name, l, w1, w2, b1, b2, b3, b4, b5, b6, b7, b8)        \
  EXTERN_C const GUID DECLSPEC_SELECTANY name = {                              \
      l, w1, w2, {b1, b2, b3, b4, b5, b6, b7, b8}}
MY_DEFINE_GUID(CGID_DocHostCommandHandler, 0xf38bc242, 0xb950, 0x11d1, 0x89,
               0x18, 0x00, 0xc0, 0x4f, 0xc2, 0xc8, 0x36);

#define DISPID_EXTERNAL_GOLANGREQUEST (1)

CWebBrowserHost::CWebBrowserHost(
    std::shared_ptr<CWebBrowserContainer> container,
    const std::string &strInitialHtml, const std::string& id)
    : m_pContainer(container), m_cRef(1), m_hwnd(NULL), m_pWebBrowser2(nullptr),
      m_nId(0), m_bEnableForward(FALSE), m_bEnableBack(FALSE),
      m_dwAmbientDLControl(DLCTL_DLIMAGES | DLCTL_VIDEOS | DLCTL_BGSOUNDS),
      m_pEventSink(nullptr), m_strInitialHtml(strInitialHtml), m_pHtmlMoniker(nullptr), m_strID(id) {}

CWebBrowserHost::~CWebBrowserHost() {
  if (m_pHtmlMoniker) {
    m_pHtmlMoniker->Release();
    m_pHtmlMoniker = nullptr;
  }
}

STDMETHODIMP CWebBrowserHost::QueryInterface(REFIID riid, void **ppvObject) {
  *ppvObject = NULL;

  if (IsEqualIID(riid, IID_IUnknown) || IsEqualIID(riid, IID_IOleClientSite))
    *ppvObject = static_cast<IOleClientSite *>(this);
  else if (IsEqualIID(riid, IID_IOleWindow) ||
           IsEqualIID(riid, IID_IOleInPlaceSite))
    *ppvObject = static_cast<IOleInPlaceSite *>(this);
  else if (IsEqualIID(riid, IID_IDispatch))
    *ppvObject = static_cast<IDispatch *>(this);
  else if (IsEqualIID(riid, IID_IOleCommandTarget))
    *ppvObject = static_cast<IOleCommandTarget *>(this);
  else if (IsEqualIID(riid, IID_IDocHostUIHandler))
    *ppvObject = static_cast<IDocHostUIHandler *>(this);
  else
    return E_NOINTERFACE;

  AddRef();

  return S_OK;
}

STDMETHODIMP_(ULONG) CWebBrowserHost::AddRef() {
  return InterlockedIncrement(&m_cRef);
}

STDMETHODIMP_(ULONG) CWebBrowserHost::Release() {
  return InterlockedDecrement(&m_cRef);
}

STDMETHODIMP CWebBrowserHost::SaveObject() { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::GetMoniker(DWORD dwAssign, DWORD dwWhichMoniker,
                                         IMoniker **ppmk) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::GetContainer(IOleContainer **ppContainer) {
  *ppContainer = NULL;

  return E_NOINTERFACE;
}

STDMETHODIMP CWebBrowserHost::ShowObject() { return S_OK; }

STDMETHODIMP CWebBrowserHost::OnShowWindow(BOOL fShow) { return S_OK; }

STDMETHODIMP CWebBrowserHost::RequestNewObjectLayout() { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::GetWindow(HWND *phwnd) {
  *phwnd = m_hwnd;

  return S_OK;
}

STDMETHODIMP CWebBrowserHost::ContextSensitiveHelp(BOOL fEnterMode) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::CanInPlaceActivate() { return S_OK; }

STDMETHODIMP CWebBrowserHost::OnInPlaceActivate() { return S_OK; }

STDMETHODIMP CWebBrowserHost::OnUIActivate() { return S_OK; }

STDMETHODIMP CWebBrowserHost::GetWindowContext(
    IOleInPlaceFrame **ppFrame, IOleInPlaceUIWindow **ppDoc, LPRECT lprcPosRect,
    LPRECT lprcClipRect, LPOLEINPLACEFRAMEINFO lpFrameInfo) {
  *ppFrame = NULL;
  *ppDoc = NULL;

  GetClientRect(m_hwnd, lprcPosRect);
  GetClientRect(m_hwnd, lprcClipRect);

  return S_OK;
}

STDMETHODIMP CWebBrowserHost::Scroll(SIZE scrollExtant) { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::OnUIDeactivate(BOOL fUndoable) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::OnInPlaceDeactivate() { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::DiscardUndoState() { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::DeactivateAndUndo() { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::OnPosRectChange(LPCRECT lprcPosRect) {
  return S_OK;
}

STDMETHODIMP CWebBrowserHost::GetTypeInfoCount(UINT *pctinfo) {
  *pctinfo = 0;

  return S_OK;
}

STDMETHODIMP CWebBrowserHost::GetTypeInfo(UINT iTInfo, LCID lcid,
                                          ITypeInfo **ppTInfo) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::GetIDsOfNames(REFIID riid, LPOLESTR *rgszNames,
                                            UINT cNames, LCID lcid,
                                            DISPID *rgDispId) {
  UINT i;
  HRESULT hr;
  hr = NOERROR;
  for (i = 0; i < cNames; i++) {
    if (2 == ::CompareStringW(lcid, NORM_IGNOREWIDTH, L"golangRequest", -1,
                              rgszNames[i], -1)) {
      *rgDispId = DISPID_EXTERNAL_GOLANGREQUEST;
    } else {
      hr = ResultFromScode(DISP_E_UNKNOWNNAME);
      *rgDispId = DISPID_UNKNOWN;
    }
  }
  return hr;
}

STDMETHODIMP CWebBrowserHost::Invoke(DISPID dispIdMember, REFIID riid,
                                     LCID lcid, WORD wFlags,
                                     DISPPARAMS *pDispParams,
                                     VARIANT *pVarResult, EXCEPINFO *pExcepInfo,
                                     UINT *puArgErr) {
  if (dispIdMember == DISPID_AMBIENT_DLCONTROL) {
    pVarResult->vt = VT_I4;
    pVarResult->lVal = m_dwAmbientDLControl;
    return S_OK;
  }
  if (dispIdMember == DISPID_AMBIENT_USERMODE) {
    V_VT(pVarResult) = VT_BOOL;
    V_BOOL(pVarResult) = VARIANT_TRUE;
    return S_OK;
  }
  if (dispIdMember == DISPID_EXTERNAL_GOLANGREQUEST) {
    if (wFlags & DISPATCH_PROPERTYGET) {
      if (pVarResult) {
        ::VariantInit(pVarResult);
        V_VT(pVarResult) = VT_BOOL;
        V_BOOL(pVarResult) = true;
      }
    }
    if (wFlags & DISPATCH_METHOD) {
      IDispatch *lpArgDisp = pDispParams->rgvarg[0].pdispVal;
      DISPID dispIdProp;
      HRESULT hr;
      LPOLESTR lpProp;
      DISPPARAMS dispparamsNoArgs;
      std::string strPath;
      std::string strArg;
      lpProp = const_cast<LPWSTR>(L"path");
      dispparamsNoArgs.cArgs = 0;
      dispparamsNoArgs.cNamedArgs = 0;
      hr = lpArgDisp->GetIDsOfNames(IID_NULL, &lpProp, 1, LOCALE_USER_DEFAULT,
                                    &dispIdProp);
      if (hr == S_OK) {
        VARIANT varProp;
        ::VariantInit(&varProp);
        hr = lpArgDisp->Invoke(dispIdProp, IID_NULL, LOCALE_USER_DEFAULT,
                               DISPATCH_PROPERTYGET, &dispparamsNoArgs,
                               &varProp, nullptr, nullptr);
        if (hr == S_OK) {
          strPath = exciton::util::ToUTF8String(varProp.bstrVal);
        }
        ::VariantClear(&varProp);
      }

      lpProp = const_cast<LPWSTR>(L"arg");
      hr = lpArgDisp->GetIDsOfNames(IID_NULL, &lpProp, 1, LOCALE_USER_DEFAULT,
                                    &dispIdProp);
      if (hr == S_OK) {
        VARIANT varProp;
        ::VariantInit(&varProp);
        hr = lpArgDisp->Invoke(dispIdProp, IID_NULL, LOCALE_USER_DEFAULT,
                               DISPATCH_PROPERTYGET, &dispparamsNoArgs,
                               &varProp, nullptr, nullptr);
        if (hr == S_OK) {
          strArg = exciton::util::ToUTF8String(varProp.bstrVal);
        }
        ::VariantClear(&varProp);
      }
      LOG_DEBUG("[%d] CWebBrowserHost::Invoke: call external: path=%s, args=%s",
                __LINE__, strPath.c_str(), strArg.c_str());
      Driver::Current().emitEvent(strPath, strArg);
    }
    return S_OK;
  }

  return DISP_E_MEMBERNOTFOUND;
}

STDMETHODIMP CWebBrowserHost::QueryStatus(const GUID *pguidCmdGroup,
                                          ULONG cCmds, OLECMD prgCmds[],
                                          OLECMDTEXT *pCmdText) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::Exec(const GUID *pguidCmdGroup, DWORD nCmdID,
                                   DWORD nCmdexecopt, VARIANT *pvaIn,
                                   VARIANT *pvaOut) {
  if (pguidCmdGroup &&
      IsEqualGUID(*pguidCmdGroup, CGID_DocHostCommandHandler)) {
    switch (nCmdID) {
    case OLECMDID_SHOWSCRIPTERROR: {
      IHTMLDocument2 *pDoc = nullptr;
      HRESULT hr;
      static const LPCWSTR awszPropNames[] = {L"errorLine", L"errorCharacter",
                                              L"errorCode", L"errorMessage",
                                              L"errorUrl"};

      pvaOut->vt = VT_BOOL;
      pvaOut->boolVal = VARIANT_FALSE;
      std::string errorMsg = "[";
      hr = pvaIn->punkVal->QueryInterface(IID_PPV_ARGS(&pDoc));
      if (hr == S_OK) {
        IHTMLWindow2 *pWindow = nullptr;
        hr = pDoc->get_parentWindow(&pWindow);
        if (hr == S_OK) {
          IHTMLEventObj *pEventObj = nullptr;
          hr = pWindow->get_event(&pEventObj);
          if (hr == S_OK) {
            for (int i = 0;
                 i < sizeof(awszPropNames) / sizeof(awszPropNames[0]); i++) {
              LPCWSTR lpwszPropName = awszPropNames[i];
              BSTR bstrPropName = ::SysAllocString(lpwszPropName);
              DISPID dispId;
              hr = pEventObj->GetIDsOfNames(IID_NULL, &bstrPropName, 1,
                                            LOCALE_SYSTEM_DEFAULT, &dispId);
              if (hr == S_OK) {
                VARIANT vEventInfo;
                DISPPARAMS params;
                params.cArgs = 0;
                params.cNamedArgs = 0;
                ::VariantInit(&vEventInfo);
                hr = pEventObj->Invoke(dispId, IID_NULL, LOCALE_SYSTEM_DEFAULT,
                                       DISPATCH_PROPERTYGET, &params,
                                       &vEventInfo, nullptr, nullptr);
                if (hr == S_OK) {
                  if (i != 0)
                    errorMsg += ", ";
                  errorMsg +=
                      exciton::util::FormatString("[%d]%S", i, lpwszPropName);
                  switch (vEventInfo.vt) {
                  case VT_I4:
                    errorMsg +=
                        exciton::util::FormatString("[%d]", vEventInfo.intVal);
                    break;
                  case VT_BSTR: {
                    std::wstring ws(vEventInfo.bstrVal,
                                    ::SysStringLen(vEventInfo.bstrVal));
                    errorMsg += exciton::util::FormatString("[%S]", ws.c_str());
                    break;
                  }
                  default:
                    errorMsg += exciton::util::FormatString("[UNKNOWN(%d)]",
                                                            vEventInfo.vt);
                  }
                }
                ::VariantClear(&vEventInfo);
              }
              ::SysFreeString(bstrPropName);
            }
            pEventObj->Release();
          }
          pWindow->Release();
        }
        pDoc->Release();
      }
      errorMsg += "]";
      LOG_ERROR("[%d] CWebBrowserHost::Exec: OLECMDID_SHOWSCRIPTERROR: %s",
                __LINE__, errorMsg.c_str());
      return S_OK;
    }
    default:
      return OLECMDERR_E_NOTSUPPORTED;
    }
  }

  return OLECMDERR_E_UNKNOWNGROUP;
}

STDMETHODIMP CWebBrowserHost::ShowContextMenu(DWORD dwID, POINT *ppt,
                                              IUnknown *pcmdtReserved,
                                              IDispatch *pdispReserved) {
  return S_FALSE;
}

STDMETHODIMP CWebBrowserHost::GetHostInfo(DOCHOSTUIINFO *pInfo) {
  pInfo->cbSize = sizeof(DOCHOSTUIINFO);
  pInfo->dwFlags = DOCHOSTUIFLAG_NO3DBORDER;
  pInfo->dwDoubleClick = DOCHOSTUIDBLCLK_DEFAULT;
  pInfo->pchHostCss = NULL;
  pInfo->pchHostNS = NULL;

  return S_OK;
}

STDMETHODIMP CWebBrowserHost::ShowUI(DWORD dwID,
                                     IOleInPlaceActiveObject *pActiveObject,
                                     IOleCommandTarget *pCommandTarget,
                                     IOleInPlaceFrame *pFrame,
                                     IOleInPlaceUIWindow *pDoc) {
  return S_FALSE;
}

STDMETHODIMP CWebBrowserHost::HideUI(VOID) { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::UpdateUI(VOID) { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::EnableModeless(BOOL fEnable) { return E_NOTIMPL; }

STDMETHODIMP CWebBrowserHost::OnDocWindowActivate(BOOL fActivate) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::OnFrameWindowActivate(BOOL fActivate) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::ResizeBorder(LPCRECT prcBorder,
                                           IOleInPlaceUIWindow *pUIWindow,
                                           BOOL fFrameWindow) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::TranslateAccelerator(LPMSG lpMsg,
                                                   const GUID *pguidCmdGroup,
                                                   DWORD nCmdID) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::GetOptionKeyPath(LPOLESTR *pchKey, DWORD dw) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::GetDropTarget(IDropTarget *pDropTarget,
                                            IDropTarget **ppDropTarget) {
  return E_NOTIMPL;
}

STDMETHODIMP CWebBrowserHost::GetExternal(IDispatch **ppDispatch) {
  return QueryInterface(IID_PPV_ARGS(ppDispatch));
}

STDMETHODIMP CWebBrowserHost::TranslateUrl(DWORD dwTranslate, OLECHAR *pchURLIn,
                                           OLECHAR **ppchURLOut) {
  return S_FALSE;
}

STDMETHODIMP CWebBrowserHost::FilterDataObject(IDataObject *pDO,
                                               IDataObject **ppDORet) {
  return E_NOTIMPL;
}

BOOL CWebBrowserHost::Create(HWND hwndParent, LPCWSTR lpszUrl, int nId) {
  IOleObject *pOleObject;
  HRESULT hr;
  MSG msg;
  RECT rc;

  hr = CoCreateInstance(CLSID_WebBrowser, NULL, CLSCTX_INPROC_SERVER,
                        IID_PPV_ARGS(&m_pWebBrowser2));
  if (FAILED(hr))
    return FALSE;

  m_nId = nId;

  m_hwnd = CreateWindowEx(
      WS_EX_CLIENTEDGE, g_szHostClassName, NULL,
      WS_CHILD | WS_CLIPCHILDREN | WS_VISIBLE | WS_CLIPSIBLINGS, 0, 0, 0, 0,
      hwndParent, (HMENU)(UINT_PTR)m_nId, NULL, NULL);
  SetWindowLongPtr(m_hwnd, GWLP_USERDATA, (LONG_PTR)this);

  m_pWebBrowser2->QueryInterface(IID_PPV_ARGS(&pOleObject));
  pOleObject->SetClientSite(static_cast<IOleClientSite *>(this));

  SetRectEmpty(&rc);
  hr = pOleObject->DoVerb(OLEIVERB_INPLACEACTIVATE, &msg,
                          static_cast<IOleClientSite *>(this), 0, m_hwnd, &rc);
  if (FAILED(hr)) {
    pOleObject->Release();
    m_pWebBrowser2->Release();
    m_pWebBrowser2 = NULL;
    return FALSE;
  }

  m_pEventSink = new CEventSink(*m_pContainer);
  m_pEventSink->Create(this, pOleObject);

  Navigate(lpszUrl);

  pOleObject->Release();

  return TRUE;
}

BOOL CWebBrowserHost::Destroy() {
  if (m_pEventSink != NULL) {
    m_pEventSink->Destroy();
    m_pEventSink->Release();
  }

  if (m_pWebBrowser2 != NULL) {
    IOleObject *pOleObject;
    IOleInPlaceObject *pOleInPlaceObject;

    m_pWebBrowser2->QueryInterface(IID_PPV_ARGS(&pOleObject));
    m_pWebBrowser2->QueryInterface(IID_PPV_ARGS(&pOleInPlaceObject));
    pOleInPlaceObject->InPlaceDeactivate();
    pOleObject->Close(OLECLOSE_NOSAVE);
    pOleObject->Release();
    pOleInPlaceObject->Release();

    m_pWebBrowser2->Release();
  }

  DestroyWindow(m_hwnd);

  return TRUE;
}

LRESULT CWebBrowserHost::WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam,
                                    LPARAM lParam) {
  switch (uMsg) {

  case WM_SIZE: {
    RECT rc = {0, 0, LOWORD(lParam), HIWORD(lParam)};
    IOleInPlaceObject *pOleInPlaceObject;
    m_pWebBrowser2->put_Width(LOWORD(lParam));
    m_pWebBrowser2->put_Height(HIWORD(lParam));
    m_pWebBrowser2->QueryInterface(IID_PPV_ARGS(&pOleInPlaceObject));
    pOleInPlaceObject->SetObjectRects(&rc, &rc);
    pOleInPlaceObject->Release();
    picojson::object arg;
    arg.emplace("width", static_cast<int64_t>(LOWORD(lParam)));
    arg.emplace("hdight", static_cast<int64_t>(HIWORD(lParam)));
    picojson::value val(arg);
    auto json = val.serialize();
    auto name = exciton::util::FormatString("/window/%s/resize", m_strID.c_str());
    Driver::Current().emitEvent(name, json);
    return 0;
  }

  default:
    break;
  }

  return DefWindowProc(hwnd, uMsg, wParam, lParam);
}

void CWebBrowserHost::SetWindowSize(int x, int y, int nWidth, int nHeight) {
  MoveWindow(m_hwnd, x, y, nWidth, nHeight, TRUE);
}

void CWebBrowserHost::Show(BOOL bShow) {
  if (bShow) {
    ShowWindow(m_hwnd, SW_SHOW);
  } else {
    ShowWindow(m_hwnd, SW_HIDE);
  }
}

BOOL CWebBrowserHost::Navigate(LPCWSTR lpszUrl) {
  HRESULT hr;
  BSTR bstrUrl;
  VARIANT varFlags, varTargetFrameName, varPostData, varHeaders;

  if (lpszUrl == NULL)
    bstrUrl = SysAllocString(L"");
  else
    bstrUrl = SysAllocString(lpszUrl);

  VariantInit(&varFlags);
  VariantInit(&varTargetFrameName);
  VariantInit(&varPostData);
  VariantInit(&varHeaders);

  hr = m_pWebBrowser2->Navigate(bstrUrl, &varFlags, &varTargetFrameName,
                                &varPostData, &varHeaders);

  SysFreeString(bstrUrl);

  return hr == S_OK;
}

HRESULT CWebBrowserHost::QueryBrowserInterface(REFIID riid, void **ppvObject) {
  return m_pWebBrowser2->QueryInterface(riid, ppvObject);
}

int CWebBrowserHost::GetId() { return m_nId; }

void CWebBrowserHost::AmbientPropertyChange(int nId) {
  DWORD dwAmbientDLControl[] = {DLCTL_DLIMAGES,        DLCTL_VIDEOS,
                                DLCTL_BGSOUNDS,        DLCTL_NO_SCRIPTS,
                                DLCTL_NO_JAVA,         DLCTL_NO_RUNACTIVEXCTLS,
                                DLCTL_NO_DLACTIVEXCTLS};
  IOleControl *pOleControl;

  m_dwAmbientDLControl ^= dwAmbientDLControl[nId];

  m_pWebBrowser2->QueryInterface(IID_PPV_ARGS(&pOleControl));
  pOleControl->OnAmbientPropertyChange(DISPID_AMBIENT_DLCONTROL);
  pOleControl->Release();
}

BOOL CWebBrowserHost::IsContainAmbientProperty(int nId) {
  DWORD dwAmbientDLControl[] = {DLCTL_DLIMAGES,        DLCTL_VIDEOS,
                                DLCTL_BGSOUNDS,        DLCTL_NO_SCRIPTS,
                                DLCTL_NO_JAVA,         DLCTL_NO_RUNACTIVEXCTLS,
                                DLCTL_NO_DLACTIVEXCTLS};

  return m_dwAmbientDLControl & dwAmbientDLControl[nId];
}

void CWebBrowserHost::ForwardOrBack(BOOL bForward) {
  if (bForward)
    m_pWebBrowser2->GoForward();
  else
    m_pWebBrowser2->GoBack();
}

void CWebBrowserHost::GetTravelState(LPBOOL lpbEnableForward,
                                     LPBOOL lpbEnableBack) {
  *lpbEnableForward = m_bEnableForward;
  *lpbEnableBack = m_bEnableBack;
}

void CWebBrowserHost::SetTravelState(BOOL bEnableForward, BOOL bEnableBack) {
  m_bEnableForward = bEnableForward;
  m_bEnableBack = bEnableBack;
}

HRESULT CWebBrowserHost::TranslateAccelerator(LPMSG lpMsg) {
  HRESULT hr = S_OK;
  IOleInPlaceActiveObject *pOleInPlaceActiveObject;

  hr = m_pWebBrowser2->QueryInterface(IID_PPV_ARGS(&pOleInPlaceActiveObject));
  if (hr == S_OK) {
    hr = pOleInPlaceActiveObject->TranslateAccelerator(lpMsg);
    pOleInPlaceActiveObject->Release();
  }

  return hr;
}

void CWebBrowserHost::Exec(OLECMDID nCmdID, OLECMDEXECOPT nCmdexecopt,
                           VARIANT *pvaIn, VARIANT *pvaOut) {
  m_pWebBrowser2->ExecWB(nCmdID, nCmdexecopt, pvaIn, pvaOut);
}

void CWebBrowserHost::ExecDocument(const GUID *pguid, DWORD nCmdID,
                                   DWORD nCmdexecopt, VARIANT *pvaIn,
                                   VARIANT *pvaOut) {
  IDispatch *pDocument;
  IOleCommandTarget *pOleCommandTarget;

  m_pWebBrowser2->get_Document(&pDocument);

  pDocument->QueryInterface(IID_PPV_ARGS(&pOleCommandTarget));
  pOleCommandTarget->Exec(pguid, nCmdID, nCmdexecopt, pvaIn, pvaOut);
  pOleCommandTarget->Release();

  pDocument->Release();
}

void CWebBrowserHost::OnDocumentComplate(IDispatch *pDisp, const std::wstring& strURL) {
  if (strURL != L"about:blank") {
    return;
  }
  if (m_strInitialHtml.empty()) {
    return;
  }
  auto initHTML = m_strInitialHtml;
  m_strInitialHtml.clear();
  if (m_pHtmlMoniker) {
    m_pHtmlMoniker->Release();
    m_pHtmlMoniker = nullptr;
  }
  m_pHtmlMoniker = new HtmlMoniker();
  m_pHtmlMoniker->SetHtml(initHTML);
  m_pHtmlMoniker->SetBaseURL(strURL);

  IDispatch* pDocDisp;
  HRESULT hr = m_pWebBrowser2->get_Document(&pDocDisp);
  if (hr == S_OK && pDocDisp) {
    IHTMLDocument2* pDoc;
    hr = pDocDisp->QueryInterface(IID_PPV_ARGS(&pDoc));
    if (hr == S_OK && pDoc) {
      IPersistMoniker* pPersistMoniker;
      hr = pDoc->QueryInterface(IID_PPV_ARGS(&pPersistMoniker));
      if (hr == S_OK && pPersistMoniker) {
        IMoniker* pMoniker;
        hr = m_pHtmlMoniker->QueryInterface(IID_PPV_ARGS(&pMoniker));
        if (hr == S_OK && pMoniker) {
          hr = pPersistMoniker->Load(TRUE, pMoniker, nullptr, STGM_READ);
          pMoniker->Release();
        }
        pPersistMoniker->Release();
      }
      pDoc->Release();
    }
    pDocDisp->Release();
  }
}

void CWebBrowserHost::EvaluateJavasScript(const std::wstring &evalFuncName,
                                          const std::wstring &jsonArg,
                                          VARIANT *pRetValue) {
  IDispatch *pDocDisp;
  HRESULT hr;

  hr = m_pWebBrowser2->get_Document(&pDocDisp);
  if (hr == S_OK) {
    IHTMLDocument2 *pDocument;
    hr = pDocDisp->QueryInterface(IID_PPV_ARGS(&pDocument));
    if (hr == S_OK) {
      IDispatch *pScript;
      hr = pDocument->get_Script(&pScript);
      if (hr == S_OK) {
        DISPID dispIdExciton;
        BSTR bstrExciton = ::SysAllocString(L"exciton");
        hr = pScript->GetIDsOfNames(IID_NULL, &bstrExciton, 1,
                                    LOCALE_USER_DEFAULT, &dispIdExciton);
        if (hr == S_OK) {
          DISPPARAMS dispparamsNoArgs;
          VARIANT varExciton;
          dispparamsNoArgs.cArgs = 0;
          dispparamsNoArgs.cNamedArgs = 0;
          ::VariantInit(&varExciton);
          hr = pScript->Invoke(dispIdExciton, IID_NULL, LOCALE_USER_DEFAULT,
                               DISPATCH_PROPERTYGET, &dispparamsNoArgs,
                               &varExciton, nullptr, nullptr);
          if ((hr == S_OK) && (varExciton.vt == VT_DISPATCH)) {
            IDispatch *pExciton = varExciton.pdispVal;
            DISPID dispIdReqBE;
            BSTR bstrReqBE = ::SysAllocString(
                pRetValue ? L"requestBrowerEventSync" : L"requestBrowserEvent");
            hr = pExciton->GetIDsOfNames(IID_NULL, &bstrReqBE, 1,
                                         LOCALE_USER_DEFAULT, &dispIdReqBE);
            if (hr == S_OK) {
              VARIANT varArgs[2];
              DISPPARAMS dispparams;
              memset(&dispparams, 0, sizeof dispparams);
              ::VariantInit(&varArgs[0]);
              ::VariantInit(&varArgs[1]);
              // varArgs is argument stack
              varArgs[0].vt = VT_BSTR;
              varArgs[0].bstrVal = ::SysAllocString(jsonArg.c_str());
              varArgs[1].vt = VT_BSTR;
              varArgs[1].bstrVal = ::SysAllocString(evalFuncName.c_str());
              dispparams.cArgs = 2;
              dispparams.rgvarg = varArgs;
              dispparams.cNamedArgs = 0;

              EXCEPINFO excepInfo;
              memset(&excepInfo, 0, sizeof excepInfo);
              UINT nArgErr = (UINT)-1; // initialize to invalid arg
              // Call JavaScript function
              hr = pExciton->Invoke(dispIdReqBE, IID_NULL, 0, DISPATCH_METHOD,
                                    &dispparams, pRetValue, &excepInfo,
                                    &nArgErr);
              if (hr == S_OK) {
                // TODO: return value
              }

              ::VariantClear(&varArgs[1]);
              ::VariantClear(&varArgs[0]);
            }
            ::SysFreeString(bstrReqBE);
          }
          ::VariantClear(&varExciton);
        }
        ::SysFreeString(bstrExciton);
        pScript->Release();
      }
      pDocument->Release();
    }
    pDocDisp->Release();
  }
}

void CWebBrowserHost::PutFullscreen(bool bEnter) {
  VARIANT_BOOL b = bEnter ? VARIANT_TRUE : VARIANT_FALSE;
  // HRESULT hr = m_pWebBrowser2->put_FullScreen(b);
  HRESULT hr = m_pWebBrowser2->put_TheaterMode(b);
  if (hr != S_OK) {
    LOG_ERROR("[%d] CWebBrowserHost::PutFullscreen() failed:[0x%08x]", __LINE__,
              hr);
  }
}