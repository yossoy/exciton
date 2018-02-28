#pragma once
#ifdef __cplusplus
#include <windows.h>
#include <functional>
#include <map>
#include <set>
#include <memory>
#include <string>
#include <string_view>

#include "myjson.h"

class CWebBrowserHost;
class CWebBrowserContainer;
using NativeEventHandler = std::function<void(
    const picojson::value &argument,
    const std::map<std::string, std::string> &parameter, int responceNo)>;
using DelayProcHandler = std::function<void()>;
using ObjectKey = std::string;

class Driver final {
public:
  struct HostHolder {
    CWebBrowserHost *ptr_;
    explicit HostHolder(CWebBrowserHost *p);
    HostHolder();
    ~HostHolder();
  };
private:
  //TODO: tab support
  CRITICAL_SECTION m_cs;
  std::vector<DelayProcHandler> delayProcs_;
//  std::map<ObjectKey, std::shared_ptr<CWebBrowserContainer>> browsers_;
  std::map<std::string, NativeEventHandler> handlers_;
//   std::map<int, std::string> deferId2Name_;
//   std::map<std::string, int> deferName2Id_;
//   std::map<int, picojson::value> deferValues_;
   std::map<ObjectKey, HostHolder> hosts_;
  std::set<std::string> deferEventNames_;
  HINSTANCE hInstance_;
  DWORD mainThreadId_;
  std::wstring productName_;
  std::wstring productVersion_;

public:
  static Driver &Current();
  HINSTANCE InstanceHandle() const { return hInstance_; }
  std::map<ObjectKey, HostHolder>& Hosts() { return hosts_; }
  const std::map<ObjectKey, HostHolder>& Hosts() const { return hosts_; }
  void PushDelayProc(DelayProcHandler proc);
  const std::wstring& GetProductName(void);
  const std::wstring& GetProductVersion(void);

public:
  Driver();
  ~Driver();

public:
  bool Emit(std::string_view name);
  bool Emit(std::string_view name, std::string_view jsonEncodedArgument);
  void addEventHandler(const std::string &name, NativeEventHandler handler);
  void addDeferEventHandler(const std::string &name,
                            NativeEventHandler handler);
  int emitEvent(const std::string &name, const std::string &argument);
  int emitEvent(const void *bytes, int length);
  void responceEventResult(int responceNo, picojson::value result);
  void responceEventBoolResult(int responceNo, bool result);
  void responceEventJsonResult(int responceNo, const std::string& result);
public:
  void notifyUpdateMenu(const std::string& menuId);
public:
  void Run();
  void Quit();

private:
  int emitEvent(const picojson::value &value);
  void procDelayEvent(/*int eventNo, int valueNo*/);
  void initAppVersionInfo(void);
};
extern "C" {
#endif

extern void Driver_Run(void);
extern void Driver_Terminate(void);
extern int Driver_EmitEvent(void *bytes, int length);
extern char* Driver_GetProductName(void);
extern char* Driver_GetProductVersion(void);

#ifdef __cplusplus
};
#endif
