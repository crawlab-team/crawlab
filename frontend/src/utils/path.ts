export const getPrimaryPath = (path: string): string => {
  const arr = path.split('/');
  if (arr.length <= 1) {
    return path;
  } else {
    return `/${arr[1]}`;
  }
};
