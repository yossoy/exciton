package web

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _assetsData4029148de0a3b24c533ab030ac54e17cbacc31b7 = "<!-- Generated by IcoMoon.io -->\n<svg version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" width=\"32\" height=\"32\" viewBox=\"0 0 32 32\">\n<title>question</title>\n<path d=\"M14 22h4v4h-4zM22 8c1.105 0 2 0.895 2 2v6l-6 4h-4v-2l6-4v-2h-10v-4h12zM16 3c-3.472 0-6.737 1.352-9.192 3.808s-3.808 5.72-3.808 9.192c0 3.472 1.352 6.737 3.808 9.192s5.72 3.808 9.192 3.808c3.472 0 6.737-1.352 9.192-3.808s3.808-5.72 3.808-9.192c0-3.472-1.352-6.737-3.808-9.192s-5.72-3.808-9.192-3.808zM16 0v0c8.837 0 16 7.163 16 16s-7.163 16-16 16c-8.837 0-16-7.163-16-16s7.163-16 16-16z\"></path>\n</svg>\n"
var _assetsData811b89eace18a3615f392e27d8492a6c890bdd96 = "!function(e){var t={};function n(o){if(t[o])return t[o].exports;var i=t[o]={i:o,l:!1,exports:{}};return e[o].call(i.exports,i,i.exports,n),i.l=!0,i.exports}n.m=e,n.c=t,n.d=function(e,t,o){n.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:o})},n.r=function(e){\"undefined\"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:\"Module\"}),Object.defineProperty(e,\"__esModule\",{value:!0})},n.t=function(e,t){if(1&t&&(e=n(e)),8&t)return e;if(4&t&&\"object\"==typeof e&&e&&e.__esModule)return e;var o=Object.create(null);if(n.r(o),Object.defineProperty(o,\"default\",{enumerable:!0,value:e}),2&t&&\"string\"!=typeof e)for(var i in e)n.d(o,i,function(t){return e[t]}.bind(null,i));return o},n.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return n.d(t,\"a\",t),t},n.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},n.p=\"\",n(n.s=7)}([function(e,t,n){var o;!function(){var i=window.CustomEvent;function a(e){for(;e;){if(\"dialog\"===e.localName)return e;e=e.parentElement}return null}function l(e){e&&e.blur&&e!==document.body&&e.blur()}function s(e,t){for(var n=0;n<e.length;++n)if(e[n]===t)return!0;return!1}function r(e){return!(!e||!e.hasAttribute(\"method\"))&&\"dialog\"===e.getAttribute(\"method\").toLowerCase()}function d(e){if(this.dialog_=e,this.replacedStyleTop_=!1,this.openAsModal_=!1,e.hasAttribute(\"role\")||e.setAttribute(\"role\",\"dialog\"),e.show=this.show.bind(this),e.showModal=this.showModal.bind(this),e.close=this.close.bind(this),\"returnValue\"in e||(e.returnValue=\"\"),\"MutationObserver\"in window){new MutationObserver(this.maybeHideModal.bind(this)).observe(e,{attributes:!0,attributeFilter:[\"open\"]})}else{var t,n=!1,o=function(){n?this.downgradeModal():this.maybeHideModal(),n=!1}.bind(this),i=function(i){if(i.target===e){var a=\"DOMNodeRemoved\";n|=i.type.substr(0,a.length)===a,window.clearTimeout(t),t=window.setTimeout(o,0)}};[\"DOMAttrModified\",\"DOMNodeRemoved\",\"DOMNodeRemovedFromDocument\"].forEach(function(t){e.addEventListener(t,i)})}Object.defineProperty(e,\"open\",{set:this.setOpen.bind(this),get:e.hasAttribute.bind(e,\"open\")}),this.backdrop_=document.createElement(\"div\"),this.backdrop_.className=\"backdrop\",this.backdrop_.addEventListener(\"click\",this.backdropClick_.bind(this))}i&&\"object\"!=typeof i||((i=function(e,t){t=t||{};var n=document.createEvent(\"CustomEvent\");return n.initCustomEvent(e,!!t.bubbles,!!t.cancelable,t.detail||null),n}).prototype=window.Event.prototype),d.prototype={get dialog(){return this.dialog_},maybeHideModal:function(){this.dialog_.hasAttribute(\"open\")&&document.body.contains(this.dialog_)||this.downgradeModal()},downgradeModal:function(){this.openAsModal_&&(this.openAsModal_=!1,this.dialog_.style.zIndex=\"\",this.replacedStyleTop_&&(this.dialog_.style.top=\"\",this.replacedStyleTop_=!1),this.backdrop_.parentNode&&this.backdrop_.parentNode.removeChild(this.backdrop_),u.dm.removeDialog(this))},setOpen:function(e){e?this.dialog_.hasAttribute(\"open\")||this.dialog_.setAttribute(\"open\",\"\"):(this.dialog_.removeAttribute(\"open\"),this.maybeHideModal())},backdropClick_:function(e){if(this.dialog_.hasAttribute(\"tabindex\"))this.dialog_.focus();else{var t=document.createElement(\"div\");this.dialog_.insertBefore(t,this.dialog_.firstChild),t.tabIndex=-1,t.focus(),this.dialog_.removeChild(t)}var n=document.createEvent(\"MouseEvents\");n.initMouseEvent(e.type,e.bubbles,e.cancelable,window,e.detail,e.screenX,e.screenY,e.clientX,e.clientY,e.ctrlKey,e.altKey,e.shiftKey,e.metaKey,e.button,e.relatedTarget),this.dialog_.dispatchEvent(n),e.stopPropagation()},focus_:function(){var e=this.dialog_.querySelector(\"[autofocus]:not([disabled])\");if(!e&&this.dialog_.tabIndex>=0&&(e=this.dialog_),!e){var t=[\"button\",\"input\",\"keygen\",\"select\",\"textarea\"].map(function(e){return e+\":not([disabled])\"});t.push('[tabindex]:not([disabled]):not([tabindex=\"\"])'),e=this.dialog_.querySelector(t.join(\", \"))}l(document.activeElement),e&&e.focus()},updateZIndex:function(e,t){if(e<t)throw new Error(\"dialogZ should never be < backdropZ\");this.dialog_.style.zIndex=e,this.backdrop_.style.zIndex=t},show:function(){this.dialog_.open||(this.setOpen(!0),this.focus_())},showModal:function(){if(this.dialog_.hasAttribute(\"open\"))throw new Error(\"Failed to execute 'showModal' on dialog: The element is already open, and therefore cannot be opened modally.\");if(!document.body.contains(this.dialog_))throw new Error(\"Failed to execute 'showModal' on dialog: The element is not in a Document.\");if(!u.dm.pushDialog(this))throw new Error(\"Failed to execute 'showModal' on dialog: There are too many open modal dialogs.\");(function(e){for(;e&&e!==document.body;){var t=window.getComputedStyle(e),n=function(e,n){return!(void 0===t[e]||t[e]===n)};if(t.opacity<1||n(\"zIndex\",\"auto\")||n(\"transform\",\"none\")||n(\"mixBlendMode\",\"normal\")||n(\"filter\",\"none\")||n(\"perspective\",\"none\")||\"isolate\"===t.isolation||\"fixed\"===t.position||\"touch\"===t.webkitOverflowScrolling)return!0;e=e.parentElement}return!1})(this.dialog_.parentElement)&&console.warn(\"A dialog is being shown inside a stacking context. This may cause it to be unusable. For more information, see this link: https://github.com/GoogleChrome/dialog-polyfill/#stacking-context\"),this.setOpen(!0),this.openAsModal_=!0,u.needsCentering(this.dialog_)?(u.reposition(this.dialog_),this.replacedStyleTop_=!0):this.replacedStyleTop_=!1,this.dialog_.parentNode.insertBefore(this.backdrop_,this.dialog_.nextSibling),this.focus_()},close:function(e){if(!this.dialog_.hasAttribute(\"open\"))throw new Error(\"Failed to execute 'close' on dialog: The element does not have an 'open' attribute, and therefore cannot be closed.\");this.setOpen(!1),void 0!==e&&(this.dialog_.returnValue=e);var t=new i(\"close\",{bubbles:!1,cancelable:!1});this.dialog_.dispatchEvent(t)}};var u={reposition:function(e){var t=document.body.scrollTop||document.documentElement.scrollTop,n=t+(window.innerHeight-e.offsetHeight)/2;e.style.top=Math.max(t,n)+\"px\"},isInlinePositionSetByStylesheet:function(e){for(var t=0;t<document.styleSheets.length;++t){var n=document.styleSheets[t],o=null;try{o=n.cssRules}catch(e){}if(o)for(var i=0;i<o.length;++i){var a=o[i],l=null;try{l=document.querySelectorAll(a.selectorText)}catch(e){}if(l&&s(l,e)){var r=a.style.getPropertyValue(\"top\"),d=a.style.getPropertyValue(\"bottom\");if(r&&\"auto\"!==r||d&&\"auto\"!==d)return!0}}}return!1},needsCentering:function(e){return\"absolute\"===window.getComputedStyle(e).position&&(!(\"auto\"!==e.style.top&&\"\"!==e.style.top||\"auto\"!==e.style.bottom&&\"\"!==e.style.bottom)&&!u.isInlinePositionSetByStylesheet(e))},forceRegisterDialog:function(e){if((window.HTMLDialogElement||e.showModal)&&console.warn(\"This browser already supports <dialog>, the polyfill may not work correctly\",e),\"dialog\"!==e.localName)throw new Error(\"Failed to register dialog: The element is not a dialog.\");new d(e)},registerDialog:function(e){e.showModal||u.forceRegisterDialog(e)},DialogManager:function(){this.pendingDialogStack=[];var e=this.checkDOM_.bind(this);this.overlay=document.createElement(\"div\"),this.overlay.className=\"_dialog_overlay\",this.overlay.addEventListener(\"click\",function(t){this.forwardTab_=void 0,t.stopPropagation(),e([])}.bind(this)),this.handleKey_=this.handleKey_.bind(this),this.handleFocus_=this.handleFocus_.bind(this),this.zIndexLow_=1e5,this.zIndexHigh_=100150,this.forwardTab_=void 0,\"MutationObserver\"in window&&(this.mo_=new MutationObserver(function(t){var n=[];t.forEach(function(e){for(var t,o=0;t=e.removedNodes[o];++o)t instanceof Element&&(\"dialog\"===t.localName&&n.push(t),n=n.concat(t.querySelectorAll(\"dialog\")))}),n.length&&e(n)}))}};if(u.DialogManager.prototype.blockDocument=function(){document.documentElement.addEventListener(\"focus\",this.handleFocus_,!0),document.addEventListener(\"keydown\",this.handleKey_),this.mo_&&this.mo_.observe(document,{childList:!0,subtree:!0})},u.DialogManager.prototype.unblockDocument=function(){document.documentElement.removeEventListener(\"focus\",this.handleFocus_,!0),document.removeEventListener(\"keydown\",this.handleKey_),this.mo_&&this.mo_.disconnect()},u.DialogManager.prototype.updateStacking=function(){for(var e,t=this.zIndexHigh_,n=0;e=this.pendingDialogStack[n];++n)e.updateZIndex(--t,--t),0===n&&(this.overlay.style.zIndex=--t);var o=this.pendingDialogStack[0];o?(o.dialog.parentNode||document.body).appendChild(this.overlay):this.overlay.parentNode&&this.overlay.parentNode.removeChild(this.overlay)},u.DialogManager.prototype.containedByTopDialog_=function(e){for(;e=a(e);){for(var t,n=0;t=this.pendingDialogStack[n];++n)if(t.dialog===e)return 0===n;e=e.parentElement}return!1},u.DialogManager.prototype.handleFocus_=function(e){if(!this.containedByTopDialog_(e.target)&&(e.preventDefault(),e.stopPropagation(),l(e.target),void 0!==this.forwardTab_)){var t=this.pendingDialogStack[0];return t.dialog.compareDocumentPosition(e.target)&Node.DOCUMENT_POSITION_PRECEDING&&(this.forwardTab_?t.focus_():document.documentElement.focus()),!1}},u.DialogManager.prototype.handleKey_=function(e){if(this.forwardTab_=void 0,27===e.keyCode){e.preventDefault(),e.stopPropagation();var t=new i(\"cancel\",{bubbles:!1,cancelable:!0}),n=this.pendingDialogStack[0];n&&n.dialog.dispatchEvent(t)&&n.dialog.close()}else 9===e.keyCode&&(this.forwardTab_=!e.shiftKey)},u.DialogManager.prototype.checkDOM_=function(e){this.pendingDialogStack.slice().forEach(function(t){-1!==e.indexOf(t.dialog)?t.downgradeModal():t.maybeHideModal()})},u.DialogManager.prototype.pushDialog=function(e){var t=(this.zIndexHigh_-this.zIndexLow_)/2-1;return!(this.pendingDialogStack.length>=t)&&(1===this.pendingDialogStack.unshift(e)&&this.blockDocument(),this.updateStacking(),!0)},u.DialogManager.prototype.removeDialog=function(e){var t=this.pendingDialogStack.indexOf(e);-1!==t&&(this.pendingDialogStack.splice(t,1),0===this.pendingDialogStack.length&&this.unblockDocument(),this.updateStacking())},u.dm=new u.DialogManager,u.formSubmitter=null,u.useValue=null,void 0===window.HTMLDialogElement){var c=document.createElement(\"form\");if(c.setAttribute(\"method\",\"dialog\"),\"dialog\"!==c.method){var m=Object.getOwnPropertyDescriptor(HTMLFormElement.prototype,\"method\");if(m){var p=m.get;m.get=function(){return r(this)?\"dialog\":p.call(this)};var h=m.set;m.set=function(e){return\"string\"==typeof e&&\"dialog\"===e.toLowerCase()?this.setAttribute(\"method\",e):h.call(this,e)},Object.defineProperty(HTMLFormElement.prototype,\"method\",m)}}document.addEventListener(\"click\",function(e){if(u.formSubmitter=null,u.useValue=null,!e.defaultPrevented){var t=e.target;if(t&&r(t.form)){if(!(\"submit\"===t.type&&[\"button\",\"input\"].indexOf(t.localName)>-1)){if(\"input\"!==t.localName||\"image\"!==t.type)return;u.useValue=e.offsetX+\",\"+e.offsetY}a(t)&&(u.formSubmitter=t)}}},!1);var f=HTMLFormElement.prototype.submit;HTMLFormElement.prototype.submit=function(){if(!r(this))return f.call(this);var e=a(this);e&&e.close()},document.addEventListener(\"submit\",function(e){var t=e.target;if(r(t)){e.preventDefault();var n=a(t);if(n){var o=u.formSubmitter;o&&o.form===t?n.close(u.useValue||o.value):n.close(),u.formSubmitter=null}}},!0)}u.forceRegisterDialog=u.forceRegisterDialog,u.registerDialog=u.registerDialog,\"amd\"in n(2)?void 0===(o=function(){return u}.call(t,n,t,e))||(e.exports=o):\"object\"==typeof e&&\"object\"==typeof e.exports?e.exports=u:window.dialogPolyfill=u}()},function(e,t,n){\"use strict\";n.r(t);n(6);const o=1,i=4,a=5,l=6,s=9,r=10,d=11,u=12,c=13,m=14,p=19,h=20,f=21,g=23,b=24,w=26,v=28,y=29;function M(e,t){let n;for(let o of t)n=n?n.subMenu.itemAtIndex(o):e.itemAtIndex(o);return n}class D{constructor(e){this.id=null,this.items=[],this.hostItem=null,this.title=\"\",this.menuData=e}initWithMenuTemplate(e,t){for(var n of t){const t=new x(this.menuData);n.separator?t.separator=!0:t.initWithMenuTemplate(e,n),this.addMenuItem(t)}}addMenuItem(e){this.items.push(e)}getSubMenu(e){if(e<this.items.length){const t=this.items[e];if(t.subMenu)return t.subMenu}return null}itemAtIndex(e){return e<this.items.length?this.items[e]:null}getNodeAtIndex(e){const t=this.itemAtIndex(e);return t&&t.subMenu?t.subMenu.getNode():null}getNode(){if(!this.items||0==this.items.length)return null;const e=document.createElement(\"ul\");for(let t of this.items){let n;if(t.separator)e.lastChild&&(e.lastChild.style.marginBottom=\"10px\");else if(t.subMenu){const e=t.subMenu.getNode();if(e){(n=document.createElement(\"dl\")).classList.add(\"column\");const o=document.createTextNode(t.title),i=document.createElement(\"dt\");i.appendChild(o),n.appendChild(i);const a=document.createElement(\"dd\");a.appendChild(e),n.appendChild(a),a.style.display=\"none\"}}else if(\"\"!==t.title){const e=document.createElement(\"a\"),o=document.createTextNode(t.title);t.enabled&&t.handler&&(e.onclick=t.handler,e.setAttribute(\"href\",\"#\")),e.appendChild(o),n=e}if(n){const t=document.createElement(\"li\");t.classList.add(\"menuItem\"),t.appendChild(n),e.appendChild(t)}}return e}}const _={about:{command:function(e,t){e.showAboutDialog()},label:\"About...\"},front:{command:null,label:\"Bring All to Front\"},cut:{command:function(e,t){},label:\"Cut\"},copy:{command:function(e,t){},label:\"Copy\"},paste:{command:function(e,t){},label:\"Paste\"},delete:{command:function(e,t){},label:\"Delete\"},selectall:{command:function(e,t){},label:\"Select All\"},minimize:{command:null,label:\"Minimize\"},zoom:{command:function(e,t){},label:\"Zoom\"},togglefullscreen:{command:function(e,t){},label:\"Toggle Full Screen\"},viewsource:{command:null,label:\"View Source\"},back:{command:function(e,t){},label:\"Back\"},forward:{command:function(e,t){},label:\"Forward\"}};function k(e,t,n,o){const i=[\"bubbles\",\"cancelBubble\",\"cancelable\",\"composed\",\"defaultPrevented\",\"eventPhase\",\"timeStamp\",\"type\",\"isTrusted\",\"detail\",\"altKey\",\"button\",\"buttons\",\"clientX\",\"clientY\",\"ctrlKey\",\"metaKey\",\"movementX\",\"movementY\",\"region\",\"@relatedTarget\",\"screenX\",\"screenY\",\"shiftKey\"],a={};for(let s of i){const e=t[s];e&&(a[s]=e)}const l={menuId:o.id,elementId:n,appId:e.ID};return a.currentTarget=l,a.target=l,a}class E{constructor(){this.menu=null,this.id=\"\"}getAppMenuNode(e){for(let t of this.menu.items){const n=document.createElement(\"div\");n.classList.add(\"dropdown\");const o=document.createElement(\"button\");if(o.classList.add(\"dropbtn\"),o.appendChild(document.createTextNode(t.title)),n.appendChild(o),t.subMenu){const e=t.subMenu.getNode();if(e){const t=document.createElement(\"div\");t.classList.add(\"dropdown-content\"),t.appendChild(e),n.appendChild(t)}}e.appendChild(n)}return e}getPopupMenuNode(){const e=document.createElement(\"div\");return e.classList.add(\"popupMenu\"),e.appendChild(this.menu.getNode()),e}populateWithAppMenuTemplate(e,t){const n=new D(this);this.menu=n,n.initWithMenuTemplate(e,t)}polulateWithDiffset(e,t){const n=t.items,E=[];let C,S,I;console.log(\"DiffSet:\",JSON.stringify(n));for(let T of n){const t=T.t,n=T.k,A=T.v;switch(t){case o:if(\"menu\"===A){const e=new D;if(E.length>0||this.menu){const t=new x;t.setSubMenu(e),C=t,E.push(t)}else C=e,E.push(e),this.menu=e,this.menu.id=this.id}else if(\"menuitem\"===A){const e=new x;E.push(e),C=e}else{if(\"hr\"!==A)throw\"unsupported tag: \"+A;{const e=new x;e.separator=!0,E.push(e),C=e}}break;case i:C=A?\"number\"==typeof A?E[A]:M(this.menu,A):this.menu;break;case a:S=A?\"number\"==typeof A?E[A]:M(this.menu,A):this.menu;break;case l:I=A?\"number\"==typeof A?E[A]:M(this.menu,A):this.menu;break;case s:if(C instanceof x){const e=C;\"label\"===n&&(e.title=A,e.subMenu&&(e.subMenu.title=A))}else if(\"type\"!==n)throw\"invalid attribute: \"+n+\"/\"+A;break;case r:console.log(\"ditDelAttributeValue\");break;case c:if(!(C&&C instanceof x))throw\"invalid target: \"+C;{const t=C;if(\"menuRole\"===n){const n=_[A];if(!n){console.warn(\"unsupported role name:\",A);break}n.label&&(t.title=n.label,t.subMenu&&(t.subMenu.title=n.label)),n.command&&(t.handler=(t=>{t.preventDefault(),t.stopPropagation(),n.command(e,t)}),t.enabled=!0)}else\"menuAcclerator\"==n?console.log(\"menuAcclerator not implement yet.\"):console.warn(\"unknown dataset name:\",n)}break;case m:console.warn(\"ditDelDataSet: not implement yet.\");break;case p:{let e;if(C?C instanceof D?e=C:C instanceof x&&(e=C.subMenu):e=this.menu,!e)throw\"ditAppendChild: invalid arg: \"+C;if(S instanceof x)e.addMenuItem(S);else if(S!=e)throw\"ditAppendChild: invalid arg1: \"+S}break;case h:console.warn(\"ditInsertBefore: not implement yet.\");break;case f:console.warn(\"ditRemoveChild: not implement yet.\");break;case g:if(!(C instanceof x))throw\"ditAddEventListener: invalid target: \"+C;if(\"click\"!==n)throw\"ditAddEventListener: unsupported event\";{const t=C,n=A.id;t.handler=(t=>{t.preventDefault(),t.stopPropagation();const o=k(e,t,n,this.menu);console.log(\"fakeEvent ==> \",o),e.callNativeMethod(\"/menu/\"+this.menu.id+\"/html/\"+n+\"/click\",o)}),t.enabled=!0}break;case b:console.warn(\"ditRemoveEventListener: not implement yet.\");break;case w:{const e=A;if(C instanceof D)C!=this.menu&&(C.id=e,C.hostItem&&(C.hostItem.id=e));else{if(!(C instanceof x))throw\"node is invalid\";C.id=e}}break;case d:case u:case v:case y:break;default:throw\"Unsupported diff type:\"+t}}}}class x{constructor(e){this.id=\"\",this.subMenu=null,this.title=\"\",this.cmdId=-1,this.enabled=!1,this.separator=!1,this.handler=null,this.menuData=e}initWithMenuTemplate(e,t){const n=t.label,o=t.subMenu;if(this.id=t.id,o){const t=new D(this.menuData);t.initWithMenuTemplate(e,o),this.setSubMenu(t)}const i=t.role;if(i){const t=_[i];if(!t)return void console.warn(\"unsupported role name:\",i);t.label&&(this.title=t.label,this.subMenu&&(this.subMenu.title=t.label)),t.command&&(this.handler=(n=>{n.preventDefault(),n.stopPropagation(),t.command(e,n)}),this.enabled=!0)}else this.subMenu||(this.handler=(t=>{t.preventDefault(),t.stopPropagation(),e.callNativeMethod(\"/menu/\"+this.menuData.id,\"emit\",this.id)}),this.enabled=!0);n&&(this.title=n)}setSubMenu(e){this.subMenu=e,e.hostItem=this,this.enabled=!0}}const C={};function S(e){const t=e.menu;if(!t)throw\"parameter[menu] not found\";const n=C[t];if(!n)throw\"invalid menu\";return n}var I={newMenu:function(e,t){const n=t.parameter.menu;if(!n)throw\"parameter[menu] not found\";const o=new E;o.id=n,C[n]=o,e.responceValue(!0,t.respCallbackNo)},newAppMenu:function(e,t){const n=t.parameter.menu,o=t.argument;if(!n)throw\"parameter[menu] not found\";const i=new E;i.id=n,i.populateWithAppMenuTemplate(e,o),C[n]=i,e.responceValue(!0,t.respCallbackNo)},updateDiffSetHandler:function(e,t){const n=t.parameter,o=t.argument,i=n.menu;if(!i)throw\"parameter[menu] not found\";const a=C[i];if(!a)throw\"invalid menu\";a.polulateWithDiffset(e,o),e.responceValue(!0,t.respCallbackNo)},setApplicationMenu:function(e,t){const n=S(t.parameter),o=document.getElementById(\"menubar\");for(;o.firstChild;)o.removeChild(o.firstChild);n.getAppMenuNode(o)},getMenuData:S,popupContextMenu:function(e,t){S(t.parameter);const n=t.argument.position,o=\"win\"+t.argument.windowId,i=document.getElementById(o);console.log(\"popupContextMenu is not supported ==> pos: \",n,\",  win:\",i)}},T=n(0),A=n.n(T),N=window.exciton,O=window.location,B=\"https\"==O.protocol?\"wss://\":\"ws://\";console.log(N),B+=O.host+O.pathname+\"app/\"+N.ID+\"/ws\",console.log(B);var L=new WebSocket(B);L.onopen=function(){N.callNativeMethod(\"\",\"init\",null)},L.onmessage=function(e){const t=JSON.parse(e.data);console.log(\"onmessage: \",t);const n=t.data,o=\"/app/:app/window/:window\",i=\"/app/:app/menu/:menu\";if(console.log(\"onmessage ===> \",e),t.sync){if(n.target===o&&\"new\"===n.name)return N.newWindow(n);if(n.target===i&&\"new\"===n.name)return I.newMenu(N,n);if(n.target===i&&\"newApplicationMenu\"===n.name)return I.newAppMenu(N,n);if(n.target===o){const e=\"win\"+n.parameter.window,t=document.getElementById(e);console.log(\"call child event: \"+n.name+\", winid = \"+e,n.argument,t,e);const o=t.contentWindow.exciton.requestBrowerEventSync(n.name,JSON.stringify(n.argument));let i;o&&(i=JSON.parse(o)),N.responceValue(i,n.respCallbackNo)}else if(n.target===i)switch(n.name){case\"updateDiffSetHandler\":I.updateDiffSetHandler(N,n)}else{if(\"/app/:app\"!==n.target)throw\"invalid event: \"+n.name;switch(n.name){case\"showMessageBox\":N.showMessageBox(n);break;case\"showOpenDialog\":N.showOpenDialog(n)}}}else if(n.target===o){const e=\"win\"+n.parameter.window,t=document.getElementById(e);console.log(\"call child event: \"+n.name,n.argument),t.contentWindow.exciton.requestBrowserEvent(n.name,JSON.stringify(n.argument))}else{if(n.target!==i)throw\"invalid event: \"+n.name;switch(n.name){case\"setApplicationMenu\":I.setApplicationMenu(N,n);break;case\"popupContextMenu\":I.popupContextMenu(N,n);break;default:throw\"invalid menu event:\"+menuevt}}},N.newWindow=function(e){console.log(\"window/new\",e);const t=document.createElement(\"iframe\");t.classList.add(\"page\"),t.setAttribute(\"title\",e.argument.title),t.setAttribute(\"src\",e.argument.url),t.setAttribute(\"frameborder\",0),t.id=\"win\"+e.parameter.window;const n=document.getElementById(\"pagecontainer\");for(;n.firstChild;)n.removeChild(n.firstChild);n.appendChild(t),N.responceValue(!0,e.respCallbackNo)},N.callWindowMethod=function(e){console.log(\"callWindowMethod\",e),N.callNativeMethod(e.path,e.name,JSON.parse(e.arg))},N.responceValue=function(e,t){var n={target:\"\",name:\"responceEventResult\",argument:e,respCallbackNo:t};console.log(\"responceValue\",n),N.callnative(n)},N.callNativeMethod=function(e,t,n){var o={target:\"/app/\"+N.ID+e,name:t,argument:n,respCallbackNo:-1};N.callnative(o)},N.callnative=function(e){L.send(JSON.stringify(e))},N.showAboutDialog=function(){N.showMessageBoxCore(\"\",\"About...\",\"TODO: App name\",\"\",[\"OK\"],0,null)};let P=!1;N.showMessageBoxCore=function(e,t,n,o,i,a,l){const s=document.getElementById(\"messageBox\"),r=document.getElementById(\"messageBoxIcon\");r.src=e,r.style.display=\"\"===e?\"none\":\"inline\",document.getElementById(\"messageBoxTitle\").innerText=t,document.getElementById(\"messageBoxContent\").innerText=n;const d=document.getElementById(\"messageBoxDetail\");d.innerText=o,d.style.display=\"\"===o?\"none\":\"block\";const u=document.getElementById(\"messageBoxButtons\");for(;u.firstChild;)u.removeChild(u.firstChild);for(let c=0;c<i.length;c++){const e=document.createElement(\"button\");e.type=\"submit\",e.value=c,e.innerText=i[c],e.autofocus=c==a,u.appendChild(e)}P||(A.a.registerDialog(s),P=!0),l&&s.addEventListener(\"close\",e=>{l(e,parseInt(s.returnValue))},{once:!0}),s.showModal()},N.showMessageBox=function(e){console.log(\"window/new\",e);let t=\"\";const n=e.argument.type;switch(n){case 0:break;case 1:t=\"/exciton/web/assets/info.svg\";break;case 2:t=\"/exciton/web/assets/warning.svg\";break;case 3:t=\"/exciton/web/assets/error.svg\";break;case 4:t=\"/exciton/web/assets/question.svg\";break;default:console.warn(\"invalid icon type: \"+n)}let o=e.argument.buttons,i=e.argument.defaultId;0==o.length&&(4==n?(o=[\"YES\",\"NO\"],i=1):(o=[\"OK\"],i=0));const a=e.argument.title,l=e.argument.message,s=e.argument.detail;N.showMessageBoxCore(t,a,l,s,o,i,(t,n)=>{N.responceValue(n,e.respCallbackNo)})};let F=!1;N.showOpenDialog=function(e){const t=document.getElementById(\"fileOpenDialog\"),n=e.argument.title?e.argument.title:\"Open File\";document.getElementById(\"fileOpenDialogTitle\").innerText=n;const o=e.argument.buttonLabel?e.argument.buttonLabel:\"OK\";document.getElementById(\"fileOpenOK\").innerText=o;let i=null;if(e.argument.filters){let t=[];for(let n of e.argument.filters)for(let e of n.extensions)t.push(e);i=t.join(\" \")}const a=document.getElementById(\"selFile\");i?a.setAttribute(\"accept\",i):a.removeAttribute(\"accept\");const l=e.argument.properties;a.multiple=0!=(4&l),F||(A.a.registerDialog(t),F=!0,t.addEventListener(\"close\",n=>{if(\"ok\"===t.returnValue){const t=document.getElementById(\"fileUploadForm\");document.getElementById(\"openDialogResponceNo\").value=e.respCallbackNo;const n=new XMLHttpRequest,o=new FormData(t);n.open(\"POST\",\"/webFileOpenDialog\"),n.send(o)}else console.log(\"file open\")})),t.showModal()}},function(e,t){e.exports=function(){throw new Error(\"define cannot be used indirect\")}},,,,function(e,t,n){},function(e,t,n){e.exports=n(1)}]);"
var _assetsData22d2a248b11d3fc073fe48ee846b045753f35b62 = "<!-- Generated by IcoMoon.io -->\n<svg version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" width=\"32\" height=\"32\" viewBox=\"0 0 32 32\">\n<title>error</title>\n<path d=\"M16 0c-8.837 0-16 7.163-16 16s7.163 16 16 16 16-7.163 16-16-7.163-16-16-16zM16 29c-7.18 0-13-5.82-13-13s5.82-13 13-13 13 5.82 13 13-5.82 13-13 13z\"></path>\n<path d=\"M21 8l-5 5-5-5-3 3 5 5-5 5 3 3 5-5 5 5 3-3-5-5 5-5z\"></path>\n</svg>\n"
var _assetsDatac16fe42b604e35d48eef2497efd9a481bf67a431 = "/*! normalize.css v8.0.0 | MIT License | github.com/necolas/normalize.css */html{-webkit-text-size-adjust:100%;line-height:1.15}h1{font-size:2em;margin:.67em 0}hr{box-sizing:content-box;height:0;overflow:visible}pre{font-family:monospace,monospace;font-size:1em}a{background-color:transparent}abbr[title]{border-bottom:none;text-decoration:underline;text-decoration:underline dotted}b,strong{font-weight:bolder}code,kbd,samp{font-family:monospace,monospace;font-size:1em}small{font-size:80%}sub,sup{font-size:75%;line-height:0;position:relative;vertical-align:baseline}sub{bottom:-.25em}sup{top:-.5em}img{border-style:none}button,input,optgroup,select,textarea{font-family:inherit;font-size:100%;line-height:1.15;margin:0}button,input{overflow:visible}button,select{text-transform:none}[type=button],[type=reset],[type=submit],button{-webkit-appearance:button}[type=button]::-moz-focus-inner,[type=reset]::-moz-focus-inner,[type=submit]::-moz-focus-inner,button::-moz-focus-inner{border-style:none;padding:0}[type=button]:-moz-focusring,[type=reset]:-moz-focusring,[type=submit]:-moz-focusring,button:-moz-focusring{outline:1px dotted ButtonText}fieldset{padding:.35em .75em .625em}legend{box-sizing:border-box;color:inherit;display:table;max-width:100%;padding:0;white-space:normal}progress{vertical-align:baseline}textarea{overflow:auto}[type=checkbox],[type=radio]{box-sizing:border-box;padding:0}[type=number]::-webkit-inner-spin-button,[type=number]::-webkit-outer-spin-button{height:auto}[type=search]{-webkit-appearance:textfield;outline-offset:-2px}[type=search]::-webkit-search-decoration{-webkit-appearance:none}::-webkit-file-upload-button{-webkit-appearance:button;font:inherit}details{display:block}summary{display:list-item}[hidden],template{display:none}dialog{background:#fff;border:solid;color:#000;display:block;height:-moz-fit-content;height:-webkit-fit-content;height:fit-content;left:0;margin:auto;padding:1em;position:absolute;right:0;width:-moz-fit-content;width:-webkit-fit-content;width:fit-content}dialog:not([open]){display:none}dialog+.backdrop{background:rgba(0,0,0,.1)}._dialog_overlay,dialog+.backdrop{bottom:0;left:0;position:fixed;right:0;top:0}dialog.fixed{position:fixed;top:50%;transform:translateY(-50%)}html{font-family:system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif}body{display:flex;flex-direction:column;height:100vh;margin:0;padding:0;width:100vw}#menubar{flex:initial}#pagecontainer{display:flex;flex:auto}.page{flex:auto}.navbar{background-color:#333;overflow:hidden}.navbar a{color:#fff;float:left;padding:14px 16px;text-align:center;text-decoration:none}.dropdown{float:left;overflow:hidden}.dropdown .dropbtn{background-color:inherit;border:none;color:#fff;font:inherit;margin:0;outline:none;padding:14px 16px}.dropdown:hover .dropbtn,.navbar a:hover{background-color:red}.dropdown-content{background-color:#f9f9f9;column-rule:1px solid #ccc;column-width:240px;display:none;left:0;position:absolute;width:100%;z-index:1}.dropdown:hover .dropdown-content{display:block;padding:16px}li{list-style-type:none}li,ul{margin:0;padding:0}.column{margin:0 0 16px}dd dl.column{margin-left:8px}.column dt{font-weight:700}.column dd{margin-left:0}.column a,li.menuItem a{color:#000;display:block;float:none;padding:8px;text-align:left;text-decoration:none}li{break-after:auto;break-before:auto;break-inside:avoid-column}.column a:hover{background-color:#ddd}.row:after{clear:both;content:\"\";display:table}dialog{border:0;border-radius:.6rem;box-shadow:0 0 1em #000;min-width:250px;padding:0}dialog[open]{animation:slide-up .4s ease-out}dialog h3{background-color:#333;border-bottom:1px solid #fff;border-top-left-radius:.6rem;border-top-right-radius:.6em;color:#fff;margin:0;padding:.6em 1em}dialog div.dialogContent{margin:0;padding:.6em}dialog footer{border-top:1px solid #333;display:flex;margin:0;padding:.64em 1em}dialog footer form button{flex:auto}.dialogContentDetail{display:none;font-size:80%}"
var _assetsData0115ca79e54f87ebaca3eef672d910710a7b2459 = "<!DOCTYPE html>\n<html>\n\n<head>\n    <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n    <meta charset=\"UTF-8\">\n    <link type=\"text/css\" rel=\"stylesheet\" href=\"/exciton/web/assets/webroot.css\"></link>\n</head>\n\n<body>\n    <nav id=\"menubar\" class=\"navbar\"></nav>\n    <div id=\"pagecontainer\"></div>\n    <script>window.exciton = { ID: {{.ID}} };</script>\n    <script src=\"/exciton/web/assets/exciton-ws.js\"></script>\n\n    <dialog id=\"fileOpenDialog\">\n        <h3 id=\"fileOpenDialogTitle\"></h3>\n        <div class=\"dialogContent\">\n            <form id=\"fileUploadForm\" enctype=\"multipart/form-data\">\n            <input type=\"file\" id=\"selFile\" name=\"selFile\"></input>\n            <input type=\"hidden\" id=\"openDialogResponceNo\" name=\"openDialogResponceNo\" value=\"0\"></input>\n            </form>\n        </div>\n        <footer>\n        <form method=\"dialog\">\n            <button value=\"ok\" id=\"fileOpenOK\">Upload</button>\n            <button value=\"cancel\">Cancel</button>\n        </form>\n        </footer>\n    </dialog>\n    <dialog id=\"messageBox\">\n    <h3 id=\"messageBoxTitle\"></h3>\n    <div class=\"dialogContent\">\n        <p><img width=\"32\" height=\"32\" id=\"messageBoxIcon\"><span id=\"messageBoxContent\"></span></p>\n        <p id=\"messageBoxDetail\" class=\"dialogContentDetail\"></p>\n    </div>\n    <footer>\n        <form method=\"dialog\" id=\"messageBoxButtons\">\n            <!--input class=\"dialogCloseButton\" type=\"submit\" value=\"Close\"-->\n        </form>\n    </dialog>\n</body>\n\n</html>"
var _assetsDataaef3c0e72af327412e4a31eb57713d72a7622b08 = "<!-- Generated by IcoMoon.io -->\n<svg version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" width=\"32\" height=\"32\" viewBox=\"0 0 32 32\">\n<title>warning</title>\n<path d=\"M16 2.899l13.409 26.726h-26.819l13.409-26.726zM16 0c-0.69 0-1.379 0.465-1.903 1.395l-13.659 27.222c-1.046 1.86-0.156 3.383 1.978 3.383h27.166c2.134 0 3.025-1.522 1.978-3.383h0l-13.659-27.222c-0.523-0.93-1.213-1.395-1.903-1.395v0z\"></path>\n<path d=\"M18 26c0 1.105-0.895 2-2 2s-2-0.895-2-2c0-1.105 0.895-2 2-2s2 0.895 2 2z\"></path>\n<path d=\"M16 22c-1.105 0-2-0.895-2-2v-6c0-1.105 0.895-2 2-2s2 0.895 2 2v6c0 1.105-0.895 2-2 2z\"></path>\n</svg>\n"
var _assetsData67b16a7f0b419a9f7e7cf60118d68cd4fb55175c = "<!-- Generated by IcoMoon.io -->\n<svg version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" width=\"32\" height=\"32\" viewBox=\"0 0 32 32\">\n<title>info</title>\n<path d=\"M14 9.5c0-0.825 0.675-1.5 1.5-1.5h1c0.825 0 1.5 0.675 1.5 1.5v1c0 0.825-0.675 1.5-1.5 1.5h-1c-0.825 0-1.5-0.675-1.5-1.5v-1z\"></path>\n<path d=\"M20 24h-8v-2h2v-6h-2v-2h6v8h2z\"></path>\n<path d=\"M16 0c-8.837 0-16 7.163-16 16s7.163 16 16 16 16-7.163 16-16-7.163-16-16-16zM16 29c-7.18 0-13-5.82-13-13s5.82-13 13-13 13 5.82 13 13-5.82 13-13 13z\"></path>\n</svg>\n"

// assetsData returns go-assets FileSystem
var assetsData = assets.NewFileSystem(map[string][]string{"/": []string{"webroot.css", "webroot.gohtml", "question.svg", "warning.svg", "exciton-ws.js", "info.svg", "error.svg"}}, map[string]*assets.File{
	"/question.svg": &assets.File{
		Path:     "/question.svg",
		FileMode: 0x1ed,
		Mtime:    time.Unix(1533560732, 1533560732000000000),
		Data:     []byte(_assetsData4029148de0a3b24c533ab030ac54e17cbacc31b7),
	}, "/exciton-ws.js": &assets.File{
		Path:     "/exciton-ws.js",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1554419629, 1554419629832411486),
		Data:     []byte(_assetsData811b89eace18a3615f392e27d8492a6c890bdd96),
	}, "/error.svg": &assets.File{
		Path:     "/error.svg",
		FileMode: 0x1ed,
		Mtime:    time.Unix(1533560732, 1533560732000000000),
		Data:     []byte(_assetsData22d2a248b11d3fc073fe48ee846b045753f35b62),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1533898704, 1533898704000000000),
		Data:     nil,
	}, "/webroot.css": &assets.File{
		Path:     "/webroot.css",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1554419629, 1554419629832463834),
		Data:     []byte(_assetsDatac16fe42b604e35d48eef2497efd9a481bf67a431),
	}, "/webroot.gohtml": &assets.File{
		Path:     "/webroot.gohtml",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1533684338, 1533684338000000000),
		Data:     []byte(_assetsData0115ca79e54f87ebaca3eef672d910710a7b2459),
	}, "/warning.svg": &assets.File{
		Path:     "/warning.svg",
		FileMode: 0x1ed,
		Mtime:    time.Unix(1533560732, 1533560732000000000),
		Data:     []byte(_assetsDataaef3c0e72af327412e4a31eb57713d72a7622b08),
	}, "/info.svg": &assets.File{
		Path:     "/info.svg",
		FileMode: 0x1ed,
		Mtime:    time.Unix(1533560732, 1533560732000000000),
		Data:     []byte(_assetsData67b16a7f0b419a9f7e7cf60118d68cd4fb55175c),
	}}, "")
