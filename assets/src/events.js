'use strict';
const tEvent = {
  super: null,
  props: [
    'bubbles', 'cancelBubble', 'cancelable', 'composed', '@currentTarget',
    'defaultPrevented', 'eventPhase', '@target', 'timeStamp', 'type',
    'isTrusted'
  ]
};
const tUIEvent = {
  super: tEvent,
  props: ['detail', '@view']
};
const tBeforeUnloadEvent = {
  super: tEvent,
  props: []
};
const tMouseEvent = {
  super: tUIEvent,
  props: [
    'altKey', 'button', 'buttons', 'clientX', 'clientY', 'ctrlKey', 'metaKey',
    'movementX', 'movementY', 'region', '@relatedTarget', 'screenX', 'screenY',
    'shiftKey'
  ]
};
const tPopStateEvent = {
  super: tEvent,
  props: []
};
const tWheelEvent = {
  super: tMouseEvent,
  props: ['deltaX', 'deltaY', 'deltaZ', 'deltaMode']
};
const tPageTransitionEvent = {
  super: tEvent,
  props: ['persisted']
};
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
const tFocusEvent = {
  super: tUIEvent,
  props: ['@relatedTarget']
};
const tCompositionEvent = {
  super: tUIEvent,
  props: ['data', 'locale']
};
const tDragEvent = {
  super: tMouseEvent,
  props: [/*'dataTransfer'*/]
};
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

function translateEvent(e) {
  // TODO: autogenerate on build time or runtime
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
            if (p.charCodeAt(0) == 64) {
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
export default translateEvent;
