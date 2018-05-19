(function(nsobj) {
'use strict';
nsobj.system = this;

const ExcitonEventData = 'exciton-event-data';

// event emitter
function EventEmitter() {
  const target = document.createDocumentFragment();
  function delegate(method) {
    this[method] = target[method].bind(target);
  }
  [
    "addEventListener",
    "dispatchEvent",
    "removeEventListener"
  ].forEach(delegate, this);
}

function translateEvent(e) {
  // TODO: autogenerate on build time or runtime
  const tEvent = {
    super: null,
    props: [
      'bubbles', 'cancelBubble', 'cancelable', 'composed', '@currentTarget',
      'defaultPrevented', 'eventPhase', '@target', 'timeStamp', 'type',
      'isTrusted'
    ]
  };
  const tUIEvent = {super: tEvent, props: ['detail', '@view']};
  const tBeforeUnloadEvent = {super: tEvent, props: []};
  const tMouseEvent = {
    super: tUIEvent,
    props: [
      'altKey', 'button', 'buttons', 'clientX', 'clientY', 'ctrlKey', 'metaKey',
      'movementX', 'movementY', 'region', '@relatedTarget', 'screenX',
      'screenY', 'shiftKey'
    ]
  };
  const tPopStateEvent = {super: tEvent, props: []};
  const tWheelEvent = {
    super: tMouseEvent,
    props: ['deltaX', 'deltaY', 'deltaZ', 'deltaMode']
  };
  const tPageTransitionEvent = {super: tEvent, props: ['persisted']};
  const tProgressEvent = {
    super: tEvent,
    props: ['lengthComputable', 'loaded', 'total']
  };
  const tKeyboardEvent = {
    super: tUIEvent,
    props: [
      'altKey', 'code', 'ctrlKey', 'isComposing', 'key', 'locale', 'location',
      'metaKey', 'repeat', 'shiftKey'
    ]
  };
  const tFocusEvent = {super: tUIEvent, props: ['@relatedTarget']};
  const tCompositionEvent = {super: tUIEvent, props: ['data', 'locale']};
  const tDragEvent = {super: tMouseEvent, props: [/*'dataTransfer'*/]};
  const tHashChangeEvent = {
    super: null,
    props: ['@target', 'type', 'bubbles', 'cancelable', 'oldURL', 'newURL']
  };
  const tOfflineAudioCompletionEvent = {
    super: tEvent,
    props: [/*'renderedBuffer'*/]
  };

  const nEvents = [
    {
      type: tEvent,
      events: [
        'afterprint',     'beforeprint',
        'canplay',        'canplaythrough',
        'durationchange', 'languagechange',
        'loadeddata',     'loadedmetadata',
        'noupdate',       'cached',
        'change',         'checking',
        'reset',          'DOMContentLoaded',
        'downloading',    'emptied',
        'ended',          'error',
        'input',          'invalid',
        'obsolete',       'offline',
        'online',         'pause',
        'play',           'playing',
        'seeked',         'seeking',
        'stalled',        'submit',
        'suspend',        'waiting',
        'ratechange',     'readystatechange',
        'selectstart',    'selectionchange',
        'timeupdate',     'updateready',
        'volumechange'
      ]
    },
    {
      type: tMouseEvent,
      events: [
        'contextmenu', 'dblclick', 'mousedown', 'mouseenter', 'mouseleave',
        'mousemove', 'mouseout', 'mouseover', 'mouseup', 'click', 'show'
      ]
    },
    {
      type: tWheelEvent,
      events: ['wheel']

    },
    {
      type: tKeyboardEvent,
      events: ['keydown', 'keypress', 'keyup']

    },
    {
      type: tCompositionEvent,
      events: ['compositionend', 'compositionstart', 'compositionupdate']
    },
    {type: tFocusEvent, events: ['focusin', 'focusout', 'blur', 'focus']}, {
      type: tProgressEvent,
      events: ['progress']

    },
    {type: tBeforeUnloadEvent, events: ['beforeunload']}, {
      type: tUIEvent,
      events: ['abort', 'load', 'resize', 'scroll', 'select', 'unload']
    },
    {
      type: tDragEvent,
      events: [
        'dragend', 'dragenter', 'dragleave', 'dragover', 'dragstart', 'drag',
        'drop'
      ]
    },
    {type: tHashChangeEvent, events: ['hashchange']},
    {type: tOfflineAudioCompletionEvent, events: ['complete']},
    {type: tPageTransitionEvent, events: ['pagehide', 'pageshow']},
    {type: tPopStateEvent, events: ['popstate']}
  ];

  const eventType = e.type;
  for (let nei = 0; nei < nEvents.length; nei++) {
    const ne = nEvents[nei];
    for (let eti = 0; eti < ne.events.length; eti++) {
      const et = ne.events[eti];
      if (et == eventType) {
        const result = {};
        let eot = ne.type;
        while (eot != null) {
          for (let pi = 0; pi < eot.props.length; pi++) {
            const p = eot.props[pi];
            if (p.startsWith('@')) {
              const rp = p.substr(1);
              const pv = e[rp];
              if ((typeof pv) !== 'undefined' && pv != null) {
                // console.log(rp, pv);
                if (rp === 'view') {
                  result[rp] = { windowId: pv.exciton.ID }
                } else {
                  result[rp] = {
                    windowId: pv.ownerDocument.defaultView.exciton.ID,
                    elementId: pv.dataset['excitonId']
                  };
                }
                // console.log("==>", result[p]);
              }
            } else {
              const pv = e[p];
              if ((typeof pv) !== 'undefined') {
                result[p] = pv;
              }
            }
          }
          eot = eot.super;
        }
        return result;
      }
    }
  }
  return {};
}

function callNativeMethod(method, arg) {
  nsobj.callnative({path: '/window/' + nsobj.ID + '/' + method, arg: JSON.stringify(arg)});
}

const exciton = new EventEmitter();

exciton.addEventListener('requestAnimationFrame', function(e) {
  const timestamp = e.detail;
  window.requestAnimationFrame(function(timestamp) {
    callNativeMethod('onRequestAnimationFrame', timestamp);
  });
}, false);

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
const ditSetRootItem = 25;
const ditNodeUUID = 26;

function resolvePathNode(root, items) {
  let ret = root;
  for (let i = 0; i < items.length; i++) {
    ret = ret.childNodes[items[i]];
  }
  return ret;
}

function addEventCallback(n, name, itemv) {
  const id = itemv.id;
  const pd = itemv.pd;
  const sp = itemv.sp;
  if (n[ExcitonEventData] === undefined) {
    n[ExcitonEventData] = {}
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
    callNativeMethod('html/' + id + '/' + name, goevent);
  };
  n.addEventListener(name, f);
  n[ExcitonEventData][name] = f;
}

function removeEventCallback(n, name) {
  const m = n[ExcitonEventData];
  if (m !== undefined) {
    const f = m[name];
    delete m[name];
    n.removeEventListenr(name, f);
  }
}

exciton.addEventListener('updateDiffData', function(e) {
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
        } else if (typeof(item.v) === 'number') {
          curNode = creNodes[item.v];
        } else {
          curNode = resolvePathNode(rootObj, item.v);
        }
        break;
      case ditSelectArg1Node:
        if (item.v === null || item.v === undefined) {
          arg1Node = rootObj;
        } else if (typeof(item.v) === 'number') {
          arg1Node = creNodes[item.v];
        } else {
          arg1Node = resolvePathNode(rootObj, item.v);
        }
        break;
      case ditSelectArg2Node:
        if (item.v === null || item.v === undefined) {
          arg2Node = rootObj;
        } if (typeof(item.v) === 'number') {
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
        addEventCallback(curNode, item.k, item.v);
        break;
      case ditRemoveEventListener:
        curNode.removeEventCallback(curNode, item.v);
        break;
      case ditSetRootItem:
        // TODO:
        //document.body = curNode;
        break;
      case ditNodeUUID:
        if (!curNode.dataset) {
          curNode.dataset = {};
        }
        curNode.dataset['excitonId'] = item.v;
        break;
      default:
        console.log('invalid type', item.t);
        break;
    }
  }
  console.log('updateDiffData', diff)
}, false);

nsobj.requestBrowserEvent = function(method, jsonArg) {
  // console.log('jsonArg', jsonArg);
  const arg = JSON.parse(jsonArg);
  const e = new CustomEvent(method, {detail: arg});
  exciton.dispatchEvent(e);
};

nsobj.requestBrowerEventSync = function(method, jsonArg) {
  const arg = JSON.parse(jsonArg);
  console.log('requestBrowerEventSync', arg);
  switch (arg.cmd) {
    case 'getProp':
      const elem =
          document.querySelector('[data-exciton-id="' + arg.elemId + '"]');
      if (elem == null) {
        throw new Error('invalid target: ' + arg.elemId);
      }
      return JSON.stringify(elem[arg.propName]);
    default:
      throw new Error('invalid command: ' + arg.cmd);
  }
};

callNativeMethod('ready', null);
})(window.exciton)