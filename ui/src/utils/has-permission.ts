/**
 * 判断权限
 * @param permissions 带判断权限
 * @param mode 模式 some:任意 every:所有
 * @returns
 */
export default function hasPermission(
  permissions: string[],
  storePermissions: string | string[] | null,
  mode: 'every' | 'some'
): boolean {
  // 未传入需判断的权限时 直接通过
  if (!permissions || permissions.length === 0) {
    return true
  }

  if (storePermissions == '*') {
    //所有权限标识
    return true
  }
  if (mode == 'every') {
    return permissions.every((e) => (storePermissions || []).includes(e))
  }
  if (mode == 'some') {
    return permissions.some((e) => (storePermissions || []).includes(e))
  }
  return false
}
