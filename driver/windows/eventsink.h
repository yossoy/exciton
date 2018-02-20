#pragma once
class CWebBrowserHost;
class CWebBrowserContainer;

class CEventSink : public IDispatch {
public:
  STDMETHODIMP QueryInterface(REFIID riid, void **ppvObject);
  STDMETHODIMP_(ULONG) AddRef();
  STDMETHODIMP_(ULONG) Release();

  STDMETHODIMP GetTypeInfoCount(UINT *pctinfo);
  STDMETHODIMP GetTypeInfo(UINT iTInfo, LCID lcid, ITypeInfo **ppTInfo);
  STDMETHODIMP GetIDsOfNames(REFIID riid, LPOLESTR *rgszNames, UINT cNames,
                             LCID lcid, DISPID *rgDispId);
  STDMETHODIMP Invoke(DISPID dispIdMember, REFIID riid, LCID lcid, WORD wFlags,
                      DISPPARAMS *pDispParams, VARIANT *pVarResult,
                      EXCEPINFO *pExcepInfo, UINT *puArgErr);

  CEventSink(CWebBrowserContainer &container);
  CEventSink(const CEventSink &) = delete;
  CEventSink(CEventSink &&) = delete;
  ~CEventSink();
  BOOL Create(CWebBrowserHost *pWebBrowserHost, IUnknown *pUnknown);
  void Destroy();

private:
  void OnBeforeNavigate2(IDispatch *pDisp, VARIANT *url, VARIANT *Flags,
                         VARIANT *TargetFrameName, VARIANT *PostData,
                         VARIANT *Headers, VARIANT_BOOL *Cancel);
  void OnDocumentComplete(IDispatch *pDisp, VARIANT *url);
  // void OnNewWindow3(IDispatch **ppDisp, VARIANT_BOOL *Cancel, DWORD dwFlags,
  // BSTR bstrUrlContext, BSTR bstrUrl);
  void OnNavigateComplete2(IDispatch *pDispatch, VARIANT *URL);
  void OnCommandStateChange(long Command, VARIANT_BOOL Enable);
  void OnStatusTextChange(BSTR Text);
  void OnTitleChange(BSTR Text);

private:
  CWebBrowserContainer &m_container;
  LONG m_cRef;
  CWebBrowserHost *m_pWebBrowserHost;
  IConnectionPoint *m_pConnectionPoint;
  DWORD m_dwCookie;
};