// baidu tongji
export const initBaiduTonji = () => {
  if (localStorage.getItem('useStats') !== '0') {
    window._hmt = window._hmt || [];
    (function () {
      const hm = document.createElement('script');
      hm.src = 'https://hm.baidu.com/hm.js?c35e3a563a06caee2524902c81975add';
      const s = document.getElementsByTagName('script')[0];
      s?.parentNode?.insertBefore(hm, s);
    })();
  }
};
