var clicked = 0;
export function onClickClient1(c, e) {
    console.log(c, e);
    c.clickCount++;
    e.currentTarget.innerHTML = '<b>clicked!</b>: ' + c.clickCount;
}
export function mountComponent(n, c) {
    c.clickCount = 0;
    c.callNativeEvent('on-mount', 'called')
}

export function unmountComponent(c) {
    console.log('mountComponent', c);
}

export function clientFunc1(c, add) {
    return c.clickCount + add;
}
