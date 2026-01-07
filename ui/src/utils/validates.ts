type Validate = {
  validator: any
  trigger: string
}

export const mobilePattern = /^\+?[1-9]\d{1,14}$/
//export const passowrdPattern = /^.*(?=.{6,})(?=.*\d)(?=.*[A-Za-z]).*$/
export const passowrdPattern = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?])(?!.*\s)/
export const usernamePattern = /^[A-Za-z0-9._-]{1,20}$/
export const pointIDPattern = /^[A-Za-z0-9_-]{0,35}$/
export const deviceIDPattern = /^[A-Za-z]\w{0,35}$/
export const gatewayIDPattern = /^[\w-]{0,35}$/
export const emailPattern = /^\w+([+.-]\w+)*@\w+([.-]\w+)*\.\w+([.-]\w+)*$/
export const variableNamePattern = /^[A-Z_a-z]\w*$/

export default {
  /**
   * 特殊字符校验
   */
  get specialChar(): Validate {
    return {
      validator: (rule: any, value: any, callback: any) => {
        if (value == null || value == '') {
          callback()
          return
        }
        if (value.includes(' ')) {
          callback(new Error('不能含有空格'))
          return
        }

        const reg = new RegExp(
          `[${'`'}${'"'}${'\\\\'}~!@#$^&*()=|{}':;',\\[\\].<>/?~！@#￥……&*（）——|{}【】‘；：”“'。，、？]`
        )
        if (!reg.test(value)) {
          callback()
        } else {
          callback(new Error('不能含有特殊字符'))
        }
      },
      trigger: 'blur',
    }
  },
  get space(): Validate {
    return {
      validator: (rule: any, value: any, callback: any) => {
        const reg = new RegExp('(^\\s+)|(\\s+$)')
        if (!reg.test(value)) {
          callback()
        } else {
          callback(new Error('字符串两端不能包含空格'))
        }
      },
      trigger: 'blur',
    }
  },
  get email(): Validate {
    return {
      validator: (rule: any, value: any, callback: any) => {
        if (value == null || value == '') {
          callback()
          return
        }
        const reg = new RegExp(
          '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$'
        )

        if (reg.test(value)) {
          callback()
        } else {
          callback(new Error('邮箱格式不正确'))
        }
      },
      trigger: 'blur',
    }
  },
  get mobile(): Validate {
    return {
      validator: (rule: any, value: any, callback: any) => {
        if (value == null || value == '') {
          callback()
          return
        }
        const reg = /^\+?[1-9]\d{0,3}[\s-]?\d{6,14}$/
        if (reg.test(value)) {
          callback()
        } else {
          callback(new Error('手机号格式不正确'))
        }
      },
      trigger: 'blur',
    }
  },
}
