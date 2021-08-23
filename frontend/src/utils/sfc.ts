const vue = window['Vue'];
const {loadModule: sfcLoadModule} = window['vue3-sfc-loader'];

export const getLoadModuleOptions = () => {
  return {
    moduleCache: {
      vue,
    },
    pathResolve({ refPath, relPath }: {refPath?: string; relPath?: string}) {
      // self
      if ( relPath === '.' ) {
        return refPath;
      }

      // relPath is a module name ?
      if ( relPath?.[0] !== '.' && relPath?.[0] !== '/' ) {
        return relPath;
      }

      return String(new URL(relPath, refPath === undefined ? window.location.toString() : refPath));
    },
    async getFile(url: string) {
      const res = await fetch(url);
      if (!res.ok) {
        throw Object.assign(new Error(res.statusText + ' ' + url), {res});
      }
      return {
        getContentData: (asBinary: boolean) => asBinary ? res.arrayBuffer() : res.text(),
      };
    },
    addStyle(textContent: string) {
      const style = Object.assign(document.createElement('style'), {textContent});
      const ref = document.head.getElementsByTagName('style')[0] || null;
      document.head.insertBefore(style, ref);
    },
  };
};

export const loadModule = sfcLoadModule;
