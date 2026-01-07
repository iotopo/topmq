function has(browser: string) {
  const ua = navigator.userAgent
  if (browser === 'ie') {
    const isIE = ua.includes('compatible') && ua.includes('MSIE')
    if (isIE) {
      const reIE = new RegExp('MSIE (\\d+\\.\\d+);')
      reIE.test(ua)
      return Number.parseFloat(RegExp['$1'])
    } else {
      return false
    }
  } else {
    return ua.includes(browser)
  }
}
export const csvService = {
  _isIE11() {
    let iev = 0
    const ieold = /MSIE (\d+\.\d+);/.test(navigator.userAgent)
    const trident = !!navigator.userAgent.match(/Trident\/7.0/)
    const rv = navigator.userAgent.indexOf('rv:11.0')

    if (ieold) {
      iev = Number(RegExp.$1)
    }
    if (navigator.appVersion.includes('MSIE 10')) {
      iev = 10
    }
    if (trident && rv !== -1) {
      iev = 11
    }

    return iev === 11
  },

  _isEdge() {
    return /Edge/.test(navigator.userAgent)
  },

  _getDownloadUrl(text: string) {
    const BOM = '\uFEFF'
    // Add BOM to text for open in excel correctly
    // if (window.Blob && window.URL && window.URL.createObjectURL) {
    // } else {
    //   return `data:attachment/csv;charset=utf-8,${BOM}${encodeURIComponent(
    //     text
    //   )}`
    // }
    const csvData = new Blob([BOM + text], {
      type: 'text/csv',
    })
    return URL.createObjectURL(csvData)
  },

  download(filename: string, text: string) {
    // if (has('ie') && has('ie') < 10) {
    //   // has module unable identify ie11 and Edge
    //   const oWin = window.top.open('about:blank', '_blank')
    //   oWin.document.charset = 'utf-8'
    //   oWin.document.write(text)
    //   oWin.document.close()
    //   oWin.document.execCommand('SaveAs', filename)
    //   oWin.close()
    // } else if (has('ie') === 10 || this._isIE11() || this._isEdge()) {
    //   const BOM = '\uFEFF'
    //   const csvData = new Blob([BOM + text], {
    //     type: 'text/csv',
    //   })
    //   navigator.msSaveBlob(csvData, filename)
    // } else {
    //   const link = document.createElement('a')
    //   link.download = filename
    //   link.href = this._getDownloadUrl(text)
    //   document.body.append(link)
    //   link.click()
    //   link.remove()
    // }
    const link = document.createElement('a')
    link.download = filename
    link.href = this._getDownloadUrl(text)
    document.body.append(link)
    link.click()
    link.remove()
  },
  upload(size: number) {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = '.csv'
    input.style.display = 'none'
    document.body.append(input)

    return new Promise<string[][]>((resolve, reject) => {
      input.addEventListener('change', () => {
        if (!input.files) {
          return
        }
        const file = input.files[0]

        if (file.size > size << 20) {
          reject(`文件大小不可超过${size}M`)
          return
        }

        const reader = new FileReader()
        reader.addEventListener('load', (e) => {
          const target = e.currentTarget
          if (!target) {
            return
          }
          const content = (target as any).result
          if (typeof content === 'string') {
            let lines = content.split('\r\n')
            // let headerLine = lines[0];
            // let columns = headerLine.split(',');
            lines = lines.slice(1)
            const rows: any[] = []
            lines.forEach((line) => {
              if (!line.trim()) {
                return
              }
              const row = line.split(',')
              rows.push(row)
            })
            resolve(rows)
          } else {
            reject('invalid content')
          }
        })
        reader.readAsText(file)
      })
      input.click()
    })
  },
}
