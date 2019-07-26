import Vue from 'vue'

import 'normalize.css/normalize.css' // A modern alternative to CSS resets

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import locale from 'element-ui/lib/locale/lang/en' // lang i18n

import '@/styles/index.scss' // global css

import 'font-awesome/scss/font-awesome.scss'// FontAwesome
import { library } from '@fortawesome/fontawesome-svg-core'
import { fab } from '@fortawesome/free-brands-svg-icons'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { FontAwesomeIcon, FontAwesomeLayers, FontAwesomeLayersText } from '@fortawesome/vue-fontawesome'

import 'codemirror/lib/codemirror.js'
import { codemirror } from 'vue-codemirror-lite'
import 'codemirror/lib/codemirror.css'
import 'codemirror/theme/darcula.css'

import App from './App'
import store from './store'
import router from './router'

import '@/icons' // icon
import '@/permission' // permission control

import request from './api/request'
import i18n from './i18n'
import utils from './utils'

// code mirror
Vue.use(codemirror)

// element-ui
Vue.use(ElementUI, { locale })

// font-awesome
library.add(fab)
library.add(far)
library.add(fas)
Vue.component('font-awesome-icon', FontAwesomeIcon)
Vue.component('font-awesome-layers', FontAwesomeLayers)
Vue.component('font-awesome-layers-text', FontAwesomeLayersText)

Vue.config.productionTip = false

// 百度统计
if (localStorage.getItem('useStats') !== '0') {
  window._hmt = window._hmt || [];
  (function () {
    let hm = document.createElement('script')
    hm.src = 'https://hm.baidu.com/hm.js?c35e3a563a06caee2524902c81975add'
    let s = document.getElementsByTagName('script')[0]
    s.parentNode.insertBefore(hm, s)
  })()
}

// inject request api
Vue.prototype.$request = request

// inject utils
Vue.prototype.$utils = utils

// inject stats
Vue.prototype.$st = utils.stats

const app = new Vue({
  el: '#app',
  i18n,
  router,
  store,
  render: h => h(App)
})
export default app
