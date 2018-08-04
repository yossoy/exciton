'use strict';
import styles from './webroot.css';
import menu from './menu';
import dialogPolyfill from 'dialog-polyfill';
var nsobj = window.exciton;


var l = window.location;
var ws_url = (l.protocol == 'https') ? 'wss://' : 'ws://';
console.log(nsobj);
ws_url += (l.host + l.pathname + 'exciton/' + nsobj.ID + '/ws');
console.log(ws_url);

var sock = new WebSocket(ws_url);
sock.onopen = function () {
    // send appid for communication
    nsobj.callNativeMethod('/app/init', null);
};
sock.onmessage = function (e) {
    const ed = JSON.parse(e.data);
    const d = ed.data;
    const winPrefix = '/exciton/:appid/window/:id/';
    const menuPrefix = '/exciton/:appid/menu/:id/';
    console.log(ed);
    if (ed.sync) {
        if (d.name == (winPrefix + 'new')) {
            return nsobj.newWindow(d);
        }
        if (d.name == (menuPrefix + 'new')) {
            return menu.newMenu(nsobj, d);
        }
        if (d.name.startsWith(winPrefix)) {
            const winevnt = d.name.slice(winPrefix.length);
            const winid = 'win' + d.parameter['id'];
            const w = document.getElementById(winid);
            console.log('call child event: ' + winevnt, d.argument);
            const resultStr = w.contentWindow.exciton.requestBrowerEventSync(winevnt, JSON.stringify(d.argument));
            const result = JSON.parse(resultStr);
            nsobj.responceValue(result, d.respCallbackNo);
        } else if (d.name.startsWith(menuPrefix)) {
            const menuevt = d.name.slice(menuPrefix.length);
            switch (menuevt) {
                case 'updateDiffSetHandler':
                menu.updateDiffSetHandler(nsobj, d);
                break;
            }
        } else {
            throw 'invalid event: ' + d.name;
        }
    } else {
        if (d.name.startsWith(winPrefix)) {
            const winevnt = d.name.slice(winPrefix.length);
            const winid = 'win' + d.parameter['id'];
            const w = document.getElementById(winid);
            console.log('call child event: ' + winevnt, d.argument);
            w.contentWindow.exciton.requestBrowserEvent(winevnt, JSON.stringify(d.argument));
        } else if (d.name.startsWith(menuPrefix)) {
            const menuevt = d.name.slice(menuPrefix.length);
            switch (menuevt) {
                case 'setApplicationMenu':
                    menu.setApplicationMenu(nsobj, d);
                    break;
                default:
                    throw 'invalid menu event:' + menuevt;
            }
        } else {
            throw 'invalid event: ' + d.name;
        }
    }
};

nsobj.newWindow = function (dd) {
    console.log('window/new', dd);
    const iframe = document.createElement('iframe');
    iframe.classList.add('page');
    iframe.setAttribute('title', dd.argument.title);
    iframe.setAttribute('src', dd.argument.url);
    iframe.setAttribute('frameborder', 0);
    iframe.id = 'win' + dd.parameter['id'];
    const p = document.getElementById('pagecontainer');
    while (p.firstChild) {
        //or hide children?
        p.removeChild(p.firstChild);
    }
    p.appendChild(iframe);
    nsobj.responceValue(true, dd.respCallbackNo);
};

nsobj.callWindowMethod = function (d) {
    console.log('callWindowMethod', d);
    nsobj.callNativeMethod(d.path, JSON.parse(d.arg));
};

nsobj.responceValue = function (val, respNo) {
    var data = {
        name: '/responceEventResult',
        argument: val, //TODO: error result
        respCallbackNo: respNo,
    };
    console.log('responceValue', data);
    nsobj.callnative(data);
};

nsobj.callNativeMethod = function (method, arg) {
    var data = {
        name: '/exciton/' + nsobj.ID + method,
        argument: arg,
        respCallbackNo: -1,
    };
    nsobj.callnative(data);
};


nsobj.callnative = function (data) {
    sock.send(JSON.stringify(data));
};

nsobj.showAboutDialog = function() {
    const dlg = document.getElementById('aboutDialog');
    dialogPolyfill.registerDialog(dlg);
    dlg.showModal();
    console.log('roleCmdAbout!');
};
