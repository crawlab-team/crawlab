export const getRoutePathByDepth = (path: string, depth?: number) => {
  if (!depth) depth = 1;
  const arr = path.split('/');
  if (!arr[0]) depth += 1;
  return arr.slice(0, depth).join('/');
};
