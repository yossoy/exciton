#pragma once
#include <oaidl.h>
#include <string>
#include <windows.h>

class HtmlMoniker : public IMoniker {
  ULONG refCount_;
  IStream *htmlStream_;
  std::string html_;
  std::wstring baseURL_;

public:
  HtmlMoniker(void);
  virtual ~HtmlMoniker(void);

public:
  void SetHtml(std::string_view html);
  void SetBaseURL(std::wstring_view baseURL);

public:
public:
  // IUnknown
  STDMETHODIMP QueryInterface(REFIID riid, void **ppvObject);
  ULONG STDMETHODCALLTYPE AddRef(void);
  ULONG STDMETHODCALLTYPE Release(void);

  // IMoniker
  STDMETHODIMP BindToStorage(IBindCtx *pbc, IMoniker *pmkToLeft, REFIID riid,
                             void **ppvObj);
  STDMETHODIMP GetDisplayName(IBindCtx *pbc, IMoniker *pmkToLeft,
                              LPOLESTR *ppszDisplayName);
  STDMETHODIMP BindToObject(IBindCtx *, IMoniker *, REFIID, void **) {
    return E_NOTIMPL;
  }
  STDMETHODIMP Reduce(IBindCtx *, DWORD, IMoniker **, IMoniker **) {
    return E_NOTIMPL;
  }
  STDMETHODIMP ComposeWith(IMoniker *, BOOL, IMoniker **) { return E_NOTIMPL; }
  STDMETHODIMP Enum(BOOL, IEnumMoniker **) { return E_NOTIMPL; }
  STDMETHODIMP IsEqual(IMoniker *) { return E_NOTIMPL; }
  STDMETHODIMP Hash(DWORD *) { return E_NOTIMPL; }
  STDMETHODIMP IsRunning(IBindCtx *, IMoniker *, IMoniker *) {
    return E_NOTIMPL;
  }
  STDMETHODIMP GetTimeOfLastChange(IBindCtx *, IMoniker *, FILETIME *) {
    return E_NOTIMPL;
  }
  STDMETHODIMP Inverse(IMoniker **) { return E_NOTIMPL; }
  STDMETHODIMP CommonPrefixWith(IMoniker *, IMoniker **) { return E_NOTIMPL; }
  STDMETHODIMP RelativePathTo(IMoniker *, IMoniker **) { return E_NOTIMPL; }
  STDMETHODIMP ParseDisplayName(IBindCtx *, IMoniker *, LPOLESTR, ULONG *,
                                IMoniker **) {
    return E_NOTIMPL;
  }
  STDMETHODIMP IsSystemMoniker(DWORD *pdwMksys) {
    if (!pdwMksys) {
      return E_POINTER;
    }
    *pdwMksys = MKSYS_NONE;
    return S_OK;
  }

  // IPersistStream methods
  STDMETHODIMP Save(IStream *, BOOL) { return E_NOTIMPL; }
  STDMETHODIMP IsDirty() { return E_NOTIMPL; }
  STDMETHODIMP Load(IStream *) { return E_NOTIMPL; }
  STDMETHODIMP GetSizeMax(ULARGE_INTEGER *) { return E_NOTIMPL; }

  // IPersist
  STDMETHODIMP GetClassID(CLSID *) { return E_NOTIMPL; }
};