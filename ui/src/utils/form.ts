const isBlob = (value: any) => value instanceof Blob
const isFile = (value: any) => value instanceof File
const isBoolean = (value: any) => typeof value === 'boolean'
const isNull = (value: any) => value === null
const isUndefined = (value: any) => value === undefined
const isArray = (value: any) => Array.isArray(value)
const isObject = (value: any) => !isArray(value) && typeof value === 'object'
const isNumber = (value: any) =>
  typeof value === 'number' && !Number.isNaN(value)

interface Options {
  arrayIndexes: boolean
  excludeNull: boolean
  useDotSeparator: boolean
  useBrackets: boolean
  booleanAsNumbers: boolean
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const processData = (
  value: any,
  options: Options,
  formData: FormData,
  parent?: string
) => {
  const processedKey = parent || ''

  if (isNull(value) || isUndefined(value)) {
    if (!options.excludeNull) {
      formData.append(processedKey, '')
    }
  } else if (isFile(value)) {
    formData.append(processedKey, value)
  } else if (isBlob(value)) {
    formData.append(processedKey, value)
  } else if (isArray(value)) {
    value.forEach((item: any, index: number) => {
      let computedKey = processedKey
      if (options.useBrackets) {
        computedKey += `[${options.arrayIndexes ? index : ''}]`
      }
      processData(item, options, formData, computedKey)
    })
  } else if (isObject(value)) {
    Object.entries(value).forEach(([key, data]) => {
      let computedKey = key
      if (parent) {
        computedKey = options.useDotSeparator
          ? `${parent}.${key}`
          : `${parent}[${key}]`
      }
      processData(data, options, formData, computedKey)
    })
  } else if (isBoolean(value)) {
    if (options.booleanAsNumbers) {
      formData.append(processedKey, `${Number(value)}`)
    } else {
      formData.append(processedKey, value ? 'true' : 'false')
    }
  } else if (isNumber(value)) {
    formData.append(processedKey, `${value}`)
  } else {
    formData.append(processedKey, value)
  }
}

const defaultOptions: Options = {
  arrayIndexes: true,
  excludeNull: true,
  useDotSeparator: true,
  useBrackets: false,
  booleanAsNumbers: false,
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const objectToFormData = (
  payload: any,
  options: Partial<Options> = {},
  formData: FormData = new FormData()
) => {
  if (!payload) return formData

  options = Object.assign({}, defaultOptions, options)

  processData(payload, options as Options, formData)

  return formData
}

export default objectToFormData
