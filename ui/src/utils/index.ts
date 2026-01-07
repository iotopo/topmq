export const randomString = function (length: number) {
  let str = ''
  for (; str.length < length; str += Math.random().toString(36).slice(2));
  return str.slice(0, Math.max(0, length))
}

export function getCssVarValue(name: string) {
  return getComputedStyle(document.documentElement).getPropertyValue(name)
}

// base64 encode
export function b2a(s: string) {
  return window.btoa(window.encodeURIComponent(s))
}
// base64 decode
export function a2b(s: string) {
  return window.decodeURIComponent(window.atob(s))
}

export const baseUrl = '/' // document.baseURI.slice(location.origin.length)
// console.log('document.baseURI:', document.baseURI)

// 添加一个工具函数来创建防抖函数, 在指定时间周期内只能执行一次，默认 300ms
export function debounce(fn: () => void, t: number = 300) {
  let ts = 0
  return function (this: any, ...args: any[]) {
    if (Date.now() - ts < t) {
      return
    }
    ts = Date.now()
    fn()
  }
}

// 创建一个延迟执行的函数，每次执行时，如果之前有定时器，则清除之前定时器
export function delayer(fn: (...args: any[]) => void, delay: number = 300) {
  if (delay < 100) {
    delay = 100
  }
  let timer: any = null
  return function (this: any, ...args: any[]) {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}

export function getBrowserEngine() {
  const userAgent = navigator.userAgent.toLowerCase();
  let engine: {
    type: 'Blink' | 'Blink (Edge)' | 'WebKit' | 'Gecko' | 'Trident' | '';
    version: string | null;
  } = { type: '', version: '' };

  // 检测 Blink 内核（Chrome、Edge、Opera 等现代浏览器）
  if (/chrome\/\d+/.test(userAgent) && !/edge\/\d+/.test(userAgent) && !/edg\/\d+/.test(userAgent)) {
    const match = userAgent.match(/chrome\/(\d+\.\d+\.\d+\.\d+)/);
    engine = { type: 'Blink', version: match ? match[1] : null };
  }
  // 检测 Edge 浏览器（基于 Blink，但需单独识别）
  else if (/edge\/\d+/.test(userAgent) || /edg\/\d+/.test(userAgent)) {
    const match = userAgent.match(/(edge|edg)\/(\d+\.\d+\.\d+\.\d+)/);
    engine = { type: 'Blink (Edge)', version: match ? match[2] : null };
  }
  // 检测 WebKit 内核（Safari 等）
  else if (/safari\/\d+/.test(userAgent) && !/chrome\/\d+/.test(userAgent)) {
    const match = userAgent.match(/version\/(\d+\.\d+\.\d+)/);
    engine = { type: 'WebKit', version: match ? match[1] : null };
  }
  // 检测 Gecko 内核（Firefox 等）
  else if (/firefox\/\d+/.test(userAgent)) {
    const match = userAgent.match(/firefox\/(\d+\.\d+)/);
    engine = { type: 'Gecko', version: match ? match[1] : null };
  }
  // // 检测 Trident 内核（IE 浏览器）
  // else if (/trident\/\d+/.test(userAgent)) {
  //   const match = userAgent.match(/trident\/(\d+\.\d+)/);
  //   engine = { type: 'Trident', version: match ? match[1] : null };
  // }

  return engine;
}
