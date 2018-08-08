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
    const dialogPrefix = '/exciton/:appid/dialog/:id/';
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
        } else if (d.name.startsWith(dialogPrefix)) {
            const dlgevt = d.name.slice(dialogPrefix.length);
            switch (dlgevt) {
                case 'showMessageBox':
                    nsobj.showMessageBox(d);
                    break;
                case 'showOpenDialog':
                    nsobj.showOpenDialog(d);
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

nsobj.showAboutDialog = function () {
    //TODO: app icon
    nsobj.showMessageBoxCore('', 'About...', 'TODO: App name', '', ['OK'], 0, null);
};

let mesasgeBoxRegisterd = false;

nsobj.showMessageBoxCore = function (iconSrc, title, message, detail, buttons, defaultId, respCallback) {
    const dlg = document.getElementById('messageBox');
    const icon = document.getElementById('messageBoxIcon');
    icon.src = iconSrc;
    icon.style.display = (iconSrc === '') ? 'none' : 'inline';
    document.getElementById('messageBoxTitle').innerText = title;
    document.getElementById('messageBoxContent').innerText = message;
    const detailElem = document.getElementById('messageBoxDetail');
    detailElem.innerText = detail;
    detailElem.style.display = (detail === '') ? 'none' : 'block';
    const buttonBase = document.getElementById('messageBoxButtons');
    while (buttonBase.firstChild) {
        buttonBase.removeChild(buttonBase.firstChild);
    }
    for (let i = 0; i < buttons.length; i++) {
        const b = document.createElement('button');
        b.type = 'submit';
        b.value = i;
        b.innerText = buttons[i];
        b.autofocus = (i == defaultId);
        buttonBase.appendChild(b);
    }
    if (!mesasgeBoxRegisterd) {
        dialogPolyfill.registerDialog(dlg);
        mesasgeBoxRegisterd = true;
    }
    if (respCallback) {
        dlg.addEventListener('close', (e) => {
            respCallback(e, parseInt(dlg.returnValue));
        }, {
            once: true
        });
    }

    dlg.showModal();
};

nsobj.showMessageBox = function (dd) {
    console.log('window/new', dd);
    let iconSrc = '';
    const type = dd.argument['type'];
    switch (type) {
        case 0: // none
            break;
        case 1: // info
            iconSrc = '/exciton/web/assets/info.svg'
            break;
        case 2: // warning
            iconSrc = '/exciton/web/assets/warning.svg'
            break;
        case 3: // error
            iconSrc = '/exciton/web/assets/error.svg'
            break;
        case 4: // question
            iconSrc = '/exciton/web/assets/question.svg'
            break;
        default:
            console.warn('invalid icon type: ' + type);
            break;
    }
    let buttons = dd.argument['buttons'];
    let defaultId = dd.argument['defaultId'];
    if (buttons.length == 0) {
        if (type == 4 /*question*/ ) {
            buttons = ['YES', 'NO'];
            defaultId = 1;
        } else {
            buttons = ['OK'];
            defaultId = 0;
        }
    }
    const title = dd.argument['title'];
    const message = dd.argument['message'];
    const detail = dd.argument['detail'];
    const ret = nsobj.showMessageBoxCore(iconSrc, title, message, detail, buttons, defaultId, (e, returnValue) => {
        nsobj.responceValue(returnValue, dd.respCallbackNo);
    });
};

let fileOpenDialogRegisterd = false;
nsobj.showOpenDialog = function (dd) {
    const dlg = document.getElementById('fileOpenDialog');
    const title = (dd.argument['title']) ? dd.argument['title'] : "Open File";
    document.getElementById('fileOpenDialogTitle').innerText = title;
    const buttonLabel = (dd.argument['buttonLabel']) ? dd.argument['buttonLabel'] : "OK";
    document.getElementById('fileOpenOK').innerText = buttonLabel;
    let accept = null;
    if (dd.argument['filters']) {
        let exts = [];
        for (let filter of dd.argument['filters']) {
            for (let ext of filter['extensions']) {
                exts.push(ext);
            }
        }
        accept = exts.join(' ')
    }
    const infile = document.getElementById('selFile');
    if (accept) {
        infile.setAttribute('accept', accept);
    } else {
        infile.removeAttribute('accept');
    }
    const prop = dd.argument['properties'];
    infile.multiple = (prop & 4) != 0;
    if (!fileOpenDialogRegisterd) {
        dialogPolyfill.registerDialog(dlg);
        fileOpenDialogRegisterd = true;
        dlg.addEventListener('close', (e) => {
            if (dlg.returnValue === 'ok') {
                const form = document.getElementById('fileUploadForm');
                document.getElementById('openDialogResponceNo').value = dd.respCallbackNo;
                const XHR = new XMLHttpRequest();
                const FD = new FormData(form);
                XHR.open('POST', '/webFileOpenDialog');
                XHR.send(FD);
                //TODO: error
            } else {
                console.log('file open');
                //TODO: cancel proc
            }
        });
    }
    dlg.showModal();
};