/**
 * Created by jiachenpan on 16/11/18.
 */

export function isValidUsername (str) {
  if (!str) return false
  if (str.length > 100) return false
  return true
  // const validMap = ['admin', 'editor']
  // return validMap.indexOf(str.trim()) >= 0
}

export function isExternal (path) {
  return /^(https?:|mailto:|tel:)/.test(path)
}
