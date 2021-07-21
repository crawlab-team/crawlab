import getos from 'getos';
import {OS_LINUX, OS_MAC, OS_WINDOWS} from '@/constants/os';

let os: OS;

getos((e, _os) => {
  if (e) {
    console.error(e);
    return;
  }

  switch (_os.os) {
    case 'win32':
      return OS_WINDOWS;
    case 'darwin':
      return OS_MAC;
    default:
      return OS_LINUX;
  }
});

export const getOS = (): OS => {
  return os;
};

export const isWindows = (): boolean => {
  return getOS() === OS_WINDOWS;
};

export const getOSPathSeparator = () => {
  return isWindows() ? '\\' : '/';
};
