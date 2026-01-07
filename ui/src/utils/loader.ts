export function loadJS(url: string, onload?: () => void) {
  return new Promise<void>((resolve) => {
    const js = document.createElement('script')
    js.addEventListener('load', () => {
      if (onload) {
        onload()
      }
      js.remove()
      resolve()
    })
    js.setAttribute('src', url)
    document.querySelectorAll('head')[0].append(js)
  })
}
