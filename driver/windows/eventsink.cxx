#include <windows.h>

#include <exdispid.h>
#include <shlobj.h>

#include "browsercontainer.h"
#include "browserhost.h"
#include "eventsink.h"

CEventSink::CEventSink(CWebBrowserContainer &container)
    : m_container(container), m_cRef(1), m_pWebBrowserHost(nullptr),
      m_pConnectionPoint(nullptr), m_dwCookie(0) {}

CEventSink::~CEventSink() {}

STDMETHODIMP_(ULONG) CEventSink::AddRef() {
  return InterlockedIncrement(&m_cRef);
}

STDMETHODIMP_(ULONG) CEventSink::Release() {
  if (InterlockedDecrement(&m_cRef) == 0) {
    delete this;
    return 0;
  }

  return m_cRef;
}

STDMETHODIMP CEventSink::QueryInterface(REFIID riid, void **ppvObject) {
  *ppvObject = NULL;

  if (IsEqualIID(riid, IID_IUnknown) || IsEqualIID(riid, IID_IDispatch) ||
      IsEqualIID(riid, DIID_DWebBrowserEvents2))
    *ppvObject = static_cast<IDispatch *>(this);
  else
    return E_NOINTERFACE;

  AddRef();

  return S_OK;
}

STDMETHODIMP CEventSink::GetTypeInfoCount(UINT *pctinfo) {
  *pctinfo = 0;

  return S_OK;
}

STDMETHODIMP CEventSink::GetTypeInfo(UINT iTInfo, LCID lcid,
                                     ITypeInfo **ppTInfo) {
  return E_NOTIMPL;
}

STDMETHODIMP CEventSink::GetIDsOfNames(REFIID riid, LPOLESTR *rgszNames,
                                       UINT cNames, LCID lcid,
                                       DISPID *rgDispId) {
  return E_NOTIMPL;
}

STDMETHODIMP CEventSink::Invoke(DISPID dispIdMember, REFIID riid, LCID lcid,
                                WORD wFlags, DISPPARAMS *pDispParams,
                                VARIANT *pVarResult, EXCEPINFO *pExcepInfo,
                                UINT *puArgErr) {
  switch (dispIdMember) {
  case DISPID_BEFORENAVIGATE2:
    OnBeforeNavigate2(
        pDispParams->rgvarg[6].pdispVal, pDispParams->rgvarg[5].pvarVal,
        pDispParams->rgvarg[4].pvarVal, pDispParams->rgvarg[3].pvarVal,
        pDispParams->rgvarg[2].pvarVal, pDispParams->rgvarg[1].pvarVal,
        pDispParams->rgvarg[0].pboolVal);
    break;
  case DISPID_NAVIGATECOMPLETE2:
    OnNavigateComplete2(pDispParams->rgvarg[1].pdispVal,
                        pDispParams->rgvarg[0].pvarVal);
    break;
  //   case DISPID_NEWWINDOW3:
  //     OnNewWindow3(pDispParams->rgvarg[4].ppdispVal,
  //                  pDispParams->rgvarg[3].pboolVal,
  // 				 pDispParams->rgvarg[2].lVal,
  //                  pDispParams->rgvarg[1].bstrVal,
  //                  pDispParams->rgvarg[0].bstrVal);
  //     break;
  case DISPID_COMMANDSTATECHANGE:
    OnCommandStateChange(pDispParams->rgvarg[1].lVal,
                         pDispParams->rgvarg[0].boolVal);
    break;
  case DISPID_TITLECHANGE:
    OnTitleChange(pDispParams->rgvarg[0].bstrVal);
    break;
  case DISPID_STATUSTEXTCHANGE:
    OnStatusTextChange(pDispParams->rgvarg[0].bstrVal);
    break;
  case DISPID_DOCUMENTCOMPLETE:
    OnDocumentComplete(pDispParams->rgvarg[1].pdispVal,
                       pDispParams->rgvarg[0].pvarVal);
    break;
  default:
    return DISP_E_MEMBERNOTFOUND;
  }

  return S_OK;
}

BOOL CEventSink::Create(CWebBrowserHost *pWebBrowserHost, IUnknown *pUnknown) {
  HRESULT hr;
  IConnectionPointContainer *pConnectionPointContainer;

  m_pWebBrowserHost = pWebBrowserHost;

  hr = pUnknown->QueryInterface(IID_PPV_ARGS(&pConnectionPointContainer));
  if (FAILED(hr))
    return FALSE;

  pConnectionPointContainer->FindConnectionPoint(DIID_DWebBrowserEvents2,
                                                 &m_pConnectionPoint);
  pConnectionPointContainer->Release();

  m_pConnectionPoint->Advise(this, &m_dwCookie);

  return TRUE;
}

void CEventSink::Destroy() {
  if (m_pConnectionPoint != NULL) {
    m_pConnectionPoint->Unadvise(m_dwCookie);
    m_pConnectionPoint->Release();
  }
}

void CEventSink::OnBeforeNavigate2(IDispatch *pDisp, VARIANT *url,
                                   VARIANT *Flags, VARIANT *TargetFrameName,
                                   VARIANT *PostData, VARIANT *Headers,
                                   VARIANT_BOOL *Cancel) {}

void CEventSink::OnDocumentComplete(IDispatch *pDisp, VARIANT *url) {
  CWebBrowserHost *pWebBrowserHost;
  pWebBrowserHost = m_container.GetActiveBrowser();
  pWebBrowserHost->OnDocumentComplate(pDisp);
}

void CEventSink::OnNavigateComplete2(IDispatch *pDispatch, VARIANT *URL) {
  CWebBrowserHost *pWebBrowserHost;
  IWebBrowser2 *pWebBrowser2;
  BSTR bstr;

  pWebBrowserHost = m_container.GetActiveBrowser();
  pWebBrowserHost->QueryBrowserInterface(IID_PPV_ARGS(&pWebBrowser2));
  pWebBrowser2->get_LocationURL(&bstr);

  m_container.OnNavigateComplete2(bstr);

  SysFreeString(bstr);
  pWebBrowser2->Release();
  pWebBrowserHost->Release();
}

// void CEventSink::OnNewWindow3(IDispatch **ppDisp, VARIANT_BOOL *Cancel, DWORD
// dwFlags, BSTR bstrUrlContext, BSTR bstrUrl)
// {
// 	*Cancel = VARIANT_TRUE;

// 	g_pWebBrowserContainer->OnNewWindow3(bstrUrl);
// }

void CEventSink::OnCommandStateChange(long Command, VARIANT_BOOL Enable) {
  BOOL bEnableForward, bEnableBack;

  m_pWebBrowserHost->GetTravelState(&bEnableForward, &bEnableBack);

  if (Command == CSC_NAVIGATEFORWARD) {
    if (Enable == VARIANT_TRUE)
      bEnableForward = TRUE;
    else
      bEnableForward = FALSE;
  } else if (Command == CSC_NAVIGATEBACK) {
    if (Enable == VARIANT_TRUE)
      bEnableBack = TRUE;
    else
      bEnableBack = FALSE;
  } else
    return;

  m_pWebBrowserHost->SetTravelState(bEnableForward, bEnableBack);

  m_container.OnCommandStateChange(Command, Enable);
}

void CEventSink::OnStatusTextChange(BSTR Text) {
  ::OutputDebugStringW(static_cast<LPCWSTR>(Text));
  //	g_pWebBrowserContainer->OnStatusTextChange(Text);
}

void CEventSink::OnTitleChange(BSTR Text) { m_container.OnTitleChange(Text); }
