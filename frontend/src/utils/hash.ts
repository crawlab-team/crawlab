import md5 from 'md5';

export const getMd5 = (text: string): string => {
  return md5(text).toString();
};
