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
#include "dialog.h"
#include "driver.h"
#include "myjson.h"
#include "util.h"
#include "log.h"

namespace {
std::string
getIdFromParam(const std::map<std::string, std::string> &parameter) {
  auto fiter = parameter.find("id");
  if (fiter == parameter.end()) {
    return "";
  }
  return (*fiter).second;
}

HWND resolveParentWindow(const std::string &id) {
  Driver &d = Driver::Current();
  if (!id.empty()) {
    auto fiter = d.Hosts().find(id);
    if (fiter != d.Hosts().end()) {
      CWebBrowserHost *p = (*fiter).second.ptr_;
      HWND hWnd = NULL;
      if (S_OK == p->GetWindow(&hWnd)) {
        return hWnd;
      }
      auto container = p->GetHostContainer();
      return container->GetWindow();
    }
  }
  auto pActiveContainer = CWebBrowserContainer::GetActiveContainer();
  return pActiveContainer->GetWindow();
}

const int nMsgBoxIdOfs = 1000;

HRESULT CALLBACK messageBoxCallbackProc(HWND hwnd, UINT uNotification,
                                        WPARAM wParam, LPARAM lParam,
                                        LONG_PTR dwRefData) {
#if 0
	//TODO:
	if (uNotification == TDN_HYPERLINK_CLICKED) {
        link is open default browser?
        or 
        new tab?
        ShellExecuteEx((LPWSTR)lParam)?

    }
		g_pWebBrowserContainer->Navigate((LPWSTR)lParam, TRUE);
#endif

  return S_OK;
}

void showMessageBox(const picojson::value &argument,
                    const std::map<std::string, std::string> &parameter,
                    int responceNo) {
  Driver &d = Driver::Current();
  std::map<int, int> mapBtnId2Idx;
  TASKDIALOGCONFIG config;
  std::vector<std::wstring> vecButtonLabels;
  std::vector<TASKDIALOG_BUTTON> vecButtons;
  ZeroMemory(&config, sizeof(config));
  config.cbSize = sizeof(TASKDIALOGCONFIG);

  // icon
  auto type =
      static_cast<MESSAGE_BOX_TYPE>(argument.get("type").get<int64_t>());
  switch (type) {
  case MESSAGE_BOX_TYPE_INFO:
    config.pszMainIcon = TD_INFORMATION_ICON;
    break;
  case MESSAGE_BOX_TYPE_WARNING:
    config.pszMainIcon = TD_WARNING_ICON;
    break;
  case MESSAGE_BOX_TYPE_ERROR:
    config.pszMainIcon = TD_ERROR_ICON;
    break;
  case MESSAGE_BOX_TYPE_QUESTION:
    // no question icon in TaskDialog style?
  default:
    break;
  }

  // parentWindow
  {
    auto &parentId = argument.get("windowId").get<std::string>();
    config.hwndParent = resolveParentWindow(parentId);
  }

  // buttons
  {
    auto &buttons = argument.get("buttons").get<picojson::array>();
    std::vector<std::string> strbuttons;
    strbuttons.reserve(buttons.size());
    for (auto &button : buttons) {
      strbuttons.push_back(button.get<std::string>());
    }
    // Detect default buttons only
    if (strbuttons.empty()) {
      if (MESSAGE_BOX_TYPE_QUESTION == type) {
        config.dwCommonButtons = TDCBF_YES_BUTTON | TDCBF_NO_BUTTON;
        mapBtnId2Idx[IDYES] = 0;
        mapBtnId2Idx[IDNO] = 1;
      } else {
        config.dwCommonButtons = TDCBF_OK_BUTTON;
        mapBtnId2Idx[IDOK] = 0;
      }
    } else {
      DWORD dwCmnBtn = 0;
      size_t nCmnCnt = 0;
      std::map<int, int> mapCmnBttn;
      for (size_t idx = 0; idx < strbuttons.size(); idx++) {
        static const struct l2id {
          const char *l;
          DWORD id;
        } al2id[] = {
            {"OK", TDCBF_OK_BUTTON},       {"YES", TDCBF_YES_BUTTON},
            {"NO", TDCBF_NO_BUTTON},       {"Cancel", TDCBF_CANCEL_BUTTON},
            {"Retry", TDCBF_RETRY_BUTTON}, {"Close", TDCBF_CLOSE_BUTTON},
        };
        auto &label = strbuttons[idx];
        bool match = false;
        for (auto &l2id : al2id) {
          if (strcasecmp(label.c_str(), l2id.l) == 0) {
            dwCmnBtn |= l2id.id;
            mapCmnBttn[l2id.id] = idx;
            match = true;
            break;
          }
        }
        if (!match) {
          break;
        }
        nCmnCnt++;
      }
      if (nCmnCnt == strbuttons.size()) {
        strbuttons.clear();
        config.dwCommonButtons = dwCmnBtn;
        mapBtnId2Idx = mapCmnBttn;
      } else {
        config.dwCommonButtons = 0;
        vecButtonLabels.resize(strbuttons.size());
        vecButtons.resize(strbuttons.size());
        for (size_t idx = 0; idx < strbuttons.size(); idx++) {
          vecButtonLabels[idx] = exciton::util::ToUTF16String(strbuttons[idx]);
          vecButtons[idx].pszButtonText = vecButtonLabels[idx].c_str();
          vecButtons[idx].nButtonID = nMsgBoxIdOfs + idx;
        }
        config.pButtons = vecButtons.data();
        config.cButtons = static_cast<UINT>(vecButtons.size());
      }
    }
  }

  // auto idstr = getIdFromParam(parameter);
  int defaultIdx = static_cast<int>(argument.get("defaultId").get<int64_t>());
  int i;

  // title, message, etc.
  auto title =
      exciton::util::ToUTF16String(argument.get("title").get<std::string>());
  config.pszWindowTitle = title.c_str();
  auto message =
      exciton::util::ToUTF16String(argument.get("message").get<std::string>());
  config.pszMainInstruction = message.c_str();
  auto detail =
      exciton::util::ToUTF16String(argument.get("detail").get<std::string>());
  config.pszContent = detail.c_str();
  config.pfCallback = messageBoxCallbackProc;

  int nResult;
  auto r = TaskDialogIndirect(&config, &nResult, NULL, NULL);
  auto fiter = mapBtnId2Idx.find(nResult);
  if (fiter != mapBtnId2Idx.end()) {
    nResult = (*fiter).second;
  }
  if (r == S_OK) {
    Driver::Current().responceEventResult(
        responceNo, picojson::value(static_cast<int64_t>(nResult)));
  }
}

class CDialogEventHandler : public IFileDialogEvents,
                            public IFileDialogControlEvents {
public:
  // IUnknown methods
  IFACEMETHODIMP QueryInterface(REFIID riid, void **ppv) {
    *ppv = nullptr;

    if (IsEqualIID(riid, IID_IUnknown) ||
        IsEqualIID(riid, IID_IFileDialogEvents)) {
      *ppv = static_cast<IFileDialogEvents *>(this);
    } else if (IsEqualIID(riid, IID_IFileDialogControlEvents)) {
      *ppv = static_cast<IFileDialogControlEvents *>(this);
    } else {
      return E_NOINTERFACE;
    }
    AddRef();

    return S_OK;
  }

  IFACEMETHODIMP_(ULONG) AddRef() { return InterlockedIncrement(&_cRef); }

  IFACEMETHODIMP_(ULONG) Release() {
    long cRef = InterlockedDecrement(&_cRef);
    if (!cRef)
      delete this;
    return cRef;
  }

  // IFileDialogEvents methods
  IFACEMETHODIMP OnFileOk(IFileDialog *) { return S_OK; };
  IFACEMETHODIMP OnFolderChange(IFileDialog *) { return S_OK; };
  IFACEMETHODIMP OnFolderChanging(IFileDialog *, IShellItem *) { return S_OK; };
  IFACEMETHODIMP OnHelp(IFileDialog *) { return S_OK; };
  IFACEMETHODIMP OnSelectionChange(IFileDialog *) { return S_OK; };
  IFACEMETHODIMP OnShareViolation(IFileDialog *, IShellItem *,
                                  FDE_SHAREVIOLATION_RESPONSE *) {
    return S_OK;
  };
  IFACEMETHODIMP OnTypeChange(IFileDialog *pfd) {
    IFileSaveDialog *pfsd;
    HRESULT hr = pfd->QueryInterface(&pfsd);
    if (SUCCEEDED(hr)) {
      UINT uIndex;
      hr = pfsd->GetFileTypeIndex(&uIndex); // index of current file-type
      if (SUCCEEDED(hr)) {
        IPropertyDescriptionList *pdl = nullptr;
      }
      pfsd->Release();
    }
    return hr;
  }
  IFACEMETHODIMP OnOverwrite(IFileDialog *, IShellItem *,
                             FDE_OVERWRITE_RESPONSE *) {
    return S_OK;
  };

  // IFileDialogControlEvents methods
  IFACEMETHODIMP OnItemSelected(IFileDialogCustomize *pfdc, DWORD dwIDCtl,
                                DWORD dwIDItem) {
    return E_NOTIMPL;
  }
  IFACEMETHODIMP OnButtonClicked(IFileDialogCustomize *, DWORD) {
    return S_OK;
  };
  IFACEMETHODIMP OnCheckButtonToggled(IFileDialogCustomize *, DWORD, BOOL) {
    return S_OK;
  };
  IFACEMETHODIMP OnControlActivating(IFileDialogCustomize *, DWORD) {
    return S_OK;
  };

  CDialogEventHandler() : _cRef(1){};

private:
  ~CDialogEventHandler(){};
  long _cRef;
};

HRESULT CDialogEventHandler_CreateInstance(REFIID riid, void **ppv) {
  *ppv = nullptr;
  CDialogEventHandler *pDialogEventHandler =
      new (std::nothrow) CDialogEventHandler();
  HRESULT hr = pDialogEventHandler ? S_OK : E_OUTOFMEMORY;
  if (SUCCEEDED(hr)) {
    hr = pDialogEventHandler->QueryInterface(riid, ppv);
    pDialogEventHandler->Release();
  }
  return hr;
}

DWORD setOpenFileDialogOptions(DWORD dwOptions,
                               const picojson::value &argument) {
  auto props = argument.get("properties").get<int64_t>();
  if (props & OPEN_DIALOG_FOR_OPEN_DIRECTORY) {
    dwOptions |= FOS_PICKFOLDERS;
  }
  if (props & OPEN_DIALOG_WITH_CREATE_DIRECTORY) {
    // ?
  }
  if (props & OPEN_DIALOG_WITH_MULTIPLE_SELECTIONS) {
    dwOptions |= FOS_ALLOWMULTISELECT;
  } else {
    dwOptions &= ~FOS_ALLOWMULTISELECT;
  }
  if (props & OPEN_DIALOG_WITH_SHOW_HIDDEN_FILES) {
    dwOptions |= FOS_FORCESHOWHIDDEN;
  }
  return dwOptions;
}

void setupAllowedFileTypes(IFileDialog *pfd,
                           const picojson::array &args) {
  std::vector<COMDLG_FILTERSPEC> vecFileSpecs;
  std::vector<std::wstring> vecStrs;
  for (auto &filter : args) {
    LOG_DEBUG("[%d] filter = '%s'", __LINE__,  filter.get("name").to_str().c_str());
    auto name =
        exciton::util::ToUTF16String(filter.get("name").get<std::string>());
    auto extensions = filter.get("extensions").get<picojson::array>();
    std::wstring exts;
    for (auto &ext : extensions) {
      if (!exts.empty()) {
        exts += L';';
      }
      exts += L"*." + exciton::util::ToUTF16String(ext.get<std::string>());
    }
    name += L"(" + exts + L")";
    COMDLG_FILTERSPEC spec;
    vecStrs.push_back(name);
    spec.pszName = vecStrs.back().c_str();
    vecStrs.push_back(exts);
    spec.pszSpec = vecStrs.back().c_str();
    vecFileSpecs.push_back(spec);
  }
  pfd->SetFileTypes(vecFileSpecs.size(), vecFileSpecs.data());
}

void createFileDialog(bool forOpen, const picojson::value &argument,
                      int responceNo) {
  IFileDialog *pfd = nullptr;
  IFileDialogEvents *pfde = nullptr;
  HRESULT hr;
  bool succeeded = false;
  bool allowMultipleFiles = false;
  bool openFolder = false;
  HWND hWndParent = NULL;

  if (forOpen) {
    hr = ::CoCreateInstance(CLSID_FileOpenDialog, nullptr, CLSCTX_INPROC_SERVER,
                            IID_PPV_ARGS(&pfd));
  } else {
    hr = ::CoCreateInstance(CLSID_FileSaveDialog, nullptr, CLSCTX_INPROC_SERVER,
                            IID_PPV_ARGS(&pfd));
  }
  if (FAILED(hr)) {
    goto ERROR1;
  }

  hr = CDialogEventHandler_CreateInstance(IID_PPV_ARGS(&pfde));
  if (FAILED(hr))
    goto ERROR2;
  DWORD dwCookie;
  hr = pfd->Advise(pfde, &dwCookie);
  if (FAILED(hr))
    goto ERROR3;
  if (forOpen) {
    DWORD dwFlags;
    hr = pfd->GetOptions(&dwFlags);
    if (FAILED(hr))
      goto ERROR4;
    dwFlags = setOpenFileDialogOptions(dwFlags, argument);
    hr = pfd->SetOptions(dwFlags);
    if (FAILED(hr))
      goto ERROR4;
    if (dwFlags & FOS_ALLOWMULTISELECT) {
      allowMultipleFiles = true;
    }
    if (dwFlags & FOS_PICKFOLDERS) {
      openFolder = true;
    }
  }
  if (argument.contains("title")) {
    auto title =
        exciton::util::ToUTF16String(argument.get("title").get<std::string>());
    if (!title.empty()) {
      pfd->SetTitle(title.c_str());
    }
  }
  if (argument.contains("buttonLabel")) {
    auto label = exciton::util::ToUTF16String(
        argument.get("buttonLabel").get<std::string>());
    if (!label.empty()) {
      pfd->SetOkButtonLabel(label.c_str());
    }
  }
  if (argument.contains("defaultPath")) {
    auto defPath = exciton::util::ToUTF16String(
        argument.get("defaultPath").get<std::string>());
    if (!defPath.empty()) {
      IShellItem *psi = nullptr;
      IShellItem *ppi = nullptr;
      LPWSTR lpItemName = nullptr;
      SHCreateItemFromParsingName(defPath.c_str(), nullptr, IID_PPV_ARGS(&psi));
      psi->GetParent(&ppi);
      psi->GetDisplayName(SIGDN_NORMALDISPLAY, &lpItemName);
      pfd->SetFolder(ppi);
      pfd->SetFileName(lpItemName);
      ::CoTaskMemFree(lpItemName);
      ppi->Release();
      psi->Release();
    }
  }
  if (argument.contains("filters")) {
     setupAllowedFileTypes(pfd, argument.get("filters").get<picojson::array>());
  }

  {
    auto &parentId = argument.get("windowId").get<std::string>();
    hWndParent = resolveParentWindow(parentId);
  }
  hr = pfd->Show(hWndParent);
  if (SUCCEEDED(hr)) {
    picojson::value result;
    if (allowMultipleFiles) {
      IFileOpenDialog *lofd = nullptr;
      picojson::array results;
      hr = pfd->QueryInterface(IID_PPV_ARGS(&lofd));
      if (SUCCEEDED(hr)) {
        IShellItemArray *pItemArray = nullptr;
        hr = lofd->GetResults(&pItemArray);
        if (SUCCEEDED(hr)) {
          DWORD dwItemCount = 0;
          hr = pItemArray->GetCount(&dwItemCount);
          for (DWORD dwIdx = 0; dwIdx < dwItemCount; dwIdx++) {
            IShellItem *psi = nullptr;
            LPWSTR lpszName = nullptr;
            hr = pItemArray->GetItemAt(dwIdx, &psi);
            hr = psi->GetDisplayName(SIGDN_FILESYSPATH, &lpszName);
            auto file =
                exciton::util::ToUTF8String(const_cast<LPCWSTR>(lpszName));
            results.push_back(picojson::value(file));
            ::CoTaskMemFree(lpszName);
            psi->Release();
          }
          result = picojson::value(results);
          succeeded = true;
          pItemArray->Release();
        }
      }
    } else {
      IShellItem *psi = nullptr;
      hr = pfd->GetResult(&psi);
      if (SUCCEEDED(hr)) {
        LPWSTR lpszName = nullptr;
        hr = psi->GetDisplayName(SIGDN_FILESYSPATH, &lpszName);
        auto fname = exciton::util::ToUTF8String(const_cast<LPCWSTR>(lpszName));
        if (forOpen) {
          picojson::array a;
          a.emplace_back(fname);
          result = picojson::value(a);
        } else {
          result = picojson::value(fname);
        }
        ::CoTaskMemFree(lpszName);
        psi->Release();
        succeeded = true;
      }
    }
    if (succeeded) {
      Driver::Current().responceEventResult(responceNo, result);
    }
  }

ERROR4:
  pfd->Unadvise(dwCookie);
ERROR3:
  pfde->Release();
ERROR2:
  pfd->Release();
ERROR1:
  if (!succeeded) {
    // TODO: error response
  }
}

void showOpenDialog(const picojson::value &argument,
                    const std::map<std::string, std::string> &parameter,
                    int responceNo) {
  Driver &d = Driver::Current();
  //  auto idstr = getIdFromParam(parameter);
  createFileDialog(true, argument, responceNo);
}

void showSaveDialog(const picojson::value &argument,
                    const std::map<std::string, std::string> &parameter,
                    int responceNo) {
  Driver &d = Driver::Current();
  auto idstr = getIdFromParam(parameter);
  createFileDialog(false, argument, responceNo);
}
} // namespace

void Dialog_Init() {
  auto &d = Driver::Current();
  d.addDeferEventHandler("/dialog/:id/showMessageBox", showMessageBox);
  d.addDeferEventHandler("/dialog/:id/showOpenDialog", showOpenDialog);
  d.addDeferEventHandler("/dialog/:id/showSaveDialog", showSaveDialog);
}
