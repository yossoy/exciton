'use strict';
const ExcitonEventData = 'exciton-event-data';
const ExcitonComponentID = 'exciton-component-id';

const ditNone = 0;
const ditCreateNode = 1;
const ditCreateNodeWithNS = 2;
const ditCreateTextNode = 3;
const ditSelectCurNode = 4;
const ditSelectArg1Node = 5;
const ditSelectArg2Node = 6;
const ditPropertyValue = 7;
const ditDelProperty = 8;
const ditAttributeValue = 9;
const ditDelAttributeValue = 10;
const ditAddClassList = 11;
const ditDelClassList = 12;
const ditAddDataSet = 13;
const ditDelDataSet = 14;
const ditAddStyle = 15;
const ditDelStyle = 16;
const ditNodeValue = 17;
const ditInnerHTML = 18;
const ditAppendChild = 19;
const ditInsertBefore = 20;
const ditRemoveChild = 21;
const ditReplaceNode = 22;
const ditAddEventListener = 23;
const ditRemoveEventListener = 24;
const ditSetRootItem = 25
const ditNodeUUID = 26;
const ditAddClientEvent = 27;
const ditMountComponent = 28;
const ditUnmountComponent = 29;
const ditAttributeValueNS = 30;
const ditDelAttributeValueNS = 31;


function resolvePathNode(root, items) {
  let ret = root;
  for (let i = 0; i < items.length; i++) {
    ret = ret.childNodes[items[i]];
  }
  return ret;
}
import translateEvent from './events';
import {METHODS} from 'http';

function addEventCallback(nsobj, n, name, itemv) {
  const id = itemv.id;
  const pd = itemv.pd;
  const sp = itemv.sp;
  if (n[ExcitonEventData] === undefined) {
    n[ExcitonEventData] = {};
  }
  const f = function(e) {
    if (pd) {
      e.preventDefault();
    }
    if (sp) {
      e.stopPropagation();
    }
    const goevent = translateEvent(e);
    if (e.target) {
      console.log('dataset:', e.target.dataset);
    }
    console.log('called', goevent);
    nsobj.callNativeMethod('html/' + id + '/' + name, goevent);
  };
  n.addEventListener(name, f);
  n[ExcitonEventData][name] = f;
}

function findExcitonComponent(nsobj, n) {
  while (n) {
    const cid = n[ExcitonComponentID];
    if (cid) {
      return nsobj.components[cid];
    }
    n = n.parentNode;
  }
  return null;
}

function addClientEventCallback(nsobj, n, name, itemv) {
  const id = itemv.id;
  const csp = itemv.sp;
  const shn = itemv.sh;
  const sas = itemv.sas;
  const saslen = (sas) ? sas.length : 0;
  if (n[ExcitonEventData] === undefined) {
    n[ExcitonEventData] = {};
  }
  let f = undefined;
  if (csp === '') {
    const gf = window[shn];
    if (gf && typeof (gf) === 'function') {
      f = gf;
    }
  } else if (csp === '*exciton*') {
    const ef = nsobj[shn];
    if (ef && typeof (ef) === 'function') {
      f = (saslen == 0) ? ef : (e) => ef(e, ...sas);
    }
  } else {
    const mf = nsobj.findModuleFunction(csp, shn);
    if (mf) {
      switch (mf.length) {
        case saslen + 1:  // event only
          f = (saslen == 0) ? mf : (e) => mf(e, ...sas);
          break;
        case saslen + 2:  // component and event
          f = (e) => {
            const cid = findExcitonComponent(nsobj, e.target);
            if (!cid) {
              throw 'invalid event target';
            }
            // TODO: check c.
            return (saslen == 0) ? mf(cid, e) : mf(cid, e, ...sas);
          };
        default:
          break;
      }
    }
  }
  if (f) {
    n.addEventListener(name, f);
    n[ExcitonEventData][name] = f;
  } else {
    throw 'invalid event: ' + csp + ' / ' + name;
  }
}

function mountUnmountComponent(nsobj, n, itemv, mounted) {
  const classId = itemv.classId;
  const id = itemv.id;
  const localJSKey = itemv.localJSKey;
  var instanceData = {
    classId: classId,
    id: id,
    localJSKey: localJSKey,
    callNativeEvent: (method, arg) => {
      nsobj.callnative({
        path: '/components/' + nsobj.ID + '/' + id + method,
        arg: JSON.stringify(arg)
      });
    }
  };
  console.log('locaJSKey = ' + localJSKey);
  if (localJSKey !== '') {
    const m = nsobj.modules[localJSKey];
    if (m) {
      const f =
          mounted ? m.exports['mountComponent'] : m.exports['unmountComponent'];
      if (f && typeof (f) === 'function' && f.length == 2) {
        f(n, instanceData)
      }
    }
  }
  if (mounted) {
    nsobj.components[id] = instanceData;
    n[ExcitonComponentID] = id;
  } else {
    delete nsobj.components[id];
  }
}

function updateDiffData(nsobj, e) {
  const rootObj = document.getElementById(nsobj.ID);
  const diff = e.detail;
  let curNode = null;
  let arg1Node = null;
  let arg2Node = null;
  let creNodes = [];
  for (let ii = 0; ii < diff.items.length; ii++) {
    const item = diff.items[ii];
    switch (item.t) {
      case ditCreateNode:
        curNode = document.createElement(item.v);
        creNodes.push(curNode);
        break;
      case ditCreateNodeWithNS:
        curNode = document.createElementNS(item.k, item.v);
        creNodes.push(curNode);
        break;
      case ditCreateTextNode:
        curNode = document.createTextNode(item.v);
        creNodes.push(curNode);
        break;
      case ditSelectCurNode:
        if (item.v === null || item.v === undefined) {
          curNode = rootObj;
        } else if (typeof (item.v) === 'number') {
          curNode = creNodes[item.v];
        } else {
          curNode = resolvePathNode(rootObj, item.v);
        }
        break;
      case ditSelectArg1Node:
        if (item.v === null || item.v === undefined) {
          arg1Node = rootObj;
        } else if (typeof (item.v) === 'number') {
          arg1Node = creNodes[item.v];
        } else {
          arg1Node = resolvePathNode(rootObj, item.v);
        }
        break;
      case ditSelectArg2Node:
        if (item.v === null || item.v === undefined) {
          arg2Node = rootObj;
        }
        if (typeof (item.v) === 'number') {
          arg2Node = creNodes[item.v];
        } else {
          arg2Node = resolvePathNode(rootObj, item.v);
        }
        break;
      case ditPropertyValue:
        curNode[item.k] = item.v;
        break;
      case ditDelProperty:
        delete curNode[item.v];
        break;
      case ditAttributeValue:
        curNode.setAttribute(item.k, item.v);
        break;
      case ditDelAttributeValue:
        curNode.removeAttribute(item.v);
        break;
      case ditAddClassList:
        curNode.classList.add(item.v);
        break;
      case ditDelClassList:
        curNode.classList.remove(item.v);
        break;
      case ditAddDataSet:
        if (!curNode.dataset) {
          curNode.dataset = {};
        }
        curNode.dataset[item.k] = item.v;
        break;
      case ditDelDataSet:
        delete curNode.dataset[item.v];
        break;
      case ditAddStyle:
        curNode.style.setProperty(item.k, item.v);
        break;
      case ditDelStyle:
        curNode.style.removeProperty(item.v);
        break;
      case ditNodeValue:
        curNode.nodeValue = item.v;
        break;
      case ditInnerHTML:
        curNode.innerHTML = item.v;
        break;
      case ditAppendChild:
        curNode.appendChild(arg1Node);
        break;
      case ditInsertBefore:
        curNode.insertBefore(arg1Node, arg2Node);
        break;
      case ditRemoveChild:
        curNode.removeChild(arg1Node);
        break;
      case ditReplaceNode:
        curNode.replaceChild(arg1Node, arg2Node);
        break;
      case ditAddEventListener:
        addEventCallback(nsobj, curNode, item.k, item.v);
        break;
      case ditRemoveEventListener:
        curNode.removeEventCallback(curNode, item.v);
        break;
      case ditSetRootItem:
        break;
      case ditNodeUUID:
        if (!curNode.dataset) {
          curNode.dataset = {};
        }
        curNode.dataset['excitonId'] = item.v;
        break;
      case ditAddClientEvent:
        addClientEventCallback(nsobj, curNode, item.k, item.v);
        break;
      case ditMountComponent:
        mountUnmountComponent(nsobj, curNode, item.v, true);
        break;
      case ditUnmountComponent:
        mountUnmountComponent(nsobj, null, item.v, false);
        break;
      case ditAttributeValueNS:
        curNode.setAttributeNS(item.v.ns, item.k, item.v.v);
        break;
      case ditDelAttributeValueNS:
        curNode.removeAttributeNS(item.v, item.k);
        break;
      default:
        console.log('invalid type', item.t);
        break;
    }
  }
  console.log('updateDiffData', diff)
}

export default updateDiffData;
