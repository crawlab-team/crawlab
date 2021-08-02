import {createApp} from 'vue';
import ElementPlus from 'element-plus';
import App from '@/App.vue';
import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';
import {FontAwesomeIcon} from '@fortawesome/vue-fontawesome';
import {library} from '@fortawesome/fontawesome-svg-core';
import {fab} from '@fortawesome/free-brands-svg-icons';
import {far} from '@fortawesome/free-regular-svg-icons';
import {fas} from '@fortawesome/free-solid-svg-icons';
import 'normalize.css/normalize.css';
import 'font-awesome/css/font-awesome.min.css';
import 'element-plus/lib/theme-chalk/index.css';
import '@/styles/index.scss';
import {initBaiduTonji} from '@/admin/baidu';

library.add(fab, far, fas);

// baidu tongji
initBaiduTonji();

// remove loading placeholder
document.querySelector('#loading-placeholder')?.remove();

createApp(App)
  .use(store)
  .use(router)
  .use(ElementPlus)
  .use(i18n)
  .component('font-awesome-icon', FontAwesomeIcon)
  .mount('#app');
