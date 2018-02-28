#include <windows.h>

#include <mshtmdid.h>
#include <mshtmhst.h>
#include <shlobj.h>
#include <shlwapi.h>

#include "browsercontainer.h"
#include "browserhost.h"
#include "driver.h"
#include "log.h"
#include "menumgr.h"
#include "myjson.h"
#include "util.h"
#include "window.h"

namespace {
std::string
getIdFromParam(const std::map<std::string, std::string> &parameter) {
  auto fiter = parameter.find("id");
  if (fiter == parameter.end()) {
    return "";
  }
  return (*fiter).second;
}

void evaluateJavaScript(const std::string &id, const std::string &funcName,
                        const picojson::value &argument, int responceNo = -1) {
  std::string json = argument.serialize();
  std::wstring wFuncName = exciton::util::ToUTF16String(funcName);
  std::wstring jsonArg = exciton::util::ToUTF16String(json);
  Driver::Current().PushDelayProc([=]() {
    Driver &d = Driver::Current();
    auto fiter = d.Hosts().find(id);
    if (fiter == d.Hosts().end()) {
      LOG_ERROR("requestAnimationFrame: window not found!: %s\n", id.c_str());
      return;
    }
    CWebBrowserHost *p = (*fiter).second.ptr_;
    p->AddRef();
    VARIANT vResult;
    VARIANT *pvResult = nullptr;
    if (0 <= responceNo) {
      ::VariantInit(&vResult);
      pvResult = &vResult;
    }
    p->EvaluateJavasScript(wFuncName, jsonArg, pvResult);
    p->Release();
    if (0 <= responceNo) {
      std::string retValue;
      if (vResult.vt == VT_BSTR) {
        retValue = exciton::util::ToUTF8String(vResult.bstrVal);
      } else {
        LOG_ERROR("requestAnimationFrame: invalid result type: %d\n",
                  vResult.vt);
        retValue = "undefined";
      }
      ::VariantClear(&vResult);
      d.responceEventJsonResult(responceNo, retValue);
    }
  });
}

void newWindow(const picojson::value &argument,
               const std::map<std::string, std::string> &parameter,
               int responceNo) {
  Driver &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR("parameter['id'] not found\n");
    Driver::Current().responceEventBoolResult(responceNo, false);
    return;
  }
  LOG_INFO("newWindow called: %s\n", id.c_str());
  auto container = std::make_shared<CWebBrowserContainer>();
  std::string html = argument.get("html").get<std::string>();
  auto width = argument.get("size").get("width").get<int64_t>();
  auto height = argument.get("size").get("height").get<int64_t>();
  if (!container->NewWindow(d.InstanceHandle(), width, height)) {
    LOG_ERROR("Container: newWindow failed\n");
    Driver::Current().responceEventBoolResult(responceNo, false);
    return;
  }
  d.PushDelayProc([=]() {
    Driver &d = Driver::Current();
    // TODO: tab browse
    auto p = container->NewBrowser(html);
    d.Hosts()[id] = Driver::HostHolder(p);
    p->Release();
    d.responceEventBoolResult(responceNo, true);
  });
}

void requestAnimationFrame(const picojson::value &argument,
                           const std::map<std::string, std::string> &parameter,
                           int responceNo) {
  Driver &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR("parameter['id'] not found\n");
    return;
  }
  LOG_INFO("requestAnimationFrame called: %s\n", id.c_str());
  evaluateJavaScript(id, "requestAnimationFrame", argument);
}

void updateDiffSetHandler(const picojson::value &argument,
                          const std::map<std::string, std::string> &parameter,
                          int responceNo) {
  Driver &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR("parameter['id'] not found\n");
    return;
  }
  LOG_INFO("updateDiffSetHandler called: %s\n", id.c_str());
  evaluateJavaScript(id, "updateDiffData", argument);
}

void browserSync(const picojson::value &argument,
                 const std::map<std::string, std::string> &parameter,
                 int responceNo) {
  Driver &d = Driver::Current();
  auto id = getIdFromParam(parameter);
  if (id.empty()) {
    LOG_ERROR("parameter['id'] not found\n");
    return;
  }
  LOG_INFO("updateDiffSetHandler called: %s\n", id.c_str());
  evaluateJavaScript(id, "updateDiffData", argument, responceNo);
}
} // namespace

void Window_Init() {
  auto &d = Driver::Current();
  d.addDeferEventHandler("/window/:id/new", newWindow);
  d.addDeferEventHandler("/window/:id/requestAnimationFrame",
                         requestAnimationFrame);
  d.addDeferEventHandler("/window/:id/updateDiffSetHandler",
                         updateDiffSetHandler);
  d.addDeferEventHandler("/window/:id/browserSync", browserSync);

  if (!CMenuMgr::InitClass()) {
    LOG_ERROR("menumgr::initclass failed\n");
  }
}
