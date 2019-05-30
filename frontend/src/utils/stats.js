export default {
  sendPv (page) {
    window._hmt.push(['_trackPageview', page])
  },
  sendEv (ev) {
    window._hmt.push(['_trackCustomEvent', ev])
  }
}
