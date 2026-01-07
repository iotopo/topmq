/** 变量名称转换 */
/**
 * constant   -> I_AM_24_YEARS_OLD
 * camel      -> iAm24YearsOld
 * pascal     -> IAm24YearsOld
 * hyphen     -> i-am-24-years-old
 * snake      -> i_am_24_years_old
 * underscore -> _i_am_24_years_old
 * unknown -> I-Am_24YearsOLD
 */
export type VariableNameMethodsType =
  | 'camel'
  | 'pascal'
  | 'hyphen'
  | 'constant'
  | 'snake'
  | 'underscore'
  | 'unknown'

export type VariableNameTransferType =
  | 'camel'
  | 'pascal'
  | 'hyphen'
  | 'constant'
  | 'snake'
  | 'underscore'

/** Determine the style of variable naming */
export const judgeStyle = (varName: string): VariableNameMethodsType => {
  const isFullUpperCase = isUpperCase(varName) // is all upper case letters [try to judge is constant]
  const isFullLowerCase = isLowerCase(varName)
  const hasHyphen = hasSomeWords(varName, '-')
  const hasUnderLine = hasSomeWords(varName, '_')
  const startWithUnderLine = hasSomeWords(varName[0], '_')
  const startWithUpperCase = isUpperCase(varName[0])

  if (isFullUpperCase && !hasHyphen) {
    return 'constant'
  } else if (!hasUnderLine && hasHyphen && isFullLowerCase) {
    return 'hyphen'
  } else if (
    !hasHyphen &&
    hasUnderLine &&
    isFullLowerCase &&
    !startWithUnderLine
  ) {
    return 'snake'
  } else if (
    isFullLowerCase &&
    startWithUnderLine &&
    hasUnderLine &&
    !hasHyphen
  ) {
    return 'underscore'
  } else if (
    !hasUnderLine &&
    !hasHyphen &&
    !startWithUpperCase &&
    !startWithUnderLine
  ) {
    return 'camel'
  } else if (
    !hasUnderLine &&
    !hasHyphen &&
    startWithUpperCase &&
    !startWithUnderLine
  ) {
    return 'pascal'
  } else {
    return 'unknown'
  }
}

/** Determines whether the string consists of all uppercase letters */
export const isUpperCase = (str: string): boolean => {
  return str === str.toUpperCase()
}

/** Determines whether the string consists of all lowercase letters */
export const isLowerCase = (str: string): boolean => {
  return str === str.toLowerCase()
}

/** Determines whether certain constructions are present in a string */
export const hasSomeWords = (
  str: string,
  words: string | string[]
): boolean => {
  if (typeof words === 'string') {
    return str.search(words) !== -1
  } else {
    for (const word of words) {
      if (str.search(word) !== -1) return true
    }
    return false
  }
}

/** Transfer variable name style */
export const varNameTransfer = (
  varName: string,
  style: VariableNameTransferType
): string => {
  const varNameSlice: string[] = variableNameDivider(varName, ['_', '-'])
  let transferVar = ''
  if (style === 'constant') {
    const upperVarNameSlice = varNameSlice.map((item: string) =>
      item.toUpperCase()
    )
    transferVar = upperVarNameSlice.join('_')
  } else if (style === 'camel') {
    varNameSlice.forEach((item: string, index: number) => {
      const strArr = item.split('')
      if (strArr[0] && index !== 0) strArr[0] = strArr[0].toUpperCase()
      transferVar += strArr.join('')
    })
  } else if (style === 'pascal') {
    varNameSlice.forEach((item: string) => {
      const strArr = item.split('')
      if (strArr[0]) strArr[0] = strArr[0].toUpperCase()
      transferVar += strArr.join('')
    })
  } else if (style === 'hyphen') {
    transferVar = varNameSlice.join('-')
  } else if (style === 'snake') {
    transferVar = varNameSlice.join('_')
  } else if (style === 'underscore') {
    transferVar = '_'
    transferVar += varNameSlice.join('_')
  }

  return transferVar
}

/**
 * @param str need cut string
 * @param cutter cut option like '-' or ['_', '-']
 * @returns string | string[]
 */
export const variableNameDivider = (
  str: string,
  cutter: string[]
): string[] => {
  const result: string[] = []
  let pushStr = ''
  const strSplit = str.split('')
  strSplit.forEach((item: string) => {
    if (cutter.includes(item)) {
      if (pushStr) {
        result.push(pushStr)
      }
      pushStr = ''
    } else if (item.toUpperCase() === item && Number.isNaN(Number(item))) {
      pushStr && result.push(pushStr)
      pushStr = item.toLowerCase()
    } else {
      pushStr += item.toLowerCase()
    }
  })
  if (pushStr) result.push(pushStr)
  return result
}
