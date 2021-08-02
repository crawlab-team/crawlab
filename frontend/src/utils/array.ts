export const isDuplicated = <T = any>(array?: T[]) => {
  if (!array) return false;
  return array.length > new Set(array).size;
};
