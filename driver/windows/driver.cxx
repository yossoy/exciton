#include <windows.h>

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
Driver::HostHolder::~HostHolder() {
  if (ptr_)
    ptr_->Release();
}

namespace {
  Driver s_driver;
}
Driver &Driver::Current() {
  //static Driver driver;
  return s_driver;
}
Driver::Driver() : mainThreadId_(0) { ::InitializeCriticalSection(&m_cs); }
Driver::~Driver() { ::DeleteCriticalSection(&m_cs); }

void Driver::PushDelayProc(DelayProcHandler proc) {
  ::EnterCriticalSection(&m_cs);
  delayProcs_.push_back(proc);
  ::LeaveCriticalSection(&m_cs);
  ::PostThreadMessageW(mainThreadId_, WM_DEFER_CALL, 0, 0);
}

bool Driver::Emit(std::string_view name) {
  ::OutputDebugStringA(std::string(name).c_str());
  ::OutputDebugStringA("\n");
}

bool Driver::Emit(std::string_view name, std::string_view jsonEncodedArgument) {
  ::OutputDebugStringA(std::string(name).c_str());
  ::OutputDebugStringA("\n");
  ::OutputDebugStringA(std::string(jsonEncodedArgument).c_str());
  ::OutputDebugStringA("\n");
}
void Driver::addEventHandler(const std::string &name,
                             NativeEventHandler handler) {
  handlers_[name] = handler;
}
void Driver::addDeferEventHandler(const std::string &name,
                                  NativeEventHandler handler) {
  handlers_[name] = handler;
  //   int id = static_cast<int>(deferId2Name_.size());
  //   deferId2Name_[id] = name;
  //   deferName2Id_[name] = id;
  deferEventNames_.insert(name);
}

int Driver::emitEvent(const void *bytes, int length) {
  std::string jsonStr(reinterpret_cast<const char *>(bytes), length);
  picojson::value jsonValue;
  auto err = picojson::parse(jsonValue, jsonStr);
  if (!err.empty()) {
    LOG_ERROR("emitEvent: error: %s", err.c_str());
    return FALSE;
  }
  auto name = jsonValue.get("name").get<std::string>();
  if (deferEventNames_.find(name) != deferEventNames_.end()) {
    PushDelayProc([this, jsonValue]() { emitEvent(jsonValue); });
    return TRUE;
  }

  return emitEvent(jsonValue);
}

int Driver::emitEvent(const std::string &name, const std::string &argument) {
  LOG_INFO("[%d] Driver::emitEvent(%s, %s)", __LINE__, name.c_str(), argument.c_str());
  std::stringstream ss;
  ss << "{\"name\":\"";
  ss << name;
  ss << "\",\"argument\":";
  ss << argument;
  ss << ",\"respCallbackNo\":-1}";
  std::string json = ss.str();
  requestEventEmit(const_cast<char *>(json.data()), json.length());
  return TRUE;
}

int Driver::emitEvent(const picojson::value &jsonValue) {
  auto name = jsonValue.get("name").get<std::string>();
  auto fiter = handlers_.find(name);
  if (fiter == handlers_.end()) {
    LOG_ERROR("[%d] Driver::emitEvent: event not found[%s]", __LINE__, name.c_str());
    return FALSE;
  }
  auto parameter = jsonValue.get("parameter").get<picojson::object>();
  auto argument = jsonValue.get("argument");
  auto callback = jsonValue.get("respCallbackNo").get<double>();
  std::map<std::string, std::string> mapParam;
  for (auto &kv : parameter) {
    if (kv.second.is<std::string>()) {
      mapParam[kv.first] = kv.second.get<std::string>();
    }
  }
  ((*fiter).second)(argument, mapParam, static_cast<int>(callback));

  return TRUE;
}

void Driver::procDelayEvent(void) {
  std::vector<DelayProcHandler> procs;
  ::EnterCriticalSection(&m_cs);
  procs.swap(delayProcs_);
  ::LeaveCriticalSection(&m_cs);
  for (auto proc : procs) {
    proc();
  }
}

void Driver::responceEventResult(int responceNo, picojson::value result) {
  std::string json = result.serialize();
  ::responceEventResult(responceNo, const_cast<char *>(json.c_str()),
                        json.length());
}
void Driver::responceEventBoolResult(int responceNo, bool result) {
  const char *ret = result ? "true" : "false";
  ::responceEventResult(responceNo, const_cast<char *>(ret), strlen(ret));
}
void Driver::responceEventJsonResult(int responceNo,
                                     const std::string &result) {
  ::responceEventResult(responceNo, const_cast<char *>(result.c_str()),
                        result.length());
}
void Driver::notifyUpdateMenu(const std::string &menuId) {
  std::string id = menuId;
  PushDelayProc([this, id]() {
    std::set<std::shared_ptr<CWebBrowserContainer>> containers;
    for (auto &kv : hosts_) {
      containers.insert(kv.second.ptr_->GetHostContainer());
    }
  });
}

void Driver::Run() {
  hInstance_ = ::GetModuleHandle(nullptr);
  MSG msg;

  mainThreadId_ = ::GetCurrentThreadId();

  ::OleInitialize(NULL);

  while (GetMessage(&msg, NULL, 0, 0) > 0) {
    if (msg.message == WM_QUIT_MESSAGE) {
      ::PostQuitMessage(0);
    }
    if (msg.message == WM_DEFER_CALL) {
      procDelayEvent(/*msg.wParam, msg.lParam*/);
    }
    auto ac = CWebBrowserContainer::GetActiveContainer();
    if (ac && ac->TranslateAccelerator(&msg)) {
      continue;
    }
    auto pActiveContainer = CWebBrowserContainer::GetActiveContainer();
    if (pActiveContainer) {
      if (pActiveContainer->TranslateAccelerator(&msg)) {
        continue;
      }
    }
    TranslateMessage(&msg);
    DispatchMessage(&msg);
  }

  ::OleUninitialize();
}

void Driver::Quit() {
  ::PostThreadMessage(mainThreadId_, WM_QUIT_MESSAGE, 0, 0);
}

void Driver_Run(void) { Driver::Current().Run(); }

void Driver_Terminate(void) { Driver::Current().Quit(); }

int Driver_EmitEvent(void *bytes, int length) {
  return Driver::Current().emitEvent(bytes, length);
}
