import * as vue from '@vue/runtime-dom';
import {getRequestBaseUrl} from '@/utils/request';
import useRequest from '@/services/request';

const {loadModule: sfcLoadModule} = window['vue3-sfc-loader'];

const {
  getRaw,
} = useRequest();

const getLoadModuleOptions = (): any => {
  return {
    moduleCache: {
      vue,
    },
    pathResolve({refPath, relPath}: { refPath?: string; relPath?: string }) {
      // self
      if (relPath === '.') {
        return refPath;
      }

      // relPath is a module name ?
      if (relPath?.toString()?.[0] !== '.' && relPath?.toString()?.[0] !== '/') {
        return relPath;
      }

      return String(new URL(relPath.toString(), refPath === undefined ? window.location.toString() : refPath.toString()));
    },
    async getFile(url: string) {
      const res = await getRaw(url.toString());
      return {
        getContentData: async (_: boolean) => res.data,
      };
    },
    addStyle(textContent: string) {
      const style = Object.assign(document.createElement('style'), {textContent});
      const ref = document.head.getElementsByTagName('style')[0] || null;
      document.head.insertBefore(style, ref);
    },
  };
};

export const loadModule = (path: string) => sfcLoadModule(`${getRequestBaseUrl()}${path}`, getLoadModuleOptions());
