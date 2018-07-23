'use strict';
import styles from './default.css';
const nsobj = window.exciton;
nsobj.system = this;

import ee from 'event-emitter';
import updateDiffData from './diff';

class Exciton {};
ee(Exciton.prototype);

const exciton = new Exciton();

nsobj.callNativeMethod = (method, arg) => {
  nsobj.callnative(
      {path: '/window/' + nsobj.ID + '/' + method, arg: JSON.stringify(arg)});
};

nsobj.requestBrowserEvent = function(method, jsonArg) {
  const arg = JSON.parse(jsonArg);
  exciton.emit(method, {detail: arg});
};

nsobj.findModuleFunction = function(localJSKey, funcName) {
  if (localJSKey !== '') {
    const m = this.modules[localJSKey];
    if (m) {
      const mf = m.exports[funcName];
      if (mf && typeof (mf) === 'function') {
        return mf;
      }
    }
  }
  return null;
}

window.addEventListener('popstate', function(e) {
  const s = e.state;
  nsobj.callNativeMethod('changeRoute', {'route': s.redirectRoute});
});

nsobj.redirectTo = function(route) {
  window.history.pushState({'redirectRoute': route}, route, route);
  nsobj.callNativeMethod('changeRoute', {'route': route});
};

nsobj.onClickRedirectTo = function(e, route) {
  nsobj.redirectTo(route);
  e.preventDefault();
  e.stopPropagation();
  return false;
};

nsobj.requestBrowerEventSync = function(method, jsonArg) {
  const arg = JSON.parse(jsonArg);
  console.log('requestBrowerEventSync', arg);
  switch (arg.cmd) {
    case 'getProp':
      const elem =
          document.querySelector('[data-exciton-id="' + arg.target + '"]');
      if (elem == null) {
        throw new Error('invalid target: ' + arg.target);
      }
      return JSON.stringify(elem[arg.argument]);
    case 'callClientFunction':
      const funcName = arg.argument.funcName;
      const args = arg.argument.arguments;
      let f = null;
      if (arg.target == null) {
        f = window[funcName];
      } else {
        const mf = nsobj.findModuleFunction(arg.target.localJSKey, funcName);
        if (mf != null) {
          const cid = nsobj.components[arg.target.id];
          if (cid != null && mf.length >= 1) {
            f = (...args) => {
              return mf(cid, ...args);
            };
          } else {
            mf = f;
          }
        }
      }
      if (!f) {
        throw new Error('function not found: ' + funcName);
      }
      const r = f(...args);
      return JSON.stringify(r);
    default:
      throw new Error('invalid command: ' + arg.cmd);
  }
};

exciton.on('requestAnimationFrame', function(e) {
  const timestamp = e.detail;
  window.requestAnimationFrame(function(timestamp) {
    nsobj.callNativeMethod('onRequestAnimationFrame', timestamp);
  });
});

exciton.on('updateDiffData', (e) => {
  updateDiffData(nsobj, e);
});

exciton.on('redirectTo', (e) => {
  nsobj.redirectTo(e.detail);
});

class Module {
  constructor(id, w) {
    this.id = id;
    this.exports = {};
    this.loaded = false;
    this.wrapper = w;
  }

  require(id) {
    if (id == 'exciton') {
      // special case
      return {id: this.id, exports: {id: this.id}};
    }
    return nsobj.require(id)
  }
};

nsobj.registerModule = function(id, w) {
  if (id in nsobj.modules) {
    throw 'multiple register component: ' + id;
  }
  nsobj.modules[id] = new Module(id, w);
};

nsobj.require = function(id) {
  const m = nsobj.modules[id];
  if (m) {
    if (m.loaded) {
      return m.exports;
    }
    m.wrapper(m.exports, (id) => {
      return m.require(id);
    }, m, '', '');
    m.loaded = true;
    return m.exports;
  }
  throw 'not contains key: ' + id;
};

const loadComponentsScripts = function() {
  for (let id in nsobj.modules) {
    nsobj.require(id);
  }
  if (!window.history.state) {
    window.history.replaceState({'redirectRoute': '/'}, '/', '/');
  }
  window.removeEventListener('load', loadComponentsScripts, false);
  nsobj.callNativeMethod('ready', null);
};

window.addEventListener('load', loadComponentsScripts, false);
