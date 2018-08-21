var clicked = 0;
export function onClickClient1(c, e) {
    console.log(c, e);
    c.clickCount++;
    e.currentTarget.innerHTML = '<b>clicked!</b>: ' + c.clickCount;
}
export function mountComponent(n, inst) {
    inst.clickCount = 0;
    inst.callNativeEvent('/on-mount', 'called')
}

export function unmountComponent(inst) {
    console.log('mountComponent', inst);
}

export function clientFunc1(c, add) {
    return c.clickCount + add;
}
