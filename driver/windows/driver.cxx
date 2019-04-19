#include <windows.h>

#include <shlwapi.h>
#include <winnls.h>

#include <set>
#include <sstream>

#include "_cgo_export.h"

#include "browsercontainer.h"
#include "browserhost.h"

#include "driver.h"
#include "log.h"

#define WM_QUIT_MESSAGE (WM_USER + 100)
#define WM_DEFER_CALL (WM_USER + 101)

Driver::HostHolder::HostHolder(CWebBrowserHost *p) : ptr_(p) {}
Driver::HostHolder::HostHolder(void) : ptr_(nullptr) {}
Driver::HostHolder::~HostHolder()
{
  if (ptr_)
    ptr_->Release();
}

Driver &Driver::Current()
{
  static Driver s_driver;
  return s_driver;
}

Driver::Driver() : mainThreadId_(0) { ::InitializeCriticalSection(&m_cs); }
Driver::~Driver() { ::DeleteCriticalSection(&m_cs); }

void Driver::PushDelayProc(DelayProcHandler proc)
{
  ::EnterCriticalSection(&m_cs);
  delayProcs_.push_back(proc);
  ::LeaveCriticalSection(&m_cs);
  ::PostThreadMessageW(mainThreadId_, WM_DEFER_CALL, 0, 0);
}

bool Driver::Emit(std::string_view name)
{
  LOG_DEBUG("Driver::Emit(%s)", std::string(name).c_str());
  return false;
}

bool Driver::Emit(std::string_view name, std::string_view jsonEncodedArgument)
{
  LOG_DEBUG("Driver::Emit(%s, %s)", std::string(name).c_str(),
            std::string(jsonEncodedArgument).c_str());
  return false;
}
void Driver::addEventHandler(const std::string& path, const std::string &name,
                             NativeEventHandler handler)
{
  handlers_[path][name] = handler;
}
void Driver::addDeferEventHandler(const std::string& path, const std::string &name,
                                  NativeEventHandler handler)
{
  handlers_[path][name] = handler;
  deferEventNames_[path].insert(name);
}

int Driver::emitEvent(const void *bytes, int length)
{
  std::string jsonStr(reinterpret_cast<const char *>(bytes), length);
  picojson::value jsonValue;
  auto err = picojson::parse(jsonValue, jsonStr);
  if (!err.empty())
  {
    LOG_ERROR("emitEvent: error: %s", err.c_str());
    return FALSE;
  }
  auto path = jsonValue.get("target").get<std::string>();
  auto name = jsonValue.get("name").get<std::string>();
  auto fpath = deferEventNames_.find(path);
  if (fpath != deferEventNames_.end()) {
    auto fname = (*fpath).second.find(name);
    if (fname != (*fpath).second.end()) {
      PushDelayProc([this, jsonValue]() { emitEvent(jsonValue); });
      return TRUE;
    }
  }
  return emitEvent(jsonValue);
}

int Driver::responseEvent(int respNo, const void *bytes, int length)
{
  std::string jsonStr(reinterpret_cast<const char *>(bytes), length);
  picojson::value jsonValue;
  auto err = picojson::parse(jsonValue, jsonStr);
  ::EnterCriticalSection(&m_cs);
  NativeResponseHandler callback;
  {
    auto fresponse = responses_.find(respNo);
    if (fresponse != responses_.end()) {
      callback = (*fresponse).second;
      responses_.erase(fresponse);
    }
  }
  ::LeaveCriticalSection(&m_cs);
  if (!callback) {
    LOG_ERROR("responseEvent: handler not found:[%d]\n", respNo);
    return FALSE;
  }
  auto result = jsonValue.get("result");
  auto error = jsonValue.get("error").get<std::string>();
  callback(result, error);
  return TRUE;
}

int Driver::emitEvent(const std::string& target, const std::string &name, const std::string &argument)
{
  LOG_INFO("[%d] Driver::emitEvent(%s, %s)", __LINE__, name.c_str(),
           argument.c_str());
  std::stringstream ss;
  ss << "{\"target\":\"";
  ss << target,
  ss << "\",\"name\":\"";
  ss << name;
  ss << "\",\"argument\":";
  ss << argument;
  ss << ",\"respCallbackNo\":-1}";
  std::string json = ss.str();
  requestEventEmit(const_cast<char *>(json.data()), json.length());
  return TRUE;
}

int Driver::emitEvent(const picojson::value &jsonValue)
{
  auto target = jsonValue.get("target").get<std::string>();
  auto name = jsonValue.get("name").get<std::string>();
  auto ftarget = handlers_.find(target);
  if (ftarget == handlers_.end()) {
    LOG_ERROR("[%d] Driver::emitEvent: event not found[%s/%s]", __LINE__,
              target.c_str(), name.c_str());
    return FALSE;

  }
  auto fiter = (*ftarget).second.find(name);
  if (fiter == (*ftarget).second.end())
  {
    LOG_ERROR("[%d] Driver::emitEvent: event not found[%s/%s]", __LINE__,
              target.c_str(), name.c_str());
    return FALSE;
  }
  auto parameter = jsonValue.get("parameter").get<picojson::object>();
  auto argument = jsonValue.get("argument");
  auto callback = jsonValue.get("respCallbackNo").get<double>();
  std::map<std::string, std::string> mapParam;
  for (auto &kv : parameter)
  {
    if (kv.second.is<std::string>())
    {
      mapParam[kv.first] = kv.second.get<std::string>();
    }
  }
  ((*fiter).second)(argument, mapParam, static_cast<int>(callback));

  return TRUE;
}

void Driver::procDelayEvent(void)
{
  std::vector<DelayProcHandler> procs;
  ::EnterCriticalSection(&m_cs);
  procs.swap(delayProcs_);
  ::LeaveCriticalSection(&m_cs);
  for (auto proc : procs)
  {
    proc();
  }
}

void Driver::responceEventResult(int responceNo, picojson::value result)
{
  std::string json = result.serialize();
  ::responceEventResult(responceNo, const_cast<char *>(json.c_str()),
                        json.length());
}
void Driver::responceEventBoolResult(int responceNo, bool result)
{
  const char *ret = result ? "true" : "false";
  ::responceEventResult(responceNo, const_cast<char *>(ret), strlen(ret));
}
void Driver::responceEventJsonResult(int responceNo,
                                     const std::string &result)
{
  ::responceEventResult(responceNo, const_cast<char *>(result.c_str()),
                        result.length());
}
void Driver::notifyUpdateMenu(const std::string &menuId)
{
  std::string id = menuId;
  PushDelayProc([this, id]() {
    std::set<std::shared_ptr<CWebBrowserContainer>> containers;
    for (auto &kv : hosts_)
    {
      containers.insert(kv.second.ptr_->GetHostContainer());
    }
  });
}

void Driver::Run()
{
  hInstance_ = ::GetModuleHandle(nullptr);
  MSG msg;

  mainThreadId_ = ::GetCurrentThreadId();

  ::OleInitialize(NULL);

  while (GetMessage(&msg, NULL, 0, 0) > 0)
  {
    if (msg.message == WM_QUIT_MESSAGE)
    {
      ::PostQuitMessage(0);
    }
    if (msg.message == WM_DEFER_CALL)
    {
      procDelayEvent(/*msg.wParam, msg.lParam*/);
    }
    auto ac = CWebBrowserContainer::GetActiveContainer();
    if (ac && ac->TranslateAccelerator(&msg))
    {
      continue;
    }
    auto pActiveContainer = CWebBrowserContainer::GetActiveContainer();
    if (pActiveContainer)
    {
      if (pActiveContainer->TranslateAccelerator(&msg))
      {
        continue;
      }
    }
    TranslateMessage(&msg);
    DispatchMessage(&msg);
  }

  ::OleUninitialize();
}

void Driver::Quit(void)
{
  ::PostThreadMessage(mainThreadId_, WM_QUIT_MESSAGE, 0, 0);
}

void Driver::initAppVersionInfo()
{
  // TODO: cache data?
  WCHAR strFilePath[MAX_PATH];
  GetModuleFileNameW(hInstance_, strFilePath, MAX_PATH);
  DWORD dwDummy = 0;
  DWORD dwVersionInfoSize = ::GetFileVersionInfoSizeW(strFilePath, &dwDummy);
  if (dwVersionInfoSize > 0)
  {
    std::unique_ptr<unsigned char[]> pVersionInfos(
        new unsigned char[dwVersionInfoSize]);
    if (::GetFileVersionInfoW(strFilePath, 0, dwVersionInfoSize,
                              pVersionInfos.get()))
    {
      LPVOID pvVersion;
      UINT uVersionLen;
      LPVOID pvName;
      UINT uNameLen;
      if (VerQueryValueW(pVersionInfos.get(),
                         L"\\StringFileInfo\\000004b0\\ProductVersion",
                         &pvVersion, &uVersionLen))
      {
        productVersion_ = std::wstring((const WCHAR *)pvVersion, uVersionLen);
      }
      if (VerQueryValueW(pVersionInfos.get(),
                         L"\\StringFileInfo\\000004b0\\ProductName", &pvVersion,
                         &uVersionLen))
      {
        productName_ = std::wstring((const WCHAR *)pvVersion, uVersionLen);
      }
    }
  }
  if (productVersion_.empty())
  {
    productVersion_ = L"1.0.0.0";
  }
  if (productName_.empty())
  {
    productName_ = ::PathFindFileNameW(strFilePath);
  }
}

const std::wstring &Driver::GetProductName(void)
{
  if (productName_.empty())
  {
    initAppVersionInfo();
  }
  return productName_;
}
const std::wstring &Driver::GetProductVersion(void)
{
  if (productVersion_.empty())
  {
    initAppVersionInfo();
  }
  return productVersion_;
}

char *Driver_GetProductName(void)
{
  auto &appName = Driver::Current().GetProductName();
  auto utf8AppName = exciton::util::ToUTF8String(appName.c_str());
  return ::strdup(utf8AppName.c_str());
}
char *Driver_GetProductVersion(void)
{
  auto &appVersion = Driver::Current().GetProductVersion();
  auto utf8AppVersion = exciton::util::ToUTF8String(appVersion.c_str());
  return ::strdup(utf8AppVersion.c_str());
}

void Driver_Run(void) { Driver::Current().Run(); }

void Driver_Terminate(void) { Driver::Current().Quit(); }

int Driver_EmitEvent(void *bytes, int length)
{
  LOG_DEBUG("Driver_EventEmit(%p, %d)\n", bytes, length);
  return Driver::Current().emitEvent(bytes, length);
}

int Driver_ResponseEvent(int respNo, void* bytes, int length)
{
  LOG_DEBUG("Driver_ResponseEvent(%d, %p, %d)\n", respNo, bytes, length);
  return Driver::Current().responseEvent(respNo, bytes, length);
}


struct ResFileItem Driver_GetResFile(int resNo)
{
  HINSTANCE hInst = Driver::Current().InstanceHandle();
  HRSRC rsc = ::FindResource(hInst, MAKEINTRESOURCE(resNo), RT_RCDATA);
  ResFileItem ret;
  ret.ptr = NULL;
  ret.size = 0;
  if (!rsc)
  {
    return ret;
  }
  HGLOBAL rh = ::LoadResource(NULL, rsc);
  if (!rh)
  {
    return ret;
  }
  ret.ptr = LockResource(rh);
  ret.size = SizeofResource(NULL, rsc);
  return ret;
}

const char* Driver_GetPreferrdLanguage()
{
  ULONG ulNumLanguages = 0UL;
  ULONG cchLanguagesBuffer = 0UL;
  if (!GetUserPreferredUILanguages(MUI_LANGUAGE_NAME, &ulNumLanguages, NULL, &cchLanguagesBuffer)) {
    return NULL;
  }
  if (0 == cchLanguagesBuffer) {
    return NULL;
  }
  std::unique_ptr<WCHAR[]> pwszLanguagesBuffer(new WCHAR[cchLanguagesBuffer]);
  if (!GetUserPreferredUILanguages(MUI_LANGUAGE_NAME, &ulNumLanguages, pwszLanguagesBuffer.get(), &cchLanguagesBuffer)) {
    return NULL;
  }
  for (ULONG i = 0; i < (cchLanguagesBuffer - 1); i++) {
    WCHAR* p = pwszLanguagesBuffer.get() + i;
    if (*p == 0) {
      *p = TEXT(';');
    }
  }
  auto utf8Languages = exciton::util::ToUTF8String(pwszLanguagesBuffer.get());
  return strdup(utf8Languages.c_str());
}
