import { ElMessage } from 'element-plus'

export function uploadFile(options: any) {
  const opts = Object.assign(
    {
      accept: '.txt',
      size: 5 << 20, // 5M
      multiple: false,
    },
    options || {}
  )
  const onload = opts.onload

  const input = document.createElement('input')
  input.type = 'file'
  // input.accept = 'aplication/zip,application/x-zip,application/x-zip-compressed'
  input.accept = opts.accept
  input.multiple = opts.multiple
  input.style.display = 'none'
  document.body.append(input)

  input.addEventListener('change', () => {
    if (!input.files || input.files.length === 0) {
      return
    }
    for (let i = 0; i < input.files.length; i++) {
      const file = input.files[i]

      if (file.size > opts.size) {
        ElMessage.warning(
          `文件大小不可超过${Number(opts.size >> 20).toFixed(0)}MB`
        )
        return
      }

      if (typeof onload == 'function') {
        onload(file)
      }
      if (!opts.multiple) {
        break
      }
    }
  })
  input.click()
}

export function uploadZip(onload: any) {
  // NOTE: 上传压缩包最大支持到 500 M
  uploadFile({
    accept: 'aplication/zip,application/x-zip,application/x-zip-compressed',
    size: 500 << 20, // 500 M
    onload,
  })
}

export function uploadText(options: any) {
  const opts = Object.assign(
    {
      accept: '.txt',
      size: 1024 * 1024,
    },
    options || {}
  )

  const input = document.createElement('input')
  input.type = 'file'
  input.accept = opts.accept
  input.style.display = 'none'
  document.body.append(input)

  return new Promise((resolve, reject) => {
    input.addEventListener('change', () => {
      if (!input.files || input.files.length === 0) {
        return
      }
      const file = input.files[0]

      if (file.size > opts.size) {
        ElMessage.warning(
          `文件大小不可超过${Number(opts.size / 1024 / 1024).toFixed(0)}M`
        )
        return
      }

      if (typeof opts.onload == 'function') {
        opts.onload(file)
      } else {
        const reader = new FileReader()
        reader.addEventListener('load', (e: any) => {
          const content = e.currentTarget.result
          if (typeof content === 'string') {
            resolve(content)
          } else {
            reject(new Error('无法读取文件内容'))
          }
        })
        reader.readAsText(file)
      }
    })
    input.click()
  })
}
