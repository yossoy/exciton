#include "htmlmoniker.h"

namespace {
IStream *createStreamFromText(const std::string &str) {
  IStream *stream;
  if (S_OK != ::CreateStreamOnHGlobal(NULL, TRUE, &stream)) {
    return nullptr;
  }

  ULONG written;
  if ((S_OK != stream->Write(str.data(), str.size(), &written)) || written != str.size()) {
      stream->Release();
      return nullptr;
    }

  LARGE_INTEGER zero = {0};
  stream->Seek(zero, STREAM_SEEK_SET, nullptr);
  stream->AddRef();
  return stream;
}

static LPOLESTR oleStrDup(const std::wstring &str) {
  size_t cb = sizeof(WCHAR) * (str.size() + 1);
  LPOLESTR ret = (LPOLESTR)CoTaskMemAlloc(cb);
  if (ret) {
    memcpy(ret, str.data(), cb);
  }
  return ret;
}

} // namespace

HtmlMoniker::HtmlMoniker() : refCount_(1), htmlStream_(nullptr) {}

HtmlMoniker::~HtmlMoniker() {
  if (htmlStream_) {
    htmlStream_->Release();
  }
}

void HtmlMoniker::SetHtml(std::string_view html) {
  html_ = html;
  if (htmlStream_) {
    htmlStream_->Release();
  }
  htmlStream_ = createStreamFromText(html_);
}

void HtmlMoniker::SetBaseURL(std::wstring_view baseURL) { baseURL_ = baseURL; }

STDMETHODIMP HtmlMoniker::BindToStorage(IBindCtx *, IMoniker *, REFIID riid,
                                        void **ppvObj) {
  LARGE_INTEGER seek = {0};
  htmlStream_->Seek(seek, STREAM_SEEK_SET, nullptr);
  return htmlStream_->QueryInterface(riid, ppvObj);
}

STDMETHODIMP HtmlMoniker::GetDisplayName(IBindCtx *, IMoniker *,
                                         LPOLESTR *ppszDisplayName) {
  if (!ppszDisplayName) {
    return E_POINTER;
  }
  *ppszDisplayName = oleStrDup(baseURL_);
  return *ppszDisplayName ? S_OK : E_OUTOFMEMORY;
}

STDMETHODIMP HtmlMoniker::QueryInterface(REFIID riid, void **ppv) {
  if (IsEqualIID(riid, IID_IUnknown) || IsEqualIID(riid, IID_IMoniker))
    *ppv = static_cast<IMoniker *>(this);
  else if (IsEqualIID(riid, IID_IPersistStream))
    *ppv = static_cast<IPersistStream *>(this);
  else if (IsEqualIID(riid, IID_IPersist))
    *ppv = static_cast<IPersist *>(this);
  else
    return E_NOINTERFACE;
  AddRef();

  return S_OK;
}

ULONG STDMETHODCALLTYPE HtmlMoniker::AddRef() {
  return InterlockedIncrement(&refCount_);
}

ULONG STDMETHODCALLTYPE HtmlMoniker::Release() {
  LONG res = InterlockedDecrement(&refCount_);
  if (0 == res)
    delete this;
  return res;
}
