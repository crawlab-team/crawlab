import Vue from 'vue'

import 'normalize.css/normalize.css' // A modern alternative to CSS resets

import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import locale from 'element-ui/lib/locale/lang/en' // lang i18n

import '@/styles/index.scss' // global css

import 'font-awesome/scss/font-awesome.scss'// FontAwesome

import 'codemirror/lib/codemirror.css'

import App from './App'
import store from './store'
import router from './router'

import '@/icons' // icon
import '@/permission' // permission control

import request from './api/request'

Vue.use(ElementUI, { locale })

Vue.config.productionTip = false

Vue.prototype.$request = request

const app = new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
export default app
