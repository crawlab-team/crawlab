export const isValidUsername = (str: string): boolean => {
  if (!str) return false;
  if (str.length > 100) return false;
  return true;
  // const validMap = ['admin', 'editor']
  // return validMap.indexOf(str.trim()) >= 0
};

export const isExternal = (path: string) => {
  return /^(https?:|mailto:|tel:)/.test(path);
};
