export const plainClone = <T = any>(obj: T): T => {
  if (obj === undefined || obj === null) return obj;
  return JSON.parse(JSON.stringify(obj));
};

export const cloneArray = <T = any>(arr: T[]): T[] => {
  return Array.from(arr);
};
