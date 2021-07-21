export const sendPv = (page: any) => {
  if (localStorage.getItem('useStats') !== '0') {
    window._hmt?.push(['_trackPageview', page]);
  }
};

export const sendEv = (category: string, eventName: string, optLabel: string, optValue: string) => {
  if (localStorage.getItem('useStats') !== '0') {
    window._hmt?.push(['_trackEvent', category, eventName, optLabel, optValue]);
  }
};
