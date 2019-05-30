export default {
  sendPv (page) {
    if (localStorage.getItem('useStats') !== '0') {
      window._hmt.push(['_trackPageview', page])
    }
  },
  sendEv (category, eventName, optLabel, optValue) {
    if (localStorage.getItem('useStats') !== '0') {
      window._hmt.push(['_trackEvent', category, eventName, optLabel, optValue])
    }
  }
}
