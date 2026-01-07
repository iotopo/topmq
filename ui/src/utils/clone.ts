import Clone from 'clone'

/**
 * 深拷贝
 * @param x
 * @param errOrDef
 * @returns
 */
// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function clone<T>(x: T, errOrDef = true): T {
  return Clone(x)
}
